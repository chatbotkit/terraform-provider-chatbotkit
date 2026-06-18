package provider

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

// maxDirectUploadBytes is the largest payload uploaded directly through the API
// upload endpoints. Anything larger uses the two-stage presigned flow. The
// value stays comfortably under the platform's ~4.5MB request body limit.
const maxDirectUploadBytes = 4 * 1024 * 1024

// restBaseURL derives the REST API base (".../v1") from the configured GraphQL
// endpoint (".../graphql"). Both are served from the same host, so we simply
// swap the trailing path segment.
func (c *Client) restBaseURL() string {
	base := strings.TrimSuffix(c.BaseURL, "/")
	base = strings.TrimSuffix(base, "/graphql")
	return base + "/v1"
}

// fileUploadEndpoint returns the upload URL for a file's content.
func (c *Client) fileUploadEndpoint(fileID string) string {
	return fmt.Sprintf("%s/file/%s/upload", c.restBaseURL(), url.PathEscape(fileID))
}

// spaceStorageUploadEndpoint returns the upload URL for an object in space
// storage. The path may contain multiple segments; each is escaped while the
// slashes separating them are preserved.
func (c *Client) spaceStorageUploadEndpoint(spaceID, path string) string {
	return fmt.Sprintf(
		"%s/space/%s/storage/upload/%s",
		c.restBaseURL(),
		url.PathEscape(spaceID),
		escapeStoragePath(path),
	)
}

// spaceStorageDeleteEndpoint returns the delete URL for an object in space
// storage.
func (c *Client) spaceStorageDeleteEndpoint(spaceID, path string) string {
	return fmt.Sprintf(
		"%s/space/%s/storage/delete/%s",
		c.restBaseURL(),
		url.PathEscape(spaceID),
		escapeStoragePath(path),
	)
}

// escapeStoragePath escapes each segment of a storage path while keeping the
// "/" separators intact.
func escapeStoragePath(path string) string {
	segments := strings.Split(strings.TrimPrefix(path, "/"), "/")
	for i, segment := range segments {
		segments[i] = url.PathEscape(segment)
	}
	return strings.Join(segments, "/")
}

// uploadRequest mirrors the presigned upload descriptor returned by the API for
// direct-to-storage uploads of large files.
type uploadRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

// uploadResponse is the JSON body returned by the upload endpoints.
type uploadResponse struct {
	ID            *string        `json:"id"`
	Path          *string        `json:"path"`
	UploadRequest *uploadRequest `json:"uploadRequest"`
}

// sha256Hex returns the lowercase hex-encoded SHA-256 digest of data.
func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// detectContentType guesses the MIME type of content when the caller did not
// specify one. The file name/path extension is preferred (it reliably
// distinguishes text formats such as CSV, Markdown, and JSON), falling back to
// content sniffing when the extension is unknown or unavailable.
func detectContentType(name string, data []byte) string {
	if ext := filepath.Ext(name); ext != "" {
		if t := mime.TypeByExtension(ext); t != "" {
			return t
		}
	}
	return http.DetectContentType(data)
}

// uploadContent uploads raw bytes to a REST upload endpoint, transparently
// switching to the two-stage presigned flow for payloads larger than the direct
// upload limit. name is optional and only used as a hint for large uploads.
func (c *Client) uploadContent(ctx context.Context, endpoint string, data []byte, contentType, name string) error {
	if contentType == "" {
		contentType = detectContentType(name, data)
	}

	if len(data) <= maxDirectUploadBytes {
		_, err := c.doUploadRaw(ctx, endpoint, data, contentType)
		return err
	}

	// Two-stage presigned upload for large files: ask the API for upload
	// credentials, then PUT the bytes directly to storage.
	file := map[string]interface{}{
		"type": contentType,
		"size": len(data),
	}
	if name != "" {
		file["name"] = name
	}

	resp, err := c.doUploadJSON(ctx, endpoint, map[string]interface{}{"file": file})
	if err != nil {
		return err
	}

	if resp.UploadRequest == nil {
		return fmt.Errorf("expected a presigned upload request for large content but none was returned")
	}

	return c.doPresignedUpload(ctx, resp.UploadRequest, data)
}

// uploadContentFromURL instructs the API to fetch and store content from an
// HTTP or data URL. The platform performs the fetch server-side, so large
// remote files never transit the provider.
func (c *Client) uploadContentFromURL(ctx context.Context, endpoint, fileURL string) error {
	_, err := c.doUploadJSON(ctx, endpoint, map[string]interface{}{"file": fileURL})
	return err
}

// deleteSpaceStorageFile removes an object from space storage.
func (c *Client) deleteSpaceStorageFile(ctx context.Context, spaceID, path string) error {
	endpoint := c.spaceStorageDeleteEndpoint(spaceID, path)

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	return c.doStorageRequest(req, nil)
}

// doUploadRaw performs a raw-stream upload (the request body is the content).
func (c *Client) doUploadRaw(ctx context.Context, endpoint string, data []byte, contentType string) (*uploadResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	var result uploadResponse
	if err := c.doStorageRequest(req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// doUploadJSON performs an application/json upload request (URL, data URL, or
// presigned metadata form).
func (c *Client) doUploadJSON(ctx context.Context, endpoint string, payload map[string]interface{}) (*uploadResponse, error) {
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	var result uploadResponse
	if err := c.doStorageRequest(req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// doPresignedUpload uploads data directly to storage using credentials returned
// by the API. The request is signed via the provided URL, so it carries no API
// key — only the headers the API instructed us to send.
func (c *Client) doPresignedUpload(ctx context.Context, ur *uploadRequest, data []byte) error {
	method := ur.Method
	if method == "" {
		method = "PUT"
	}

	req, err := http.NewRequestWithContext(ctx, method, ur.URL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create presigned request: %w", err)
	}

	for key, value := range ur.Headers {
		req.Header.Set(key, value)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute presigned upload: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("presigned upload failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return nil
}

// doStorageRequest executes a REST request and, when result is non-nil, decodes
// the JSON response body into it. Non-2xx responses are surfaced as errors.
func (c *Client) doStorageRequest(req *http.Request, result *uploadResponse) error {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

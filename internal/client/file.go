package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// File represents a ChatBotKit file
type File struct {
	ID        string                 `json:"id,omitempty"`
	Name      string                 `json:"name"`
	Type      string                 `json:"type,omitempty"`
	Source    string                 `json:"source,omitempty"`
	Meta      map[string]interface{} `json:"meta,omitempty"`
	CreatedAt int64                  `json:"createdAt,omitempty"`
	UpdatedAt int64                  `json:"updatedAt,omitempty"`
}

// FileListResponse represents the response from listing files
type FileListResponse struct {
	Items  []File `json:"items"`
	Cursor string `json:"cursor,omitempty"`
}

// CreateFile creates a new file
func (c *Client) CreateFile(ctx context.Context, file *File) (*File, error) {
	data, err := c.doRequest(ctx, "POST", "/file/upload", file)
	if err != nil {
		return nil, err
	}

	var result File
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetFile retrieves a file by ID
func (c *Client) GetFile(ctx context.Context, id string) (*File, error) {
	data, err := c.doRequest(ctx, "GET", fmt.Sprintf("/file/%s/fetch", id), nil)
	if err != nil {
		return nil, err
	}

	var result File
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateFile updates an existing file
func (c *Client) UpdateFile(ctx context.Context, id string, file *File) (*File, error) {
	data, err := c.doRequest(ctx, "POST", fmt.Sprintf("/file/%s/update", id), file)
	if err != nil {
		return nil, err
	}

	var result File
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteFile deletes a file
func (c *Client) DeleteFile(ctx context.Context, id string) error {
	_, err := c.doRequest(ctx, "POST", fmt.Sprintf("/file/%s/delete", id), nil)
	return err
}

// ListFiles lists all files
func (c *Client) ListFiles(ctx context.Context, cursor string) (*FileListResponse, error) {
	path := "/file/list"
	if cursor != "" {
		path = fmt.Sprintf("%s?cursor=%s", path, cursor)
	}

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result FileListResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

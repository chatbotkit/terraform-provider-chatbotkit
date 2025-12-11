package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// Secret represents a ChatBotKit secret
type Secret struct {
	ID        string                 `json:"id,omitempty"`
	Name      string                 `json:"name"`
	Value     string                 `json:"value,omitempty"`
	Meta      map[string]interface{} `json:"meta,omitempty"`
	CreatedAt int64                  `json:"createdAt,omitempty"`
	UpdatedAt int64                  `json:"updatedAt,omitempty"`
}

// SecretListResponse represents the response from listing secrets
type SecretListResponse struct {
	Items  []Secret `json:"items"`
	Cursor string   `json:"cursor,omitempty"`
}

// CreateSecret creates a new secret
func (c *Client) CreateSecret(ctx context.Context, secret *Secret) (*Secret, error) {
	data, err := c.doRequest(ctx, "POST", "/secret/create", secret)
	if err != nil {
		return nil, err
	}

	var result Secret
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetSecret retrieves a secret by ID
func (c *Client) GetSecret(ctx context.Context, id string) (*Secret, error) {
	data, err := c.doRequest(ctx, "GET", fmt.Sprintf("/secret/%s/fetch", id), nil)
	if err != nil {
		return nil, err
	}

	var result Secret
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateSecret updates an existing secret
func (c *Client) UpdateSecret(ctx context.Context, id string, secret *Secret) (*Secret, error) {
	data, err := c.doRequest(ctx, "POST", fmt.Sprintf("/secret/%s/update", id), secret)
	if err != nil {
		return nil, err
	}

	var result Secret
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteSecret deletes a secret
func (c *Client) DeleteSecret(ctx context.Context, id string) error {
	_, err := c.doRequest(ctx, "POST", fmt.Sprintf("/secret/%s/delete", id), nil)
	return err
}

// ListSecrets lists all secrets
func (c *Client) ListSecrets(ctx context.Context, cursor string) (*SecretListResponse, error) {
	path := "/secret/list"
	if cursor != "" {
		path = fmt.Sprintf("%s?cursor=%s", path, cursor)
	}

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result SecretListResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

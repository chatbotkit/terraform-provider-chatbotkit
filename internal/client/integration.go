package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// Integration represents a ChatBotKit integration
type Integration struct {
	ID          string                 `json:"id,omitempty"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Type        string                 `json:"type,omitempty"`
	BotID       string                 `json:"botId,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	CreatedAt   int64                  `json:"createdAt,omitempty"`
	UpdatedAt   int64                  `json:"updatedAt,omitempty"`
}

// IntegrationListResponse represents the response from listing integrations
type IntegrationListResponse struct {
	Items  []Integration `json:"items"`
	Cursor string        `json:"cursor,omitempty"`
}

// CreateIntegration creates a new integration
func (c *Client) CreateIntegration(ctx context.Context, integration *Integration) (*Integration, error) {
	data, err := c.doRequest(ctx, "POST", "/integration/create", integration)
	if err != nil {
		return nil, err
	}

	var result Integration
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetIntegration retrieves an integration by ID
func (c *Client) GetIntegration(ctx context.Context, id string) (*Integration, error) {
	data, err := c.doRequest(ctx, "GET", fmt.Sprintf("/integration/%s/fetch", id), nil)
	if err != nil {
		return nil, err
	}

	var result Integration
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateIntegration updates an existing integration
func (c *Client) UpdateIntegration(ctx context.Context, id string, integration *Integration) (*Integration, error) {
	data, err := c.doRequest(ctx, "POST", fmt.Sprintf("/integration/%s/update", id), integration)
	if err != nil {
		return nil, err
	}

	var result Integration
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteIntegration deletes an integration
func (c *Client) DeleteIntegration(ctx context.Context, id string) error {
	_, err := c.doRequest(ctx, "POST", fmt.Sprintf("/integration/%s/delete", id), nil)
	return err
}

// ListIntegrations lists all integrations
func (c *Client) ListIntegrations(ctx context.Context, cursor string) (*IntegrationListResponse, error) {
	path := "/integration/list"
	if cursor != "" {
		path = fmt.Sprintf("%s?cursor=%s", path, cursor)
	}

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result IntegrationListResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

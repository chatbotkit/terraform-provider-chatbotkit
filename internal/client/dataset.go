package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// Dataset represents a ChatBotKit dataset
type Dataset struct {
	ID          string                 `json:"id,omitempty"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	CreatedAt   int64                  `json:"createdAt,omitempty"`
	UpdatedAt   int64                  `json:"updatedAt,omitempty"`
}

// DatasetListResponse represents the response from listing datasets
type DatasetListResponse struct {
	Items  []Dataset `json:"items"`
	Cursor string    `json:"cursor,omitempty"`
}

// CreateDataset creates a new dataset
func (c *Client) CreateDataset(ctx context.Context, dataset *Dataset) (*Dataset, error) {
	data, err := c.doRequest(ctx, "POST", "/dataset/create", dataset)
	if err != nil {
		return nil, err
	}

	var result Dataset
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetDataset retrieves a dataset by ID
func (c *Client) GetDataset(ctx context.Context, id string) (*Dataset, error) {
	data, err := c.doRequest(ctx, "GET", fmt.Sprintf("/dataset/%s/fetch", id), nil)
	if err != nil {
		return nil, err
	}

	var result Dataset
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateDataset updates an existing dataset
func (c *Client) UpdateDataset(ctx context.Context, id string, dataset *Dataset) (*Dataset, error) {
	data, err := c.doRequest(ctx, "POST", fmt.Sprintf("/dataset/%s/update", id), dataset)
	if err != nil {
		return nil, err
	}

	var result Dataset
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteDataset deletes a dataset
func (c *Client) DeleteDataset(ctx context.Context, id string) error {
	_, err := c.doRequest(ctx, "POST", fmt.Sprintf("/dataset/%s/delete", id), nil)
	return err
}

// ListDatasets lists all datasets
func (c *Client) ListDatasets(ctx context.Context, cursor string) (*DatasetListResponse, error) {
	path := "/dataset/list"
	if cursor != "" {
		path = fmt.Sprintf("%s?cursor=%s", path, cursor)
	}

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result DatasetListResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

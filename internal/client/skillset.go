package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// Skillset represents a ChatBotKit skillset
type Skillset struct {
	ID          string                 `json:"id,omitempty"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	CreatedAt   int64                  `json:"createdAt,omitempty"`
	UpdatedAt   int64                  `json:"updatedAt,omitempty"`
}

// SkillsetListResponse represents the response from listing skillsets
type SkillsetListResponse struct {
	Items  []Skillset `json:"items"`
	Cursor string     `json:"cursor,omitempty"`
}

// CreateSkillset creates a new skillset
func (c *Client) CreateSkillset(ctx context.Context, skillset *Skillset) (*Skillset, error) {
	data, err := c.doRequest(ctx, "POST", "/skillset/create", skillset)
	if err != nil {
		return nil, err
	}

	var result Skillset
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetSkillset retrieves a skillset by ID
func (c *Client) GetSkillset(ctx context.Context, id string) (*Skillset, error) {
	data, err := c.doRequest(ctx, "GET", fmt.Sprintf("/skillset/%s/fetch", id), nil)
	if err != nil {
		return nil, err
	}

	var result Skillset
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateSkillset updates an existing skillset
func (c *Client) UpdateSkillset(ctx context.Context, id string, skillset *Skillset) (*Skillset, error) {
	data, err := c.doRequest(ctx, "POST", fmt.Sprintf("/skillset/%s/update", id), skillset)
	if err != nil {
		return nil, err
	}

	var result Skillset
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteSkillset deletes a skillset
func (c *Client) DeleteSkillset(ctx context.Context, id string) error {
	_, err := c.doRequest(ctx, "POST", fmt.Sprintf("/skillset/%s/delete", id), nil)
	return err
}

// ListSkillsets lists all skillsets
func (c *Client) ListSkillsets(ctx context.Context, cursor string) (*SkillsetListResponse, error) {
	path := "/skillset/list"
	if cursor != "" {
		path = fmt.Sprintf("%s?cursor=%s", path, cursor)
	}

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result SkillsetListResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

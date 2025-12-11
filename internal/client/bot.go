package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// Bot represents a ChatBotKit bot
type Bot struct {
	ID              string                 `json:"id,omitempty"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description,omitempty"`
	Model           string                 `json:"model,omitempty"`
	DatasetID       string                 `json:"datasetId,omitempty"`
	SkillsetID      string                 `json:"skillsetId,omitempty"`
	Backstory       string                 `json:"backstory,omitempty"`
	Temperature     float64                `json:"temperature,omitempty"`
	Instructions    string                 `json:"instructions,omitempty"`
	Moderation      bool                   `json:"moderation,omitempty"`
	Privacy         bool                   `json:"privacy,omitempty"`
	Meta            map[string]interface{} `json:"meta,omitempty"`
	CreatedAt       int64                  `json:"createdAt,omitempty"`
	UpdatedAt       int64                  `json:"updatedAt,omitempty"`
}

// BotListResponse represents the response from listing bots
type BotListResponse struct {
	Items  []Bot  `json:"items"`
	Cursor string `json:"cursor,omitempty"`
}

// CreateBot creates a new bot
func (c *Client) CreateBot(ctx context.Context, bot *Bot) (*Bot, error) {
	data, err := c.doRequest(ctx, "POST", "/bot/create", bot)
	if err != nil {
		return nil, err
	}

	var result Bot
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// GetBot retrieves a bot by ID
func (c *Client) GetBot(ctx context.Context, id string) (*Bot, error) {
	data, err := c.doRequest(ctx, "GET", fmt.Sprintf("/bot/%s/fetch", id), nil)
	if err != nil {
		return nil, err
	}

	var result Bot
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// UpdateBot updates an existing bot
func (c *Client) UpdateBot(ctx context.Context, id string, bot *Bot) (*Bot, error) {
	data, err := c.doRequest(ctx, "POST", fmt.Sprintf("/bot/%s/update", id), bot)
	if err != nil {
		return nil, err
	}

	var result Bot
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

// DeleteBot deletes a bot
func (c *Client) DeleteBot(ctx context.Context, id string) error {
	_, err := c.doRequest(ctx, "POST", fmt.Sprintf("/bot/%s/delete", id), nil)
	return err
}

// ListBots lists all bots
func (c *Client) ListBots(ctx context.Context, cursor string) (*BotListResponse, error) {
	path := "/bot/list"
	if cursor != "" {
		path = fmt.Sprintf("%s?cursor=%s", path, cursor)
	}

	data, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result BotListResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

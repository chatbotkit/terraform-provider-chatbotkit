package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

const defaultBaseURL = "https://api.chatbotkit.com/graphql"

// Client represents the ChatBotKit GraphQL API client.
type Client struct {
	APIKey     string
	BaseURL    string
	RunAs      string
	HTTPClient *http.Client
}

// NewClient creates a new ChatBotKit API client.
func NewClient(apiKey, baseURL string) *Client {
	if apiKey == "" {
		apiKey = os.Getenv("CHATBOTKIT_API_KEY")
	}
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &Client{
		APIKey:     apiKey,
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

// GraphQLRequest represents a GraphQL request.
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// GraphQLResponse represents a GraphQL response.
type GraphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

// doRequest executes a GraphQL request.
func (c *Client) doRequest(ctx context.Context, query string, variables map[string]interface{}, result interface{}) error {
	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	// When set, act on behalf of a sub-account (partner user). This lets a single
	// api_key manage many sub-accounts by selecting one per provider configuration.
	if c.RunAs != "" {
		req.Header.Set("X-RunAs-UserId", c.RunAs)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var gqlResp GraphQLResponse
	if err := json.Unmarshal(respBody, &gqlResp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(gqlResp.Errors) > 0 {
		return fmt.Errorf("GraphQL error: %v", gqlResp.Errors[0].Message)
	}

	if result != nil {
		if err := json.Unmarshal(gqlResp.Data, result); err != nil {
			return fmt.Errorf("failed to unmarshal data: %w", err)
		}
	}

	return nil
}

// convertMapToInterface converts types.Map to map[string]interface{}.
func convertMapToInterface(ctx context.Context, m types.Map) map[string]interface{} {
	if m.IsNull() || m.IsUnknown() {
		return nil
	}
	result := make(map[string]interface{})
	m.ElementsAs(ctx, &result, false)
	return result
}

// CreateBlueprintInput represents the input for creating a blueprint.
type CreateBlueprintInput struct {
	Alias       *string                `json:"alias,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Visibility  *string                `json:"visibility,omitempty"`
}

// CreateBlueprintResponse represents the response from creating a blueprint.
type CreateBlueprintResponse struct {
	ID *string `json:"id"`
}

// CreateBlueprint creates a new blueprint.
func (c *Client) CreateBlueprint(ctx context.Context, input CreateBlueprintInput) (*CreateBlueprintResponse, error) {
	query := `
		mutation CreateBlueprint($input: BlueprintCreateRequest!) {
			createBlueprint(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateBlueprint *CreateBlueprintResponse `json:"createBlueprint"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateBlueprint, nil
}

// UpdateBlueprintInput represents the input for updating a blueprint.
type UpdateBlueprintInput struct {
	Alias       *string                `json:"alias,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Visibility  *string                `json:"visibility,omitempty"`
}

// UpdateBlueprintResponse represents the response from updating a blueprint.
type UpdateBlueprintResponse struct {
	ID *string `json:"id"`
}

// UpdateBlueprint updates an existing blueprint.
func (c *Client) UpdateBlueprint(ctx context.Context, id string, input UpdateBlueprintInput) (*UpdateBlueprintResponse, error) {
	query := `
		mutation UpdateBlueprint($blueprintId: ID!, $input: BlueprintUpdateRequest!) {
			updateBlueprint(blueprintId: $blueprintId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"blueprintId": id,
		"input":       input,
	}

	var response struct {
		UpdateBlueprint *UpdateBlueprintResponse `json:"updateBlueprint"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateBlueprint, nil
}

// DeleteBlueprintResponse represents the response from deleting a blueprint.
type DeleteBlueprintResponse struct {
	ID *string `json:"id"`
}

// DeleteBlueprint deletes a blueprint.
func (c *Client) DeleteBlueprint(ctx context.Context, id string) (*DeleteBlueprintResponse, error) {
	query := `
		mutation DeleteBlueprint($blueprintId: ID!) {
			deleteBlueprint(blueprintId: $blueprintId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"blueprintId": id,
	}

	var response struct {
		DeleteBlueprint *DeleteBlueprintResponse `json:"deleteBlueprint"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteBlueprint, nil
}

// GetBlueprintResponse represents the response from fetching a blueprint.
type GetBlueprintResponse struct {
	ID          *string                `json:"id"`
	Alias       *string                `json:"alias,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Visibility  *string                `json:"visibility,omitempty"`
	CreatedAt   *string                `json:"createdAt,omitempty"`
	UpdatedAt   *string                `json:"updatedAt,omitempty"`
}

// GetBlueprint fetches a blueprint by ID.
func (c *Client) GetBlueprint(ctx context.Context, id string) (*GetBlueprintResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetBlueprint($cursor: ID) {
			blueprints(first: 1, after: $cursor) {
				edges {
					node {
						id
						alias
						description
						meta
						name
						visibility
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		Blueprints struct {
			Edges []struct {
				Node *GetBlueprintResponse `json:"node"`
			} `json:"edges"`
		} `json:"blueprints"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.Blueprints.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("blueprint with ID %s not found", id)
}

// CreateBotInput represents the input for creating a bot.
type CreateBotInput struct {
	Alias       *string                `json:"alias,omitempty"`
	Backstory   *string                `json:"backstory,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	DatasetId   *string                `json:"datasetId,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Model       *string                `json:"model,omitempty"`
	Moderation  *bool                  `json:"moderation,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Privacy     *bool                  `json:"privacy,omitempty"`
	SkillsetId  *string                `json:"skillsetId,omitempty"`
	Visibility  *string                `json:"visibility,omitempty"`
}

// CreateBotResponse represents the response from creating a bot.
type CreateBotResponse struct {
	ID *string `json:"id"`
}

// CreateBot creates a new bot.
func (c *Client) CreateBot(ctx context.Context, input CreateBotInput) (*CreateBotResponse, error) {
	query := `
		mutation CreateBot($input: BotCreateRequest!) {
			createBot(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateBot *CreateBotResponse `json:"createBot"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateBot, nil
}

// UpdateBotInput represents the input for updating a bot.
type UpdateBotInput struct {
	Alias       *string                `json:"alias,omitempty"`
	Backstory   *string                `json:"backstory,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	DatasetId   *string                `json:"datasetId,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Model       *string                `json:"model,omitempty"`
	Moderation  *bool                  `json:"moderation,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Privacy     *bool                  `json:"privacy,omitempty"`
	SkillsetId  *string                `json:"skillsetId,omitempty"`
	Visibility  *string                `json:"visibility,omitempty"`
}

// UpdateBotResponse represents the response from updating a bot.
type UpdateBotResponse struct {
	ID *string `json:"id"`
}

// UpdateBot updates an existing bot.
func (c *Client) UpdateBot(ctx context.Context, id string, input UpdateBotInput) (*UpdateBotResponse, error) {
	query := `
		mutation UpdateBot($botId: ID!, $input: BotUpdateRequest!) {
			updateBot(botId: $botId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"botId": id,
		"input": input,
	}

	var response struct {
		UpdateBot *UpdateBotResponse `json:"updateBot"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateBot, nil
}

// DeleteBotResponse represents the response from deleting a bot.
type DeleteBotResponse struct {
	ID *string `json:"id"`
}

// DeleteBot deletes a bot.
func (c *Client) DeleteBot(ctx context.Context, id string) (*DeleteBotResponse, error) {
	query := `
		mutation DeleteBot($botId: ID!) {
			deleteBot(botId: $botId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"botId": id,
	}

	var response struct {
		DeleteBot *DeleteBotResponse `json:"deleteBot"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteBot, nil
}

// GetBotResponse represents the response from fetching a bot.
type GetBotResponse struct {
	ID          *string                `json:"id"`
	Alias       *string                `json:"alias,omitempty"`
	Backstory   *string                `json:"backstory,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	DatasetId   *string                `json:"datasetId,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Model       *string                `json:"model,omitempty"`
	Moderation  *bool                  `json:"moderation,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Privacy     *bool                  `json:"privacy,omitempty"`
	SkillsetId  *string                `json:"skillsetId,omitempty"`
	Visibility  *string                `json:"visibility,omitempty"`
	CreatedAt   *string                `json:"createdAt,omitempty"`
	UpdatedAt   *string                `json:"updatedAt,omitempty"`
}

// GetBot fetches a bot by ID.
func (c *Client) GetBot(ctx context.Context, id string) (*GetBotResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetBot($cursor: ID) {
			bots(first: 1, after: $cursor) {
				edges {
					node {
						id
						alias
						backstory
						blueprintId
						datasetId
						description
						meta
						model
						moderation
						name
						privacy
						skillsetId
						visibility
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		Bots struct {
			Edges []struct {
				Node *GetBotResponse `json:"node"`
			} `json:"edges"`
		} `json:"bots"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.Bots.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("bot with ID %s not found", id)
}

// CreateContextInput represents the input for creating a context.
type CreateContextInput struct {
	BlueprintId *string                `json:"blueprintId,omitempty"`
	BotId       *string                `json:"botId,omitempty"`
	ContactId   *string                `json:"contactId,omitempty"`
	DatasetId   *string                `json:"datasetId,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Payload     map[string]interface{} `json:"payload,omitempty"`
	SkillsetId  *string                `json:"skillsetId,omitempty"`
}

// CreateContextResponse represents the response from creating a context.
type CreateContextResponse struct {
	ID *string `json:"id"`
}

// CreateContext creates a new context.
func (c *Client) CreateContext(ctx context.Context, input CreateContextInput) (*CreateContextResponse, error) {
	query := `
		mutation CreateContext($input: ContextCreateRequest!) {
			createContext(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateContext *CreateContextResponse `json:"createContext"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateContext, nil
}

// UpdateContextInput represents the input for updating a context.
type UpdateContextInput struct {
	BlueprintId *string                `json:"blueprintId,omitempty"`
	BotId       *string                `json:"botId,omitempty"`
	ContactId   *string                `json:"contactId,omitempty"`
	DatasetId   *string                `json:"datasetId,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Payload     map[string]interface{} `json:"payload,omitempty"`
	SkillsetId  *string                `json:"skillsetId,omitempty"`
}

// UpdateContextResponse represents the response from updating a context.
type UpdateContextResponse struct {
	ID *string `json:"id"`
}

// UpdateContext updates an existing context.
func (c *Client) UpdateContext(ctx context.Context, id string, input UpdateContextInput) (*UpdateContextResponse, error) {
	query := `
		mutation UpdateContext($contextId: ID!, $input: ContextUpdateRequest!) {
			updateContext(contextId: $contextId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"contextId": id,
		"input":     input,
	}

	var response struct {
		UpdateContext *UpdateContextResponse `json:"updateContext"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateContext, nil
}

// DeleteContextResponse represents the response from deleting a context.
type DeleteContextResponse struct {
	ID *string `json:"id"`
}

// DeleteContext deletes a context.
func (c *Client) DeleteContext(ctx context.Context, id string) (*DeleteContextResponse, error) {
	query := `
		mutation DeleteContext($contextId: ID!) {
			deleteContext(contextId: $contextId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"contextId": id,
	}

	var response struct {
		DeleteContext *DeleteContextResponse `json:"deleteContext"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteContext, nil
}

// GetContextResponse represents the response from fetching a context.
type GetContextResponse struct {
	ID          *string                `json:"id"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	BotId       *string                `json:"botId,omitempty"`
	ContactId   *string                `json:"contactId,omitempty"`
	DatasetId   *string                `json:"datasetId,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Payload     map[string]interface{} `json:"payload,omitempty"`
	SkillsetId  *string                `json:"skillsetId,omitempty"`
	CreatedAt   *string                `json:"createdAt,omitempty"`
	UpdatedAt   *string                `json:"updatedAt,omitempty"`
}

// GetContext fetches a context by ID.
func (c *Client) GetContext(ctx context.Context, id string) (*GetContextResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetContext($cursor: ID) {
			contexts(first: 1, after: $cursor) {
				edges {
					node {
						id
						blueprintId
						botId
						contactId
						datasetId
						description
						meta
						name
						payload
						skillsetId
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		Contexts struct {
			Edges []struct {
				Node *GetContextResponse `json:"node"`
			} `json:"edges"`
		} `json:"contexts"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.Contexts.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("context with ID %s not found", id)
}

// CreateDatasetInput represents the input for creating a dataset.
type CreateDatasetInput struct {
	Alias               *string                `json:"alias,omitempty"`
	BlueprintId         *string                `json:"blueprintId,omitempty"`
	Description         *string                `json:"description,omitempty"`
	MatchInstruction    *string                `json:"matchInstruction,omitempty"`
	Meta                map[string]interface{} `json:"meta,omitempty"`
	MismatchInstruction *string                `json:"mismatchInstruction,omitempty"`
	Name                *string                `json:"name,omitempty"`
	RecordMaxTokens     *int64                 `json:"recordMaxTokens,omitempty"`
	Reranker            *string                `json:"reranker,omitempty"`
	SearchMaxRecords    *int64                 `json:"searchMaxRecords,omitempty"`
	SearchMaxTokens     *int64                 `json:"searchMaxTokens,omitempty"`
	SearchMinScore      *float64               `json:"searchMinScore,omitempty"`
	Separators          *string                `json:"separators,omitempty"`
	Store               *string                `json:"store,omitempty"`
	Visibility          *string                `json:"visibility,omitempty"`
}

// CreateDatasetResponse represents the response from creating a dataset.
type CreateDatasetResponse struct {
	ID *string `json:"id"`
}

// CreateDataset creates a new dataset.
func (c *Client) CreateDataset(ctx context.Context, input CreateDatasetInput) (*CreateDatasetResponse, error) {
	query := `
		mutation CreateDataset($input: DatasetCreateRequest!) {
			createDataset(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateDataset *CreateDatasetResponse `json:"createDataset"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateDataset, nil
}

// UpdateDatasetInput represents the input for updating a dataset.
type UpdateDatasetInput struct {
	Alias               *string                `json:"alias,omitempty"`
	BlueprintId         *string                `json:"blueprintId,omitempty"`
	Description         *string                `json:"description,omitempty"`
	MatchInstruction    *string                `json:"matchInstruction,omitempty"`
	Meta                map[string]interface{} `json:"meta,omitempty"`
	MismatchInstruction *string                `json:"mismatchInstruction,omitempty"`
	Name                *string                `json:"name,omitempty"`
	RecordMaxTokens     *int64                 `json:"recordMaxTokens,omitempty"`
	Reranker            *string                `json:"reranker,omitempty"`
	SearchMaxRecords    *int64                 `json:"searchMaxRecords,omitempty"`
	SearchMaxTokens     *int64                 `json:"searchMaxTokens,omitempty"`
	SearchMinScore      *float64               `json:"searchMinScore,omitempty"`
	Separators          *string                `json:"separators,omitempty"`
	Store               *string                `json:"store,omitempty"`
	Visibility          *string                `json:"visibility,omitempty"`
}

// UpdateDatasetResponse represents the response from updating a dataset.
type UpdateDatasetResponse struct {
	ID *string `json:"id"`
}

// UpdateDataset updates an existing dataset.
func (c *Client) UpdateDataset(ctx context.Context, id string, input UpdateDatasetInput) (*UpdateDatasetResponse, error) {
	query := `
		mutation UpdateDataset($datasetId: ID!, $input: DatasetUpdateRequest!) {
			updateDataset(datasetId: $datasetId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"datasetId": id,
		"input":     input,
	}

	var response struct {
		UpdateDataset *UpdateDatasetResponse `json:"updateDataset"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateDataset, nil
}

// DeleteDatasetResponse represents the response from deleting a dataset.
type DeleteDatasetResponse struct {
	ID *string `json:"id"`
}

// DeleteDataset deletes a dataset.
func (c *Client) DeleteDataset(ctx context.Context, id string) (*DeleteDatasetResponse, error) {
	query := `
		mutation DeleteDataset($datasetId: ID!) {
			deleteDataset(datasetId: $datasetId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"datasetId": id,
	}

	var response struct {
		DeleteDataset *DeleteDatasetResponse `json:"deleteDataset"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteDataset, nil
}

// GetDatasetResponse represents the response from fetching a dataset.
type GetDatasetResponse struct {
	ID                  *string                `json:"id"`
	Alias               *string                `json:"alias,omitempty"`
	BlueprintId         *string                `json:"blueprintId,omitempty"`
	Description         *string                `json:"description,omitempty"`
	MatchInstruction    *string                `json:"matchInstruction,omitempty"`
	Meta                map[string]interface{} `json:"meta,omitempty"`
	MismatchInstruction *string                `json:"mismatchInstruction,omitempty"`
	Name                *string                `json:"name,omitempty"`
	RecordMaxTokens     *int64                 `json:"recordMaxTokens,omitempty"`
	Reranker            *string                `json:"reranker,omitempty"`
	SearchMaxRecords    *int64                 `json:"searchMaxRecords,omitempty"`
	SearchMaxTokens     *int64                 `json:"searchMaxTokens,omitempty"`
	SearchMinScore      *float64               `json:"searchMinScore,omitempty"`
	Separators          *string                `json:"separators,omitempty"`
	Store               *string                `json:"store,omitempty"`
	Visibility          *string                `json:"visibility,omitempty"`
	CreatedAt           *string                `json:"createdAt,omitempty"`
	UpdatedAt           *string                `json:"updatedAt,omitempty"`
}

// GetDataset fetches a dataset by ID.
func (c *Client) GetDataset(ctx context.Context, id string) (*GetDatasetResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetDataset($cursor: ID) {
			datasets(first: 1, after: $cursor) {
				edges {
					node {
						id
						alias
						blueprintId
						description
						matchInstruction
						meta
						mismatchInstruction
						name
						recordMaxTokens
						reranker
						searchMaxRecords
						searchMaxTokens
						searchMinScore
						separators
						store
						visibility
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		Datasets struct {
			Edges []struct {
				Node *GetDatasetResponse `json:"node"`
			} `json:"edges"`
		} `json:"datasets"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.Datasets.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("dataset with ID %s not found", id)
}

// CreateDiscordIntegrationInput represents the input for creating a discordintegration.
type CreateDiscordIntegrationInput struct {
	Alias             *string                `json:"alias,omitempty"`
	AllowFrom         *string                `json:"allowFrom,omitempty"`
	AppId             *string                `json:"appId,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	BotToken          *string                `json:"botToken,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Handle            *string                `json:"handle,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	PublicKey         *string                `json:"publicKey,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
}

// CreateDiscordIntegrationResponse represents the response from creating a discordintegration.
type CreateDiscordIntegrationResponse struct {
	ID *string `json:"id"`
}

// CreateDiscordIntegration creates a new discordintegration.
func (c *Client) CreateDiscordIntegration(ctx context.Context, input CreateDiscordIntegrationInput) (*CreateDiscordIntegrationResponse, error) {
	query := `
		mutation CreateDiscordIntegration($input: DiscordIntegrationCreateRequest!) {
			createDiscordIntegration(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateDiscordIntegration *CreateDiscordIntegrationResponse `json:"createDiscordIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateDiscordIntegration, nil
}

// UpdateDiscordIntegrationInput represents the input for updating a discordintegration.
type UpdateDiscordIntegrationInput struct {
	Alias             *string                `json:"alias,omitempty"`
	AllowFrom         *string                `json:"allowFrom,omitempty"`
	AppId             *string                `json:"appId,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	BotToken          *string                `json:"botToken,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Handle            *string                `json:"handle,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	PublicKey         *string                `json:"publicKey,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
}

// UpdateDiscordIntegrationResponse represents the response from updating a discordintegration.
type UpdateDiscordIntegrationResponse struct {
	ID *string `json:"id"`
}

// UpdateDiscordIntegration updates an existing discordintegration.
func (c *Client) UpdateDiscordIntegration(ctx context.Context, id string, input UpdateDiscordIntegrationInput) (*UpdateDiscordIntegrationResponse, error) {
	query := `
		mutation UpdateDiscordIntegration($discordIntegrationId: ID!, $input: DiscordIntegrationUpdateRequest!) {
			updateDiscordIntegration(discordIntegrationId: $discordIntegrationId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"discordIntegrationId": id,
		"input":                input,
	}

	var response struct {
		UpdateDiscordIntegration *UpdateDiscordIntegrationResponse `json:"updateDiscordIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateDiscordIntegration, nil
}

// DeleteDiscordIntegrationResponse represents the response from deleting a discordintegration.
type DeleteDiscordIntegrationResponse struct {
	ID *string `json:"id"`
}

// DeleteDiscordIntegration deletes a discordintegration.
func (c *Client) DeleteDiscordIntegration(ctx context.Context, id string) (*DeleteDiscordIntegrationResponse, error) {
	query := `
		mutation DeleteDiscordIntegration($discordIntegrationId: ID!) {
			deleteDiscordIntegration(discordIntegrationId: $discordIntegrationId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"discordIntegrationId": id,
	}

	var response struct {
		DeleteDiscordIntegration *DeleteDiscordIntegrationResponse `json:"deleteDiscordIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteDiscordIntegration, nil
}

// GetDiscordIntegrationResponse represents the response from fetching a discordintegration.
type GetDiscordIntegrationResponse struct {
	ID                *string                `json:"id"`
	Alias             *string                `json:"alias,omitempty"`
	AllowFrom         *string                `json:"allowFrom,omitempty"`
	AppId             *string                `json:"appId,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	BotToken          *string                `json:"botToken,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Handle            *string                `json:"handle,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	PublicKey         *string                `json:"publicKey,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
	CreatedAt         *string                `json:"createdAt,omitempty"`
	UpdatedAt         *string                `json:"updatedAt,omitempty"`
}

// GetDiscordIntegration fetches a discordintegration by ID.
func (c *Client) GetDiscordIntegration(ctx context.Context, id string) (*GetDiscordIntegrationResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetDiscordIntegration($cursor: ID) {
			discordIntegrations(first: 1, after: $cursor) {
				edges {
					node {
						id
						alias
						allowFrom
						appId
						blueprintId
						botId
						botToken
						contactCollection
						description
						handle
						meta
						name
						publicKey
						sessionDuration
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		DiscordIntegrations struct {
			Edges []struct {
				Node *GetDiscordIntegrationResponse `json:"node"`
			} `json:"edges"`
		} `json:"discordIntegrations"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.DiscordIntegrations.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("discordintegration with ID %s not found", id)
}

// CreateEmailIntegrationInput represents the input for creating a emailintegration.
type CreateEmailIntegrationInput struct {
	Alias             *string                `json:"alias,omitempty"`
	AllowFrom         *string                `json:"allowFrom,omitempty"`
	Attachments       *bool                  `json:"attachments,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
}

// CreateEmailIntegrationResponse represents the response from creating a emailintegration.
type CreateEmailIntegrationResponse struct {
	ID *string `json:"id"`
}

// CreateEmailIntegration creates a new emailintegration.
func (c *Client) CreateEmailIntegration(ctx context.Context, input CreateEmailIntegrationInput) (*CreateEmailIntegrationResponse, error) {
	query := `
		mutation CreateEmailIntegration($input: EmailIntegrationCreateRequest!) {
			createEmailIntegration(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateEmailIntegration *CreateEmailIntegrationResponse `json:"createEmailIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateEmailIntegration, nil
}

// UpdateEmailIntegrationInput represents the input for updating a emailintegration.
type UpdateEmailIntegrationInput struct {
	Alias             *string                `json:"alias,omitempty"`
	AllowFrom         *string                `json:"allowFrom,omitempty"`
	Attachments       *bool                  `json:"attachments,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
}

// UpdateEmailIntegrationResponse represents the response from updating a emailintegration.
type UpdateEmailIntegrationResponse struct {
	ID *string `json:"id"`
}

// UpdateEmailIntegration updates an existing emailintegration.
func (c *Client) UpdateEmailIntegration(ctx context.Context, id string, input UpdateEmailIntegrationInput) (*UpdateEmailIntegrationResponse, error) {
	query := `
		mutation UpdateEmailIntegration($emailIntegrationId: ID!, $input: EmailIntegrationUpdateRequest!) {
			updateEmailIntegration(emailIntegrationId: $emailIntegrationId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"emailIntegrationId": id,
		"input":              input,
	}

	var response struct {
		UpdateEmailIntegration *UpdateEmailIntegrationResponse `json:"updateEmailIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateEmailIntegration, nil
}

// DeleteEmailIntegrationResponse represents the response from deleting a emailintegration.
type DeleteEmailIntegrationResponse struct {
	ID *string `json:"id"`
}

// DeleteEmailIntegration deletes a emailintegration.
func (c *Client) DeleteEmailIntegration(ctx context.Context, id string) (*DeleteEmailIntegrationResponse, error) {
	query := `
		mutation DeleteEmailIntegration($emailIntegrationId: ID!) {
			deleteEmailIntegration(emailIntegrationId: $emailIntegrationId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"emailIntegrationId": id,
	}

	var response struct {
		DeleteEmailIntegration *DeleteEmailIntegrationResponse `json:"deleteEmailIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteEmailIntegration, nil
}

// GetEmailIntegrationResponse represents the response from fetching a emailintegration.
type GetEmailIntegrationResponse struct {
	ID                *string                `json:"id"`
	Alias             *string                `json:"alias,omitempty"`
	AllowFrom         *string                `json:"allowFrom,omitempty"`
	Attachments       *bool                  `json:"attachments,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
	CreatedAt         *string                `json:"createdAt,omitempty"`
	UpdatedAt         *string                `json:"updatedAt,omitempty"`
}

// GetEmailIntegration fetches a emailintegration by ID.
func (c *Client) GetEmailIntegration(ctx context.Context, id string) (*GetEmailIntegrationResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetEmailIntegration($cursor: ID) {
			emailIntegrations(first: 1, after: $cursor) {
				edges {
					node {
						id
						alias
						allowFrom
						attachments
						blueprintId
						botId
						contactCollection
						description
						meta
						name
						sessionDuration
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		EmailIntegrations struct {
			Edges []struct {
				Node *GetEmailIntegrationResponse `json:"node"`
			} `json:"edges"`
		} `json:"emailIntegrations"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.EmailIntegrations.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("emailintegration with ID %s not found", id)
}

// CreateExtractIntegrationInput represents the input for creating a extractintegration.
type CreateExtractIntegrationInput struct {
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	BotId       *string                `json:"botId,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Model       *string                `json:"model,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Request     *string                `json:"request,omitempty"`
	Schema      map[string]interface{} `json:"schema,omitempty"`
}

// CreateExtractIntegrationResponse represents the response from creating a extractintegration.
type CreateExtractIntegrationResponse struct {
	ID *string `json:"id"`
}

// CreateExtractIntegration creates a new extractintegration.
func (c *Client) CreateExtractIntegration(ctx context.Context, input CreateExtractIntegrationInput) (*CreateExtractIntegrationResponse, error) {
	query := `
		mutation CreateExtractIntegration($input: ExtractIntegrationCreateRequest!) {
			createExtractIntegration(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateExtractIntegration *CreateExtractIntegrationResponse `json:"createExtractIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateExtractIntegration, nil
}

// UpdateExtractIntegrationInput represents the input for updating a extractintegration.
type UpdateExtractIntegrationInput struct {
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	BotId       *string                `json:"botId,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Model       *string                `json:"model,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Request     *string                `json:"request,omitempty"`
	Schema      map[string]interface{} `json:"schema,omitempty"`
}

// UpdateExtractIntegrationResponse represents the response from updating a extractintegration.
type UpdateExtractIntegrationResponse struct {
	ID *string `json:"id"`
}

// UpdateExtractIntegration updates an existing extractintegration.
func (c *Client) UpdateExtractIntegration(ctx context.Context, id string, input UpdateExtractIntegrationInput) (*UpdateExtractIntegrationResponse, error) {
	query := `
		mutation UpdateExtractIntegration($extractIntegrationId: ID!, $input: ExtractIntegrationUpdateRequest!) {
			updateExtractIntegration(extractIntegrationId: $extractIntegrationId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"extractIntegrationId": id,
		"input":                input,
	}

	var response struct {
		UpdateExtractIntegration *UpdateExtractIntegrationResponse `json:"updateExtractIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateExtractIntegration, nil
}

// DeleteExtractIntegrationResponse represents the response from deleting a extractintegration.
type DeleteExtractIntegrationResponse struct {
	ID *string `json:"id"`
}

// DeleteExtractIntegration deletes a extractintegration.
func (c *Client) DeleteExtractIntegration(ctx context.Context, id string) (*DeleteExtractIntegrationResponse, error) {
	query := `
		mutation DeleteExtractIntegration($extractIntegrationId: ID!) {
			deleteExtractIntegration(extractIntegrationId: $extractIntegrationId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"extractIntegrationId": id,
	}

	var response struct {
		DeleteExtractIntegration *DeleteExtractIntegrationResponse `json:"deleteExtractIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteExtractIntegration, nil
}

// GetExtractIntegrationResponse represents the response from fetching a extractintegration.
type GetExtractIntegrationResponse struct {
	ID          *string                `json:"id"`
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	BotId       *string                `json:"botId,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Model       *string                `json:"model,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Request     *string                `json:"request,omitempty"`
	Schema      map[string]interface{} `json:"schema,omitempty"`
	CreatedAt   *string                `json:"createdAt,omitempty"`
	UpdatedAt   *string                `json:"updatedAt,omitempty"`
}

// GetExtractIntegration fetches a extractintegration by ID.
func (c *Client) GetExtractIntegration(ctx context.Context, id string) (*GetExtractIntegrationResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetExtractIntegration($cursor: ID) {
			extractIntegrations(first: 1, after: $cursor) {
				edges {
					node {
						id
						alias
						blueprintId
						botId
						description
						meta
						model
						name
						request
						schema
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		ExtractIntegrations struct {
			Edges []struct {
				Node *GetExtractIntegrationResponse `json:"node"`
			} `json:"edges"`
		} `json:"extractIntegrations"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.ExtractIntegrations.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("extractintegration with ID %s not found", id)
}

// CreateFileInput represents the input for creating a file.
type CreateFileInput struct {
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Visibility  *string                `json:"visibility,omitempty"`
}

// CreateFileResponse represents the response from creating a file.
type CreateFileResponse struct {
	ID *string `json:"id"`
}

// CreateFile creates a new file.
func (c *Client) CreateFile(ctx context.Context, input CreateFileInput) (*CreateFileResponse, error) {
	query := `
		mutation CreateFile($input: FileCreateRequest!) {
			createFile(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateFile *CreateFileResponse `json:"createFile"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateFile, nil
}

// UpdateFileInput represents the input for updating a file.
type UpdateFileInput struct {
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Visibility  *string                `json:"visibility,omitempty"`
}

// UpdateFileResponse represents the response from updating a file.
type UpdateFileResponse struct {
	ID *string `json:"id"`
}

// UpdateFile updates an existing file.
func (c *Client) UpdateFile(ctx context.Context, id string, input UpdateFileInput) (*UpdateFileResponse, error) {
	query := `
		mutation UpdateFile($fileId: ID!, $input: FileUpdateRequest!) {
			updateFile(fileId: $fileId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"fileId": id,
		"input":  input,
	}

	var response struct {
		UpdateFile *UpdateFileResponse `json:"updateFile"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateFile, nil
}

// DeleteFileResponse represents the response from deleting a file.
type DeleteFileResponse struct {
	ID *string `json:"id"`
}

// DeleteFile deletes a file.
func (c *Client) DeleteFile(ctx context.Context, id string) (*DeleteFileResponse, error) {
	query := `
		mutation DeleteFile($fileId: ID!) {
			deleteFile(fileId: $fileId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"fileId": id,
	}

	var response struct {
		DeleteFile *DeleteFileResponse `json:"deleteFile"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteFile, nil
}

// GetFileResponse represents the response from fetching a file.
type GetFileResponse struct {
	ID          *string                `json:"id"`
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Visibility  *string                `json:"visibility,omitempty"`
	CreatedAt   *string                `json:"createdAt,omitempty"`
	UpdatedAt   *string                `json:"updatedAt,omitempty"`
}

// GetFile fetches a file by ID.
func (c *Client) GetFile(ctx context.Context, id string) (*GetFileResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetFile($cursor: ID) {
			files(first: 1, after: $cursor) {
				edges {
					node {
						id
						alias
						blueprintId
						description
						meta
						name
						visibility
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		Files struct {
			Edges []struct {
				Node *GetFileResponse `json:"node"`
			} `json:"edges"`
		} `json:"files"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.Files.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("file with ID %s not found", id)
}

// CreateGooglechatIntegrationInput represents the input for creating a googlechatintegration.
type CreateGooglechatIntegrationInput struct {
	Alias             *string                `json:"alias,omitempty"`
	AllowFrom         *string                `json:"allowFrom,omitempty"`
	Attachments       *bool                  `json:"attachments,omitempty"`
	AutoRespond       *string                `json:"autoRespond,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	ProjectNumber     *string                `json:"projectNumber,omitempty"`
	ServiceAccountKey *string                `json:"serviceAccountKey,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
}

// CreateGooglechatIntegrationResponse represents the response from creating a googlechatintegration.
type CreateGooglechatIntegrationResponse struct {
	ID *string `json:"id"`
}

// CreateGooglechatIntegration creates a new googlechatintegration.
func (c *Client) CreateGooglechatIntegration(ctx context.Context, input CreateGooglechatIntegrationInput) (*CreateGooglechatIntegrationResponse, error) {
	query := `
		mutation CreateGooglechatIntegration($input: GooglechatIntegrationCreateRequest!) {
			createGooglechatIntegration(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateGooglechatIntegration *CreateGooglechatIntegrationResponse `json:"createGooglechatIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateGooglechatIntegration, nil
}

// UpdateGooglechatIntegrationInput represents the input for updating a googlechatintegration.
type UpdateGooglechatIntegrationInput struct {
	Alias             *string                `json:"alias,omitempty"`
	AllowFrom         *string                `json:"allowFrom,omitempty"`
	Attachments       *bool                  `json:"attachments,omitempty"`
	AutoRespond       *string                `json:"autoRespond,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	ProjectNumber     *string                `json:"projectNumber,omitempty"`
	ServiceAccountKey *string                `json:"serviceAccountKey,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
}

// UpdateGooglechatIntegrationResponse represents the response from updating a googlechatintegration.
type UpdateGooglechatIntegrationResponse struct {
	ID *string `json:"id"`
}

// UpdateGooglechatIntegration updates an existing googlechatintegration.
func (c *Client) UpdateGooglechatIntegration(ctx context.Context, id string, input UpdateGooglechatIntegrationInput) (*UpdateGooglechatIntegrationResponse, error) {
	query := `
		mutation UpdateGooglechatIntegration($googlechatIntegrationId: ID!, $input: GooglechatIntegrationUpdateRequest!) {
			updateGooglechatIntegration(googlechatIntegrationId: $googlechatIntegrationId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"googlechatIntegrationId": id,
		"input":                   input,
	}

	var response struct {
		UpdateGooglechatIntegration *UpdateGooglechatIntegrationResponse `json:"updateGooglechatIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateGooglechatIntegration, nil
}

// DeleteGooglechatIntegrationResponse represents the response from deleting a googlechatintegration.
type DeleteGooglechatIntegrationResponse struct {
	ID *string `json:"id"`
}

// DeleteGooglechatIntegration deletes a googlechatintegration.
func (c *Client) DeleteGooglechatIntegration(ctx context.Context, id string) (*DeleteGooglechatIntegrationResponse, error) {
	query := `
		mutation DeleteGooglechatIntegration($googlechatIntegrationId: ID!) {
			deleteGooglechatIntegration(googlechatIntegrationId: $googlechatIntegrationId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"googlechatIntegrationId": id,
	}

	var response struct {
		DeleteGooglechatIntegration *DeleteGooglechatIntegrationResponse `json:"deleteGooglechatIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteGooglechatIntegration, nil
}

// GetGooglechatIntegrationResponse represents the response from fetching a googlechatintegration.
type GetGooglechatIntegrationResponse struct {
	ID                *string                `json:"id"`
	Alias             *string                `json:"alias,omitempty"`
	AllowFrom         *string                `json:"allowFrom,omitempty"`
	Attachments       *bool                  `json:"attachments,omitempty"`
	AutoRespond       *string                `json:"autoRespond,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	ProjectNumber     *string                `json:"projectNumber,omitempty"`
	ServiceAccountKey *string                `json:"serviceAccountKey,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
	CreatedAt         *string                `json:"createdAt,omitempty"`
	UpdatedAt         *string                `json:"updatedAt,omitempty"`
}

// GetGooglechatIntegration fetches a googlechatintegration by ID.
func (c *Client) GetGooglechatIntegration(ctx context.Context, id string) (*GetGooglechatIntegrationResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetGooglechatIntegration($cursor: ID) {
			googlechatIntegrations(first: 1, after: $cursor) {
				edges {
					node {
						id
						alias
						allowFrom
						attachments
						autoRespond
						blueprintId
						botId
						contactCollection
						description
						meta
						name
						projectNumber
						serviceAccountKey
						sessionDuration
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		GooglechatIntegrations struct {
			Edges []struct {
				Node *GetGooglechatIntegrationResponse `json:"node"`
			} `json:"edges"`
		} `json:"googlechatIntegrations"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.GooglechatIntegrations.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("googlechatintegration with ID %s not found", id)
}

// CreateInstagramIntegrationInput represents the input for creating a instagramintegration.
type CreateInstagramIntegrationInput struct {
	AccessToken       *string                `json:"accessToken,omitempty"`
	Alias             *string                `json:"alias,omitempty"`
	Attachments       *bool                  `json:"attachments,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
}

// CreateInstagramIntegrationResponse represents the response from creating a instagramintegration.
type CreateInstagramIntegrationResponse struct {
	ID *string `json:"id"`
}

// CreateInstagramIntegration creates a new instagramintegration.
func (c *Client) CreateInstagramIntegration(ctx context.Context, input CreateInstagramIntegrationInput) (*CreateInstagramIntegrationResponse, error) {
	query := `
		mutation CreateInstagramIntegration($input: InstagramIntegrationCreateRequest!) {
			createInstagramIntegration(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateInstagramIntegration *CreateInstagramIntegrationResponse `json:"createInstagramIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateInstagramIntegration, nil
}

// UpdateInstagramIntegrationInput represents the input for updating a instagramintegration.
type UpdateInstagramIntegrationInput struct {
	AccessToken       *string                `json:"accessToken,omitempty"`
	Alias             *string                `json:"alias,omitempty"`
	Attachments       *bool                  `json:"attachments,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
}

// UpdateInstagramIntegrationResponse represents the response from updating a instagramintegration.
type UpdateInstagramIntegrationResponse struct {
	ID *string `json:"id"`
}

// UpdateInstagramIntegration updates an existing instagramintegration.
func (c *Client) UpdateInstagramIntegration(ctx context.Context, id string, input UpdateInstagramIntegrationInput) (*UpdateInstagramIntegrationResponse, error) {
	query := `
		mutation UpdateInstagramIntegration($instagramIntegrationId: ID!, $input: InstagramIntegrationUpdateRequest!) {
			updateInstagramIntegration(instagramIntegrationId: $instagramIntegrationId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"instagramIntegrationId": id,
		"input":                  input,
	}

	var response struct {
		UpdateInstagramIntegration *UpdateInstagramIntegrationResponse `json:"updateInstagramIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateInstagramIntegration, nil
}

// DeleteInstagramIntegrationResponse represents the response from deleting a instagramintegration.
type DeleteInstagramIntegrationResponse struct {
	ID *string `json:"id"`
}

// DeleteInstagramIntegration deletes a instagramintegration.
func (c *Client) DeleteInstagramIntegration(ctx context.Context, id string) (*DeleteInstagramIntegrationResponse, error) {
	query := `
		mutation DeleteInstagramIntegration($instagramIntegrationId: ID!) {
			deleteInstagramIntegration(instagramIntegrationId: $instagramIntegrationId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"instagramIntegrationId": id,
	}

	var response struct {
		DeleteInstagramIntegration *DeleteInstagramIntegrationResponse `json:"deleteInstagramIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteInstagramIntegration, nil
}

// GetInstagramIntegrationResponse represents the response from fetching a instagramintegration.
type GetInstagramIntegrationResponse struct {
	ID                *string                `json:"id"`
	AccessToken       *string                `json:"accessToken,omitempty"`
	Alias             *string                `json:"alias,omitempty"`
	Attachments       *bool                  `json:"attachments,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
	CreatedAt         *string                `json:"createdAt,omitempty"`
	UpdatedAt         *string                `json:"updatedAt,omitempty"`
}

// GetInstagramIntegration fetches a instagramintegration by ID.
func (c *Client) GetInstagramIntegration(ctx context.Context, id string) (*GetInstagramIntegrationResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetInstagramIntegration($cursor: ID) {
			instagramIntegrations(first: 1, after: $cursor) {
				edges {
					node {
						id
						accessToken
						alias
						attachments
						blueprintId
						botId
						contactCollection
						description
						meta
						name
						sessionDuration
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		InstagramIntegrations struct {
			Edges []struct {
				Node *GetInstagramIntegrationResponse `json:"node"`
			} `json:"edges"`
		} `json:"instagramIntegrations"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.InstagramIntegrations.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("instagramintegration with ID %s not found", id)
}

// CreateMcpserverIntegrationInput represents the input for creating a mcpserverintegration.
type CreateMcpserverIntegrationInput struct {
	Alias             *string                `json:"alias,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	OAuthConnectionId *string                `json:"oAuthConnectionId,omitempty"`
	SkillsetId        *string                `json:"skillsetId,omitempty"`
}

// CreateMcpserverIntegrationResponse represents the response from creating a mcpserverintegration.
type CreateMcpserverIntegrationResponse struct {
	ID *string `json:"id"`
}

// CreateMcpserverIntegration creates a new mcpserverintegration.
func (c *Client) CreateMcpserverIntegration(ctx context.Context, input CreateMcpserverIntegrationInput) (*CreateMcpserverIntegrationResponse, error) {
	query := `
		mutation CreateMcpserverIntegration($input: McpserverIntegrationCreateRequest!) {
			createMcpserverIntegration(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateMcpserverIntegration *CreateMcpserverIntegrationResponse `json:"createMcpserverIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateMcpserverIntegration, nil
}

// UpdateMcpserverIntegrationInput represents the input for updating a mcpserverintegration.
type UpdateMcpserverIntegrationInput struct {
	Alias             *string                `json:"alias,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	OAuthConnectionId *string                `json:"oAuthConnectionId,omitempty"`
	SkillsetId        *string                `json:"skillsetId,omitempty"`
}

// UpdateMcpserverIntegrationResponse represents the response from updating a mcpserverintegration.
type UpdateMcpserverIntegrationResponse struct {
	ID *string `json:"id"`
}

// UpdateMcpserverIntegration updates an existing mcpserverintegration.
func (c *Client) UpdateMcpserverIntegration(ctx context.Context, id string, input UpdateMcpserverIntegrationInput) (*UpdateMcpserverIntegrationResponse, error) {
	query := `
		mutation UpdateMcpserverIntegration($mcpserverIntegrationId: ID!, $input: McpserverIntegrationUpdateRequest!) {
			updateMcpserverIntegration(mcpserverIntegrationId: $mcpserverIntegrationId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"mcpserverIntegrationId": id,
		"input":                  input,
	}

	var response struct {
		UpdateMcpserverIntegration *UpdateMcpserverIntegrationResponse `json:"updateMcpserverIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateMcpserverIntegration, nil
}

// DeleteMcpserverIntegrationResponse represents the response from deleting a mcpserverintegration.
type DeleteMcpserverIntegrationResponse struct {
	ID *string `json:"id"`
}

// DeleteMcpserverIntegration deletes a mcpserverintegration.
func (c *Client) DeleteMcpserverIntegration(ctx context.Context, id string) (*DeleteMcpserverIntegrationResponse, error) {
	query := `
		mutation DeleteMcpserverIntegration($mcpserverIntegrationId: ID!) {
			deleteMcpserverIntegration(mcpserverIntegrationId: $mcpserverIntegrationId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"mcpserverIntegrationId": id,
	}

	var response struct {
		DeleteMcpserverIntegration *DeleteMcpserverIntegrationResponse `json:"deleteMcpserverIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteMcpserverIntegration, nil
}

// GetMcpserverIntegrationResponse represents the response from fetching a mcpserverintegration.
type GetMcpserverIntegrationResponse struct {
	ID                *string                `json:"id"`
	Alias             *string                `json:"alias,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	OAuthConnectionId *string                `json:"oAuthConnectionId,omitempty"`
	SkillsetId        *string                `json:"skillsetId,omitempty"`
	CreatedAt         *string                `json:"createdAt,omitempty"`
	UpdatedAt         *string                `json:"updatedAt,omitempty"`
}

// GetMcpserverIntegration fetches a mcpserverintegration by ID.
func (c *Client) GetMcpserverIntegration(ctx context.Context, id string) (*GetMcpserverIntegrationResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetMcpserverIntegration($cursor: ID) {
			mcpserverIntegrations(first: 1, after: $cursor) {
				edges {
					node {
						id
						alias
						blueprintId
						description
						meta
						name
						oAuthConnectionId
						skillsetId
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		McpserverIntegrations struct {
			Edges []struct {
				Node *GetMcpserverIntegrationResponse `json:"node"`
			} `json:"edges"`
		} `json:"mcpserverIntegrations"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.McpserverIntegrations.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("mcpserverintegration with ID %s not found", id)
}

// CreateMessengerIntegrationInput represents the input for creating a messengerintegration.
type CreateMessengerIntegrationInput struct {
	AccessToken       *string                `json:"accessToken,omitempty"`
	Alias             *string                `json:"alias,omitempty"`
	Attachments       *bool                  `json:"attachments,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
}

// CreateMessengerIntegrationResponse represents the response from creating a messengerintegration.
type CreateMessengerIntegrationResponse struct {
	ID *string `json:"id"`
}

// CreateMessengerIntegration creates a new messengerintegration.
func (c *Client) CreateMessengerIntegration(ctx context.Context, input CreateMessengerIntegrationInput) (*CreateMessengerIntegrationResponse, error) {
	query := `
		mutation CreateMessengerIntegration($input: MessengerIntegrationCreateRequest!) {
			createMessengerIntegration(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateMessengerIntegration *CreateMessengerIntegrationResponse `json:"createMessengerIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateMessengerIntegration, nil
}

// UpdateMessengerIntegrationInput represents the input for updating a messengerintegration.
type UpdateMessengerIntegrationInput struct {
	AccessToken       *string                `json:"accessToken,omitempty"`
	Alias             *string                `json:"alias,omitempty"`
	Attachments       *bool                  `json:"attachments,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
}

// UpdateMessengerIntegrationResponse represents the response from updating a messengerintegration.
type UpdateMessengerIntegrationResponse struct {
	ID *string `json:"id"`
}

// UpdateMessengerIntegration updates an existing messengerintegration.
func (c *Client) UpdateMessengerIntegration(ctx context.Context, id string, input UpdateMessengerIntegrationInput) (*UpdateMessengerIntegrationResponse, error) {
	query := `
		mutation UpdateMessengerIntegration($messengerIntegrationId: ID!, $input: MessengerIntegrationUpdateRequest!) {
			updateMessengerIntegration(messengerIntegrationId: $messengerIntegrationId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"messengerIntegrationId": id,
		"input":                  input,
	}

	var response struct {
		UpdateMessengerIntegration *UpdateMessengerIntegrationResponse `json:"updateMessengerIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateMessengerIntegration, nil
}

// DeleteMessengerIntegrationResponse represents the response from deleting a messengerintegration.
type DeleteMessengerIntegrationResponse struct {
	ID *string `json:"id"`
}

// DeleteMessengerIntegration deletes a messengerintegration.
func (c *Client) DeleteMessengerIntegration(ctx context.Context, id string) (*DeleteMessengerIntegrationResponse, error) {
	query := `
		mutation DeleteMessengerIntegration($messengerIntegrationId: ID!) {
			deleteMessengerIntegration(messengerIntegrationId: $messengerIntegrationId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"messengerIntegrationId": id,
	}

	var response struct {
		DeleteMessengerIntegration *DeleteMessengerIntegrationResponse `json:"deleteMessengerIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteMessengerIntegration, nil
}

// GetMessengerIntegrationResponse represents the response from fetching a messengerintegration.
type GetMessengerIntegrationResponse struct {
	ID                *string                `json:"id"`
	AccessToken       *string                `json:"accessToken,omitempty"`
	Alias             *string                `json:"alias,omitempty"`
	Attachments       *bool                  `json:"attachments,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
	CreatedAt         *string                `json:"createdAt,omitempty"`
	UpdatedAt         *string                `json:"updatedAt,omitempty"`
}

// GetMessengerIntegration fetches a messengerintegration by ID.
func (c *Client) GetMessengerIntegration(ctx context.Context, id string) (*GetMessengerIntegrationResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetMessengerIntegration($cursor: ID) {
			messengerIntegrations(first: 1, after: $cursor) {
				edges {
					node {
						id
						accessToken
						alias
						attachments
						blueprintId
						botId
						contactCollection
						description
						meta
						name
						sessionDuration
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		MessengerIntegrations struct {
			Edges []struct {
				Node *GetMessengerIntegrationResponse `json:"node"`
			} `json:"edges"`
		} `json:"messengerIntegrations"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.MessengerIntegrations.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("messengerintegration with ID %s not found", id)
}

// CreateMicrosoftteamsIntegrationInput represents the input for creating a microsoftteamsintegration.
type CreateMicrosoftteamsIntegrationInput struct {
	Alias                 *string                `json:"alias,omitempty"`
	AllowFrom             *string                `json:"allowFrom,omitempty"`
	BlueprintId           *string                `json:"blueprintId,omitempty"`
	BotFrameworkAppId     *string                `json:"botFrameworkAppId,omitempty"`
	BotFrameworkAppSecret *string                `json:"botFrameworkAppSecret,omitempty"`
	BotId                 *string                `json:"botId,omitempty"`
	ContactCollection     *bool                  `json:"contactCollection,omitempty"`
	Description           *string                `json:"description,omitempty"`
	Meta                  map[string]interface{} `json:"meta,omitempty"`
	Name                  *string                `json:"name,omitempty"`
	SessionDuration       *int64                 `json:"sessionDuration,omitempty"`
	TenantId              *string                `json:"tenantId,omitempty"`
}

// CreateMicrosoftteamsIntegrationResponse represents the response from creating a microsoftteamsintegration.
type CreateMicrosoftteamsIntegrationResponse struct {
	ID *string `json:"id"`
}

// CreateMicrosoftteamsIntegration creates a new microsoftteamsintegration.
func (c *Client) CreateMicrosoftteamsIntegration(ctx context.Context, input CreateMicrosoftteamsIntegrationInput) (*CreateMicrosoftteamsIntegrationResponse, error) {
	query := `
		mutation CreateMicrosoftteamsIntegration($input: MicrosoftteamsIntegrationCreateRequest!) {
			createMicrosoftteamsIntegration(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateMicrosoftteamsIntegration *CreateMicrosoftteamsIntegrationResponse `json:"createMicrosoftteamsIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateMicrosoftteamsIntegration, nil
}

// UpdateMicrosoftteamsIntegrationInput represents the input for updating a microsoftteamsintegration.
type UpdateMicrosoftteamsIntegrationInput struct {
	Alias                 *string                `json:"alias,omitempty"`
	AllowFrom             *string                `json:"allowFrom,omitempty"`
	BlueprintId           *string                `json:"blueprintId,omitempty"`
	BotFrameworkAppId     *string                `json:"botFrameworkAppId,omitempty"`
	BotFrameworkAppSecret *string                `json:"botFrameworkAppSecret,omitempty"`
	BotId                 *string                `json:"botId,omitempty"`
	ContactCollection     *bool                  `json:"contactCollection,omitempty"`
	Description           *string                `json:"description,omitempty"`
	Meta                  map[string]interface{} `json:"meta,omitempty"`
	Name                  *string                `json:"name,omitempty"`
	SessionDuration       *int64                 `json:"sessionDuration,omitempty"`
	TenantId              *string                `json:"tenantId,omitempty"`
}

// UpdateMicrosoftteamsIntegrationResponse represents the response from updating a microsoftteamsintegration.
type UpdateMicrosoftteamsIntegrationResponse struct {
	ID *string `json:"id"`
}

// UpdateMicrosoftteamsIntegration updates an existing microsoftteamsintegration.
func (c *Client) UpdateMicrosoftteamsIntegration(ctx context.Context, id string, input UpdateMicrosoftteamsIntegrationInput) (*UpdateMicrosoftteamsIntegrationResponse, error) {
	query := `
		mutation UpdateMicrosoftteamsIntegration($microsoftteamsIntegrationId: ID!, $input: MicrosoftteamsIntegrationUpdateRequest!) {
			updateMicrosoftteamsIntegration(microsoftteamsIntegrationId: $microsoftteamsIntegrationId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"microsoftteamsIntegrationId": id,
		"input":                       input,
	}

	var response struct {
		UpdateMicrosoftteamsIntegration *UpdateMicrosoftteamsIntegrationResponse `json:"updateMicrosoftteamsIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateMicrosoftteamsIntegration, nil
}

// DeleteMicrosoftteamsIntegrationResponse represents the response from deleting a microsoftteamsintegration.
type DeleteMicrosoftteamsIntegrationResponse struct {
	ID *string `json:"id"`
}

// DeleteMicrosoftteamsIntegration deletes a microsoftteamsintegration.
func (c *Client) DeleteMicrosoftteamsIntegration(ctx context.Context, id string) (*DeleteMicrosoftteamsIntegrationResponse, error) {
	query := `
		mutation DeleteMicrosoftteamsIntegration($microsoftteamsIntegrationId: ID!) {
			deleteMicrosoftteamsIntegration(microsoftteamsIntegrationId: $microsoftteamsIntegrationId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"microsoftteamsIntegrationId": id,
	}

	var response struct {
		DeleteMicrosoftteamsIntegration *DeleteMicrosoftteamsIntegrationResponse `json:"deleteMicrosoftteamsIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteMicrosoftteamsIntegration, nil
}

// GetMicrosoftteamsIntegrationResponse represents the response from fetching a microsoftteamsintegration.
type GetMicrosoftteamsIntegrationResponse struct {
	ID                    *string                `json:"id"`
	Alias                 *string                `json:"alias,omitempty"`
	AllowFrom             *string                `json:"allowFrom,omitempty"`
	BlueprintId           *string                `json:"blueprintId,omitempty"`
	BotFrameworkAppId     *string                `json:"botFrameworkAppId,omitempty"`
	BotFrameworkAppSecret *string                `json:"botFrameworkAppSecret,omitempty"`
	BotId                 *string                `json:"botId,omitempty"`
	ContactCollection     *bool                  `json:"contactCollection,omitempty"`
	Description           *string                `json:"description,omitempty"`
	Meta                  map[string]interface{} `json:"meta,omitempty"`
	Name                  *string                `json:"name,omitempty"`
	SessionDuration       *int64                 `json:"sessionDuration,omitempty"`
	TenantId              *string                `json:"tenantId,omitempty"`
	CreatedAt             *string                `json:"createdAt,omitempty"`
	UpdatedAt             *string                `json:"updatedAt,omitempty"`
}

// GetMicrosoftteamsIntegration fetches a microsoftteamsintegration by ID.
func (c *Client) GetMicrosoftteamsIntegration(ctx context.Context, id string) (*GetMicrosoftteamsIntegrationResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetMicrosoftteamsIntegration($cursor: ID) {
			microsoftteamsIntegrations(first: 1, after: $cursor) {
				edges {
					node {
						id
						alias
						allowFrom
						blueprintId
						botFrameworkAppId
						botFrameworkAppSecret
						botId
						contactCollection
						description
						meta
						name
						sessionDuration
						tenantId
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		MicrosoftteamsIntegrations struct {
			Edges []struct {
				Node *GetMicrosoftteamsIntegrationResponse `json:"node"`
			} `json:"edges"`
		} `json:"microsoftteamsIntegrations"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.MicrosoftteamsIntegrations.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("microsoftteamsintegration with ID %s not found", id)
}

// CreateNotionIntegrationInput represents the input for creating a notionintegration.
type CreateNotionIntegrationInput struct {
	Alias        *string                `json:"alias,omitempty"`
	BlueprintId  *string                `json:"blueprintId,omitempty"`
	DatasetId    *string                `json:"datasetId,omitempty"`
	Description  *string                `json:"description,omitempty"`
	ExpiresIn    *int64                 `json:"expiresIn,omitempty"`
	Meta         map[string]interface{} `json:"meta,omitempty"`
	Name         *string                `json:"name,omitempty"`
	SyncSchedule *string                `json:"syncSchedule,omitempty"`
	Token        *string                `json:"token,omitempty"`
}

// CreateNotionIntegrationResponse represents the response from creating a notionintegration.
type CreateNotionIntegrationResponse struct {
	ID *string `json:"id"`
}

// CreateNotionIntegration creates a new notionintegration.
func (c *Client) CreateNotionIntegration(ctx context.Context, input CreateNotionIntegrationInput) (*CreateNotionIntegrationResponse, error) {
	query := `
		mutation CreateNotionIntegration($input: NotionIntegrationCreateRequest!) {
			createNotionIntegration(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateNotionIntegration *CreateNotionIntegrationResponse `json:"createNotionIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateNotionIntegration, nil
}

// UpdateNotionIntegrationInput represents the input for updating a notionintegration.
type UpdateNotionIntegrationInput struct {
	Alias        *string                `json:"alias,omitempty"`
	BlueprintId  *string                `json:"blueprintId,omitempty"`
	DatasetId    *string                `json:"datasetId,omitempty"`
	Description  *string                `json:"description,omitempty"`
	ExpiresIn    *int64                 `json:"expiresIn,omitempty"`
	Meta         map[string]interface{} `json:"meta,omitempty"`
	Name         *string                `json:"name,omitempty"`
	SyncSchedule *string                `json:"syncSchedule,omitempty"`
	Token        *string                `json:"token,omitempty"`
}

// UpdateNotionIntegrationResponse represents the response from updating a notionintegration.
type UpdateNotionIntegrationResponse struct {
	ID *string `json:"id"`
}

// UpdateNotionIntegration updates an existing notionintegration.
func (c *Client) UpdateNotionIntegration(ctx context.Context, id string, input UpdateNotionIntegrationInput) (*UpdateNotionIntegrationResponse, error) {
	query := `
		mutation UpdateNotionIntegration($notionIntegrationId: ID!, $input: NotionIntegrationUpdateRequest!) {
			updateNotionIntegration(notionIntegrationId: $notionIntegrationId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"notionIntegrationId": id,
		"input":               input,
	}

	var response struct {
		UpdateNotionIntegration *UpdateNotionIntegrationResponse `json:"updateNotionIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateNotionIntegration, nil
}

// DeleteNotionIntegrationResponse represents the response from deleting a notionintegration.
type DeleteNotionIntegrationResponse struct {
	ID *string `json:"id"`
}

// DeleteNotionIntegration deletes a notionintegration.
func (c *Client) DeleteNotionIntegration(ctx context.Context, id string) (*DeleteNotionIntegrationResponse, error) {
	query := `
		mutation DeleteNotionIntegration($notionIntegrationId: ID!) {
			deleteNotionIntegration(notionIntegrationId: $notionIntegrationId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"notionIntegrationId": id,
	}

	var response struct {
		DeleteNotionIntegration *DeleteNotionIntegrationResponse `json:"deleteNotionIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteNotionIntegration, nil
}

// GetNotionIntegrationResponse represents the response from fetching a notionintegration.
type GetNotionIntegrationResponse struct {
	ID           *string                `json:"id"`
	Alias        *string                `json:"alias,omitempty"`
	BlueprintId  *string                `json:"blueprintId,omitempty"`
	DatasetId    *string                `json:"datasetId,omitempty"`
	Description  *string                `json:"description,omitempty"`
	ExpiresIn    *int64                 `json:"expiresIn,omitempty"`
	Meta         map[string]interface{} `json:"meta,omitempty"`
	Name         *string                `json:"name,omitempty"`
	SyncSchedule *string                `json:"syncSchedule,omitempty"`
	Token        *string                `json:"token,omitempty"`
	CreatedAt    *string                `json:"createdAt,omitempty"`
	UpdatedAt    *string                `json:"updatedAt,omitempty"`
}

// GetNotionIntegration fetches a notionintegration by ID.
func (c *Client) GetNotionIntegration(ctx context.Context, id string) (*GetNotionIntegrationResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetNotionIntegration($cursor: ID) {
			notionIntegrations(first: 1, after: $cursor) {
				edges {
					node {
						id
						alias
						blueprintId
						datasetId
						description
						expiresIn
						meta
						name
						syncSchedule
						token
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		NotionIntegrations struct {
			Edges []struct {
				Node *GetNotionIntegrationResponse `json:"node"`
			} `json:"edges"`
		} `json:"notionIntegrations"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.NotionIntegrations.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("notionintegration with ID %s not found", id)
}

// CreatePolicyInput represents the input for creating a policy.
type CreatePolicyInput struct {
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	BotId       *string                `json:"botId,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Type        string                 `json:"type,omitempty"`
}

// CreatePolicyResponse represents the response from creating a policy.
type CreatePolicyResponse struct {
	ID *string `json:"id"`
}

// CreatePolicy creates a new policy.
func (c *Client) CreatePolicy(ctx context.Context, input CreatePolicyInput) (*CreatePolicyResponse, error) {
	query := `
		mutation CreatePolicy($input: PolicyCreateRequest!) {
			createPolicy(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreatePolicy *CreatePolicyResponse `json:"createPolicy"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreatePolicy, nil
}

// UpdatePolicyInput represents the input for updating a policy.
type UpdatePolicyInput struct {
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	BotId       *string                `json:"botId,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Type        string                 `json:"type,omitempty"`
}

// UpdatePolicyResponse represents the response from updating a policy.
type UpdatePolicyResponse struct {
	ID *string `json:"id"`
}

// UpdatePolicy updates an existing policy.
func (c *Client) UpdatePolicy(ctx context.Context, id string, input UpdatePolicyInput) (*UpdatePolicyResponse, error) {
	query := `
		mutation UpdatePolicy($policyId: ID!, $input: PolicyUpdateRequest!) {
			updatePolicy(policyId: $policyId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"policyId": id,
		"input":    input,
	}

	var response struct {
		UpdatePolicy *UpdatePolicyResponse `json:"updatePolicy"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdatePolicy, nil
}

// DeletePolicyResponse represents the response from deleting a policy.
type DeletePolicyResponse struct {
	ID *string `json:"id"`
}

// DeletePolicy deletes a policy.
func (c *Client) DeletePolicy(ctx context.Context, id string) (*DeletePolicyResponse, error) {
	query := `
		mutation DeletePolicy($policyId: ID!) {
			deletePolicy(policyId: $policyId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"policyId": id,
	}

	var response struct {
		DeletePolicy *DeletePolicyResponse `json:"deletePolicy"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeletePolicy, nil
}

// GetPolicyResponse represents the response from fetching a policy.
type GetPolicyResponse struct {
	ID          *string                `json:"id"`
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	BotId       *string                `json:"botId,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Type        string                 `json:"type,omitempty"`
	CreatedAt   *string                `json:"createdAt,omitempty"`
	UpdatedAt   *string                `json:"updatedAt,omitempty"`
}

// GetPolicy fetches a policy by ID.
func (c *Client) GetPolicy(ctx context.Context, id string) (*GetPolicyResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetPolicy($cursor: ID) {
			policies(first: 1, after: $cursor) {
				edges {
					node {
						id
						alias
						blueprintId
						botId
						config
						description
						meta
						name
						type
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		Policys struct {
			Edges []struct {
				Node *GetPolicyResponse `json:"node"`
			} `json:"edges"`
		} `json:"policies"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.Policys.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("policy with ID %s not found", id)
}

// CreatePortalInput represents the input for creating a portal.
type CreatePortalInput struct {
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Slug        *string                `json:"slug,omitempty"`
}

// CreatePortalResponse represents the response from creating a portal.
type CreatePortalResponse struct {
	ID *string `json:"id"`
}

// CreatePortal creates a new portal.
func (c *Client) CreatePortal(ctx context.Context, input CreatePortalInput) (*CreatePortalResponse, error) {
	query := `
		mutation CreatePortal($input: PortalCreateRequest!) {
			createPortal(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreatePortal *CreatePortalResponse `json:"createPortal"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreatePortal, nil
}

// UpdatePortalInput represents the input for updating a portal.
type UpdatePortalInput struct {
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Slug        *string                `json:"slug,omitempty"`
}

// UpdatePortalResponse represents the response from updating a portal.
type UpdatePortalResponse struct {
	ID *string `json:"id"`
}

// UpdatePortal updates an existing portal.
func (c *Client) UpdatePortal(ctx context.Context, id string, input UpdatePortalInput) (*UpdatePortalResponse, error) {
	query := `
		mutation UpdatePortal($portalId: ID!, $input: PortalUpdateRequest!) {
			updatePortal(portalId: $portalId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"portalId": id,
		"input":    input,
	}

	var response struct {
		UpdatePortal *UpdatePortalResponse `json:"updatePortal"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdatePortal, nil
}

// DeletePortalResponse represents the response from deleting a portal.
type DeletePortalResponse struct {
	ID *string `json:"id"`
}

// DeletePortal deletes a portal.
func (c *Client) DeletePortal(ctx context.Context, id string) (*DeletePortalResponse, error) {
	query := `
		mutation DeletePortal($portalId: ID!) {
			deletePortal(portalId: $portalId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"portalId": id,
	}

	var response struct {
		DeletePortal *DeletePortalResponse `json:"deletePortal"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeletePortal, nil
}

// GetPortalResponse represents the response from fetching a portal.
type GetPortalResponse struct {
	ID          *string                `json:"id"`
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Slug        *string                `json:"slug,omitempty"`
	CreatedAt   *string                `json:"createdAt,omitempty"`
	UpdatedAt   *string                `json:"updatedAt,omitempty"`
}

// GetPortal fetches a portal by ID.
func (c *Client) GetPortal(ctx context.Context, id string) (*GetPortalResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetPortal($cursor: ID) {
			portals(first: 1, after: $cursor) {
				edges {
					node {
						id
						alias
						blueprintId
						config
						description
						meta
						name
						slug
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		Portals struct {
			Edges []struct {
				Node *GetPortalResponse `json:"node"`
			} `json:"edges"`
		} `json:"portals"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.Portals.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("portal with ID %s not found", id)
}

// CreateSecretInput represents the input for creating a secret.
type CreateSecretInput struct {
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Description *string                `json:"description,omitempty"`
	Kind        *string                `json:"kind,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Type        *string                `json:"type,omitempty"`
	Value       *string                `json:"value,omitempty"`
	Visibility  *string                `json:"visibility,omitempty"`
}

// CreateSecretResponse represents the response from creating a secret.
type CreateSecretResponse struct {
	ID *string `json:"id"`
}

// CreateSecret creates a new secret.
func (c *Client) CreateSecret(ctx context.Context, input CreateSecretInput) (*CreateSecretResponse, error) {
	query := `
		mutation CreateSecret($input: SecretCreateRequest!) {
			createSecret(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateSecret *CreateSecretResponse `json:"createSecret"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateSecret, nil
}

// UpdateSecretInput represents the input for updating a secret.
type UpdateSecretInput struct {
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Description *string                `json:"description,omitempty"`
	Kind        *string                `json:"kind,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Type        *string                `json:"type,omitempty"`
	Value       *string                `json:"value,omitempty"`
	Visibility  *string                `json:"visibility,omitempty"`
}

// UpdateSecretResponse represents the response from updating a secret.
type UpdateSecretResponse struct {
	ID *string `json:"id"`
}

// UpdateSecret updates an existing secret.
func (c *Client) UpdateSecret(ctx context.Context, id string, input UpdateSecretInput) (*UpdateSecretResponse, error) {
	query := `
		mutation UpdateSecret($secretId: ID!, $input: SecretUpdateRequest!) {
			updateSecret(secretId: $secretId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"secretId": id,
		"input":    input,
	}

	var response struct {
		UpdateSecret *UpdateSecretResponse `json:"updateSecret"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateSecret, nil
}

// DeleteSecretResponse represents the response from deleting a secret.
type DeleteSecretResponse struct {
	ID *string `json:"id"`
}

// DeleteSecret deletes a secret.
func (c *Client) DeleteSecret(ctx context.Context, id string) (*DeleteSecretResponse, error) {
	query := `
		mutation DeleteSecret($secretId: ID!) {
			deleteSecret(secretId: $secretId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"secretId": id,
	}

	var response struct {
		DeleteSecret *DeleteSecretResponse `json:"deleteSecret"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteSecret, nil
}

// GetSecretResponse represents the response from fetching a secret.
type GetSecretResponse struct {
	ID          *string                `json:"id"`
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Description *string                `json:"description,omitempty"`
	Kind        *string                `json:"kind,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Type        *string                `json:"type,omitempty"`
	Value       *string                `json:"value,omitempty"`
	Visibility  *string                `json:"visibility,omitempty"`
	CreatedAt   *string                `json:"createdAt,omitempty"`
	UpdatedAt   *string                `json:"updatedAt,omitempty"`
}

// GetSecret fetches a secret by ID.
func (c *Client) GetSecret(ctx context.Context, id string) (*GetSecretResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetSecret($cursor: ID) {
			secrets(first: 1, after: $cursor) {
				edges {
					node {
						id
						alias
						blueprintId
						config
						description
						kind
						meta
						name
						type
						value
						visibility
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		Secrets struct {
			Edges []struct {
				Node *GetSecretResponse `json:"node"`
			} `json:"edges"`
		} `json:"secrets"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.Secrets.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("secret with ID %s not found", id)
}

// CreateSitemapIntegrationInput represents the input for creating a sitemapintegration.
type CreateSitemapIntegrationInput struct {
	Alias        *string                `json:"alias,omitempty"`
	BlueprintId  *string                `json:"blueprintId,omitempty"`
	DatasetId    *string                `json:"datasetId,omitempty"`
	Description  *string                `json:"description,omitempty"`
	ExpiresIn    *int64                 `json:"expiresIn,omitempty"`
	Glob         *string                `json:"glob,omitempty"`
	Javascript   *bool                  `json:"javascript,omitempty"`
	Meta         map[string]interface{} `json:"meta,omitempty"`
	Name         *string                `json:"name,omitempty"`
	Selectors    *string                `json:"selectors,omitempty"`
	SyncSchedule *string                `json:"syncSchedule,omitempty"`
	URL          *string                `json:"url,omitempty"`
}

// CreateSitemapIntegrationResponse represents the response from creating a sitemapintegration.
type CreateSitemapIntegrationResponse struct {
	ID *string `json:"id"`
}

// CreateSitemapIntegration creates a new sitemapintegration.
func (c *Client) CreateSitemapIntegration(ctx context.Context, input CreateSitemapIntegrationInput) (*CreateSitemapIntegrationResponse, error) {
	query := `
		mutation CreateSitemapIntegration($input: SitemapIntegrationCreateRequest!) {
			createSitemapIntegration(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateSitemapIntegration *CreateSitemapIntegrationResponse `json:"createSitemapIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateSitemapIntegration, nil
}

// UpdateSitemapIntegrationInput represents the input for updating a sitemapintegration.
type UpdateSitemapIntegrationInput struct {
	Alias        *string                `json:"alias,omitempty"`
	BlueprintId  *string                `json:"blueprintId,omitempty"`
	DatasetId    *string                `json:"datasetId,omitempty"`
	Description  *string                `json:"description,omitempty"`
	ExpiresIn    *int64                 `json:"expiresIn,omitempty"`
	Glob         *string                `json:"glob,omitempty"`
	Javascript   *bool                  `json:"javascript,omitempty"`
	Meta         map[string]interface{} `json:"meta,omitempty"`
	Name         *string                `json:"name,omitempty"`
	Selectors    *string                `json:"selectors,omitempty"`
	SyncSchedule *string                `json:"syncSchedule,omitempty"`
	URL          *string                `json:"url,omitempty"`
}

// UpdateSitemapIntegrationResponse represents the response from updating a sitemapintegration.
type UpdateSitemapIntegrationResponse struct {
	ID *string `json:"id"`
}

// UpdateSitemapIntegration updates an existing sitemapintegration.
func (c *Client) UpdateSitemapIntegration(ctx context.Context, id string, input UpdateSitemapIntegrationInput) (*UpdateSitemapIntegrationResponse, error) {
	query := `
		mutation UpdateSitemapIntegration($sitemapIntegrationId: ID!, $input: SitemapIntegrationUpdateRequest!) {
			updateSitemapIntegration(sitemapIntegrationId: $sitemapIntegrationId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"sitemapIntegrationId": id,
		"input":                input,
	}

	var response struct {
		UpdateSitemapIntegration *UpdateSitemapIntegrationResponse `json:"updateSitemapIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateSitemapIntegration, nil
}

// DeleteSitemapIntegrationResponse represents the response from deleting a sitemapintegration.
type DeleteSitemapIntegrationResponse struct {
	ID *string `json:"id"`
}

// DeleteSitemapIntegration deletes a sitemapintegration.
func (c *Client) DeleteSitemapIntegration(ctx context.Context, id string) (*DeleteSitemapIntegrationResponse, error) {
	query := `
		mutation DeleteSitemapIntegration($sitemapIntegrationId: ID!) {
			deleteSitemapIntegration(sitemapIntegrationId: $sitemapIntegrationId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"sitemapIntegrationId": id,
	}

	var response struct {
		DeleteSitemapIntegration *DeleteSitemapIntegrationResponse `json:"deleteSitemapIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteSitemapIntegration, nil
}

// GetSitemapIntegrationResponse represents the response from fetching a sitemapintegration.
type GetSitemapIntegrationResponse struct {
	ID           *string                `json:"id"`
	Alias        *string                `json:"alias,omitempty"`
	BlueprintId  *string                `json:"blueprintId,omitempty"`
	DatasetId    *string                `json:"datasetId,omitempty"`
	Description  *string                `json:"description,omitempty"`
	ExpiresIn    *int64                 `json:"expiresIn,omitempty"`
	Glob         *string                `json:"glob,omitempty"`
	Javascript   *bool                  `json:"javascript,omitempty"`
	Meta         map[string]interface{} `json:"meta,omitempty"`
	Name         *string                `json:"name,omitempty"`
	Selectors    *string                `json:"selectors,omitempty"`
	SyncSchedule *string                `json:"syncSchedule,omitempty"`
	URL          *string                `json:"url,omitempty"`
	CreatedAt    *string                `json:"createdAt,omitempty"`
	UpdatedAt    *string                `json:"updatedAt,omitempty"`
}

// GetSitemapIntegration fetches a sitemapintegration by ID.
func (c *Client) GetSitemapIntegration(ctx context.Context, id string) (*GetSitemapIntegrationResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetSitemapIntegration($cursor: ID) {
			sitemapIntegrations(first: 1, after: $cursor) {
				edges {
					node {
						id
						alias
						blueprintId
						datasetId
						description
						expiresIn
						glob
						javascript
						meta
						name
						selectors
						syncSchedule
						url
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		SitemapIntegrations struct {
			Edges []struct {
				Node *GetSitemapIntegrationResponse `json:"node"`
			} `json:"edges"`
		} `json:"sitemapIntegrations"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.SitemapIntegrations.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("sitemapintegration with ID %s not found", id)
}

// CreateSkillserverIntegrationInput represents the input for creating a skillserverintegration.
type CreateSkillserverIntegrationInput struct {
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	SkillsetId  *string                `json:"skillsetId,omitempty"`
}

// CreateSkillserverIntegrationResponse represents the response from creating a skillserverintegration.
type CreateSkillserverIntegrationResponse struct {
	ID *string `json:"id"`
}

// CreateSkillserverIntegration creates a new skillserverintegration.
func (c *Client) CreateSkillserverIntegration(ctx context.Context, input CreateSkillserverIntegrationInput) (*CreateSkillserverIntegrationResponse, error) {
	query := `
		mutation CreateSkillserverIntegration($input: SkillserverIntegrationCreateRequest!) {
			createSkillserverIntegration(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateSkillserverIntegration *CreateSkillserverIntegrationResponse `json:"createSkillserverIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateSkillserverIntegration, nil
}

// UpdateSkillserverIntegrationInput represents the input for updating a skillserverintegration.
type UpdateSkillserverIntegrationInput struct {
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	SkillsetId  *string                `json:"skillsetId,omitempty"`
}

// UpdateSkillserverIntegrationResponse represents the response from updating a skillserverintegration.
type UpdateSkillserverIntegrationResponse struct {
	ID *string `json:"id"`
}

// UpdateSkillserverIntegration updates an existing skillserverintegration.
func (c *Client) UpdateSkillserverIntegration(ctx context.Context, id string, input UpdateSkillserverIntegrationInput) (*UpdateSkillserverIntegrationResponse, error) {
	query := `
		mutation UpdateSkillserverIntegration($skillserverIntegrationId: ID!, $input: SkillserverIntegrationUpdateRequest!) {
			updateSkillserverIntegration(skillserverIntegrationId: $skillserverIntegrationId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"skillserverIntegrationId": id,
		"input":                    input,
	}

	var response struct {
		UpdateSkillserverIntegration *UpdateSkillserverIntegrationResponse `json:"updateSkillserverIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateSkillserverIntegration, nil
}

// DeleteSkillserverIntegrationResponse represents the response from deleting a skillserverintegration.
type DeleteSkillserverIntegrationResponse struct {
	ID *string `json:"id"`
}

// DeleteSkillserverIntegration deletes a skillserverintegration.
func (c *Client) DeleteSkillserverIntegration(ctx context.Context, id string) (*DeleteSkillserverIntegrationResponse, error) {
	query := `
		mutation DeleteSkillserverIntegration($skillserverIntegrationId: ID!) {
			deleteSkillserverIntegration(skillserverIntegrationId: $skillserverIntegrationId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"skillserverIntegrationId": id,
	}

	var response struct {
		DeleteSkillserverIntegration *DeleteSkillserverIntegrationResponse `json:"deleteSkillserverIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteSkillserverIntegration, nil
}

// GetSkillserverIntegrationResponse represents the response from fetching a skillserverintegration.
type GetSkillserverIntegrationResponse struct {
	ID          *string                `json:"id"`
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	SkillsetId  *string                `json:"skillsetId,omitempty"`
	CreatedAt   *string                `json:"createdAt,omitempty"`
	UpdatedAt   *string                `json:"updatedAt,omitempty"`
}

// GetSkillserverIntegration fetches a skillserverintegration by ID.
func (c *Client) GetSkillserverIntegration(ctx context.Context, id string) (*GetSkillserverIntegrationResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetSkillserverIntegration($cursor: ID) {
			skillserverIntegrations(first: 1, after: $cursor) {
				edges {
					node {
						id
						alias
						blueprintId
						description
						meta
						name
						skillsetId
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		SkillserverIntegrations struct {
			Edges []struct {
				Node *GetSkillserverIntegrationResponse `json:"node"`
			} `json:"edges"`
		} `json:"skillserverIntegrations"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.SkillserverIntegrations.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("skillserverintegration with ID %s not found", id)
}

// CreateSkillsetAbilityInput represents the input for creating a skillsetability.
type CreateSkillsetAbilityInput struct {
	BlueprintId *string                `json:"blueprintId,omitempty"`
	BotId       *string                `json:"botId,omitempty"`
	Description *string                `json:"description,omitempty"`
	FileId      *string                `json:"fileId,omitempty"`
	Instruction *string                `json:"instruction,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	SecretId    *string                `json:"secretId,omitempty"`
	SpaceId     *string                `json:"spaceId,omitempty"`
	State       *string                `json:"state,omitempty"`
}

// CreateSkillsetAbilityResponse represents the response from creating a skillsetability.
type CreateSkillsetAbilityResponse struct {
	ID *string `json:"id"`
}

// CreateSkillsetAbility creates a new skillsetability.
func (c *Client) CreateSkillsetAbility(ctx context.Context, skillsetId string, input CreateSkillsetAbilityInput) (*CreateSkillsetAbilityResponse, error) {
	query := `
		mutation CreateSkillsetAbility($skillsetId: ID!, $input: SkillsetAbilityCreateRequest!) {
			createSkillsetAbility(skillsetId: $skillsetId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"skillsetId": skillsetId,
		"input":      input,
	}

	var response struct {
		CreateSkillsetAbility *CreateSkillsetAbilityResponse `json:"createSkillsetAbility"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateSkillsetAbility, nil
}

// UpdateSkillsetAbilityInput represents the input for updating a skillsetability.
type UpdateSkillsetAbilityInput struct {
	BlueprintId *string                `json:"blueprintId,omitempty"`
	BotId       *string                `json:"botId,omitempty"`
	Description *string                `json:"description,omitempty"`
	FileId      *string                `json:"fileId,omitempty"`
	Instruction *string                `json:"instruction,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	SecretId    *string                `json:"secretId,omitempty"`
	SpaceId     *string                `json:"spaceId,omitempty"`
	State       *string                `json:"state,omitempty"`
}

// UpdateSkillsetAbilityResponse represents the response from updating a skillsetability.
type UpdateSkillsetAbilityResponse struct {
	ID *string `json:"id"`
}

// UpdateSkillsetAbility updates an existing skillsetability.
func (c *Client) UpdateSkillsetAbility(ctx context.Context, skillsetId string, abilityId string, input UpdateSkillsetAbilityInput) (*UpdateSkillsetAbilityResponse, error) {
	query := `
		mutation UpdateSkillsetAbility($skillsetId: ID!, $abilityId: ID!, $input: SkillsetAbilityUpdateRequest!) {
			updateSkillsetAbility(skillsetId: $skillsetId, abilityId: $abilityId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"skillsetId": skillsetId,
		"abilityId":  abilityId,
		"input":      input,
	}

	var response struct {
		UpdateSkillsetAbility *UpdateSkillsetAbilityResponse `json:"updateSkillsetAbility"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateSkillsetAbility, nil
}

// DeleteSkillsetAbilityResponse represents the response from deleting a skillsetability.
type DeleteSkillsetAbilityResponse struct {
	ID *string `json:"id"`
}

// DeleteSkillsetAbility deletes a skillsetability.
func (c *Client) DeleteSkillsetAbility(ctx context.Context, skillsetId string, abilityId string) (*DeleteSkillsetAbilityResponse, error) {
	query := `
		mutation DeleteSkillsetAbility($skillsetId: ID!, $abilityId: ID!) {
			deleteSkillsetAbility(skillsetId: $skillsetId, abilityId: $abilityId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"skillsetId": skillsetId,
		"abilityId":  abilityId,
	}

	var response struct {
		DeleteSkillsetAbility *DeleteSkillsetAbilityResponse `json:"deleteSkillsetAbility"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteSkillsetAbility, nil
}

// GetSkillsetAbilityResponse represents the response from fetching a skillsetability.
type GetSkillsetAbilityResponse struct {
	ID          *string                `json:"id"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	BotId       *string                `json:"botId,omitempty"`
	Description *string                `json:"description,omitempty"`
	FileId      *string                `json:"fileId,omitempty"`
	Instruction *string                `json:"instruction,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	SecretId    *string                `json:"secretId,omitempty"`
	SpaceId     *string                `json:"spaceId,omitempty"`
	State       *string                `json:"state,omitempty"`
	CreatedAt   *string                `json:"createdAt,omitempty"`
	UpdatedAt   *string                `json:"updatedAt,omitempty"`
}

// GetSkillsetAbility fetches a skillsetability by ID.
func (c *Client) GetSkillsetAbility(ctx context.Context, skillsetId string, id string) (*GetSkillsetAbilityResponse, error) {
	// Query abilities through the skillset connection
	query := `
		query GetSkillsetAbility($skillsetIds: [ID!]) {
			skillsets(first: 1, skillsetIds: $skillsetIds) {
				edges {
					node {
						id
						abilities(first: 100) {
							edges {
								node {
									id
									blueprintId
									botId
									description
									fileId
									instruction
									meta
									name
									secretId
									spaceId
									state
									createdAt
									updatedAt
								}
							}
						}
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"skillsetIds": []string{skillsetId},
	}

	var response struct {
		Skillsets struct {
			Edges []struct {
				Node struct {
					ID        string `json:"id"`
					Abilities struct {
						Edges []struct {
							Node *GetSkillsetAbilityResponse `json:"node"`
						} `json:"edges"`
					} `json:"abilities"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"skillsets"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the ability with matching ID
	for _, parentEdge := range response.Skillsets.Edges {
		for _, abilityEdge := range parentEdge.Node.Abilities.Edges {
			if abilityEdge.Node != nil && abilityEdge.Node.ID != nil && *abilityEdge.Node.ID == id {
				return abilityEdge.Node, nil
			}
		}
	}

	return nil, fmt.Errorf("skillsetability with ID %s not found in skillset %s", id, skillsetId)
}

// CreateSkillsetInput represents the input for creating a skillset.
type CreateSkillsetInput struct {
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	State       *string                `json:"state,omitempty"`
	Visibility  *string                `json:"visibility,omitempty"`
}

// CreateSkillsetResponse represents the response from creating a skillset.
type CreateSkillsetResponse struct {
	ID *string `json:"id"`
}

// CreateSkillset creates a new skillset.
func (c *Client) CreateSkillset(ctx context.Context, input CreateSkillsetInput) (*CreateSkillsetResponse, error) {
	query := `
		mutation CreateSkillset($input: SkillsetCreateRequest!) {
			createSkillset(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateSkillset *CreateSkillsetResponse `json:"createSkillset"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateSkillset, nil
}

// UpdateSkillsetInput represents the input for updating a skillset.
type UpdateSkillsetInput struct {
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	State       *string                `json:"state,omitempty"`
	Visibility  *string                `json:"visibility,omitempty"`
}

// UpdateSkillsetResponse represents the response from updating a skillset.
type UpdateSkillsetResponse struct {
	ID *string `json:"id"`
}

// UpdateSkillset updates an existing skillset.
func (c *Client) UpdateSkillset(ctx context.Context, id string, input UpdateSkillsetInput) (*UpdateSkillsetResponse, error) {
	query := `
		mutation UpdateSkillset($skillsetId: ID!, $input: SkillsetUpdateRequest!) {
			updateSkillset(skillsetId: $skillsetId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"skillsetId": id,
		"input":      input,
	}

	var response struct {
		UpdateSkillset *UpdateSkillsetResponse `json:"updateSkillset"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateSkillset, nil
}

// DeleteSkillsetResponse represents the response from deleting a skillset.
type DeleteSkillsetResponse struct {
	ID *string `json:"id"`
}

// DeleteSkillset deletes a skillset.
func (c *Client) DeleteSkillset(ctx context.Context, id string) (*DeleteSkillsetResponse, error) {
	query := `
		mutation DeleteSkillset($skillsetId: ID!) {
			deleteSkillset(skillsetId: $skillsetId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"skillsetId": id,
	}

	var response struct {
		DeleteSkillset *DeleteSkillsetResponse `json:"deleteSkillset"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteSkillset, nil
}

// GetSkillsetResponse represents the response from fetching a skillset.
type GetSkillsetResponse struct {
	ID          *string                `json:"id"`
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	State       *string                `json:"state,omitempty"`
	Visibility  *string                `json:"visibility,omitempty"`
	CreatedAt   *string                `json:"createdAt,omitempty"`
	UpdatedAt   *string                `json:"updatedAt,omitempty"`
}

// GetSkillset fetches a skillset by ID.
func (c *Client) GetSkillset(ctx context.Context, id string) (*GetSkillsetResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetSkillset($cursor: ID) {
			skillsets(first: 1, after: $cursor) {
				edges {
					node {
						id
						alias
						blueprintId
						description
						meta
						name
						state
						visibility
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		Skillsets struct {
			Edges []struct {
				Node *GetSkillsetResponse `json:"node"`
			} `json:"edges"`
		} `json:"skillsets"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.Skillsets.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("skillset with ID %s not found", id)
}

// CreateSlackIntegrationInput represents the input for creating a slackintegration.
type CreateSlackIntegrationInput struct {
	Alias             *string                `json:"alias,omitempty"`
	AllowFrom         *string                `json:"allowFrom,omitempty"`
	AutoRespond       *string                `json:"autoRespond,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	BotToken          *string                `json:"botToken,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	Ratings           *bool                  `json:"ratings,omitempty"`
	References        *bool                  `json:"references,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
	SigningSecret     *string                `json:"signingSecret,omitempty"`
	UserToken         *string                `json:"userToken,omitempty"`
	VisibleMessages   *int64                 `json:"visibleMessages,omitempty"`
}

// CreateSlackIntegrationResponse represents the response from creating a slackintegration.
type CreateSlackIntegrationResponse struct {
	ID *string `json:"id"`
}

// CreateSlackIntegration creates a new slackintegration.
func (c *Client) CreateSlackIntegration(ctx context.Context, input CreateSlackIntegrationInput) (*CreateSlackIntegrationResponse, error) {
	query := `
		mutation CreateSlackIntegration($input: SlackIntegrationCreateRequest!) {
			createSlackIntegration(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateSlackIntegration *CreateSlackIntegrationResponse `json:"createSlackIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateSlackIntegration, nil
}

// UpdateSlackIntegrationInput represents the input for updating a slackintegration.
type UpdateSlackIntegrationInput struct {
	Alias             *string                `json:"alias,omitempty"`
	AllowFrom         *string                `json:"allowFrom,omitempty"`
	AutoRespond       *string                `json:"autoRespond,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	BotToken          *string                `json:"botToken,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	Ratings           *bool                  `json:"ratings,omitempty"`
	References        *bool                  `json:"references,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
	SigningSecret     *string                `json:"signingSecret,omitempty"`
	UserToken         *string                `json:"userToken,omitempty"`
	VisibleMessages   *int64                 `json:"visibleMessages,omitempty"`
}

// UpdateSlackIntegrationResponse represents the response from updating a slackintegration.
type UpdateSlackIntegrationResponse struct {
	ID *string `json:"id"`
}

// UpdateSlackIntegration updates an existing slackintegration.
func (c *Client) UpdateSlackIntegration(ctx context.Context, id string, input UpdateSlackIntegrationInput) (*UpdateSlackIntegrationResponse, error) {
	query := `
		mutation UpdateSlackIntegration($slackIntegrationId: ID!, $input: SlackIntegrationUpdateRequest!) {
			updateSlackIntegration(slackIntegrationId: $slackIntegrationId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"slackIntegrationId": id,
		"input":              input,
	}

	var response struct {
		UpdateSlackIntegration *UpdateSlackIntegrationResponse `json:"updateSlackIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateSlackIntegration, nil
}

// DeleteSlackIntegrationResponse represents the response from deleting a slackintegration.
type DeleteSlackIntegrationResponse struct {
	ID *string `json:"id"`
}

// DeleteSlackIntegration deletes a slackintegration.
func (c *Client) DeleteSlackIntegration(ctx context.Context, id string) (*DeleteSlackIntegrationResponse, error) {
	query := `
		mutation DeleteSlackIntegration($slackIntegrationId: ID!) {
			deleteSlackIntegration(slackIntegrationId: $slackIntegrationId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"slackIntegrationId": id,
	}

	var response struct {
		DeleteSlackIntegration *DeleteSlackIntegrationResponse `json:"deleteSlackIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteSlackIntegration, nil
}

// GetSlackIntegrationResponse represents the response from fetching a slackintegration.
type GetSlackIntegrationResponse struct {
	ID                *string                `json:"id"`
	Alias             *string                `json:"alias,omitempty"`
	AllowFrom         *string                `json:"allowFrom,omitempty"`
	AutoRespond       *string                `json:"autoRespond,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	BotToken          *string                `json:"botToken,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	Ratings           *bool                  `json:"ratings,omitempty"`
	References        *bool                  `json:"references,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
	SigningSecret     *string                `json:"signingSecret,omitempty"`
	UserToken         *string                `json:"userToken,omitempty"`
	VisibleMessages   *int64                 `json:"visibleMessages,omitempty"`
	CreatedAt         *string                `json:"createdAt,omitempty"`
	UpdatedAt         *string                `json:"updatedAt,omitempty"`
}

// GetSlackIntegration fetches a slackintegration by ID.
func (c *Client) GetSlackIntegration(ctx context.Context, id string) (*GetSlackIntegrationResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetSlackIntegration($cursor: ID) {
			slackIntegrations(first: 1, after: $cursor) {
				edges {
					node {
						id
						alias
						allowFrom
						autoRespond
						blueprintId
						botId
						botToken
						contactCollection
						description
						meta
						name
						ratings
						references
						sessionDuration
						signingSecret
						userToken
						visibleMessages
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		SlackIntegrations struct {
			Edges []struct {
				Node *GetSlackIntegrationResponse `json:"node"`
			} `json:"edges"`
		} `json:"slackIntegrations"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.SlackIntegrations.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("slackintegration with ID %s not found", id)
}

// CreateSpaceInput represents the input for creating a space.
type CreateSpaceInput struct {
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	ContactId   *string                `json:"contactId,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
}

// CreateSpaceResponse represents the response from creating a space.
type CreateSpaceResponse struct {
	ID *string `json:"id"`
}

// CreateSpace creates a new space.
func (c *Client) CreateSpace(ctx context.Context, input CreateSpaceInput) (*CreateSpaceResponse, error) {
	query := `
		mutation CreateSpace($input: SpaceCreateRequest!) {
			createSpace(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateSpace *CreateSpaceResponse `json:"createSpace"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateSpace, nil
}

// UpdateSpaceInput represents the input for updating a space.
type UpdateSpaceInput struct {
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	ContactId   *string                `json:"contactId,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
}

// UpdateSpaceResponse represents the response from updating a space.
type UpdateSpaceResponse struct {
	ID *string `json:"id"`
}

// UpdateSpace updates an existing space.
func (c *Client) UpdateSpace(ctx context.Context, id string, input UpdateSpaceInput) (*UpdateSpaceResponse, error) {
	query := `
		mutation UpdateSpace($spaceId: ID!, $input: SpaceUpdateRequest!) {
			updateSpace(spaceId: $spaceId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"spaceId": id,
		"input":   input,
	}

	var response struct {
		UpdateSpace *UpdateSpaceResponse `json:"updateSpace"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateSpace, nil
}

// DeleteSpaceResponse represents the response from deleting a space.
type DeleteSpaceResponse struct {
	ID *string `json:"id"`
}

// DeleteSpace deletes a space.
func (c *Client) DeleteSpace(ctx context.Context, id string) (*DeleteSpaceResponse, error) {
	query := `
		mutation DeleteSpace($spaceId: ID!) {
			deleteSpace(spaceId: $spaceId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"spaceId": id,
	}

	var response struct {
		DeleteSpace *DeleteSpaceResponse `json:"deleteSpace"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteSpace, nil
}

// GetSpaceResponse represents the response from fetching a space.
type GetSpaceResponse struct {
	ID          *string                `json:"id"`
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	ContactId   *string                `json:"contactId,omitempty"`
	Description *string                `json:"description,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	CreatedAt   *string                `json:"createdAt,omitempty"`
	UpdatedAt   *string                `json:"updatedAt,omitempty"`
}

// GetSpace fetches a space by ID.
func (c *Client) GetSpace(ctx context.Context, id string) (*GetSpaceResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetSpace($cursor: ID) {
			spaces(first: 1, after: $cursor) {
				edges {
					node {
						id
						alias
						blueprintId
						contactId
						description
						meta
						name
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		Spaces struct {
			Edges []struct {
				Node *GetSpaceResponse `json:"node"`
			} `json:"edges"`
		} `json:"spaces"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.Spaces.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("space with ID %s not found", id)
}

// CreateSpaceSiteInput represents the input for creating a spacesite.
type CreateSpaceSiteInput struct {
	Alias       *string                `json:"alias,omitempty"`
	Description *string                `json:"description,omitempty"`
	Domain      string                 `json:"domain,omitempty"`
	Index       *string                `json:"index,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	NotFound    *string                `json:"notFound,omitempty"`
	Prefix      *string                `json:"prefix,omitempty"`
}

// CreateSpaceSiteResponse represents the response from creating a spacesite.
type CreateSpaceSiteResponse struct {
	ID *string `json:"id"`
}

// CreateSpaceSite creates a new spacesite.
func (c *Client) CreateSpaceSite(ctx context.Context, spaceId string, input CreateSpaceSiteInput) (*CreateSpaceSiteResponse, error) {
	query := `
		mutation CreateSpaceSite($spaceId: ID!, $input: SpaceSiteCreateRequest!) {
			createSpaceSite(spaceId: $spaceId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"spaceId": spaceId,
		"input":   input,
	}

	var response struct {
		CreateSpaceSite *CreateSpaceSiteResponse `json:"createSpaceSite"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateSpaceSite, nil
}

// UpdateSpaceSiteInput represents the input for updating a spacesite.
type UpdateSpaceSiteInput struct {
	Alias       *string                `json:"alias,omitempty"`
	Description *string                `json:"description,omitempty"`
	Domain      string                 `json:"domain,omitempty"`
	Index       *string                `json:"index,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	NotFound    *string                `json:"notFound,omitempty"`
	Prefix      *string                `json:"prefix,omitempty"`
}

// UpdateSpaceSiteResponse represents the response from updating a spacesite.
type UpdateSpaceSiteResponse struct {
	ID *string `json:"id"`
}

// UpdateSpaceSite updates an existing spacesite.
func (c *Client) UpdateSpaceSite(ctx context.Context, spaceId string, siteId string, input UpdateSpaceSiteInput) (*UpdateSpaceSiteResponse, error) {
	query := `
		mutation UpdateSpaceSite($spaceId: ID!, $siteId: ID!, $input: SpaceSiteUpdateRequest!) {
			updateSpaceSite(spaceId: $spaceId, siteId: $siteId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"spaceId": spaceId,
		"siteId":  siteId,
		"input":   input,
	}

	var response struct {
		UpdateSpaceSite *UpdateSpaceSiteResponse `json:"updateSpaceSite"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateSpaceSite, nil
}

// DeleteSpaceSiteResponse represents the response from deleting a spacesite.
type DeleteSpaceSiteResponse struct {
	ID *string `json:"id"`
}

// DeleteSpaceSite deletes a spacesite.
func (c *Client) DeleteSpaceSite(ctx context.Context, spaceId string, siteId string) (*DeleteSpaceSiteResponse, error) {
	query := `
		mutation DeleteSpaceSite($spaceId: ID!, $siteId: ID!) {
			deleteSpaceSite(spaceId: $spaceId, siteId: $siteId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"spaceId": spaceId,
		"siteId":  siteId,
	}

	var response struct {
		DeleteSpaceSite *DeleteSpaceSiteResponse `json:"deleteSpaceSite"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteSpaceSite, nil
}

// GetSpaceSiteResponse represents the response from fetching a spacesite.
type GetSpaceSiteResponse struct {
	ID          *string                `json:"id"`
	Alias       *string                `json:"alias,omitempty"`
	Description *string                `json:"description,omitempty"`
	Domain      string                 `json:"domain,omitempty"`
	Index       *string                `json:"index,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	NotFound    *string                `json:"notFound,omitempty"`
	Prefix      *string                `json:"prefix,omitempty"`
	CreatedAt   *string                `json:"createdAt,omitempty"`
	UpdatedAt   *string                `json:"updatedAt,omitempty"`
}

// GetSpaceSite fetches a spacesite by ID.
func (c *Client) GetSpaceSite(ctx context.Context, spaceId string, id string) (*GetSpaceSiteResponse, error) {
	// Query sites through the space connection
	query := `
		query GetSpaceSite($spaceIds: [ID!]) {
			spaces(first: 1, spaceIds: $spaceIds) {
				edges {
					node {
						id
						sites(first: 100) {
							edges {
								node {
									id
									alias
									description
									domain
									index
									meta
									name
									notFound
									prefix
									createdAt
									updatedAt
								}
							}
						}
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"spaceIds": []string{spaceId},
	}

	var response struct {
		Spaces struct {
			Edges []struct {
				Node struct {
					ID    string `json:"id"`
					Sites struct {
						Edges []struct {
							Node *GetSpaceSiteResponse `json:"node"`
						} `json:"edges"`
					} `json:"sites"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"spaces"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the ability with matching ID
	for _, parentEdge := range response.Spaces.Edges {
		for _, abilityEdge := range parentEdge.Node.Sites.Edges {
			if abilityEdge.Node != nil && abilityEdge.Node.ID != nil && *abilityEdge.Node.ID == id {
				return abilityEdge.Node, nil
			}
		}
	}

	return nil, fmt.Errorf("spacesite with ID %s not found in space %s", id, spaceId)
}

// CreateSupportIntegrationInput represents the input for creating a supportintegration.
type CreateSupportIntegrationInput struct {
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	BotId       *string                `json:"botId,omitempty"`
	Description *string                `json:"description,omitempty"`
	Email       *string                `json:"email,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
}

// CreateSupportIntegrationResponse represents the response from creating a supportintegration.
type CreateSupportIntegrationResponse struct {
	ID *string `json:"id"`
}

// CreateSupportIntegration creates a new supportintegration.
func (c *Client) CreateSupportIntegration(ctx context.Context, input CreateSupportIntegrationInput) (*CreateSupportIntegrationResponse, error) {
	query := `
		mutation CreateSupportIntegration($input: SupportIntegrationCreateRequest!) {
			createSupportIntegration(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateSupportIntegration *CreateSupportIntegrationResponse `json:"createSupportIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateSupportIntegration, nil
}

// UpdateSupportIntegrationInput represents the input for updating a supportintegration.
type UpdateSupportIntegrationInput struct {
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	BotId       *string                `json:"botId,omitempty"`
	Description *string                `json:"description,omitempty"`
	Email       *string                `json:"email,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
}

// UpdateSupportIntegrationResponse represents the response from updating a supportintegration.
type UpdateSupportIntegrationResponse struct {
	ID *string `json:"id"`
}

// UpdateSupportIntegration updates an existing supportintegration.
func (c *Client) UpdateSupportIntegration(ctx context.Context, id string, input UpdateSupportIntegrationInput) (*UpdateSupportIntegrationResponse, error) {
	query := `
		mutation UpdateSupportIntegration($supportIntegrationId: ID!, $input: SupportIntegrationUpdateRequest!) {
			updateSupportIntegration(supportIntegrationId: $supportIntegrationId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"supportIntegrationId": id,
		"input":                input,
	}

	var response struct {
		UpdateSupportIntegration *UpdateSupportIntegrationResponse `json:"updateSupportIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateSupportIntegration, nil
}

// DeleteSupportIntegrationResponse represents the response from deleting a supportintegration.
type DeleteSupportIntegrationResponse struct {
	ID *string `json:"id"`
}

// DeleteSupportIntegration deletes a supportintegration.
func (c *Client) DeleteSupportIntegration(ctx context.Context, id string) (*DeleteSupportIntegrationResponse, error) {
	query := `
		mutation DeleteSupportIntegration($supportIntegrationId: ID!) {
			deleteSupportIntegration(supportIntegrationId: $supportIntegrationId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"supportIntegrationId": id,
	}

	var response struct {
		DeleteSupportIntegration *DeleteSupportIntegrationResponse `json:"deleteSupportIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteSupportIntegration, nil
}

// GetSupportIntegrationResponse represents the response from fetching a supportintegration.
type GetSupportIntegrationResponse struct {
	ID          *string                `json:"id"`
	Alias       *string                `json:"alias,omitempty"`
	BlueprintId *string                `json:"blueprintId,omitempty"`
	BotId       *string                `json:"botId,omitempty"`
	Description *string                `json:"description,omitempty"`
	Email       *string                `json:"email,omitempty"`
	Meta        map[string]interface{} `json:"meta,omitempty"`
	Name        *string                `json:"name,omitempty"`
	CreatedAt   *string                `json:"createdAt,omitempty"`
	UpdatedAt   *string                `json:"updatedAt,omitempty"`
}

// GetSupportIntegration fetches a supportintegration by ID.
func (c *Client) GetSupportIntegration(ctx context.Context, id string) (*GetSupportIntegrationResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetSupportIntegration($cursor: ID) {
			supportIntegrations(first: 1, after: $cursor) {
				edges {
					node {
						id
						alias
						blueprintId
						botId
						description
						email
						meta
						name
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		SupportIntegrations struct {
			Edges []struct {
				Node *GetSupportIntegrationResponse `json:"node"`
			} `json:"edges"`
		} `json:"supportIntegrations"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.SupportIntegrations.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("supportintegration with ID %s not found", id)
}

// CreateTaskInput represents the input for creating a task.
type CreateTaskInput struct {
	BotId           *string                `json:"botId,omitempty"`
	ContactId       *string                `json:"contactId,omitempty"`
	Description     *string                `json:"description,omitempty"`
	MaxCalls        *int64                 `json:"maxCalls,omitempty"`
	MaxIterations   *int64                 `json:"maxIterations,omitempty"`
	MaxTime         *float64               `json:"maxTime,omitempty"`
	Meta            map[string]interface{} `json:"meta,omitempty"`
	Name            *string                `json:"name,omitempty"`
	Schedule        *string                `json:"schedule,omitempty"`
	SessionDuration *float64               `json:"sessionDuration,omitempty"`
	Timezone        *string                `json:"timezone,omitempty"`
}

// CreateTaskResponse represents the response from creating a task.
type CreateTaskResponse struct {
	ID *string `json:"id"`
}

// CreateTask creates a new task.
func (c *Client) CreateTask(ctx context.Context, input CreateTaskInput) (*CreateTaskResponse, error) {
	query := `
		mutation CreateTask($input: TaskCreateRequest!) {
			createTask(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateTask *CreateTaskResponse `json:"createTask"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateTask, nil
}

// UpdateTaskInput represents the input for updating a task.
type UpdateTaskInput struct {
	BotId           *string                `json:"botId,omitempty"`
	ContactId       *string                `json:"contactId,omitempty"`
	Description     *string                `json:"description,omitempty"`
	MaxCalls        *int64                 `json:"maxCalls,omitempty"`
	MaxIterations   *int64                 `json:"maxIterations,omitempty"`
	MaxTime         *float64               `json:"maxTime,omitempty"`
	Meta            map[string]interface{} `json:"meta,omitempty"`
	Name            *string                `json:"name,omitempty"`
	Schedule        *string                `json:"schedule,omitempty"`
	SessionDuration *float64               `json:"sessionDuration,omitempty"`
	Timezone        *string                `json:"timezone,omitempty"`
}

// UpdateTaskResponse represents the response from updating a task.
type UpdateTaskResponse struct {
	ID *string `json:"id"`
}

// UpdateTask updates an existing task.
func (c *Client) UpdateTask(ctx context.Context, id string, input UpdateTaskInput) (*UpdateTaskResponse, error) {
	query := `
		mutation UpdateTask($taskId: ID!, $input: TaskUpdateRequest!) {
			updateTask(taskId: $taskId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"taskId": id,
		"input":  input,
	}

	var response struct {
		UpdateTask *UpdateTaskResponse `json:"updateTask"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateTask, nil
}

// DeleteTaskResponse represents the response from deleting a task.
type DeleteTaskResponse struct {
	ID *string `json:"id"`
}

// DeleteTask deletes a task.
func (c *Client) DeleteTask(ctx context.Context, id string) (*DeleteTaskResponse, error) {
	query := `
		mutation DeleteTask($taskId: ID!) {
			deleteTask(taskId: $taskId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"taskId": id,
	}

	var response struct {
		DeleteTask *DeleteTaskResponse `json:"deleteTask"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteTask, nil
}

// GetTaskResponse represents the response from fetching a task.
type GetTaskResponse struct {
	ID              *string                `json:"id"`
	BotId           *string                `json:"botId,omitempty"`
	ContactId       *string                `json:"contactId,omitempty"`
	Description     *string                `json:"description,omitempty"`
	MaxCalls        *int64                 `json:"maxCalls,omitempty"`
	MaxIterations   *int64                 `json:"maxIterations,omitempty"`
	MaxTime         *float64               `json:"maxTime,omitempty"`
	Meta            map[string]interface{} `json:"meta,omitempty"`
	Name            *string                `json:"name,omitempty"`
	Schedule        *string                `json:"schedule,omitempty"`
	SessionDuration *float64               `json:"sessionDuration,omitempty"`
	Timezone        *string                `json:"timezone,omitempty"`
	CreatedAt       *string                `json:"createdAt,omitempty"`
	UpdatedAt       *string                `json:"updatedAt,omitempty"`
}

// GetTask fetches a task by ID.
func (c *Client) GetTask(ctx context.Context, id string) (*GetTaskResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetTask($cursor: ID) {
			tasks(first: 1, after: $cursor) {
				edges {
					node {
						id
						botId
						contactId
						description
						maxCalls
						maxIterations
						maxTime
						meta
						name
						schedule
						sessionDuration
						timezone
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		Tasks struct {
			Edges []struct {
				Node *GetTaskResponse `json:"node"`
			} `json:"edges"`
		} `json:"tasks"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.Tasks.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("task with ID %s not found", id)
}

// CreateTelegramIntegrationInput represents the input for creating a telegramintegration.
type CreateTelegramIntegrationInput struct {
	Alias             *string                `json:"alias,omitempty"`
	AllowFrom         *string                `json:"allowFrom,omitempty"`
	Attachments       *bool                  `json:"attachments,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	BotToken          *string                `json:"botToken,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
}

// CreateTelegramIntegrationResponse represents the response from creating a telegramintegration.
type CreateTelegramIntegrationResponse struct {
	ID *string `json:"id"`
}

// CreateTelegramIntegration creates a new telegramintegration.
func (c *Client) CreateTelegramIntegration(ctx context.Context, input CreateTelegramIntegrationInput) (*CreateTelegramIntegrationResponse, error) {
	query := `
		mutation CreateTelegramIntegration($input: TelegramIntegrationCreateRequest!) {
			createTelegramIntegration(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateTelegramIntegration *CreateTelegramIntegrationResponse `json:"createTelegramIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateTelegramIntegration, nil
}

// UpdateTelegramIntegrationInput represents the input for updating a telegramintegration.
type UpdateTelegramIntegrationInput struct {
	Alias             *string                `json:"alias,omitempty"`
	AllowFrom         *string                `json:"allowFrom,omitempty"`
	Attachments       *bool                  `json:"attachments,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	BotToken          *string                `json:"botToken,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
}

// UpdateTelegramIntegrationResponse represents the response from updating a telegramintegration.
type UpdateTelegramIntegrationResponse struct {
	ID *string `json:"id"`
}

// UpdateTelegramIntegration updates an existing telegramintegration.
func (c *Client) UpdateTelegramIntegration(ctx context.Context, id string, input UpdateTelegramIntegrationInput) (*UpdateTelegramIntegrationResponse, error) {
	query := `
		mutation UpdateTelegramIntegration($telegramIntegrationId: ID!, $input: TelegramIntegrationUpdateRequest!) {
			updateTelegramIntegration(telegramIntegrationId: $telegramIntegrationId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"telegramIntegrationId": id,
		"input":                 input,
	}

	var response struct {
		UpdateTelegramIntegration *UpdateTelegramIntegrationResponse `json:"updateTelegramIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateTelegramIntegration, nil
}

// DeleteTelegramIntegrationResponse represents the response from deleting a telegramintegration.
type DeleteTelegramIntegrationResponse struct {
	ID *string `json:"id"`
}

// DeleteTelegramIntegration deletes a telegramintegration.
func (c *Client) DeleteTelegramIntegration(ctx context.Context, id string) (*DeleteTelegramIntegrationResponse, error) {
	query := `
		mutation DeleteTelegramIntegration($telegramIntegrationId: ID!) {
			deleteTelegramIntegration(telegramIntegrationId: $telegramIntegrationId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"telegramIntegrationId": id,
	}

	var response struct {
		DeleteTelegramIntegration *DeleteTelegramIntegrationResponse `json:"deleteTelegramIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteTelegramIntegration, nil
}

// GetTelegramIntegrationResponse represents the response from fetching a telegramintegration.
type GetTelegramIntegrationResponse struct {
	ID                *string                `json:"id"`
	Alias             *string                `json:"alias,omitempty"`
	AllowFrom         *string                `json:"allowFrom,omitempty"`
	Attachments       *bool                  `json:"attachments,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	BotToken          *string                `json:"botToken,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
	CreatedAt         *string                `json:"createdAt,omitempty"`
	UpdatedAt         *string                `json:"updatedAt,omitempty"`
}

// GetTelegramIntegration fetches a telegramintegration by ID.
func (c *Client) GetTelegramIntegration(ctx context.Context, id string) (*GetTelegramIntegrationResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetTelegramIntegration($cursor: ID) {
			telegramIntegrations(first: 1, after: $cursor) {
				edges {
					node {
						id
						alias
						allowFrom
						attachments
						blueprintId
						botId
						botToken
						contactCollection
						description
						meta
						name
						sessionDuration
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		TelegramIntegrations struct {
			Edges []struct {
				Node *GetTelegramIntegrationResponse `json:"node"`
			} `json:"edges"`
		} `json:"telegramIntegrations"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.TelegramIntegrations.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("telegramintegration with ID %s not found", id)
}

// CreateTriggerIntegrationInput represents the input for creating a triggerintegration.
type CreateTriggerIntegrationInput struct {
	Alias           *string                `json:"alias,omitempty"`
	Authenticate    *bool                  `json:"authenticate,omitempty"`
	BlueprintId     *string                `json:"blueprintId,omitempty"`
	BotId           *string                `json:"botId,omitempty"`
	Description     *string                `json:"description,omitempty"`
	Meta            map[string]interface{} `json:"meta,omitempty"`
	Name            *string                `json:"name,omitempty"`
	Schedule        *string                `json:"schedule,omitempty"`
	SessionDuration *int64                 `json:"sessionDuration,omitempty"`
	Timezone        *string                `json:"timezone,omitempty"`
}

// CreateTriggerIntegrationResponse represents the response from creating a triggerintegration.
type CreateTriggerIntegrationResponse struct {
	ID *string `json:"id"`
}

// CreateTriggerIntegration creates a new triggerintegration.
func (c *Client) CreateTriggerIntegration(ctx context.Context, input CreateTriggerIntegrationInput) (*CreateTriggerIntegrationResponse, error) {
	query := `
		mutation CreateTriggerIntegration($input: TriggerIntegrationCreateRequest!) {
			createTriggerIntegration(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateTriggerIntegration *CreateTriggerIntegrationResponse `json:"createTriggerIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateTriggerIntegration, nil
}

// UpdateTriggerIntegrationInput represents the input for updating a triggerintegration.
type UpdateTriggerIntegrationInput struct {
	Alias           *string                `json:"alias,omitempty"`
	Authenticate    *bool                  `json:"authenticate,omitempty"`
	BlueprintId     *string                `json:"blueprintId,omitempty"`
	BotId           *string                `json:"botId,omitempty"`
	Description     *string                `json:"description,omitempty"`
	Meta            map[string]interface{} `json:"meta,omitempty"`
	Name            *string                `json:"name,omitempty"`
	Schedule        *string                `json:"schedule,omitempty"`
	SessionDuration *int64                 `json:"sessionDuration,omitempty"`
	Timezone        *string                `json:"timezone,omitempty"`
}

// UpdateTriggerIntegrationResponse represents the response from updating a triggerintegration.
type UpdateTriggerIntegrationResponse struct {
	ID *string `json:"id"`
}

// UpdateTriggerIntegration updates an existing triggerintegration.
func (c *Client) UpdateTriggerIntegration(ctx context.Context, id string, input UpdateTriggerIntegrationInput) (*UpdateTriggerIntegrationResponse, error) {
	query := `
		mutation UpdateTriggerIntegration($triggerIntegrationId: ID!, $input: TriggerIntegrationUpdateRequest!) {
			updateTriggerIntegration(triggerIntegrationId: $triggerIntegrationId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"triggerIntegrationId": id,
		"input":                input,
	}

	var response struct {
		UpdateTriggerIntegration *UpdateTriggerIntegrationResponse `json:"updateTriggerIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateTriggerIntegration, nil
}

// DeleteTriggerIntegrationResponse represents the response from deleting a triggerintegration.
type DeleteTriggerIntegrationResponse struct {
	ID *string `json:"id"`
}

// DeleteTriggerIntegration deletes a triggerintegration.
func (c *Client) DeleteTriggerIntegration(ctx context.Context, id string) (*DeleteTriggerIntegrationResponse, error) {
	query := `
		mutation DeleteTriggerIntegration($triggerIntegrationId: ID!) {
			deleteTriggerIntegration(triggerIntegrationId: $triggerIntegrationId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"triggerIntegrationId": id,
	}

	var response struct {
		DeleteTriggerIntegration *DeleteTriggerIntegrationResponse `json:"deleteTriggerIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteTriggerIntegration, nil
}

// GetTriggerIntegrationResponse represents the response from fetching a triggerintegration.
type GetTriggerIntegrationResponse struct {
	ID              *string                `json:"id"`
	Alias           *string                `json:"alias,omitempty"`
	Authenticate    *bool                  `json:"authenticate,omitempty"`
	BlueprintId     *string                `json:"blueprintId,omitempty"`
	BotId           *string                `json:"botId,omitempty"`
	Description     *string                `json:"description,omitempty"`
	Meta            map[string]interface{} `json:"meta,omitempty"`
	Name            *string                `json:"name,omitempty"`
	Schedule        *string                `json:"schedule,omitempty"`
	SessionDuration *int64                 `json:"sessionDuration,omitempty"`
	Timezone        *string                `json:"timezone,omitempty"`
	CreatedAt       *string                `json:"createdAt,omitempty"`
	UpdatedAt       *string                `json:"updatedAt,omitempty"`
}

// GetTriggerIntegration fetches a triggerintegration by ID.
func (c *Client) GetTriggerIntegration(ctx context.Context, id string) (*GetTriggerIntegrationResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetTriggerIntegration($cursor: ID) {
			triggerIntegrations(first: 1, after: $cursor) {
				edges {
					node {
						id
						alias
						authenticate
						blueprintId
						botId
						description
						meta
						name
						schedule
						sessionDuration
						timezone
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		TriggerIntegrations struct {
			Edges []struct {
				Node *GetTriggerIntegrationResponse `json:"node"`
			} `json:"edges"`
		} `json:"triggerIntegrations"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.TriggerIntegrations.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("triggerintegration with ID %s not found", id)
}

// CreateTwilioIntegrationInput represents the input for creating a twiliointegration.
type CreateTwilioIntegrationInput struct {
	AccountSid        *string                `json:"accountSid,omitempty"`
	Alias             *string                `json:"alias,omitempty"`
	AllowFrom         *string                `json:"allowFrom,omitempty"`
	AuthToken         *string                `json:"authToken,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
	Voice             *string                `json:"voice,omitempty"`
}

// CreateTwilioIntegrationResponse represents the response from creating a twiliointegration.
type CreateTwilioIntegrationResponse struct {
	ID *string `json:"id"`
}

// CreateTwilioIntegration creates a new twiliointegration.
func (c *Client) CreateTwilioIntegration(ctx context.Context, input CreateTwilioIntegrationInput) (*CreateTwilioIntegrationResponse, error) {
	query := `
		mutation CreateTwilioIntegration($input: TwilioIntegrationCreateRequest!) {
			createTwilioIntegration(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateTwilioIntegration *CreateTwilioIntegrationResponse `json:"createTwilioIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateTwilioIntegration, nil
}

// UpdateTwilioIntegrationInput represents the input for updating a twiliointegration.
type UpdateTwilioIntegrationInput struct {
	AccountSid        *string                `json:"accountSid,omitempty"`
	Alias             *string                `json:"alias,omitempty"`
	AllowFrom         *string                `json:"allowFrom,omitempty"`
	AuthToken         *string                `json:"authToken,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
	Voice             *string                `json:"voice,omitempty"`
}

// UpdateTwilioIntegrationResponse represents the response from updating a twiliointegration.
type UpdateTwilioIntegrationResponse struct {
	ID *string `json:"id"`
}

// UpdateTwilioIntegration updates an existing twiliointegration.
func (c *Client) UpdateTwilioIntegration(ctx context.Context, id string, input UpdateTwilioIntegrationInput) (*UpdateTwilioIntegrationResponse, error) {
	query := `
		mutation UpdateTwilioIntegration($twilioIntegrationId: ID!, $input: TwilioIntegrationUpdateRequest!) {
			updateTwilioIntegration(twilioIntegrationId: $twilioIntegrationId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"twilioIntegrationId": id,
		"input":               input,
	}

	var response struct {
		UpdateTwilioIntegration *UpdateTwilioIntegrationResponse `json:"updateTwilioIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateTwilioIntegration, nil
}

// DeleteTwilioIntegrationResponse represents the response from deleting a twiliointegration.
type DeleteTwilioIntegrationResponse struct {
	ID *string `json:"id"`
}

// DeleteTwilioIntegration deletes a twiliointegration.
func (c *Client) DeleteTwilioIntegration(ctx context.Context, id string) (*DeleteTwilioIntegrationResponse, error) {
	query := `
		mutation DeleteTwilioIntegration($twilioIntegrationId: ID!) {
			deleteTwilioIntegration(twilioIntegrationId: $twilioIntegrationId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"twilioIntegrationId": id,
	}

	var response struct {
		DeleteTwilioIntegration *DeleteTwilioIntegrationResponse `json:"deleteTwilioIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteTwilioIntegration, nil
}

// GetTwilioIntegrationResponse represents the response from fetching a twiliointegration.
type GetTwilioIntegrationResponse struct {
	ID                *string                `json:"id"`
	AccountSid        *string                `json:"accountSid,omitempty"`
	Alias             *string                `json:"alias,omitempty"`
	AllowFrom         *string                `json:"allowFrom,omitempty"`
	AuthToken         *string                `json:"authToken,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
	Voice             *string                `json:"voice,omitempty"`
	CreatedAt         *string                `json:"createdAt,omitempty"`
	UpdatedAt         *string                `json:"updatedAt,omitempty"`
}

// GetTwilioIntegration fetches a twiliointegration by ID.
func (c *Client) GetTwilioIntegration(ctx context.Context, id string) (*GetTwilioIntegrationResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetTwilioIntegration($cursor: ID) {
			twilioIntegrations(first: 1, after: $cursor) {
				edges {
					node {
						id
						accountSid
						alias
						allowFrom
						authToken
						blueprintId
						botId
						contactCollection
						description
						meta
						name
						sessionDuration
						voice
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		TwilioIntegrations struct {
			Edges []struct {
				Node *GetTwilioIntegrationResponse `json:"node"`
			} `json:"edges"`
		} `json:"twilioIntegrations"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.TwilioIntegrations.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("twiliointegration with ID %s not found", id)
}

// CreateWhatsAppIntegrationInput represents the input for creating a whatsappintegration.
type CreateWhatsAppIntegrationInput struct {
	AccessToken       *string                `json:"accessToken,omitempty"`
	Alias             *string                `json:"alias,omitempty"`
	AllowFrom         *string                `json:"allowFrom,omitempty"`
	Attachments       *bool                  `json:"attachments,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	PhoneNumberId     *string                `json:"phoneNumberId,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
}

// CreateWhatsAppIntegrationResponse represents the response from creating a whatsappintegration.
type CreateWhatsAppIntegrationResponse struct {
	ID *string `json:"id"`
}

// CreateWhatsAppIntegration creates a new whatsappintegration.
func (c *Client) CreateWhatsAppIntegration(ctx context.Context, input CreateWhatsAppIntegrationInput) (*CreateWhatsAppIntegrationResponse, error) {
	query := `
		mutation CreateWhatsAppIntegration($input: WhatsAppIntegrationCreateRequest!) {
			createWhatsAppIntegration(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateWhatsAppIntegration *CreateWhatsAppIntegrationResponse `json:"createWhatsAppIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateWhatsAppIntegration, nil
}

// UpdateWhatsAppIntegrationInput represents the input for updating a whatsappintegration.
type UpdateWhatsAppIntegrationInput struct {
	AccessToken       *string                `json:"accessToken,omitempty"`
	Alias             *string                `json:"alias,omitempty"`
	AllowFrom         *string                `json:"allowFrom,omitempty"`
	Attachments       *bool                  `json:"attachments,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	PhoneNumberId     *string                `json:"phoneNumberId,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
}

// UpdateWhatsAppIntegrationResponse represents the response from updating a whatsappintegration.
type UpdateWhatsAppIntegrationResponse struct {
	ID *string `json:"id"`
}

// UpdateWhatsAppIntegration updates an existing whatsappintegration.
func (c *Client) UpdateWhatsAppIntegration(ctx context.Context, id string, input UpdateWhatsAppIntegrationInput) (*UpdateWhatsAppIntegrationResponse, error) {
	query := `
		mutation UpdateWhatsAppIntegration($whatsAppIntegrationId: ID!, $input: WhatsAppIntegrationUpdateRequest!) {
			updateWhatsAppIntegration(whatsAppIntegrationId: $whatsAppIntegrationId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"whatsAppIntegrationId": id,
		"input":                 input,
	}

	var response struct {
		UpdateWhatsAppIntegration *UpdateWhatsAppIntegrationResponse `json:"updateWhatsAppIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateWhatsAppIntegration, nil
}

// DeleteWhatsAppIntegrationResponse represents the response from deleting a whatsappintegration.
type DeleteWhatsAppIntegrationResponse struct {
	ID *string `json:"id"`
}

// DeleteWhatsAppIntegration deletes a whatsappintegration.
func (c *Client) DeleteWhatsAppIntegration(ctx context.Context, id string) (*DeleteWhatsAppIntegrationResponse, error) {
	query := `
		mutation DeleteWhatsAppIntegration($whatsAppIntegrationId: ID!) {
			deleteWhatsAppIntegration(whatsAppIntegrationId: $whatsAppIntegrationId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"whatsAppIntegrationId": id,
	}

	var response struct {
		DeleteWhatsAppIntegration *DeleteWhatsAppIntegrationResponse `json:"deleteWhatsAppIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteWhatsAppIntegration, nil
}

// GetWhatsAppIntegrationResponse represents the response from fetching a whatsappintegration.
type GetWhatsAppIntegrationResponse struct {
	ID                *string                `json:"id"`
	AccessToken       *string                `json:"accessToken,omitempty"`
	Alias             *string                `json:"alias,omitempty"`
	AllowFrom         *string                `json:"allowFrom,omitempty"`
	Attachments       *bool                  `json:"attachments,omitempty"`
	BlueprintId       *string                `json:"blueprintId,omitempty"`
	BotId             *string                `json:"botId,omitempty"`
	ContactCollection *bool                  `json:"contactCollection,omitempty"`
	Description       *string                `json:"description,omitempty"`
	Meta              map[string]interface{} `json:"meta,omitempty"`
	Name              *string                `json:"name,omitempty"`
	PhoneNumberId     *string                `json:"phoneNumberId,omitempty"`
	SessionDuration   *int64                 `json:"sessionDuration,omitempty"`
	CreatedAt         *string                `json:"createdAt,omitempty"`
	UpdatedAt         *string                `json:"updatedAt,omitempty"`
}

// GetWhatsAppIntegration fetches a whatsappintegration by ID.
func (c *Client) GetWhatsAppIntegration(ctx context.Context, id string) (*GetWhatsAppIntegrationResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetWhatsAppIntegration($cursor: ID) {
			whatsAppIntegrations(first: 1, after: $cursor) {
				edges {
					node {
						id
						accessToken
						alias
						allowFrom
						attachments
						blueprintId
						botId
						contactCollection
						description
						meta
						name
						phoneNumberId
						sessionDuration
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		WhatsAppIntegrations struct {
			Edges []struct {
				Node *GetWhatsAppIntegrationResponse `json:"node"`
			} `json:"edges"`
		} `json:"whatsAppIntegrations"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.WhatsAppIntegrations.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("whatsappintegration with ID %s not found", id)
}

// CreateWidgetIntegrationInput represents the input for creating a widgetintegration.
type CreateWidgetIntegrationInput struct {
	Alias               *string                `json:"alias,omitempty"`
	Attachments         *bool                  `json:"attachments,omitempty"`
	AutoScroll          *bool                  `json:"autoScroll,omitempty"`
	BlueprintId         *string                `json:"blueprintId,omitempty"`
	BotId               *string                `json:"botId,omitempty"`
	Carousel            *bool                  `json:"carousel,omitempty"`
	ContactCollection   *bool                  `json:"contactCollection,omitempty"`
	Description         *string                `json:"description,omitempty"`
	ExportConversation  *bool                  `json:"exportConversation,omitempty"`
	Form                *bool                  `json:"form,omitempty"`
	Initial             *string                `json:"initial,omitempty"`
	Intro               *string                `json:"intro,omitempty"`
	Language            *string                `json:"language,omitempty"`
	Layout              *string                `json:"layout,omitempty"`
	Math                *bool                  `json:"math,omitempty"`
	Maximize            *bool                  `json:"maximize,omitempty"`
	MessagePeek         *bool                  `json:"messagePeek,omitempty"`
	Meta                map[string]interface{} `json:"meta,omitempty"`
	Name                *string                `json:"name,omitempty"`
	Origin              *string                `json:"origin,omitempty"`
	Placeholder         *string                `json:"placeholder,omitempty"`
	Plugins             *string                `json:"plugins,omitempty"`
	PoweredBy           *bool                  `json:"poweredBy,omitempty"`
	RestartConversation *bool                  `json:"restartConversation,omitempty"`
	SessionDuration     *int64                 `json:"sessionDuration,omitempty"`
	StartFirst          *bool                  `json:"startFirst,omitempty"`
	Stream              *bool                  `json:"stream,omitempty"`
	Theme               *string                `json:"theme,omitempty"`
	Title               *string                `json:"title,omitempty"`
	Tools               *bool                  `json:"tools,omitempty"`
	Unfurl              *bool                  `json:"unfurl,omitempty"`
	Verbose             *bool                  `json:"verbose,omitempty"`
	VoiceIn             *bool                  `json:"voiceIn,omitempty"`
	VoiceOut            *bool                  `json:"voiceOut,omitempty"`
}

// CreateWidgetIntegrationResponse represents the response from creating a widgetintegration.
type CreateWidgetIntegrationResponse struct {
	ID *string `json:"id"`
}

// CreateWidgetIntegration creates a new widgetintegration.
func (c *Client) CreateWidgetIntegration(ctx context.Context, input CreateWidgetIntegrationInput) (*CreateWidgetIntegrationResponse, error) {
	query := `
		mutation CreateWidgetIntegration($input: WidgetIntegrationCreateRequest!) {
			createWidgetIntegration(input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"input": input,
	}

	var response struct {
		CreateWidgetIntegration *CreateWidgetIntegrationResponse `json:"createWidgetIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.CreateWidgetIntegration, nil
}

// UpdateWidgetIntegrationInput represents the input for updating a widgetintegration.
type UpdateWidgetIntegrationInput struct {
	Alias               *string                `json:"alias,omitempty"`
	Attachments         *bool                  `json:"attachments,omitempty"`
	AutoScroll          *bool                  `json:"autoScroll,omitempty"`
	BlueprintId         *string                `json:"blueprintId,omitempty"`
	BotId               *string                `json:"botId,omitempty"`
	Carousel            *bool                  `json:"carousel,omitempty"`
	ContactCollection   *bool                  `json:"contactCollection,omitempty"`
	Description         *string                `json:"description,omitempty"`
	ExportConversation  *bool                  `json:"exportConversation,omitempty"`
	Form                *bool                  `json:"form,omitempty"`
	Initial             *string                `json:"initial,omitempty"`
	Intro               *string                `json:"intro,omitempty"`
	Language            *string                `json:"language,omitempty"`
	Layout              *string                `json:"layout,omitempty"`
	Math                *bool                  `json:"math,omitempty"`
	Maximize            *bool                  `json:"maximize,omitempty"`
	MessagePeek         *bool                  `json:"messagePeek,omitempty"`
	Meta                map[string]interface{} `json:"meta,omitempty"`
	Name                *string                `json:"name,omitempty"`
	Origin              *string                `json:"origin,omitempty"`
	Placeholder         *string                `json:"placeholder,omitempty"`
	Plugins             *string                `json:"plugins,omitempty"`
	PoweredBy           *bool                  `json:"poweredBy,omitempty"`
	RestartConversation *bool                  `json:"restartConversation,omitempty"`
	SessionDuration     *int64                 `json:"sessionDuration,omitempty"`
	StartFirst          *bool                  `json:"startFirst,omitempty"`
	Stream              *bool                  `json:"stream,omitempty"`
	Theme               *string                `json:"theme,omitempty"`
	Title               *string                `json:"title,omitempty"`
	Tools               *bool                  `json:"tools,omitempty"`
	Unfurl              *bool                  `json:"unfurl,omitempty"`
	Verbose             *bool                  `json:"verbose,omitempty"`
	VoiceIn             *bool                  `json:"voiceIn,omitempty"`
	VoiceOut            *bool                  `json:"voiceOut,omitempty"`
}

// UpdateWidgetIntegrationResponse represents the response from updating a widgetintegration.
type UpdateWidgetIntegrationResponse struct {
	ID *string `json:"id"`
}

// UpdateWidgetIntegration updates an existing widgetintegration.
func (c *Client) UpdateWidgetIntegration(ctx context.Context, id string, input UpdateWidgetIntegrationInput) (*UpdateWidgetIntegrationResponse, error) {
	query := `
		mutation UpdateWidgetIntegration($widgetIntegrationId: ID!, $input: WidgetIntegrationUpdateRequest!) {
			updateWidgetIntegration(widgetIntegrationId: $widgetIntegrationId, input: $input) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"widgetIntegrationId": id,
		"input":               input,
	}

	var response struct {
		UpdateWidgetIntegration *UpdateWidgetIntegrationResponse `json:"updateWidgetIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.UpdateWidgetIntegration, nil
}

// DeleteWidgetIntegrationResponse represents the response from deleting a widgetintegration.
type DeleteWidgetIntegrationResponse struct {
	ID *string `json:"id"`
}

// DeleteWidgetIntegration deletes a widgetintegration.
func (c *Client) DeleteWidgetIntegration(ctx context.Context, id string) (*DeleteWidgetIntegrationResponse, error) {
	query := `
		mutation DeleteWidgetIntegration($widgetIntegrationId: ID!) {
			deleteWidgetIntegration(widgetIntegrationId: $widgetIntegrationId) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"widgetIntegrationId": id,
	}

	var response struct {
		DeleteWidgetIntegration *DeleteWidgetIntegrationResponse `json:"deleteWidgetIntegration"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.DeleteWidgetIntegration, nil
}

// GetWidgetIntegrationResponse represents the response from fetching a widgetintegration.
type GetWidgetIntegrationResponse struct {
	ID                  *string                `json:"id"`
	Alias               *string                `json:"alias,omitempty"`
	Attachments         *bool                  `json:"attachments,omitempty"`
	AutoScroll          *bool                  `json:"autoScroll,omitempty"`
	BlueprintId         *string                `json:"blueprintId,omitempty"`
	BotId               *string                `json:"botId,omitempty"`
	Carousel            *bool                  `json:"carousel,omitempty"`
	ContactCollection   *bool                  `json:"contactCollection,omitempty"`
	Description         *string                `json:"description,omitempty"`
	ExportConversation  *bool                  `json:"exportConversation,omitempty"`
	Form                *bool                  `json:"form,omitempty"`
	Initial             *string                `json:"initial,omitempty"`
	Intro               *string                `json:"intro,omitempty"`
	Language            *string                `json:"language,omitempty"`
	Layout              *string                `json:"layout,omitempty"`
	Math                *bool                  `json:"math,omitempty"`
	Maximize            *bool                  `json:"maximize,omitempty"`
	MessagePeek         *bool                  `json:"messagePeek,omitempty"`
	Meta                map[string]interface{} `json:"meta,omitempty"`
	Name                *string                `json:"name,omitempty"`
	Origin              *string                `json:"origin,omitempty"`
	Placeholder         *string                `json:"placeholder,omitempty"`
	Plugins             *string                `json:"plugins,omitempty"`
	PoweredBy           *bool                  `json:"poweredBy,omitempty"`
	RestartConversation *bool                  `json:"restartConversation,omitempty"`
	SessionDuration     *int64                 `json:"sessionDuration,omitempty"`
	StartFirst          *bool                  `json:"startFirst,omitempty"`
	Stream              *bool                  `json:"stream,omitempty"`
	Theme               *string                `json:"theme,omitempty"`
	Title               *string                `json:"title,omitempty"`
	Tools               *bool                  `json:"tools,omitempty"`
	Unfurl              *bool                  `json:"unfurl,omitempty"`
	Verbose             *bool                  `json:"verbose,omitempty"`
	VoiceIn             *bool                  `json:"voiceIn,omitempty"`
	VoiceOut            *bool                  `json:"voiceOut,omitempty"`
	CreatedAt           *string                `json:"createdAt,omitempty"`
	UpdatedAt           *string                `json:"updatedAt,omitempty"`
}

// GetWidgetIntegration fetches a widgetintegration by ID.
func (c *Client) GetWidgetIntegration(ctx context.Context, id string) (*GetWidgetIntegrationResponse, error) {
	// Note: The GraphQL API uses connection-based queries, so we filter by ID
	query := `
		query GetWidgetIntegration($cursor: ID) {
			widgetIntegrations(first: 1, after: $cursor) {
				edges {
					node {
						id
						alias
						attachments
						autoScroll
						blueprintId
						botId
						carousel
						contactCollection
						description
						exportConversation
						form
						initial
						intro
						language
						layout
						math
						maximize
						messagePeek
						meta
						name
						origin
						placeholder
						plugins
						poweredBy
						restartConversation
						sessionDuration
						startFirst
						stream
						theme
						title
						tools
						unfurl
						verbose
						voiceIn
						voiceOut
						createdAt
						updatedAt
					}
				}
			}
		}
	`

	// For read operations, we need to iterate through results to find by ID
	// This is a simplified implementation - in production, you'd want proper pagination
	variables := map[string]interface{}{}

	var response struct {
		WidgetIntegrations struct {
			Edges []struct {
				Node *GetWidgetIntegrationResponse `json:"node"`
			} `json:"edges"`
		} `json:"widgetIntegrations"`
	}

	if err := c.doRequest(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	// Find the resource with matching ID
	for _, edge := range response.WidgetIntegrations.Edges {
		if edge.Node != nil && edge.Node.ID != nil && *edge.Node.ID == id {
			return edge.Node, nil
		}
	}

	return nil, fmt.Errorf("widgetintegration with ID %s not found", id)
}

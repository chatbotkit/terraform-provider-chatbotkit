package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Helper function to create a pointer to a string
func ptr(s string) *string {
	return &s
}

func TestNewClient(t *testing.T) {
	t.Run("creates client with provided values", func(t *testing.T) {
		client := NewClient("test-api-key", "https://custom.api.com/graphql")

		if client.APIKey != "test-api-key" {
			t.Errorf("expected APIKey to be 'test-api-key', got '%s'", client.APIKey)
		}
		if client.BaseURL != "https://custom.api.com/graphql" {
			t.Errorf("expected BaseURL to be 'https://custom.api.com/graphql', got '%s'", client.BaseURL)
		}
	})

	t.Run("uses default base URL when empty", func(t *testing.T) {
		client := NewClient("test-api-key", "")

		if client.BaseURL != defaultBaseURL {
			t.Errorf("expected BaseURL to be '%s', got '%s'", defaultBaseURL, client.BaseURL)
		}
	})
}

func TestCreateBot(t *testing.T) {
	t.Run("creates bot successfully", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request method and headers
			if r.Method != "POST" {
				t.Errorf("expected POST request, got %s", r.Method)
			}
			if r.Header.Get("Authorization") != "Bearer test-api-key" {
				t.Errorf("expected Authorization header 'Bearer test-api-key', got '%s'", r.Header.Get("Authorization"))
			}
			if r.Header.Get("Content-Type") != "application/json" {
				t.Errorf("expected Content-Type 'application/json', got '%s'", r.Header.Get("Content-Type"))
			}

			// Decode and verify request body
			var req GraphQLRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("failed to decode request body: %v", err)
			}

			// Return mock response
			botID := "bot_123"
			response := map[string]interface{}{
				"data": map[string]interface{}{
					"createBot": map[string]interface{}{
						"id": botID,
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewClient("test-api-key", server.URL)
		result, err := client.CreateBot(context.Background(), CreateBotInput{
			Name:        ptr("Test Bot"),
			Description: ptr("A test bot"),
			Backstory:   ptr("You are a helpful assistant."),
		})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result == nil {
			t.Fatal("expected result, got nil")
		}
		if result.ID == nil || *result.ID != "bot_123" {
			t.Errorf("expected ID 'bot_123', got '%v'", result.ID)
		}
	})

	t.Run("handles GraphQL error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := map[string]interface{}{
				"data": nil,
				"errors": []map[string]interface{}{
					{"message": "Invalid input"},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewClient("test-api-key", server.URL)
		_, err := client.CreateBot(context.Background(), CreateBotInput{
			Name: ptr("Test Bot"),
		})

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "GraphQL error: Invalid input" {
			t.Errorf("expected 'GraphQL error: Invalid input', got '%s'", err.Error())
		}
	})
}

func TestUpdateBot(t *testing.T) {
	t.Run("updates bot successfully", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req GraphQLRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("failed to decode request body: %v", err)
			}

			// Verify the bot ID is passed in variables
			if req.Variables["botId"] != "bot_123" {
				t.Errorf("expected botId 'bot_123', got '%v'", req.Variables["botId"])
			}

			response := map[string]interface{}{
				"data": map[string]interface{}{
					"updateBot": map[string]interface{}{
						"id": "bot_123",
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewClient("test-api-key", server.URL)
		result, err := client.UpdateBot(context.Background(), "bot_123", UpdateBotInput{
			Name:        ptr("Updated Bot"),
			Description: ptr("Updated description"),
		})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result == nil || result.ID == nil || *result.ID != "bot_123" {
			t.Errorf("expected ID 'bot_123', got '%v'", result)
		}
	})
}

func TestDeleteBot(t *testing.T) {
	t.Run("deletes bot successfully", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req GraphQLRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				t.Errorf("failed to decode request body: %v", err)
			}

			// Verify the bot ID is passed in variables
			if req.Variables["botId"] != "bot_123" {
				t.Errorf("expected botId 'bot_123', got '%v'", req.Variables["botId"])
			}

			response := map[string]interface{}{
				"data": map[string]interface{}{
					"deleteBot": map[string]interface{}{
						"id": "bot_123",
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewClient("test-api-key", server.URL)
		result, err := client.DeleteBot(context.Background(), "bot_123")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result == nil || result.ID == nil || *result.ID != "bot_123" {
			t.Errorf("expected ID 'bot_123', got '%v'", result)
		}
	})
}

func TestGetBot(t *testing.T) {
	t.Run("gets bot successfully when found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := map[string]interface{}{
				"data": map[string]interface{}{
					"bots": map[string]interface{}{
						"edges": []map[string]interface{}{
							{
								"node": map[string]interface{}{
									"id":          "bot_123",
									"name":        "Test Bot",
									"description": "A test bot",
									"backstory":   "You are a helpful assistant.",
								},
							},
						},
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewClient("test-api-key", server.URL)
		result, err := client.GetBot(context.Background(), "bot_123")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result == nil {
			t.Fatal("expected result, got nil")
		}
		if result.ID == nil || *result.ID != "bot_123" {
			t.Errorf("expected ID 'bot_123', got '%v'", result.ID)
		}
		if result.Name == nil || *result.Name != "Test Bot" {
			t.Errorf("expected Name 'Test Bot', got '%v'", result.Name)
		}
	})

	t.Run("returns error when bot not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := map[string]interface{}{
				"data": map[string]interface{}{
					"bots": map[string]interface{}{
						"edges": []map[string]interface{}{},
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewClient("test-api-key", server.URL)
		_, err := client.GetBot(context.Background(), "bot_nonexistent")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		expectedErr := "bot with ID bot_nonexistent not found"
		if err.Error() != expectedErr {
			t.Errorf("expected '%s', got '%s'", expectedErr, err.Error())
		}
	})
}

func TestCreateDataset(t *testing.T) {
	t.Run("creates dataset successfully", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := map[string]interface{}{
				"data": map[string]interface{}{
					"createDataset": map[string]interface{}{
						"id": "dataset_123",
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewClient("test-api-key", server.URL)
		result, err := client.CreateDataset(context.Background(), CreateDatasetInput{
			Name:        ptr("Test Dataset"),
			Description: ptr("A test dataset"),
		})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result == nil || result.ID == nil || *result.ID != "dataset_123" {
			t.Errorf("expected ID 'dataset_123', got '%v'", result)
		}
	})
}

func TestCreateSkillset(t *testing.T) {
	t.Run("creates skillset successfully", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := map[string]interface{}{
				"data": map[string]interface{}{
					"createSkillset": map[string]interface{}{
						"id": "skillset_123",
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewClient("test-api-key", server.URL)
		result, err := client.CreateSkillset(context.Background(), CreateSkillsetInput{
			Name:        ptr("Test Skillset"),
			Description: ptr("A test skillset"),
		})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result == nil || result.ID == nil || *result.ID != "skillset_123" {
			t.Errorf("expected ID 'skillset_123', got '%v'", result)
		}
	})
}

func TestCreateBlueprint(t *testing.T) {
	t.Run("creates blueprint successfully", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := map[string]interface{}{
				"data": map[string]interface{}{
					"createBlueprint": map[string]interface{}{
						"id": "blueprint_123",
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
		}))
		defer server.Close()

		client := NewClient("test-api-key", server.URL)
		result, err := client.CreateBlueprint(context.Background(), CreateBlueprintInput{
			Name:        ptr("Test Blueprint"),
			Description: ptr("A test blueprint"),
		})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result == nil || result.ID == nil || *result.ID != "blueprint_123" {
			t.Errorf("expected ID 'blueprint_123', got '%v'", result)
		}
	})
}

func TestDoRequest_HTTPError(t *testing.T) {
	t.Run("handles HTTP error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Internal Server Error"))
		}))
		defer server.Close()

		client := NewClient("test-api-key", server.URL)
		_, err := client.CreateBot(context.Background(), CreateBotInput{
			Name: ptr("Test Bot"),
		})

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		// The error will be about unmarshalling the response
		if err.Error() == "" {
			t.Error("expected non-empty error message")
		}
	})
}

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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
		t.Setenv("CHATBOTKIT_API_URL", "")
		t.Setenv("CBK_API_URL", "")

		client := NewClient("test-api-key", "")

		if client.BaseURL != defaultBaseURL {
			t.Errorf("expected BaseURL to be '%s', got '%s'", defaultBaseURL, client.BaseURL)
		}
	})

	t.Run("derives the GraphQL endpoint from the platform origin in the environment", func(t *testing.T) {
		for origin, want := range map[string]string{
			"http://localhost:3000":     "http://localhost:3000/api/v1/graphql",
			"http://127.0.0.1:4300/":    "http://127.0.0.1:4300/api/v1/graphql",
			"https://corp.example/cbk":  "https://corp.example/cbk/api/v1/graphql",
			"https://corp.example/cbk/": "https://corp.example/cbk/api/v1/graphql",
		} {
			t.Setenv("CHATBOTKIT_API_URL", origin)

			if client := NewClient("test-api-key", ""); client.BaseURL != want {
				t.Errorf("origin %s: expected BaseURL '%s', got '%s'", origin, want, client.BaseURL)
			}
		}
	})

	t.Run("reads the token from the environment under every supported name", func(t *testing.T) {
		names := []string{
			"CHATBOTKIT_API_TOKEN", "CBK_API_TOKEN",
			"CHATBOTKIT_API_SECRET", "CBK_API_SECRET",
			"CHATBOTKIT_API_KEY", "CBK_API_KEY",
		}

		clear := func() {
			for _, name := range names {
				t.Setenv(name, "")
			}
		}

		for _, name := range names {
			clear()
			t.Setenv(name, "from-"+name)

			if client := NewClient("", ""); client.APIKey != "from-"+name {
				t.Errorf("%s: expected the token to be read, got '%s'", name, client.APIKey)
			}
		}

		// @note names earlier in the list win: TOKEN over SECRET over KEY, and a
		// long name over its CBK_ shorthand
		for i := 0; i < len(names)-1; i++ {
			clear()
			t.Setenv(names[i], "winner")
			t.Setenv(names[i+1], "loser")

			if client := NewClient("", ""); client.APIKey != "winner" {
				t.Errorf("expected %s to win over %s, got '%s'", names[i], names[i+1], client.APIKey)
			}
		}

		clear()
		t.Setenv("CHATBOTKIT_API_TOKEN", "from-env")

		if client := NewClient("configured", ""); client.APIKey != "configured" {
			t.Errorf("expected the configured token to win over the environment, got '%s'", client.APIKey)
		}
	})

	t.Run("reads the run-as user under the CLI names and the original one", func(t *testing.T) {
		names := []string{"CHATBOTKIT_API_RUNAS_USERID", "CBK_API_RUNAS_USERID", "CHATBOTKIT_RUN_AS", "CBK_RUN_AS"}

		for i, name := range names {
			for _, other := range names {
				t.Setenv(other, "")
			}

			t.Setenv(name, "user-"+name)

			if got := runAsFromEnv(); got != "user-"+name {
				t.Errorf("%s: expected the run-as user to be read, got '%s'", name, got)
			}

			if i+1 < len(names) {
				t.Setenv(names[i+1], "loser")

				if got := runAsFromEnv(); got != "user-"+name {
					t.Errorf("expected %s to win over %s, got '%s'", name, names[i+1], got)
				}
			}
		}
	})

	t.Run("reads the platform origin from the CBK_API_URL shorthand", func(t *testing.T) {
		t.Setenv("CHATBOTKIT_API_URL", "")
		t.Setenv("CBK_API_URL", "http://localhost:3000")

		if client := NewClient("test-api-key", ""); client.BaseURL != "http://localhost:3000/api/v1/graphql" {
			t.Errorf("expected the shorthand origin to be used, got '%s'", client.BaseURL)
		}
	})

	t.Run("prefers the configured base URL over the environment", func(t *testing.T) {
		t.Setenv("CHATBOTKIT_API_URL", "http://localhost:3000")

		client := NewClient("test-api-key", "http://localhost:9000/graphql")

		if client.BaseURL != "http://localhost:9000/graphql" {
			t.Errorf("expected the configured BaseURL to win, got '%s'", client.BaseURL)
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

// fakeGraphQLRequester serves canned JSON pages keyed by the `after` cursor
// and records the variables it was asked for.
type fakeGraphQLRequester struct {
	pages map[string]string // key: after cursor ("" for the first page)
	calls []map[string]interface{}
}

func (f *fakeGraphQLRequester) doRequest(_ context.Context, _ string, variables map[string]interface{}, result interface{}) error {
	f.calls = append(f.calls, variables)

	key := ""
	if after, ok := variables["after"].(*string); ok && after != nil {
		key = *after
	}

	page, ok := f.pages[key]
	if !ok {
		return fmt.Errorf("no page for cursor %q", key)
	}

	return json.Unmarshal([]byte(page), result)
}

func abilityPage(hasNext bool, endCursor string, ids ...string) string {
	edges := make([]string, 0, len(ids))
	for _, id := range ids {
		edges = append(edges, fmt.Sprintf(`{"node":{"id":%q,"name":"n-%s","linkedSecret":{"id":"secret-%s"}}}`, id, id, id))
	}
	cursor := "null"
	if endCursor != "" {
		cursor = fmt.Sprintf("%q", endCursor)
	}
	return fmt.Sprintf(`{"skillsets":{"edges":[{"node":{"id":"skillset-1","abilities":{"pageInfo":{"hasNextPage":%t,"endCursor":%s},"edges":[%s]}}}]}}`,
		hasNext, cursor, strings.Join(edges, ","))
}

func TestFindSkillsetAbility(t *testing.T) {
	ctx := context.Background()

	t.Run("finds an ability on the first page without paginating", func(t *testing.T) {
		fake := &fakeGraphQLRequester{pages: map[string]string{
			"": abilityPage(true, "c1", "a1", "a2"),
		}}

		got, err := findSkillsetAbility(ctx, fake, "skillset-1", "a2")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ID == nil || *got.ID != "a2" || got.LinkedSecretId == nil || *got.LinkedSecretId != "secret-a2" {
			t.Fatalf("unexpected result: %+v", got)
		}
		if len(fake.calls) != 1 {
			t.Fatalf("expected a single request, got %d", len(fake.calls))
		}
		if fake.calls[0]["first"] != skillsetAbilityPageSize {
			t.Fatalf("expected page size %d, got %v", skillsetAbilityPageSize, fake.calls[0]["first"])
		}
	})

	t.Run("follows endCursor to later pages", func(t *testing.T) {
		fake := &fakeGraphQLRequester{pages: map[string]string{
			"":   abilityPage(true, "c1", "a1"),
			"c1": abilityPage(true, "c2", "a2"),
			"c2": abilityPage(false, "", "a3"),
		}}

		got, err := findSkillsetAbility(ctx, fake, "skillset-1", "a3")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ID == nil || *got.ID != "a3" {
			t.Fatalf("unexpected result: %+v", got)
		}
		if len(fake.calls) != 3 {
			t.Fatalf("expected three requests, got %d", len(fake.calls))
		}
		if after, _ := fake.calls[2]["after"].(*string); after == nil || *after != "c2" {
			t.Fatalf("expected third request to use cursor c2, got %v", fake.calls[2]["after"])
		}
	})

	t.Run("reports not found after exhausting the connection", func(t *testing.T) {
		fake := &fakeGraphQLRequester{pages: map[string]string{
			"":   abilityPage(true, "c1", "a1"),
			"c1": abilityPage(false, "", "a2"),
		}}

		_, err := findSkillsetAbility(ctx, fake, "skillset-1", "missing")
		if err == nil || !strings.Contains(err.Error(), "not found") {
			t.Fatalf("expected a not found error, got %v", err)
		}
		if len(fake.calls) != 2 {
			t.Fatalf("expected two requests, got %d", len(fake.calls))
		}
	})

	t.Run("stops when the cursor does not advance", func(t *testing.T) {
		fake := &fakeGraphQLRequester{pages: map[string]string{
			"":   abilityPage(true, "c1", "a1"),
			"c1": abilityPage(true, "c1", "a1"),
		}}

		_, err := findSkillsetAbility(ctx, fake, "skillset-1", "missing")
		if err == nil || !strings.Contains(err.Error(), "not found") {
			t.Fatalf("expected a not found error, got %v", err)
		}
		if len(fake.calls) != 2 {
			t.Fatalf("expected the loop to stop on a repeated cursor, got %d requests", len(fake.calls))
		}
	})

	t.Run("reports a missing skillset", func(t *testing.T) {
		fake := &fakeGraphQLRequester{pages: map[string]string{"": `{"skillsets":{"edges":[]}}`}}

		_, err := findSkillsetAbility(ctx, fake, "skillset-x", "a1")
		if err == nil || !strings.Contains(err.Error(), "skillset with ID skillset-x not found") {
			t.Fatalf("expected a skillset not found error, got %v", err)
		}
	})
}

package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewClient(t *testing.T) {
	baseURL := "http://localhost:8080"
	auth := AuthConfig{APIKey: "test-key"}

	client := NewClient(baseURL, auth)

	if client.baseURL != baseURL {
		t.Errorf("Expected baseURL %s, got %s", baseURL, client.baseURL)
	}

	if client.auth.APIKey != auth.APIKey {
		t.Errorf("Expected APIKey %s, got %s", auth.APIKey, client.auth.APIKey)
	}

	if client.httpClient.Timeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", client.httpClient.Timeout)
	}
}

func TestNewClientWithHTTPClient(t *testing.T) {
	baseURL := "http://localhost:8080"
	auth := AuthConfig{BearerToken: "test-token"}
	httpClient := &http.Client{Timeout: 10 * time.Second}

	client := NewClientWithHTTPClient(baseURL, auth, httpClient)

	if client.baseURL != baseURL {
		t.Errorf("Expected baseURL %s, got %s", baseURL, client.baseURL)
	}

	if client.auth.BearerToken != auth.BearerToken {
		t.Errorf("Expected BearerToken %s, got %s", auth.BearerToken, client.auth.BearerToken)
	}

	if client.httpClient != httpClient {
		t.Error("Expected custom HTTP client to be used")
	}
}

func TestHealth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			t.Errorf("Expected path /health, got %s", r.URL.Path)
		}
		if r.Header.Get("X-API-Key") != "test-key" {
			t.Errorf("Expected API key header, got %s", r.Header.Get("X-API-Key"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"database": "ok", "redis": "ok"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, AuthConfig{APIKey: "test-key"})

	health, err := client.Health(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if health.Database != "ok" {
		t.Errorf("Expected database ok, got %s", health.Database)
	}
	if health.Redis != "ok" {
		t.Errorf("Expected redis ok, got %s", health.Redis)
	}
}

func TestEvaluateScript(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/evaluate" {
			t.Errorf("Expected path /api/v1/evaluate, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Expected Bearer auth, got %s", r.Header.Get("Authorization"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"success": true,
			"result": 5,
			"output": ["info", "test"],
			"duration": "1.2ms"
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, AuthConfig{BearerToken: "test-token"})

	result, err := client.EvaluateScript(context.Background(), EvaluateScriptRequest{
		Script: "return 2 + 3",
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !result.Success {
		t.Error("Expected success true")
	}
	if result.Result != 5.0 {
		t.Errorf("Expected result 5, got %v", result.Result)
	}
	if len(result.Output) != 2 {
		t.Errorf("Expected 2 output items, got %d", len(result.Output))
	}
}

func TestAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{
			"error": {
				"code": "VALIDATION_ERROR",
				"message": "Script is required"
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, AuthConfig{})

	_, err := client.EvaluateScript(context.Background(), EvaluateScriptRequest{})
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("Expected APIError, got %T", err)
	}

	if apiErr.Code != "VALIDATION_ERROR" {
		t.Errorf("Expected code VALIDATION_ERROR, got %s", apiErr.Code)
	}
	if apiErr.Message != "Script is required" {
		t.Errorf("Expected message 'Script is required', got %s", apiErr.Message)
	}
	if apiErr.StatusCode != 400 {
		t.Errorf("Expected status code 400, got %d", apiErr.StatusCode)
	}
}

func TestCreateRule(t *testing.T) {
	ruleID := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/rules" {
			t.Errorf("Expected path /api/v1/rules, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{
			"id": "` + ruleID.String() + `",
			"name": "Test Rule",
			"lua_script": "return true",
			"priority": 0,
			"enabled": true,
			"created_at": "2023-01-01T00:00:00Z",
			"updated_at": "2023-01-01T00:00:00Z"
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, AuthConfig{APIKey: "test"})

	rule, err := client.CreateRule(context.Background(), CreateRuleRequest{
		Name:      "Test Rule",
		LuaScript: "return true",
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if rule.ID != ruleID {
		t.Errorf("Expected rule ID %s, got %s", ruleID, rule.ID)
	}
	if rule.Name != "Test Rule" {
		t.Errorf("Expected rule name 'Test Rule', got %s", rule.Name)
	}
}

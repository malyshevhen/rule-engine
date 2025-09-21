package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents the Rule Engine API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	auth       AuthConfig
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	BearerToken string
	APIKey      string
}

// NewClient creates a new Rule Engine API client
func NewClient(baseURL string, auth AuthConfig) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		auth: auth,
	}
}

// NewClientWithHTTPClient creates a new client with a custom HTTP client
func NewClientWithHTTPClient(baseURL string, auth AuthConfig, httpClient *http.Client) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
		auth:       auth,
	}
}

// doRequest performs an HTTP request with authentication
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	fullURL := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Set authentication
	if c.auth.BearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.auth.BearerToken)
	}
	if c.auth.APIKey != "" {
		req.Header.Set("X-API-Key", c.auth.APIKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// parseResponse parses a JSON response into the provided interface
func parseResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, resp.Status)
		}
		return &APIError{
			Code:       errResp.Error.Code,
			Message:    errResp.Error.Message,
			StatusCode: resp.StatusCode,
		}
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// APIError represents an error returned by the API
type APIError struct {
	Code       string
	Message    string
	StatusCode int
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (%d): %s - %s", e.StatusCode, e.Code, e.Message)
}

// Health checks the health of the service
func (c *Client) Health(ctx context.Context) (*HealthResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/health", nil)
	if err != nil {
		return nil, err
	}

	var health HealthResponse
	if err := parseResponse(resp, &health); err != nil {
		return nil, err
	}

	return &health, nil
}

// Metrics retrieves Prometheus metrics
func (c *Client) Metrics(ctx context.Context) (string, error) {
	resp, err := c.doRequest(ctx, "GET", "/metrics", nil)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}

// EvaluateScript evaluates a Lua script
func (c *Client) EvaluateScript(ctx context.Context, req EvaluateScriptRequest) (*EvaluateScriptResponse, error) {
	resp, err := c.doRequest(ctx, "POST", "/api/v1/evaluate", req)
	if err != nil {
		return nil, err
	}

	var result EvaluateScriptResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

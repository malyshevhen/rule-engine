package alerting

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_SendAlert_Disabled(t *testing.T) {
	config := Config{Enabled: false}
	svc := NewService(config)

	err := svc.SendAlert(context.Background(), "test_type", "low", "Test Title", "Test Message", nil)
	assert.NoError(t, err) // Should not error when disabled
}

func TestService_SendAlert_NoWebhookURL(t *testing.T) {
	config := Config{Enabled: true, WebhookURL: ""}
	svc := NewService(config)

	err := svc.SendAlert(context.Background(), "test_type", "low", "Test Title", "Test Message", nil)
	assert.NoError(t, err) // Should not error when no webhook URL
}

func TestService_SendAlert_Success(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var alert Alert
		err := json.NewDecoder(r.Body).Decode(&alert)
		require.NoError(t, err)

		assert.Equal(t, AlertType("test_type"), alert.Type)
		assert.Equal(t, "low", alert.Severity)
		assert.Equal(t, "Test Title", alert.Title)
		assert.Equal(t, "Test Message", alert.Message)
		assert.NotEmpty(t, alert.ID)
		assert.False(t, alert.Timestamp.IsZero())

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := Config{
		Enabled:    true,
		WebhookURL: server.URL,
	}
	svc := NewService(config)

	err := svc.SendAlert(context.Background(), "test_type", "low", "Test Title", "Test Message", map[string]interface{}{
		"key": "value",
	})

	assert.NoError(t, err)
}

func TestService_SendAlert_WithDetails(t *testing.T) {
	var receivedAlert Alert
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedAlert)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := Config{
		Enabled:    true,
		WebhookURL: server.URL,
	}
	svc := NewService(config)

	details := map[string]interface{}{
		"rule_id":   "test-rule-id",
		"error":     "test error",
		"timestamp": "2023-01-01T00:00:00Z",
	}

	err := svc.SendAlert(context.Background(), "rule_execution_failure", "high", "Rule Failed", "Rule execution failed", details)

	assert.NoError(t, err)
	assert.Equal(t, AlertType("rule_execution_failure"), receivedAlert.Type)
	assert.Equal(t, "high", receivedAlert.Severity)
	assert.Equal(t, "Rule Failed", receivedAlert.Title)
	assert.Equal(t, details, receivedAlert.Details)
}

func TestService_SendAlert_RetryOnFailure(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 2 {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	config := Config{
		Enabled:       true,
		WebhookURL:    server.URL,
		RetryAttempts: 3,
		RetryDelay:    1 * time.Millisecond, // Fast retry for test
	}
	svc := NewService(config)

	err := svc.SendAlert(context.Background(), "test_type", "medium", "Test", "Message", nil)

	assert.NoError(t, err)
	assert.Equal(t, 2, attempts) // Should succeed on second attempt
}

func TestService_SendAlert_MaxRetriesExceeded(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	config := Config{
		Enabled:       true,
		WebhookURL:    server.URL,
		RetryAttempts: 2,
		RetryDelay:    1 * time.Millisecond,
	}
	svc := NewService(config)

	err := svc.SendAlert(context.Background(), "test_type", "high", "Test", "Message", nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send alert after 2 attempts")
	assert.Equal(t, 2, attempts)
}

func TestService_SendAlert_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond) // Longer than timeout
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := Config{
		Enabled:        true,
		WebhookURL:     server.URL,
		RequestTimeout: 10 * time.Millisecond, // Very short timeout
	}
	svc := NewService(config)

	err := svc.SendAlert(context.Background(), "test_type", "critical", "Test", "Message", nil)

	assert.Error(t, err)
}

func TestNewService_DefaultConfig(t *testing.T) {
	config := Config{Enabled: true, WebhookURL: "http://example.com"}
	svc := NewService(config)

	// Should have default values
	assert.NotNil(t, svc)
	assert.NotNil(t, svc.client)
}

func TestAlertType(t *testing.T) {
	assert.Equal(t, AlertType("rule_execution_failure"), AlertTypeRuleExecutionFailure)
	assert.Equal(t, AlertType("trigger_evaluation_failure"), AlertTypeTriggerEvaluationFailure)
	assert.Equal(t, AlertType("queue_processing_failure"), AlertTypeQueueProcessingFailure)
	assert.Equal(t, AlertType("high_error_rate"), AlertTypeHighErrorRate)
}

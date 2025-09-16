package alerting

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/malyshevhen/rule-engine/internal/metrics"
)

// AlertType represents the type of alert
type AlertType string

const (
	AlertTypeRuleExecutionFailure     AlertType = "rule_execution_failure"
	AlertTypeTriggerEvaluationFailure AlertType = "trigger_evaluation_failure"
	AlertTypeQueueProcessingFailure   AlertType = "queue_processing_failure"
	AlertTypeHighErrorRate            AlertType = "high_error_rate"
)

// Alert represents an alert notification
type Alert struct {
	ID        uuid.UUID              `json:"id"`
	Type      AlertType              `json:"type"`
	Severity  string                 `json:"severity"` // "low", "medium", "high", "critical"
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// Config holds alerting configuration
type Config struct {
	WebhookURL     string
	Enabled        bool
	RetryAttempts  int
	RetryDelay     time.Duration
	RequestTimeout time.Duration
}

// Service handles sending alert notifications
type Service struct {
	config Config
	client *http.Client
}

// NewService creates a new alerting service
func NewService(config Config) *Service {
	if config.RetryAttempts == 0 {
		config.RetryAttempts = 3
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = 5 * time.Second
	}
	if config.RequestTimeout == 0 {
		config.RequestTimeout = 30 * time.Second
	}

	return &Service{
		config: config,
		client: &http.Client{
			Timeout: config.RequestTimeout,
		},
	}
}

// SendAlert sends an alert notification via webhook
func (s *Service) SendAlert(ctx context.Context, alertType string, severity, title, message string, details map[string]interface{}) error {
	if !s.config.Enabled || s.config.WebhookURL == "" {
		slog.Debug("Alerting disabled or no webhook URL configured", "type", alertType, "severity", severity)
		return nil
	}

	alert := Alert{
		ID:        uuid.New(),
		Type:      AlertType(alertType),
		Severity:  severity,
		Title:     title,
		Message:   message,
		Details:   details,
		Timestamp: time.Now(),
	}

	slog.Info("Sending alert notification",
		"alert_id", alert.ID,
		"type", alert.Type,
		"severity", alert.Severity,
		"title", alert.Title)

	err := s.sendWebhookWithRetry(ctx, alert)
	if err != nil {
		metrics.AlertSendErrorsTotal.WithLabelValues(string(alertType)).Inc()
		return err
	}

	// Record successful alert
	metrics.AlertsTotal.WithLabelValues(string(alertType), severity).Inc()
	return nil
}

// sendWebhookWithRetry sends the webhook with retry logic
func (s *Service) sendWebhookWithRetry(ctx context.Context, alert Alert) error {
	payload, err := json.Marshal(alert)
	if err != nil {
		return fmt.Errorf("failed to marshal alert: %w", err)
	}

	for attempt := 1; attempt <= s.config.RetryAttempts; attempt++ {
		if err := s.sendWebhook(ctx, payload); err != nil {
			slog.Warn("Failed to send alert webhook",
				"alert_id", alert.ID,
				"attempt", attempt,
				"max_attempts", s.config.RetryAttempts,
				"error", err)

			if attempt < s.config.RetryAttempts {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(s.config.RetryDelay):
					continue
				}
			}
			return fmt.Errorf("failed to send alert after %d attempts: %w", s.config.RetryAttempts, err)
		}

		slog.Info("Alert webhook sent successfully",
			"alert_id", alert.ID,
			"attempt", attempt)
		return nil
	}

	return nil
}

// sendWebhook sends a single webhook request
func (s *Service) sendWebhook(ctx context.Context, payload []byte) error {
	req, err := http.NewRequestWithContext(ctx, "POST", s.config.WebhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "rule-engine-alerting/1.0")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned non-2xx status: %d", resp.StatusCode)
	}

	return nil
}

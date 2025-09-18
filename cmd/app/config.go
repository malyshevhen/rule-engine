package app

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	Port               string
	DBURL              string
	NATSURL            string
	RedisURL           string
	AlertingEnabled    bool
	AlertWebhookURL    string
	AlertRetryAttempts int
}

// loadConfig loads configuration from environment variables
func loadConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Try to build DATABASE_URL from individual environment variables
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbName := os.Getenv("DB_NAME")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbSSLMode := os.Getenv("DB_SSL_MODE")

		if dbHost != "" && dbPort != "" && dbName != "" && dbUser != "" && dbPassword != "" {
			dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
			if dbSSLMode != "" {
				dbURL += "?sslmode=" + dbSSLMode
			} else {
				dbURL += "?sslmode=disable"
			}
		} else {
			dbURL = "postgres://user:password@localhost:5432/rule_engine?sslmode=disable"
		}
	}
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	// Alerting configuration
	alertingEnabled := os.Getenv("ALERTING_ENABLED") == "true"
	alertWebhookURL := os.Getenv("ALERT_WEBHOOK_URL")
	alertRetryAttempts := 3 // default
	if retryStr := os.Getenv("ALERT_RETRY_ATTEMPTS"); retryStr != "" {
		if retry, err := strconv.Atoi(retryStr); err == nil && retry > 0 {
			alertRetryAttempts = retry
		}
	}

	return Config{
		Port:               port,
		DBURL:              dbURL,
		NATSURL:            natsURL,
		RedisURL:           redisURL,
		AlertingEnabled:    alertingEnabled,
		AlertWebhookURL:    alertWebhookURL,
		AlertRetryAttempts: alertRetryAttempts,
	}
}

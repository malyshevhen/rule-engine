package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// ParseJSONBody parses JSON request body into the provided interface
func ParseJSONBody(r *http.Request, v any) error {
	defer func() {
		if err := r.Body.Close(); err != nil {
			slog.Error("Failed to close request body", "error", err)
		}
	}()
	return json.NewDecoder(r.Body).Decode(v)
}

// GetQueryParam gets a query parameter from the request
func GetQueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// GetHeader gets a header from the request
func GetHeader(r *http.Request, key string) string {
	return r.Header.Get(key)
}

package api

import (
	"encoding/json"
	"net/http"
)

// ParseJSONBody parses JSON request body into the provided interface
func ParseJSONBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
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

package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// JSONResponse sends a JSON response with the given status code and data
func JSONResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Failed to encode response", "error", err)
		// Avoid recursive call, send simple error
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// ErrorResponse sends an error response
func ErrorResponse(w http.ResponseWriter, status int, code, message string) {
	errorObj := map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	}
	JSONResponse(w, status, errorObj)
}

// SuccessResponse sends a success response with data
func SuccessResponse(w http.ResponseWriter, data any) {
	JSONResponse(w, http.StatusOK, data)
}

// CreatedResponse sends a created response with data
func CreatedResponse(w http.ResponseWriter, data any) {
	JSONResponse(w, http.StatusCreated, data)
}

package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	jsonpatch "github.com/evanphx/json-patch/v5"
)

// ApplyJSONPatch applies JSON Patch operations to the given JSON data and returns the modified data.
// It handles parsing, validation, and application of patch operations with proper error responses.
func ApplyJSONPatch(r *http.Request, originalJSON []byte, resourceType string, resourceID string) ([]byte, error) {
	// Parse JSON Patch operations
	var patchOps PatchRequest
	if err := ParseJSONBody(r, &patchOps); err != nil {
		slog.Error("Failed to parse JSON patch body", resourceType+"_id", resourceID, "error", err)
		return nil, fmt.Errorf("invalid JSON patch request body: %w", err)
	}

	// Validate patch operations
	if err := validatePatchOperations(patchOps); err != nil {
		return nil, err
	}

	// Convert patch operations to JSON Patch format
	patchJSON, err := json.Marshal(patchOps)
	if err != nil {
		return nil, fmt.Errorf("invalid patch operations")
	}

	// Apply the patch
	patch, err := jsonpatch.DecodePatch(patchJSON)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON Patch format: %w", err)
	}

	modifiedJSON, err := patch.Apply(originalJSON)
	if err != nil {
		// Provide more specific error messages based on the type of patch error
		if strings.Contains(err.Error(), "path") {
			return nil, fmt.Errorf("invalid patch path: %w", err)
		} else if strings.Contains(err.Error(), "operation") {
			return nil, fmt.Errorf("invalid patch operation: %w", err)
		} else {
			return nil, fmt.Errorf("patch application failed: %w", err)
		}
	}

	return modifiedJSON, nil
}

// validatePatchOperations validates the basic structure of patch operations
func validatePatchOperations(patchOps PatchRequest) error {
	if len(patchOps) == 0 {
		return fmt.Errorf("at least one patch operation is required")
	}

	// Basic validation of patch operations
	for i, op := range patchOps {
		if strings.TrimSpace(op.Path) == "" {
			return fmt.Errorf("patch operation %d: path cannot be empty", i+1)
		}
		// Validate path format (should start with /)
		if !strings.HasPrefix(op.Path, "/") {
			return fmt.Errorf("patch operation %d: path must start with '/'", i+1)
		}
	}

	return nil
}

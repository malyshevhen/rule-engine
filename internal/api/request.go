package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
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

// Global validator instance
var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom validation functions
	if err := validate.RegisterValidation("lua_script_length", validateLuaScriptLength); err != nil {
		panic("Failed to register lua_script_length validation: " + err.Error())
	}
	if err := validate.RegisterValidation("rule_name_length", validateRuleNameLength); err != nil {
		panic("Failed to register rule_name_length validation: " + err.Error())
	}
}

// validateLuaScriptLength validates Lua script length
func validateLuaScriptLength(fl validator.FieldLevel) bool {
	script := fl.Field().String()
	return len(strings.TrimSpace(script)) > 0 && len(script) <= apiConfig.MaxLuaScriptLength
}

// validateRuleNameLength validates rule name length
func validateRuleNameLength(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	return len(strings.TrimSpace(name)) > 0 && len(name) <= apiConfig.MaxRuleNameLength
}

// ValidateStruct validates a struct using the validator tags
func ValidateStruct(s any) error {
	return validate.Struct(s)
}

// ValidateAndParseJSON parses JSON and validates the struct
func ValidateAndParseJSON(r *http.Request, v any) error {
	if err := ParseJSONBody(r, v); err != nil {
		return err
	}

	if err := ValidateStruct(v); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return fmt.Errorf("validation failed: %s", formatValidationErrors(validationErrors))
		}
		return err
	}

	return nil
}

// formatValidationErrors formats validation errors into a readable string
func formatValidationErrors(errors validator.ValidationErrors) string {
	var messages []string
	for _, err := range errors {
		switch err.Tag() {
		case "required":
			messages = append(messages, fmt.Sprintf("%s is required", err.Field()))
		case "lua_script_length":
			messages = append(messages, fmt.Sprintf("Lua script must be between 1 and %d characters", apiConfig.MaxLuaScriptLength))
		case "rule_name_length":
			messages = append(messages, fmt.Sprintf("Rule name must be between 1 and %d characters", apiConfig.MaxRuleNameLength))
		case "oneof":
			messages = append(messages, fmt.Sprintf("%s must be one of: %s", err.Field(), err.Param()))
		default:
			messages = append(messages, fmt.Sprintf("%s is invalid", err.Field()))
		}
	}
	return strings.Join(messages, "; ")
}

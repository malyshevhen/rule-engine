package client

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/google/uuid"
)

// CreateRule creates a new rule
func (c *Client) CreateRule(ctx context.Context, req CreateRuleRequest) (*RuleInfo, error) {
	resp, err := c.doRequest(ctx, "POST", "/api/v1/rules", req)
	if err != nil {
		return nil, err
	}

	var rule RuleInfo
	if err := parseResponse(resp, &rule); err != nil {
		return nil, err
	}

	return &rule, nil
}

// GetRule retrieves a rule by ID
func (c *Client) GetRule(ctx context.Context, id uuid.UUID) (*RuleInfo, error) {
	resp, err := c.doRequest(ctx, "GET", fmt.Sprintf("/api/v1/rules/%s", id.String()), nil)
	if err != nil {
		return nil, err
	}

	var rule RuleInfo
	if err := parseResponse(resp, &rule); err != nil {
		return nil, err
	}

	return &rule, nil
}

// ListRules retrieves a paginated list of rules
func (c *Client) ListRules(ctx context.Context, limit, offset int) (*PaginatedRulesResponse, error) {
	params := url.Values{}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		params.Set("offset", strconv.Itoa(offset))
	}

	path := "/api/v1/rules"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result PaginatedRulesResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateRule updates a rule using JSON Patch
func (c *Client) UpdateRule(ctx context.Context, id uuid.UUID, patches PatchRequest) (*RuleInfo, error) {
	resp, err := c.doRequest(ctx, "PATCH", fmt.Sprintf("/api/v1/rules/%s", id.String()), patches)
	if err != nil {
		return nil, err
	}

	var rule RuleInfo
	if err := parseResponse(resp, &rule); err != nil {
		return nil, err
	}

	return &rule, nil
}

// DeleteRule deletes a rule by ID
func (c *Client) DeleteRule(ctx context.Context, id uuid.UUID) error {
	resp, err := c.doRequest(ctx, "DELETE", fmt.Sprintf("/api/v1/rules/%s", id.String()), nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return parseResponse(resp, nil)
	}

	return nil
}

// AddActionToRule adds an action to a rule
func (c *Client) AddActionToRule(ctx context.Context, ruleID uuid.UUID, req AddActionToRuleRequest) error {
	resp, err := c.doRequest(ctx, "POST", fmt.Sprintf("/api/v1/rules/%s/actions", ruleID.String()), req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return parseResponse(resp, nil)
	}

	return nil
}

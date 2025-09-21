package client

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/google/uuid"
)

// CreateTrigger creates a new trigger
func (c *Client) CreateTrigger(ctx context.Context, req CreateTriggerRequest) (*TriggerInfo, error) {
	resp, err := c.doRequest(ctx, "POST", "/api/v1/triggers", req)
	if err != nil {
		return nil, err
	}

	var trigger TriggerInfo
	if err := parseResponse(resp, &trigger); err != nil {
		return nil, err
	}

	return &trigger, nil
}

// GetTrigger retrieves a trigger by ID
func (c *Client) GetTrigger(ctx context.Context, id uuid.UUID) (*TriggerInfo, error) {
	resp, err := c.doRequest(ctx, "GET", fmt.Sprintf("/api/v1/triggers/%s", id.String()), nil)
	if err != nil {
		return nil, err
	}

	var trigger TriggerInfo
	if err := parseResponse(resp, &trigger); err != nil {
		return nil, err
	}

	return &trigger, nil
}

// ListTriggers retrieves a paginated list of triggers
func (c *Client) ListTriggers(ctx context.Context, limit, offset int) (*PaginatedTriggersResponse, error) {
	params := url.Values{}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		params.Set("offset", strconv.Itoa(offset))
	}

	path := "/api/v1/triggers"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result PaginatedTriggersResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteTrigger deletes a trigger by ID
func (c *Client) DeleteTrigger(ctx context.Context, id uuid.UUID) error {
	resp, err := c.doRequest(ctx, "DELETE", fmt.Sprintf("/api/v1/triggers/%s", id.String()), nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return parseResponse(resp, nil)
	}

	return nil
}

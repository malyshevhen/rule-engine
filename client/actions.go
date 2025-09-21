package client

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/google/uuid"
)

// CreateAction creates a new action
func (c *Client) CreateAction(ctx context.Context, req CreateActionRequest) (*ActionInfo, error) {
	resp, err := c.doRequest(ctx, "POST", "/api/v1/actions", req)
	if err != nil {
		return nil, err
	}

	var action ActionInfo
	if err := parseResponse(resp, &action); err != nil {
		return nil, err
	}

	return &action, nil
}

// GetAction retrieves an action by ID
func (c *Client) GetAction(ctx context.Context, id uuid.UUID) (*ActionInfo, error) {
	resp, err := c.doRequest(ctx, "GET", fmt.Sprintf("/api/v1/actions/%s", id.String()), nil)
	if err != nil {
		return nil, err
	}

	var action ActionInfo
	if err := parseResponse(resp, &action); err != nil {
		return nil, err
	}

	return &action, nil
}

// ListActions retrieves a paginated list of actions
func (c *Client) ListActions(ctx context.Context, limit, offset int) (*PaginatedActionsResponse, error) {
	params := url.Values{}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		params.Set("offset", strconv.Itoa(offset))
	}

	path := "/api/v1/actions"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	resp, err := c.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result PaginatedActionsResponse
	if err := parseResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteAction deletes an action by ID
func (c *Client) DeleteAction(ctx context.Context, id uuid.UUID) error {
	resp, err := c.doRequest(ctx, "DELETE", fmt.Sprintf("/api/v1/actions/%s", id.String()), nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return parseResponse(resp, nil)
	}

	return nil
}

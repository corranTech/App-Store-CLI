package asc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetGameCenterDetails retrieves the list of Game Center details.
func (c *Client) GetGameCenterDetails(ctx context.Context, opts ...GCDetailsOption) (*GameCenterDetailsResponse, error) {
	query := &gcDetailsQuery{}
	for _, opt := range opts {
		opt(query)
	}

	path := "/v1/gameCenterDetails"
	if query.nextURL != "" {
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("game-center-details: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildGCDetailsQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var response GameCenterDetailsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetGameCenterDetail retrieves a Game Center detail by ID.
func (c *Client) GetGameCenterDetail(ctx context.Context, detailID string) (*GameCenterDetailResponse, error) {
	path := fmt.Sprintf("/v1/gameCenterDetails/%s", strings.TrimSpace(detailID))
	data, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var response GameCenterDetailResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetGameCenterDetailGameCenterGroup retrieves the related Game Center group.
func (c *Client) GetGameCenterDetailGameCenterGroup(ctx context.Context, detailID string) (*GameCenterGroupResponse, error) {
	path := fmt.Sprintf("/v1/gameCenterDetails/%s/gameCenterGroup", strings.TrimSpace(detailID))
	data, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var response GameCenterGroupResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetGameCenterGroupGameCenterDetails retrieves Game Center details for a group.
func (c *Client) GetGameCenterGroupGameCenterDetails(ctx context.Context, groupID string, opts ...GCDetailsOption) (*GameCenterDetailsResponse, error) {
	query := &gcDetailsQuery{}
	for _, opt := range opts {
		opt(query)
	}

	groupID = strings.TrimSpace(groupID)
	if query.nextURL == "" && groupID == "" {
		return nil, fmt.Errorf("groupID is required")
	}

	path := fmt.Sprintf("/v1/gameCenterGroups/%s/gameCenterDetails", groupID)
	if query.nextURL != "" {
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("game-center-group-details: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildGCDetailsQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var response GameCenterDetailsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

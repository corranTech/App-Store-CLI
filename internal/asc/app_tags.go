package asc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// AppTagAttributes describes an app tag resource.
type AppTagAttributes struct {
	Name             string `json:"name,omitempty"`
	VisibleInAppStore bool   `json:"visibleInAppStore,omitempty"`
}

// AppTagsResponse is the response from app tags endpoints.
type AppTagsResponse = Response[AppTagAttributes]

// AppTagResponse is the response from app tag detail/updates.
type AppTagResponse = SingleResponse[AppTagAttributes]

// AppTagUpdateAttributes describes fields for updating an app tag.
type AppTagUpdateAttributes struct {
	VisibleInAppStore *bool `json:"visibleInAppStore,omitempty"`
}

// AppTagUpdateData is the data portion of an app tag update request.
type AppTagUpdateData struct {
	Type       ResourceType           `json:"type"`
	ID         string                 `json:"id"`
	Attributes *AppTagUpdateAttributes `json:"attributes,omitempty"`
}

// AppTagUpdateRequest is a request to update an app tag.
type AppTagUpdateRequest struct {
	Data AppTagUpdateData `json:"data"`
}

// GetAppTags retrieves the list of app tags for an app.
func (c *Client) GetAppTags(ctx context.Context, appID string, opts ...AppTagsOption) (*AppTagsResponse, error) {
	query := &appTagsQuery{}
	for _, opt := range opts {
		opt(query)
	}

	appID = strings.TrimSpace(appID)
	path := fmt.Sprintf("/v1/apps/%s/appTags", appID)
	if query.nextURL != "" {
		// Validate nextURL to prevent credential exfiltration
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("appTags: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildAppTagsQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response AppTagsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// UpdateAppTag updates an app tag by ID.
func (c *Client) UpdateAppTag(ctx context.Context, tagID string, attrs AppTagUpdateAttributes) (*AppTagResponse, error) {
	tagID = strings.TrimSpace(tagID)
	payload := AppTagUpdateRequest{
		Data: AppTagUpdateData{
			Type: ResourceTypeAppTags,
			ID:   tagID,
		},
	}
	if attrs.VisibleInAppStore != nil {
		payload.Data.Attributes = &attrs
	}

	body, err := BuildRequestBody(payload)
	if err != nil {
		return nil, err
	}

	data, err := c.do(ctx, "PATCH", fmt.Sprintf("/v1/appTags/%s", tagID), body)
	if err != nil {
		return nil, err
	}

	var response AppTagResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

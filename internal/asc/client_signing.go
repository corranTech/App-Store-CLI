package asc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// GetBundleIDs retrieves the list of bundle IDs.
func (c *Client) GetBundleIDs(ctx context.Context, opts ...BundleIDsOption) (*BundleIDsResponse, error) {
	query := &bundleIDsQuery{}
	for _, opt := range opts {
		opt(query)
	}

	path := "/v1/bundleIds"
	if query.nextURL != "" {
		// Validate nextURL to prevent credential exfiltration
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("bundleIds: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildBundleIDsQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response BundleIDsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetBundleID retrieves a single bundle ID by ID.
func (c *Client) GetBundleID(ctx context.Context, id string) (*BundleIDResponse, error) {
	id = strings.TrimSpace(id)
	path := fmt.Sprintf("/v1/bundleIds/%s", id)
	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response BundleIDResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// CreateBundleID creates a new bundle ID.
func (c *Client) CreateBundleID(ctx context.Context, attrs BundleIDCreateAttributes) (*BundleIDResponse, error) {
	request := BundleIDCreateRequest{
		Data: BundleIDCreateData{
			Type:       ResourceTypeBundleIds,
			Attributes: attrs,
		},
	}

	body, err := BuildRequestBody(request)
	if err != nil {
		return nil, err
	}

	data, err := c.do(ctx, "POST", "/v1/bundleIds", body)
	if err != nil {
		return nil, err
	}

	var response BundleIDResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// UpdateBundleID updates an existing bundle ID.
func (c *Client) UpdateBundleID(ctx context.Context, id string, attrs BundleIDUpdateAttributes) (*BundleIDResponse, error) {
	id = strings.TrimSpace(id)
	request := BundleIDUpdateRequest{
		Data: BundleIDUpdateData{
			Type:       ResourceTypeBundleIds,
			ID:         id,
			Attributes: &attrs,
		},
	}

	body, err := BuildRequestBody(request)
	if err != nil {
		return nil, err
	}

	data, err := c.do(ctx, "PATCH", fmt.Sprintf("/v1/bundleIds/%s", id), body)
	if err != nil {
		return nil, err
	}

	var response BundleIDResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// DeleteBundleID deletes a bundle ID by ID.
func (c *Client) DeleteBundleID(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	path := fmt.Sprintf("/v1/bundleIds/%s", id)
	_, err := c.do(ctx, "DELETE", path, nil)
	return err
}

// GetBundleIDCapabilities retrieves capabilities for a bundle ID.
func (c *Client) GetBundleIDCapabilities(ctx context.Context, bundleID string, opts ...BundleIDCapabilitiesOption) (*BundleIDCapabilitiesResponse, error) {
	bundleID = strings.TrimSpace(bundleID)
	query := &bundleIDCapabilitiesQuery{}
	for _, opt := range opts {
		opt(query)
	}

	path := fmt.Sprintf("/v1/bundleIds/%s/bundleIdCapabilities", bundleID)
	if query.nextURL != "" {
		// Validate nextURL to prevent credential exfiltration
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("bundleIdCapabilities: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildBundleIDCapabilitiesQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response BundleIDCapabilitiesResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// CreateBundleIDCapability adds a capability to a bundle ID.
func (c *Client) CreateBundleIDCapability(ctx context.Context, bundleID string, attrs BundleIDCapabilityCreateAttributes) (*BundleIDCapabilityResponse, error) {
	bundleID = strings.TrimSpace(bundleID)
	request := BundleIDCapabilityCreateRequest{
		Data: BundleIDCapabilityCreateData{
			Type:       ResourceTypeBundleIdCapabilities,
			Attributes: attrs,
			Relationships: &BundleIDCapabilityRelationships{
				BundleID: &Relationship{
					Data: ResourceData{
						Type: ResourceTypeBundleIds,
						ID:   bundleID,
					},
				},
			},
		},
	}

	body, err := BuildRequestBody(request)
	if err != nil {
		return nil, err
	}

	data, err := c.do(ctx, "POST", "/v1/bundleIdCapabilities", body)
	if err != nil {
		return nil, err
	}

	var response BundleIDCapabilityResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// DeleteBundleIDCapability deletes a bundle ID capability by ID.
func (c *Client) DeleteBundleIDCapability(ctx context.Context, capabilityID string) error {
	capabilityID = strings.TrimSpace(capabilityID)
	path := fmt.Sprintf("/v1/bundleIdCapabilities/%s", capabilityID)
	_, err := c.do(ctx, "DELETE", path, nil)
	return err
}

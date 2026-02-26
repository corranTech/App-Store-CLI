package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const appDataUsagesInclude = "category,purpose,dataProtection"

// DataUsageTuple is the normalized tuple used to create/manage app data usages.
type DataUsageTuple struct {
	Category       string `json:"category,omitempty"`
	Purpose        string `json:"purpose,omitempty"`
	DataProtection string `json:"dataProtection"`
}

// AppDataUsage models one appDataUsages resource.
type AppDataUsage struct {
	ID             string `json:"id"`
	Category       string `json:"category,omitempty"`
	Purpose        string `json:"purpose,omitempty"`
	DataProtection string `json:"dataProtection,omitempty"`
}

// AppDataUsagesPublishState captures publication state for app privacy data usages.
type AppDataUsagesPublishState struct {
	ID        string `json:"id"`
	Published bool   `json:"published"`
}

func decodeAppDataUsageResource(resource jsonAPIResource) AppDataUsage {
	usage := AppDataUsage{
		ID: strings.TrimSpace(resource.ID),
	}
	if ref := firstRelationshipRef(resource, "category"); ref != nil {
		usage.Category = strings.TrimSpace(ref.ID)
	}
	if ref := firstRelationshipRef(resource, "purpose"); ref != nil {
		usage.Purpose = strings.TrimSpace(ref.ID)
	}
	if ref := firstRelationshipRef(resource, "dataProtection"); ref != nil {
		usage.DataProtection = strings.TrimSpace(ref.ID)
	}
	if usage.DataProtection == "" {
		usage.DataProtection = stringAttr(
			resource.Attributes,
			"appDataUsageDataProtection",
			"appDataUsageDataProtectionId",
		)
	}
	return usage
}

func decodeAppDataUsages(resources []jsonAPIResource) []AppDataUsage {
	if len(resources) == 0 {
		return []AppDataUsage{}
	}
	result := make([]AppDataUsage, 0, len(resources))
	for _, resource := range resources {
		result = append(result, decodeAppDataUsageResource(resource))
	}
	return result
}

func decodeAppDataUsagesPublishStateResource(resource jsonAPIResource) AppDataUsagesPublishState {
	return AppDataUsagesPublishState{
		ID:        strings.TrimSpace(resource.ID),
		Published: boolAttr(resource.Attributes, "published"),
	}
}

func normalizeDataUsageTuple(tuple DataUsageTuple) (DataUsageTuple, error) {
	tuple.Category = strings.TrimSpace(tuple.Category)
	tuple.Purpose = strings.TrimSpace(tuple.Purpose)
	tuple.DataProtection = strings.TrimSpace(tuple.DataProtection)
	if tuple.DataProtection == "" {
		return DataUsageTuple{}, fmt.Errorf("data protection is required")
	}
	return tuple, nil
}

// ListAppDataUsages lists data usage tuples for a specific app.
func (c *Client) ListAppDataUsages(ctx context.Context, appID string) ([]AppDataUsage, error) {
	appID = strings.TrimSpace(appID)
	if appID == "" {
		return nil, fmt.Errorf("app id is required")
	}
	query := url.Values{}
	query.Set("include", appDataUsagesInclude)
	query.Set("limit", "2000")
	path := queryPath("/apps/"+url.PathEscape(appID)+"/dataUsages", query)

	responseBody, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var payload jsonAPIListPayload
	if err := json.Unmarshal(responseBody, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse app data usages response: %w", err)
	}
	return decodeAppDataUsages(payload.Data), nil
}

// CreateAppDataUsage creates one data usage tuple for an app.
func (c *Client) CreateAppDataUsage(ctx context.Context, appID string, tuple DataUsageTuple) (*AppDataUsage, error) {
	appID = strings.TrimSpace(appID)
	if appID == "" {
		return nil, fmt.Errorf("app id is required")
	}
	normalized, err := normalizeDataUsageTuple(tuple)
	if err != nil {
		return nil, err
	}

	relationships := map[string]any{
		"app": map[string]any{
			"data": map[string]string{
				"type": "apps",
				"id":   appID,
			},
		},
		"dataProtection": map[string]any{
			"data": map[string]string{
				"type": "appDataUsageDataProtections",
				"id":   normalized.DataProtection,
			},
		},
	}
	if normalized.Category != "" {
		relationships["category"] = map[string]any{
			"data": map[string]string{
				"type": "appDataUsageCategories",
				"id":   normalized.Category,
			},
		}
	}
	if normalized.Purpose != "" {
		relationships["purpose"] = map[string]any{
			"data": map[string]string{
				"type": "appDataUsagePurposes",
				"id":   normalized.Purpose,
			},
		}
	}

	requestBody := map[string]any{
		"data": map[string]any{
			"type":          "appDataUsages",
			"relationships": relationships,
		},
	}
	responseBody, err := c.doRequest(ctx, http.MethodPost, "/appDataUsages", requestBody)
	if err != nil {
		return nil, err
	}
	var payload struct {
		Data jsonAPIResource `json:"data"`
	}
	if err := json.Unmarshal(responseBody, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse create app data usage response: %w", err)
	}
	usage := decodeAppDataUsageResource(payload.Data)
	return &usage, nil
}

// DeleteAppDataUsage deletes one appDataUsages resource.
func (c *Client) DeleteAppDataUsage(ctx context.Context, appDataUsageID string) error {
	appDataUsageID = strings.TrimSpace(appDataUsageID)
	if appDataUsageID == "" {
		return fmt.Errorf("app data usage id is required")
	}
	_, err := c.doRequest(ctx, http.MethodDelete, "/appDataUsages/"+url.PathEscape(appDataUsageID), nil)
	return err
}

// GetAppDataUsagesPublishState fetches publication state for app data usages.
func (c *Client) GetAppDataUsagesPublishState(ctx context.Context, appID string) (*AppDataUsagesPublishState, error) {
	appID = strings.TrimSpace(appID)
	if appID == "" {
		return nil, fmt.Errorf("app id is required")
	}
	path := "/apps/" + url.PathEscape(appID) + "/dataUsagePublishState"
	responseBody, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var payload struct {
		Data jsonAPIResource `json:"data"`
	}
	if err := json.Unmarshal(responseBody, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse data usage publish state response: %w", err)
	}
	state := decodeAppDataUsagesPublishStateResource(payload.Data)
	return &state, nil
}

// SetAppDataUsagesPublished updates publication state for app data usages.
func (c *Client) SetAppDataUsagesPublished(ctx context.Context, publishStateID string, published bool) (*AppDataUsagesPublishState, error) {
	publishStateID = strings.TrimSpace(publishStateID)
	if publishStateID == "" {
		return nil, fmt.Errorf("publish state id is required")
	}

	requestBody := map[string]any{
		"data": map[string]any{
			"type": "appDataUsagesPublishState",
			"id":   publishStateID,
			"attributes": map[string]bool{
				"published": published,
			},
		},
	}
	path := "/appDataUsagesPublishState/" + url.PathEscape(publishStateID)
	responseBody, err := c.doRequest(ctx, http.MethodPatch, path, requestBody)
	if err != nil {
		return nil, err
	}
	var payload struct {
		Data jsonAPIResource `json:"data"`
	}
	if err := json.Unmarshal(responseBody, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse publish state update response: %w", err)
	}
	state := decodeAppDataUsagesPublishStateResource(payload.Data)
	return &state, nil
}

// PublishAppDataUsages sets app data usage publish state to published=true.
func (c *Client) PublishAppDataUsages(ctx context.Context, appID string) (*AppDataUsagesPublishState, error) {
	state, err := c.GetAppDataUsagesPublishState(ctx, appID)
	if err != nil {
		return nil, err
	}
	if state.Published {
		return state, nil
	}
	return c.SetAppDataUsagesPublished(ctx, state.ID, true)
}

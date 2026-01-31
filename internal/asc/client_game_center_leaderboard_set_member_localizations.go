package asc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetGameCenterLeaderboardSetMemberLocalizations retrieves leaderboard set member localizations.
func (c *Client) GetGameCenterLeaderboardSetMemberLocalizations(ctx context.Context, opts ...GCLeaderboardSetMemberLocalizationsOption) (*GameCenterLeaderboardSetMemberLocalizationsResponse, error) {
	query := &gcLeaderboardSetMemberLocalizationsQuery{}
	for _, opt := range opts {
		opt(query)
	}

	path := "/v1/gameCenterLeaderboardSetMemberLocalizations"
	if query.nextURL != "" {
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("game-center-leaderboard-set-member-localizations: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildGCLeaderboardSetMemberLocalizationsQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var response GameCenterLeaderboardSetMemberLocalizationsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetGameCenterLeaderboardSetMemberLocalization retrieves a leaderboard set member localization by ID.
func (c *Client) GetGameCenterLeaderboardSetMemberLocalization(ctx context.Context, localizationID string) (*GameCenterLeaderboardSetMemberLocalizationResponse, error) {
	path := fmt.Sprintf("/v1/gameCenterLeaderboardSetMemberLocalizations/%s", strings.TrimSpace(localizationID))
	data, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var response GameCenterLeaderboardSetMemberLocalizationResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetGameCenterLeaderboardSetMemberLocalizationLeaderboard retrieves the leaderboard for a member localization.
func (c *Client) GetGameCenterLeaderboardSetMemberLocalizationLeaderboard(ctx context.Context, localizationID string) (*GameCenterLeaderboardResponse, error) {
	path := fmt.Sprintf("/v1/gameCenterLeaderboardSetMemberLocalizations/%s/gameCenterLeaderboard", strings.TrimSpace(localizationID))
	data, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var response GameCenterLeaderboardResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetGameCenterLeaderboardSetMemberLocalizationLeaderboardSet retrieves the leaderboard set for a member localization.
func (c *Client) GetGameCenterLeaderboardSetMemberLocalizationLeaderboardSet(ctx context.Context, localizationID string) (*GameCenterLeaderboardSetResponse, error) {
	path := fmt.Sprintf("/v1/gameCenterLeaderboardSetMemberLocalizations/%s/gameCenterLeaderboardSet", strings.TrimSpace(localizationID))
	data, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var response GameCenterLeaderboardSetResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

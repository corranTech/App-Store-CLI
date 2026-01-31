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

// GetGameCenterDetailsAchievementReleases retrieves achievement releases for a Game Center detail.
func (c *Client) GetGameCenterDetailsAchievementReleases(ctx context.Context, gcDetailID string, opts ...GCAchievementReleasesOption) (*GameCenterAchievementReleasesResponse, error) {
	query := &gcAchievementReleasesQuery{}
	for _, opt := range opts {
		opt(query)
	}

	path := fmt.Sprintf("/v1/gameCenterDetails/%s/achievementReleases", strings.TrimSpace(gcDetailID))
	if query.nextURL != "" {
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("game-center-details-achievement-releases: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildGCAchievementReleasesQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var response GameCenterAchievementReleasesResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetGameCenterDetailsLeaderboardReleases retrieves leaderboard releases for a Game Center detail.
func (c *Client) GetGameCenterDetailsLeaderboardReleases(ctx context.Context, gcDetailID string, opts ...GCLeaderboardReleasesOption) (*GameCenterLeaderboardReleasesResponse, error) {
	query := &gcLeaderboardReleasesQuery{}
	for _, opt := range opts {
		opt(query)
	}

	path := fmt.Sprintf("/v1/gameCenterDetails/%s/leaderboardReleases", strings.TrimSpace(gcDetailID))
	if query.nextURL != "" {
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("game-center-details-leaderboard-releases: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildGCLeaderboardReleasesQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var response GameCenterLeaderboardReleasesResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetGameCenterDetailsLeaderboardSetReleases retrieves leaderboard set releases for a Game Center detail.
func (c *Client) GetGameCenterDetailsLeaderboardSetReleases(ctx context.Context, gcDetailID string, opts ...GCLeaderboardSetReleasesOption) (*GameCenterLeaderboardSetReleasesResponse, error) {
	query := &gcLeaderboardSetReleasesQuery{}
	for _, opt := range opts {
		opt(query)
	}

	path := fmt.Sprintf("/v1/gameCenterDetails/%s/leaderboardSetReleases", strings.TrimSpace(gcDetailID))
	if query.nextURL != "" {
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("game-center-details-leaderboard-set-releases: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildGCLeaderboardSetReleasesQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var response GameCenterLeaderboardSetReleasesResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetGameCenterDetailsAchievementsV2 retrieves v2 achievements for a Game Center detail.
func (c *Client) GetGameCenterDetailsAchievementsV2(ctx context.Context, gcDetailID string, opts ...GCAchievementsOption) (*GameCenterAchievementsResponse, error) {
	query := &gcAchievementsQuery{}
	for _, opt := range opts {
		opt(query)
	}

	path := fmt.Sprintf("/v1/gameCenterDetails/%s/gameCenterAchievementsV2", strings.TrimSpace(gcDetailID))
	if query.nextURL != "" {
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("game-center-details-achievements-v2: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildGCAchievementsQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var response GameCenterAchievementsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetGameCenterDetailsLeaderboardsV2 retrieves v2 leaderboards for a Game Center detail.
func (c *Client) GetGameCenterDetailsLeaderboardsV2(ctx context.Context, gcDetailID string, opts ...GCLeaderboardsOption) (*GameCenterLeaderboardsResponse, error) {
	query := &gcLeaderboardsQuery{}
	for _, opt := range opts {
		opt(query)
	}

	path := fmt.Sprintf("/v1/gameCenterDetails/%s/gameCenterLeaderboardsV2", strings.TrimSpace(gcDetailID))
	if query.nextURL != "" {
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("game-center-details-leaderboards-v2: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildGCLeaderboardsQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var response GameCenterLeaderboardsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetGameCenterDetailsLeaderboardSetsV2 retrieves v2 leaderboard sets for a Game Center detail.
func (c *Client) GetGameCenterDetailsLeaderboardSetsV2(ctx context.Context, gcDetailID string, opts ...GCLeaderboardSetsOption) (*GameCenterLeaderboardSetsResponse, error) {
	query := &gcLeaderboardSetsQuery{}
	for _, opt := range opts {
		opt(query)
	}

	path := fmt.Sprintf("/v1/gameCenterDetails/%s/gameCenterLeaderboardSetsV2", strings.TrimSpace(gcDetailID))
	if query.nextURL != "" {
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("game-center-details-leaderboard-sets-v2: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildGCLeaderboardSetsQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var response GameCenterLeaderboardSetsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetGameCenterDetailsClassicMatchmakingRequests retrieves classic matchmaking request metrics.
func (c *Client) GetGameCenterDetailsClassicMatchmakingRequests(ctx context.Context, gcDetailID string, opts ...GCMatchmakingMetricsOption) (*GameCenterMetricsResponse, error) {
	query := &gcMatchmakingMetricsQuery{}
	for _, opt := range opts {
		opt(query)
	}

	path := fmt.Sprintf("/v1/gameCenterDetails/%s/metrics/classicMatchmakingRequests", strings.TrimSpace(gcDetailID))
	if query.nextURL != "" {
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("game-center-details-classic-matchmaking-requests: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildGCMatchmakingQueueRequestsQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var response GameCenterMetricsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetGameCenterDetailsRuleBasedMatchmakingRequests retrieves rule-based matchmaking request metrics.
func (c *Client) GetGameCenterDetailsRuleBasedMatchmakingRequests(ctx context.Context, gcDetailID string, opts ...GCMatchmakingMetricsOption) (*GameCenterMetricsResponse, error) {
	query := &gcMatchmakingMetricsQuery{}
	for _, opt := range opts {
		opt(query)
	}

	path := fmt.Sprintf("/v1/gameCenterDetails/%s/metrics/ruleBasedMatchmakingRequests", strings.TrimSpace(gcDetailID))
	if query.nextURL != "" {
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("game-center-details-rule-based-matchmaking-requests: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildGCMatchmakingQueueRequestsQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var response GameCenterMetricsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

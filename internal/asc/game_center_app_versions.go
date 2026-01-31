package asc

import (
	"net/url"
	"strconv"
	"strings"
)

// GameCenterAppVersionAttributes represents a Game Center app version resource.
type GameCenterAppVersionAttributes struct {
	Enabled bool `json:"enabled,omitempty"`
}

// GameCenterAppVersionsResponse is the response from app version list endpoints.
type GameCenterAppVersionsResponse = Response[GameCenterAppVersionAttributes]

// GameCenterAppVersionResponse is the response from app version detail endpoints.
type GameCenterAppVersionResponse = SingleResponse[GameCenterAppVersionAttributes]

// GCAppVersionsOption is a functional option for Game Center app version queries.
type GCAppVersionsOption func(*gcAppVersionsQuery)

type gcAppVersionsQuery struct {
	listQuery
	enabled []string
}

// WithGCAppVersionsLimit sets the max number of app versions to return.
func WithGCAppVersionsLimit(limit int) GCAppVersionsOption {
	return func(q *gcAppVersionsQuery) {
		if limit > 0 {
			q.limit = limit
		}
	}
}

// WithGCAppVersionsNextURL uses a next page URL directly.
func WithGCAppVersionsNextURL(next string) GCAppVersionsOption {
	return func(q *gcAppVersionsQuery) {
		if strings.TrimSpace(next) != "" {
			q.nextURL = strings.TrimSpace(next)
		}
	}
}

// WithGCAppVersionsEnabled filters by enabled status.
func WithGCAppVersionsEnabled(enabled bool) GCAppVersionsOption {
	return func(q *gcAppVersionsQuery) {
		q.enabled = []string{strconv.FormatBool(enabled)}
	}
}

func buildGCAppVersionsQuery(query *gcAppVersionsQuery) string {
	values := url.Values{}
	addCSV(values, "filter[enabled]", query.enabled)
	addLimit(values, query.limit)
	return values.Encode()
}

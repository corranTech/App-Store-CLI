package asc

import (
	"net/url"
	"strings"
)

// GameCenterDetailsResponse is the response from Game Center detail list endpoints.
type GameCenterDetailsResponse = Response[GameCenterDetailAttributes]

// GCDetailsOption is a functional option for Game Center detail queries.
type GCDetailsOption func(*gcDetailsQuery)

type gcDetailsQuery struct {
	listQuery
}

// WithGCDetailsLimit sets the max number of details to return.
func WithGCDetailsLimit(limit int) GCDetailsOption {
	return func(q *gcDetailsQuery) {
		if limit > 0 {
			q.limit = limit
		}
	}
}

// WithGCDetailsNextURL uses a next page URL directly.
func WithGCDetailsNextURL(next string) GCDetailsOption {
	return func(q *gcDetailsQuery) {
		if strings.TrimSpace(next) != "" {
			q.nextURL = strings.TrimSpace(next)
		}
	}
}

func buildGCDetailsQuery(query *gcDetailsQuery) string {
	values := url.Values{}
	addLimit(values, query.limit)
	return values.Encode()
}

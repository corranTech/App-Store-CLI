package asc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// ActorAttributes describes an actor resource.
type ActorAttributes struct {
	ActorType     string `json:"actorType,omitempty"`
	UserFirstName string `json:"userFirstName,omitempty"`
	UserLastName  string `json:"userLastName,omitempty"`
	UserEmail     string `json:"userEmail,omitempty"`
	APIKeyID      string `json:"apiKeyId,omitempty"`
}

// ActorsResponse is the response from actors endpoint.
type ActorsResponse = Response[ActorAttributes]

// ActorResponse is the response from actor detail endpoint.
type ActorResponse = SingleResponse[ActorAttributes]

// GetActors retrieves actors filtered by IDs.
func (c *Client) GetActors(ctx context.Context, opts ...ActorsOption) (*ActorsResponse, error) {
	query := &actorsQuery{}
	for _, opt := range opts {
		opt(query)
	}

	path := "/v1/actors"
	if query.nextURL != "" {
		// Validate nextURL to prevent credential exfiltration
		if err := validateNextURL(query.nextURL); err != nil {
			return nil, fmt.Errorf("actors: %w", err)
		}
		path = query.nextURL
	} else if queryString := buildActorsQuery(query); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response ActorsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetActor retrieves a single actor by ID.
func (c *Client) GetActor(ctx context.Context, actorID string, fields []string) (*ActorResponse, error) {
	actorID = strings.TrimSpace(actorID)
	path := fmt.Sprintf("/v1/actors/%s", actorID)
	if queryString := buildActorsFieldsQuery(fields); queryString != "" {
		path += "?" + queryString
	}

	data, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response ActorResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

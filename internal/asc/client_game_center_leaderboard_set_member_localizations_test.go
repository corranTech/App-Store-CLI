package asc

import (
	"context"
	"net/http"
	"net/url"
	"testing"
)

func TestGetGameCenterLeaderboardSetMemberLocalizations_WithFilters(t *testing.T) {
	response := jsonResponse(http.StatusOK, `{"data":[]}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", req.Method)
		}
		if req.URL.Path != "/v1/gameCenterLeaderboardSetMemberLocalizations" {
			t.Fatalf("expected path /v1/gameCenterLeaderboardSetMemberLocalizations, got %s", req.URL.Path)
		}
		values := req.URL.Query()
		if values.Get("limit") != "5" {
			t.Fatalf("expected limit=5, got %q", values.Get("limit"))
		}
		if values.Get("filter[gameCenterLeaderboardSet]") != "set-1" {
			t.Fatalf("expected filter[gameCenterLeaderboardSet]=set-1, got %q", values.Get("filter[gameCenterLeaderboardSet]"))
		}
		if values.Get("filter[gameCenterLeaderboard]") != "lb-1,lb-2" {
			t.Fatalf("expected filter[gameCenterLeaderboard]=lb-1,lb-2, got %q", values.Get("filter[gameCenterLeaderboard]"))
		}
		assertAuthorized(t, req)
	}, response)

	opts := []GCLeaderboardSetMemberLocalizationsOption{
		WithGCLeaderboardSetMemberLocalizationsLimit(5),
		WithGCLeaderboardSetMemberLocalizationsLeaderboardSetIDs([]string{"set-1"}),
		WithGCLeaderboardSetMemberLocalizationsLeaderboardIDs([]string{"lb-1", "lb-2"}),
	}

	if _, err := client.GetGameCenterLeaderboardSetMemberLocalizations(context.Background(), opts...); err != nil {
		t.Fatalf("GetGameCenterLeaderboardSetMemberLocalizations() error: %v", err)
	}
}

func TestGetGameCenterLeaderboardSetMemberLocalizations_UsesNextURL(t *testing.T) {
	next := "https://api.appstoreconnect.apple.com/v1/gameCenterLeaderboardSetMemberLocalizations?cursor=next"
	response := jsonResponse(http.StatusOK, `{"data":[]}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.URL.String() != next {
			t.Fatalf("expected URL %q, got %q", next, req.URL.String())
		}
		assertAuthorized(t, req)
	}, response)

	if _, err := client.GetGameCenterLeaderboardSetMemberLocalizations(context.Background(), WithGCLeaderboardSetMemberLocalizationsNextURL(next)); err != nil {
		t.Fatalf("GetGameCenterLeaderboardSetMemberLocalizations() error: %v", err)
	}
}

func TestGetGameCenterLeaderboardSetMemberLocalization(t *testing.T) {
	response := jsonResponse(http.StatusOK, `{"data":{"type":"gameCenterLeaderboardSetMemberLocalizations","id":"loc-1","attributes":{"name":"Seasonal","locale":"en-US"}}}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", req.Method)
		}
		if req.URL.Path != "/v1/gameCenterLeaderboardSetMemberLocalizations/loc-1" {
			t.Fatalf("expected path /v1/gameCenterLeaderboardSetMemberLocalizations/loc-1, got %s", req.URL.Path)
		}
		assertAuthorized(t, req)
	}, response)

	if _, err := client.GetGameCenterLeaderboardSetMemberLocalization(context.Background(), "loc-1"); err != nil {
		t.Fatalf("GetGameCenterLeaderboardSetMemberLocalization() error: %v", err)
	}
}

func TestGetGameCenterLeaderboardSetMemberLocalizationLeaderboard(t *testing.T) {
	response := jsonResponse(http.StatusOK, `{"data":{"type":"gameCenterLeaderboards","id":"lb-1","attributes":{"referenceName":"Leaderboard"}}}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", req.Method)
		}
		if req.URL.Path != "/v1/gameCenterLeaderboardSetMemberLocalizations/loc-1/gameCenterLeaderboard" {
			t.Fatalf("expected path /v1/gameCenterLeaderboardSetMemberLocalizations/loc-1/gameCenterLeaderboard, got %s", req.URL.Path)
		}
		assertAuthorized(t, req)
	}, response)

	if _, err := client.GetGameCenterLeaderboardSetMemberLocalizationLeaderboard(context.Background(), "loc-1"); err != nil {
		t.Fatalf("GetGameCenterLeaderboardSetMemberLocalizationLeaderboard() error: %v", err)
	}
}

func TestGetGameCenterLeaderboardSetMemberLocalizationLeaderboardSet(t *testing.T) {
	response := jsonResponse(http.StatusOK, `{"data":{"type":"gameCenterLeaderboardSets","id":"set-1","attributes":{"referenceName":"Set"}}}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", req.Method)
		}
		if req.URL.Path != "/v1/gameCenterLeaderboardSetMemberLocalizations/loc-1/gameCenterLeaderboardSet" {
			t.Fatalf("expected path /v1/gameCenterLeaderboardSetMemberLocalizations/loc-1/gameCenterLeaderboardSet, got %s", req.URL.Path)
		}
		assertAuthorized(t, req)
	}, response)

	if _, err := client.GetGameCenterLeaderboardSetMemberLocalizationLeaderboardSet(context.Background(), "loc-1"); err != nil {
		t.Fatalf("GetGameCenterLeaderboardSetMemberLocalizationLeaderboardSet() error: %v", err)
	}
}

func TestGCLeaderboardSetMemberLocalizationsOptions(t *testing.T) {
	query := &gcLeaderboardSetMemberLocalizationsQuery{}
	WithGCLeaderboardSetMemberLocalizationsLimit(10)(query)
	if query.limit != 10 {
		t.Fatalf("expected limit 10, got %d", query.limit)
	}
	WithGCLeaderboardSetMemberLocalizationsNextURL("next")(query)
	if query.nextURL != "next" {
		t.Fatalf("expected nextURL set, got %q", query.nextURL)
	}
	WithGCLeaderboardSetMemberLocalizationsLeaderboardSetIDs([]string{" set-1 ", "", "set-2"})(query)
	if len(query.leaderboardSetIDs) != 2 {
		t.Fatalf("expected 2 set IDs, got %d", len(query.leaderboardSetIDs))
	}
	WithGCLeaderboardSetMemberLocalizationsLeaderboardIDs([]string{" lb-1 ", "lb-2"})(query)
	if len(query.leaderboardIDs) != 2 {
		t.Fatalf("expected 2 leaderboard IDs, got %d", len(query.leaderboardIDs))
	}
	values, err := url.ParseQuery(buildGCLeaderboardSetMemberLocalizationsQuery(query))
	if err != nil {
		t.Fatalf("parse query: %v", err)
	}
	if values.Get("limit") != "10" {
		t.Fatalf("expected limit=10, got %q", values.Get("limit"))
	}
	if values.Get("filter[gameCenterLeaderboardSet]") != "set-1,set-2" {
		t.Fatalf("expected filter[gameCenterLeaderboardSet]=set-1,set-2, got %q", values.Get("filter[gameCenterLeaderboardSet]"))
	}
	if values.Get("filter[gameCenterLeaderboard]") != "lb-1,lb-2" {
		t.Fatalf("expected filter[gameCenterLeaderboard]=lb-1,lb-2, got %q", values.Get("filter[gameCenterLeaderboard]"))
	}
}

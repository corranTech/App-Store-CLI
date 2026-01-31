package asc

import (
	"context"
	"net/http"
	"net/url"
	"testing"
)

func TestGetGameCenterDetails_WithLimit(t *testing.T) {
	response := jsonResponse(http.StatusOK, `{"data":[]}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", req.Method)
		}
		if req.URL.Path != "/v1/gameCenterDetails" {
			t.Fatalf("expected path /v1/gameCenterDetails, got %s", req.URL.Path)
		}
		if req.URL.Query().Get("limit") != "25" {
			t.Fatalf("expected limit=25, got %q", req.URL.Query().Get("limit"))
		}
		assertAuthorized(t, req)
	}, response)

	if _, err := client.GetGameCenterDetails(context.Background(), WithGCDetailsLimit(25)); err != nil {
		t.Fatalf("GetGameCenterDetails() error: %v", err)
	}
}

func TestGetGameCenterDetails_UsesNextURL(t *testing.T) {
	next := "https://api.appstoreconnect.apple.com/v1/gameCenterDetails?cursor=next"
	response := jsonResponse(http.StatusOK, `{"data":[]}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.URL.String() != next {
			t.Fatalf("expected URL %q, got %q", next, req.URL.String())
		}
		assertAuthorized(t, req)
	}, response)

	if _, err := client.GetGameCenterDetails(context.Background(), WithGCDetailsNextURL(next)); err != nil {
		t.Fatalf("GetGameCenterDetails() error: %v", err)
	}
}

func TestGetGameCenterDetail(t *testing.T) {
	response := jsonResponse(http.StatusOK, `{"data":{"type":"gameCenterDetails","id":"detail-1","attributes":{"arcadeEnabled":true}}}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", req.Method)
		}
		if req.URL.Path != "/v1/gameCenterDetails/detail-1" {
			t.Fatalf("expected path /v1/gameCenterDetails/detail-1, got %s", req.URL.Path)
		}
		assertAuthorized(t, req)
	}, response)

	if _, err := client.GetGameCenterDetail(context.Background(), "detail-1"); err != nil {
		t.Fatalf("GetGameCenterDetail() error: %v", err)
	}
}

func TestGetGameCenterDetailGameCenterGroup(t *testing.T) {
	response := jsonResponse(http.StatusOK, `{"data":{"type":"gameCenterGroups","id":"group-1"}}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", req.Method)
		}
		if req.URL.Path != "/v1/gameCenterDetails/detail-1/gameCenterGroup" {
			t.Fatalf("expected path /v1/gameCenterDetails/detail-1/gameCenterGroup, got %s", req.URL.Path)
		}
		assertAuthorized(t, req)
	}, response)

	if _, err := client.GetGameCenterDetailGameCenterGroup(context.Background(), "detail-1"); err != nil {
		t.Fatalf("GetGameCenterDetailGameCenterGroup() error: %v", err)
	}
}

func TestGetGameCenterGroupGameCenterDetails_WithLimit(t *testing.T) {
	response := jsonResponse(http.StatusOK, `{"data":[]}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", req.Method)
		}
		if req.URL.Path != "/v1/gameCenterGroups/group-1/gameCenterDetails" {
			t.Fatalf("expected path /v1/gameCenterGroups/group-1/gameCenterDetails, got %s", req.URL.Path)
		}
		if req.URL.Query().Get("limit") != "30" {
			t.Fatalf("expected limit=30, got %q", req.URL.Query().Get("limit"))
		}
		assertAuthorized(t, req)
	}, response)

	if _, err := client.GetGameCenterGroupGameCenterDetails(context.Background(), "group-1", WithGCDetailsLimit(30)); err != nil {
		t.Fatalf("GetGameCenterGroupGameCenterDetails() error: %v", err)
	}
}

func TestGetGameCenterGroupGameCenterDetails_UsesNextURL(t *testing.T) {
	next := "https://api.appstoreconnect.apple.com/v1/gameCenterGroups/group-1/gameCenterDetails?cursor=next"
	response := jsonResponse(http.StatusOK, `{"data":[]}`)
	client := newTestClient(t, func(req *http.Request) {
		if req.URL.String() != next {
			t.Fatalf("expected URL %q, got %q", next, req.URL.String())
		}
		assertAuthorized(t, req)
	}, response)

	if _, err := client.GetGameCenterGroupGameCenterDetails(context.Background(), "", WithGCDetailsNextURL(next)); err != nil {
		t.Fatalf("GetGameCenterGroupGameCenterDetails() error: %v", err)
	}
}

func TestGCDetailsOptions(t *testing.T) {
	query := &gcDetailsQuery{}
	WithGCDetailsLimit(8)(query)
	if query.limit != 8 {
		t.Fatalf("expected limit 8, got %d", query.limit)
	}
	WithGCDetailsNextURL("next")(query)
	if query.nextURL != "next" {
		t.Fatalf("expected nextURL set, got %q", query.nextURL)
	}
	values, err := url.ParseQuery(buildGCDetailsQuery(query))
	if err != nil {
		t.Fatalf("parse query: %v", err)
	}
	if values.Get("limit") != "8" {
		t.Fatalf("expected limit=8, got %q", values.Get("limit"))
	}
}

package cmdtest

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
	appeventscli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/app_events"
)

func newAppEventsTestClient(t *testing.T, transport roundTripFunc) *asc.Client {
	t.Helper()

	tmpDir := t.TempDir()
	keyPath := filepath.Join(tmpDir, "key.p8")
	writeECDSAPEM(t, keyPath)

	httpClient := &http.Client{Transport: transport}
	client, err := asc.NewClientWithHTTPClient("KEY123", "ISS456", keyPath, httpClient)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return client
}

func TestAppEventsCreateIgnoresScheduleFlagsAndWarns(t *testing.T) {
	var captured map[string]any

	client := newAppEventsTestClient(t, roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", req.Method)
		}
		if req.URL.Path != "/v1/appEvents" {
			t.Fatalf("expected path /v1/appEvents, got %s", req.URL.Path)
		}

		if err := json.NewDecoder(req.Body).Decode(&captured); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		return jsonResponse(http.StatusCreated, `{"data":{"type":"appEvents","id":"event-1","attributes":{"referenceName":"Launch","badge":"CHALLENGE"}}}`)
	}))

	restore := appeventscli.SetClientFactory(func() (*asc.Client, error) {
		return client, nil
	})
	defer restore()

	root := RootCommand("1.2.3")
	stdout, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{
			"app-events", "create",
			"--app", "app-123",
			"--name", "Launch",
			"--event-type", "CHALLENGE",
			"--start", "2026-06-01T00:00:00Z",
			"--end", "2026-06-30T23:59:59Z",
			"--publish-start", "2026-05-15T00:00:00Z",
			"--territories", "usa, can",
			"--primary-locale", "en-US",
		}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		if err := root.Run(context.Background()); err != nil {
			t.Fatalf("run error: %v", err)
		}
	})

	if !strings.Contains(stderr, "App Store Connect currently returns HTTP 500 when app event territorySchedules are included on create") {
		t.Fatalf("expected schedule warning, got %q", stderr)
	}
	if !strings.Contains(stderr, "creating the event without a schedule") {
		t.Fatalf("expected unscheduled explanation, got %q", stderr)
	}

	var resp asc.AppEventResponse
	if err := json.Unmarshal([]byte(stdout), &resp); err != nil {
		t.Fatalf("failed to parse JSON output: %v", err)
	}
	if resp.Data.ID != "event-1" {
		t.Fatalf("expected event id event-1, got %q", resp.Data.ID)
	}

	data, ok := captured["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object, got %#v", captured["data"])
	}
	attrs, ok := data["attributes"].(map[string]any)
	if !ok {
		t.Fatalf("expected attributes object, got %#v", data["attributes"])
	}
	if _, ok := attrs["territorySchedules"]; ok {
		t.Fatalf("expected territorySchedules to be omitted, got %#v", attrs["territorySchedules"])
	}
	if attrs["referenceName"] != "Launch" {
		t.Fatalf("expected referenceName Launch, got %#v", attrs["referenceName"])
	}
	if attrs["badge"] != "CHALLENGE" {
		t.Fatalf("expected badge CHALLENGE, got %#v", attrs["badge"])
	}
	if attrs["primaryLocale"] != "en-US" {
		t.Fatalf("expected primaryLocale en-US, got %#v", attrs["primaryLocale"])
	}

	relationships, ok := data["relationships"].(map[string]any)
	if !ok {
		t.Fatalf("expected relationships object, got %#v", data["relationships"])
	}
	appRel, ok := relationships["app"].(map[string]any)
	if !ok {
		t.Fatalf("expected app relationship, got %#v", relationships["app"])
	}
	appData, ok := appRel["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected app relationship data, got %#v", appRel["data"])
	}
	if appData["id"] != "app-123" {
		t.Fatalf("expected app id app-123, got %#v", appData["id"])
	}
}

func TestAppEventsCreateScheduleFlagsStillValidateRFC3339(t *testing.T) {
	restore := appeventscli.SetClientFactory(func() (*asc.Client, error) {
		t.Fatal("did not expect client creation for invalid schedule flags")
		return nil, nil
	})
	defer restore()

	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	stdout, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{
			"app-events", "create",
			"--app", "app-123",
			"--name", "Launch",
			"--event-type", "CHALLENGE",
			"--start", "not-a-date",
			"--end", "2026-06-30T23:59:59Z",
		}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		err := root.Run(context.Background())
		if !errors.Is(err, flag.ErrHelp) {
			t.Fatalf("expected ErrHelp, got %v", err)
		}
	})

	if stdout != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
	if !strings.Contains(stderr, "--start must be in RFC3339 format") {
		t.Fatalf("expected RFC3339 validation error, got %q", stderr)
	}
}

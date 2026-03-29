package cmdtest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
)

func TestSubscriptionsIntroductoryOffersImport_CreateSuccessSummary(t *testing.T) {
	setupAuth(t)
	t.Setenv("ASC_CONFIG_PATH", filepath.Join(t.TempDir(), "nonexistent.json"))

	originalTransport := http.DefaultTransport
	t.Cleanup(func() {
		http.DefaultTransport = originalTransport
	})

	requestCount := 0
	http.DefaultTransport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		requestCount++
		if req.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", req.Method)
		}
		if req.URL.Path != "/v1/subscriptionIntroductoryOffers" {
			t.Fatalf("unexpected path: %s", req.URL.Path)
		}

		var payload map[string]any
		if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		data := payload["data"].(map[string]any)
		attrs := data["attributes"].(map[string]any)
		relationships := data["relationships"].(map[string]any)
		territory := relationships["territory"].(map[string]any)["data"].(map[string]any)["id"]

		if attrs["duration"] != "ONE_WEEK" {
			t.Fatalf("expected ONE_WEEK duration, got %#v", attrs["duration"])
		}
		if attrs["offerMode"] != "FREE_TRIAL" {
			t.Fatalf("expected FREE_TRIAL offerMode, got %#v", attrs["offerMode"])
		}
		if attrs["numberOfPeriods"] != float64(1) {
			t.Fatalf("expected numberOfPeriods 1, got %#v", attrs["numberOfPeriods"])
		}

		switch requestCount {
		case 1:
			if territory != "USA" {
				t.Fatalf("expected USA territory, got %#v", territory)
			}
		case 2:
			if territory != "AFG" {
				t.Fatalf("expected AFG territory, got %#v", territory)
			}
		default:
			t.Fatalf("unexpected request count %d", requestCount)
		}

		body := `{"data":{"type":"subscriptionIntroductoryOffers","id":"offer-1"}}`
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     http.Header{"Content-Type": []string{"application/json"}},
		}, nil
	})

	csvPath := writeTempIntroOffersCSV(t, "territory\nUSA\nAfghanistan\n")

	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	type importSummary struct {
		DryRun  bool `json:"dryRun"`
		Total   int  `json:"total"`
		Created int  `json:"created"`
		Failed  int  `json:"failed"`
	}

	stdout, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{
			"subscriptions", "offers", "introductory", "import",
			"--subscription-id", "SUB_ID",
			"--input", csvPath,
			"--offer-duration", "ONE_WEEK",
			"--offer-mode", "FREE_TRIAL",
			"--number-of-periods", "1",
		}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		if err := root.Run(context.Background()); err != nil {
			t.Fatalf("run error: %v", err)
		}
	})

	if stderr != "" {
		t.Fatalf("expected empty stderr, got %q", stderr)
	}

	var summary importSummary
	if err := json.Unmarshal([]byte(stdout), &summary); err != nil {
		t.Fatalf("parse JSON summary: %v", err)
	}
	if summary.DryRun {
		t.Fatalf("expected dryRun=false")
	}
	if summary.Total != 2 || summary.Created != 2 || summary.Failed != 0 {
		t.Fatalf("unexpected summary: %+v", summary)
	}
	if requestCount != 2 {
		t.Fatalf("expected 2 requests, got %d", requestCount)
	}
}

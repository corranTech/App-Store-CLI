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

func TestLegacyOfferCodesValuesOutputsCodesAndWarning(t *testing.T) {
	setupAuth(t)
	t.Setenv("ASC_CONFIG_PATH", filepath.Join(t.TempDir(), "nonexistent.json"))

	originalTransport := http.DefaultTransport
	t.Cleanup(func() {
		http.DefaultTransport = originalTransport
	})

	http.DefaultTransport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", req.Method)
		}
		if req.URL.Path != "/v1/subscriptionOfferCodeOneTimeUseCodes/batch-1/values" {
			t.Fatalf("expected values path, got %s", req.URL.Path)
		}

		body := "code\nSPRING-1\nSPRING-2\n"
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     http.Header{"Content-Type": []string{"text/csv"}},
		}, nil
	})

	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	stdout, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{"offer-codes", "values", "--id", "batch-1"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		if err := root.Run(context.Background()); err != nil {
			t.Fatalf("run error: %v", err)
		}
	})

	if stdout != "SPRING-1\nSPRING-2\n" {
		t.Fatalf("expected one code per line, got %q", stdout)
	}
	if !strings.Contains(stderr, `Warning: "asc offer-codes values" is deprecated. Use "asc subscriptions offer-codes values" instead.`) {
		t.Fatalf("expected deprecation warning, got %q", stderr)
	}
	assertOnlyDeprecatedCommandWarnings(t, stderr)
}

func TestLegacyWinBackOffersListOutputsJSONAndWarning(t *testing.T) {
	setupAuth(t)
	t.Setenv("ASC_CONFIG_PATH", filepath.Join(t.TempDir(), "nonexistent.json"))

	originalTransport := http.DefaultTransport
	t.Cleanup(func() {
		http.DefaultTransport = originalTransport
	})

	http.DefaultTransport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", req.Method)
		}
		if req.URL.Path != "/v1/subscriptions/sub-1/winBackOffers" {
			t.Fatalf("expected win-back offers path, got %s", req.URL.Path)
		}

		body := `{"data":[{"type":"winBackOffers","id":"offer-1","attributes":{"referenceName":"Spring 2026"}}],"links":{"next":""}}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     http.Header{"Content-Type": []string{"application/json"}},
		}, nil
	})

	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	stdout, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{"win-back-offers", "list", "--subscription", "sub-1"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		if err := root.Run(context.Background()); err != nil {
			t.Fatalf("run error: %v", err)
		}
	})

	var out struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(stdout), &out); err != nil {
		t.Fatalf("unmarshal output: %v\nstdout: %s", err, stdout)
	}
	if len(out.Data) != 1 || out.Data[0].ID != "offer-1" {
		t.Fatalf("expected legacy win-back offer output, got %+v", out.Data)
	}
	if !strings.Contains(stderr, `Warning: "asc win-back-offers list" is deprecated. Use "asc subscriptions win-back-offers list" instead.`) {
		t.Fatalf("expected deprecation warning, got %q", stderr)
	}
	assertOnlyDeprecatedCommandWarnings(t, stderr)
}

func TestLegacyPromotedPurchasesListOutputsJSONAndWarning(t *testing.T) {
	setupAuth(t)
	t.Setenv("ASC_CONFIG_PATH", filepath.Join(t.TempDir(), "nonexistent.json"))

	originalTransport := http.DefaultTransport
	t.Cleanup(func() {
		http.DefaultTransport = originalTransport
	})

	http.DefaultTransport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", req.Method)
		}
		if req.URL.Path != "/v1/apps/app-1/promotedPurchases" {
			t.Fatalf("expected promoted purchases path, got %s", req.URL.Path)
		}

		body := `{"data":[{"type":"promotedPurchases","id":"promo-1","attributes":{"enabled":true}}],"links":{"next":""}}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     http.Header{"Content-Type": []string{"application/json"}},
		}, nil
	})

	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	stdout, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{"promoted-purchases", "list", "--app", "app-1"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		if err := root.Run(context.Background()); err != nil {
			t.Fatalf("run error: %v", err)
		}
	})

	var out struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(stdout), &out); err != nil {
		t.Fatalf("unmarshal output: %v\nstdout: %s", err, stdout)
	}
	if len(out.Data) != 1 || out.Data[0].ID != "promo-1" {
		t.Fatalf("expected legacy promoted purchases output, got %+v", out.Data)
	}
	if !strings.Contains(stderr, `Warning: "asc promoted-purchases" is deprecated. Use "asc subscriptions promoted-purchases ..." for subscriptions or "asc iap promoted-purchases ..." for in-app purchases.`) {
		t.Fatalf("expected deprecation warning, got %q", stderr)
	}
	assertOnlyDeprecatedCommandWarnings(t, stderr)
}

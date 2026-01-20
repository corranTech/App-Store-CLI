package asc

import (
	"bytes"
	"net/url"
	"strings"
	"testing"
)

func TestBuildReviewQuery(t *testing.T) {
	query := buildReviewQuery([]ReviewOption{
		WithRating(5),
		WithTerritory("us"),
	})

	values, err := url.ParseQuery(query)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	if got := values.Get("filter[rating]"); got != "5" {
		t.Fatalf("expected filter[rating]=5, got %q", got)
	}

	if got := values.Get("filter[territory]"); got != "US" {
		t.Fatalf("expected filter[territory]=US, got %q", got)
	}
}

func TestBuildReviewQuery_InvalidRating(t *testing.T) {
	query := buildReviewQuery([]ReviewOption{
		WithRating(9),
	})

	values, err := url.ParseQuery(query)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	if got := values.Get("filter[rating]"); got != "" {
		t.Fatalf("expected empty filter[rating], got %q", got)
	}
}

func TestBuildRequestBody(t *testing.T) {
	body, err := BuildRequestBody(map[string]string{"hello": "world"})
	if err != nil {
		t.Fatalf("BuildRequestBody() error: %v", err)
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(body); err != nil {
		t.Fatalf("read body error: %v", err)
	}

	if !strings.Contains(buf.String(), `"hello":"world"`) {
		t.Fatalf("unexpected body: %s", buf.String())
	}
}

func TestParseError(t *testing.T) {
	payload := []byte(`{"errors":[{"code":"FORBIDDEN","title":"Forbidden","detail":"not allowed"}]}`)
	err := ParseError(payload)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "Forbidden") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

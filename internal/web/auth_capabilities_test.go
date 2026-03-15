package web

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestClientDoIrisV1RequestUsesIntegrationsHeaders(t *testing.T) {
	var (
		gotPath   string
		gotAccept string
		gotCSRF   string
		gotOrigin string
		gotReferer string
	)
	client := &Client{
		httpClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		gotPath = r.URL.Path + "?" + r.URL.RawQuery
		gotAccept = r.Header.Get("Accept")
		gotCSRF = r.Header.Get("X-CSRF-ITC")
		gotOrigin = r.Header.Get("Origin")
		gotReferer = r.Header.Get("Referer")
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"data":[]}`)),
		}, nil
	})},
	}

	if _, err := client.doIrisV1Request(context.Background(), http.MethodGet, "/apiKeys?limit=2000", nil); err != nil {
		t.Fatalf("doIrisV1Request() error: %v", err)
	}

	if gotPath != "/iris/v1/apiKeys?limit=2000" {
		t.Fatalf("expected path %q, got %q", "/iris/v1/apiKeys?limit=2000", gotPath)
	}
	if gotAccept != "application/vnd.api+json, application/json, text/csv" {
		t.Fatalf("unexpected accept header %q", gotAccept)
	}
	if gotCSRF != "[asc-ui]" {
		t.Fatalf("unexpected csrf header %q", gotCSRF)
	}
	if gotOrigin != appStoreBaseURL {
		t.Fatalf("unexpected origin header %q", gotOrigin)
	}
	if gotReferer != integrationsAPIRefererURL {
		t.Fatalf("unexpected referer header %q", gotReferer)
	}
}

func TestClientDoOlympusRequestUsesOlympusHeaders(t *testing.T) {
	var (
		gotPath      string
		gotRequested string
		gotAccept    string
	)
	client := &Client{
		httpClient: &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		gotPath = r.URL.Path
		gotRequested = r.Header.Get("X-Requested-With")
		gotAccept = r.Header.Get("Accept")
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"data":[]}`)),
		}, nil
	})},
	}

	if _, err := client.doOlympusRequest(context.Background(), http.MethodGet, "/actors/actor-1", nil); err != nil {
		t.Fatalf("doOlympusRequest() error: %v", err)
	}

	if gotPath != "/olympus/v1/actors/actor-1" {
		t.Fatalf("expected path %q, got %q", "/olympus/v1/actors/actor-1", gotPath)
	}
	if gotRequested != "xsdr2$" {
		t.Fatalf("unexpected X-Requested-With %q", gotRequested)
	}
	if gotAccept != "application/json" {
		t.Fatalf("unexpected accept header %q", gotAccept)
	}
}

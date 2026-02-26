package web

import (
	"context"
	"errors"
	"testing"

	webcore "github.com/rudrankriyam/App-Store-Connect-CLI/internal/web"
)

func TestWebAppsCreateDefersPasswordResolutionToResolveSession(t *testing.T) {
	origResolveSession := resolveSessionFn
	origPromptPassword := promptPasswordFn
	t.Cleanup(func() {
		resolveSessionFn = origResolveSession
		promptPasswordFn = origPromptPassword
	})

	promptErr := errors.New("prompt should not run before session resolution")
	resolveErr := errors.New("stop before network call")

	promptPasswordFn = func() (string, error) {
		return "", promptErr
	}

	var (
		calledResolve bool
		receivedID    string
		receivedPass  string
	)
	resolveSessionFn = func(ctx context.Context, appleID, password, twoFactorCode string, usePasswordStdin bool) (*webcore.AuthSession, string, error) {
		calledResolve = true
		receivedID = appleID
		receivedPass = password
		return nil, "", resolveErr
	}

	cmd := WebAppsCreateCommand()
	if err := cmd.FlagSet.Parse([]string{
		"--name", "My App",
		"--bundle-id", "com.example.app",
		"--sku", "SKU123",
		"--apple-id", "user@example.com",
	}); err != nil {
		t.Fatalf("parse error: %v", err)
	}

	err := cmd.Exec(context.Background(), nil)
	if !errors.Is(err, resolveErr) {
		t.Fatalf("expected resolveSession error %v, got %v", resolveErr, err)
	}
	if !calledResolve {
		t.Fatal("expected resolveSession to be called")
	}
	if receivedID != "user@example.com" {
		t.Fatalf("expected apple ID %q, got %q", "user@example.com", receivedID)
	}
	if receivedPass != "" {
		t.Fatalf("expected empty password argument, got %q", receivedPass)
	}
}

func TestWebAppsCreateResolvesSessionBeforeTimeoutContext(t *testing.T) {
	origResolveSession := resolveSessionFn
	t.Cleanup(func() {
		resolveSessionFn = origResolveSession
	})

	resolveErr := errors.New("stop before network call")
	hadDeadline := false
	resolveSessionFn = func(ctx context.Context, appleID, password, twoFactorCode string, usePasswordStdin bool) (*webcore.AuthSession, string, error) {
		_, hadDeadline = ctx.Deadline()
		return nil, "", resolveErr
	}

	cmd := WebAppsCreateCommand()
	if err := cmd.FlagSet.Parse([]string{
		"--name", "My App",
		"--bundle-id", "com.example.app",
		"--sku", "SKU123",
		"--apple-id", "user@example.com",
	}); err != nil {
		t.Fatalf("parse error: %v", err)
	}

	err := cmd.Exec(context.Background(), nil)
	if !errors.Is(err, resolveErr) {
		t.Fatalf("expected resolveSession error %v, got %v", resolveErr, err)
	}
	if hadDeadline {
		t.Fatal("expected resolveSession to run before timeout context creation")
	}
}

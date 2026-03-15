package web

import (
	"context"
	"errors"
	"flag"
	"strings"
	"testing"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
	webcore "github.com/rudrankriyam/App-Store-Connect-CLI/internal/web"
)

func TestWebAuthCapabilitiesRejectsPositionalArgs(t *testing.T) {
	cmd := WebAuthCapabilitiesCommand()
	if err := cmd.FlagSet.Parse([]string{"extra"}); err != nil {
		t.Fatalf("parse error: %v", err)
	}

	err := cmd.Exec(context.Background(), []string{"extra"})
	if !errors.Is(err, flag.ErrHelp) {
		t.Fatalf("expected usage error, got %v", err)
	}
}

func TestWebAuthCapabilitiesKeyIDBypassesLocalAuthResolution(t *testing.T) {
	origResolveAuth := resolveWebAuthCredentialsFn
	origResolveSession := resolveSessionFn
	origNewClient := newWebAuthClientFn
	origLookup := lookupWebAuthKeyFn
	t.Cleanup(func() {
		resolveWebAuthCredentialsFn = origResolveAuth
		resolveSessionFn = origResolveSession
		newWebAuthClientFn = origNewClient
		lookupWebAuthKeyFn = origLookup
	})

	resolveWebAuthCredentialsFn = func(profile string) (shared.ResolvedAuthCredentials, error) {
		t.Fatal("did not expect local auth resolution when --key-id is provided")
		return shared.ResolvedAuthCredentials{}, nil
	}
	resolveSessionFn = func(ctx context.Context, appleID, password, twoFactorCode string) (*webcore.AuthSession, string, error) {
		return &webcore.AuthSession{}, "cache", nil
	}
	newWebAuthClientFn = func(session *webcore.AuthSession) *webcore.Client {
		return &webcore.Client{}
	}
	lookupWebAuthKeyFn = func(ctx context.Context, client *webcore.Client, keyID string) (*webcore.APIKeyRoleLookup, error) {
		if keyID != "39MX87M9Y4" {
			t.Fatalf("expected key-id override, got %q", keyID)
		}
		return &webcore.APIKeyRoleLookup{
			KeyID:      "39MX87M9Y4",
			Kind:       "team",
			Roles:      []string{"APP_MANAGER"},
			RoleSource: "key",
			Active:     true,
			Lookup:     "team_keys",
		}, nil
	}

	cmd := WebAuthCapabilitiesCommand()
	if err := cmd.FlagSet.Parse([]string{"--key-id", "39MX87M9Y4", "--output", "json"}); err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := cmd.Exec(context.Background(), nil); err != nil {
		t.Fatalf("Exec() error: %v", err)
	}
}

func TestWebAuthCapabilitiesResolvesCurrentAuthKeyID(t *testing.T) {
	origResolveAuth := resolveWebAuthCredentialsFn
	origResolveSession := resolveSessionFn
	origNewClient := newWebAuthClientFn
	origLookup := lookupWebAuthKeyFn
	t.Cleanup(func() {
		resolveWebAuthCredentialsFn = origResolveAuth
		resolveSessionFn = origResolveSession
		newWebAuthClientFn = origNewClient
		lookupWebAuthKeyFn = origLookup
	})

	resolveWebAuthCredentialsFn = func(profile string) (shared.ResolvedAuthCredentials, error) {
		return shared.ResolvedAuthCredentials{
			KeyID:   "ENVKEY",
			Profile: "client",
		}, nil
	}
	resolveSessionFn = func(ctx context.Context, appleID, password, twoFactorCode string) (*webcore.AuthSession, string, error) {
		return &webcore.AuthSession{}, "cache", nil
	}
	newWebAuthClientFn = func(session *webcore.AuthSession) *webcore.Client {
		return &webcore.Client{}
	}
	lookupWebAuthKeyFn = func(ctx context.Context, client *webcore.Client, keyID string) (*webcore.APIKeyRoleLookup, error) {
		if keyID != "ENVKEY" {
			t.Fatalf("expected resolved key id, got %q", keyID)
		}
		return &webcore.APIKeyRoleLookup{
			KeyID:      keyID,
			Kind:       "team",
			Roles:      []string{"APP_MANAGER"},
			RoleSource: "key",
			Active:     true,
			Lookup:     "team_keys",
		}, nil
	}

	cmd := WebAuthCapabilitiesCommand()
	if err := cmd.FlagSet.Parse([]string{"--output", "json"}); err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := cmd.Exec(context.Background(), nil); err != nil {
		t.Fatalf("Exec() error: %v", err)
	}
}

func TestWrapWebAuthCapabilitiesErrorFormatsLookupFailures(t *testing.T) {
	err := wrapWebAuthCapabilitiesError("missing", webcore.ErrAPIKeyNotFound)
	if err == nil || !strings.Contains(err.Error(), "not found in App Store Connect web key lists") {
		t.Fatalf("unexpected not-found error: %v", err)
	}

	err = wrapWebAuthCapabilitiesError("missing", webcore.ErrAPIKeyRolesUnresolved)
	if err == nil || !strings.Contains(err.Error(), "exact roles could not be resolved") {
		t.Fatalf("unexpected unresolved error: %v", err)
	}
}

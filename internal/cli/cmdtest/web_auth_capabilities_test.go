package cmdtest

import (
	"context"
	"encoding/json"
	"testing"

	cmd "github.com/rudrankriyam/App-Store-Connect-CLI/cmd"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
	webcmd "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/web"
	webcore "github.com/rudrankriyam/App-Store-Connect-CLI/internal/web"
)

func stubWebAuthCapabilitiesLookup(t *testing.T, fn func(context.Context, *webcore.Client, string) (*webcore.APIKeyRoleLookup, error)) {
	t.Helper()

	restoreSession := webcmd.SetResolveWebSession(func(ctx context.Context, appleID, password, twoFactorCode string) (*webcore.AuthSession, string, error) {
		return &webcore.AuthSession{}, "cache", nil
	})
	restoreClient := webcmd.SetNewWebAuthClient(func(session *webcore.AuthSession) *webcore.Client {
		return &webcore.Client{}
	})
	restoreLookup := webcmd.SetLookupWebAuthKey(fn)
	t.Cleanup(restoreLookup)
	t.Cleanup(restoreClient)
	t.Cleanup(restoreSession)
}

func TestWebAuthCapabilitiesRunWithKeyIDOutputsJSON(t *testing.T) {
	restoreResolve := webcmd.SetResolveWebAuthCredentials(func(profile string) (shared.ResolvedAuthCredentials, error) {
		t.Fatal("did not expect local auth resolution when --key-id is provided")
		return shared.ResolvedAuthCredentials{}, nil
	})
	t.Cleanup(restoreResolve)

	stubWebAuthCapabilitiesLookup(t, func(ctx context.Context, client *webcore.Client, keyID string) (*webcore.APIKeyRoleLookup, error) {
		return &webcore.APIKeyRoleLookup{
			KeyID:      keyID,
			Name:       "asc_cli",
			Kind:       "team",
			Roles:      []string{"APP_MANAGER"},
			RoleSource: "key",
			Active:     true,
			Lookup:     "team_keys",
		}, nil
	})

	var code int
	stdout, stderr := captureOutput(t, func() {
		code = cmd.Run([]string{"web", "auth", "capabilities", "--key-id", "39MX87M9Y4", "--output", "json"}, "1.0.0")
	})
	if code != cmd.ExitSuccess {
		t.Fatalf("exit code = %d, want %d; stderr=%q", code, cmd.ExitSuccess, stderr)
	}
	if stderr != "" {
		t.Fatalf("expected empty stderr, got %q", stderr)
	}

	var payload struct {
		KeyID        string   `json:"keyId"`
		ResolvedFrom string   `json:"resolvedFrom"`
		Roles        []string `json:"roles"`
	}
	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("json.Unmarshal() error: %v; stdout=%q", err, stdout)
	}
	if payload.KeyID != "39MX87M9Y4" || payload.ResolvedFrom != "flag" {
		t.Fatalf("unexpected payload: %+v", payload)
	}
	if len(payload.Roles) != 1 || payload.Roles[0] != "APP_MANAGER" {
		t.Fatalf("unexpected roles: %#v", payload.Roles)
	}
}

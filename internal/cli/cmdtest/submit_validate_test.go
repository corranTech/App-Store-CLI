package cmdtest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/submit"
)

type submitValidateRoundTripFunc func(*http.Request) (*http.Response, error)

func (fn submitValidateRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func setupSubmitValidateAuth(t *testing.T) {
	t.Helper()
	tempDir := t.TempDir()
	keyPath := filepath.Join(tempDir, "AuthKey.p8")
	writeECDSAPEM(t, keyPath)
	t.Setenv("ASC_BYPASS_KEYCHAIN", "1")
	t.Setenv("ASC_KEY_ID", "TEST_KEY")
	t.Setenv("ASC_ISSUER_ID", "TEST_ISSUER")
	t.Setenv("ASC_PRIVATE_KEY_PATH", keyPath)
}

func submitValidateJSONResponse(status int, body string) (*http.Response, error) {
	return &http.Response{
		Status:     fmt.Sprintf("%d %s", status, http.StatusText(status)),
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}, nil
}

func TestSubmitValidateRequiresVersionFlag(t *testing.T) {
	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	var runErr error
	_, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{"submit", "validate", "--app", "app-1"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		runErr = root.Run(context.Background())
	})

	if runErr == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(stderr, "--version or --version-id is required") {
		t.Fatalf("expected version required error in stderr, got %q", stderr)
	}
}

func TestSubmitValidateRequiresAppFlag(t *testing.T) {
	t.Setenv("ASC_APP_ID", "")

	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	var runErr error
	_, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{"submit", "validate", "--version", "1.0.0"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		runErr = root.Run(context.Background())
	})

	if runErr == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(stderr, "--app is required") {
		t.Fatalf("expected app required error in stderr, got %q", stderr)
	}
}

func TestSubmitValidateVersionAndVersionIDMutuallyExclusive(t *testing.T) {
	setupSubmitValidateAuth(t)

	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	var runErr error
	_, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{"submit", "validate", "--app", "app-1", "--version", "1.0.0", "--version-id", "ver-1"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		runErr = root.Run(context.Background())
	})

	if runErr == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(stderr, "mutually exclusive") {
		t.Fatalf("expected mutually exclusive error in stderr, got %q", stderr)
	}
}

func TestSubmitValidateAllChecksPass(t *testing.T) {
	setupSubmitValidateAuth(t)

	originalTransport := http.DefaultTransport
	t.Cleanup(func() {
		http.DefaultTransport = originalTransport
	})

	http.DefaultTransport = submitValidateRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		path := req.URL.Path
		switch {
		// Version lookup
		case req.Method == http.MethodGet && path == "/v1/appStoreVersions/ver-1":
			return submitValidateJSONResponse(http.StatusOK, `{
				"data":{"type":"appStoreVersions","id":"ver-1","attributes":{"appVersionState":"PREPARE_FOR_SUBMISSION","platform":"IOS","versionString":"1.0.0"}}
			}`)

		// Build attached
		case req.Method == http.MethodGet && path == "/v1/appStoreVersions/ver-1/build":
			return submitValidateJSONResponse(http.StatusOK, `{
				"data":{"type":"builds","id":"build-1","attributes":{"version":"100"}}
			}`)

		// Version localizations
		case req.Method == http.MethodGet && path == "/v1/appStoreVersions/ver-1/appStoreVersionLocalizations":
			return submitValidateJSONResponse(http.StatusOK, `{
				"data":[{"type":"appStoreVersionLocalizations","id":"loc-1","attributes":{"locale":"en-US","description":"A great app","keywords":"app,great"}}]
			}`)

		// Screenshot sets for localization
		case req.Method == http.MethodGet && path == "/v1/appStoreVersionLocalizations/loc-1/appScreenshotSets":
			return submitValidateJSONResponse(http.StatusOK, `{
				"data":[{"type":"appScreenshotSets","id":"set-1","attributes":{"screenshotDisplayType":"APP_IPHONE_67"}}]
			}`)

		// Screenshots in set
		case req.Method == http.MethodGet && path == "/v1/appScreenshotSets/set-1/appScreenshots":
			return submitValidateJSONResponse(http.StatusOK, `{
				"data":[{"type":"appScreenshots","id":"ss-1","attributes":{"fileName":"screenshot1.png"}}]
			}`)

		// App infos
		case req.Method == http.MethodGet && path == "/v1/apps/app-1/appInfos":
			return submitValidateJSONResponse(http.StatusOK, `{
				"data":[{"type":"appInfos","id":"info-1"}]
			}`)

		// App info localizations
		case req.Method == http.MethodGet && path == "/v1/appInfos/info-1/appInfoLocalizations":
			return submitValidateJSONResponse(http.StatusOK, `{
				"data":[{"type":"appInfoLocalizations","id":"ailoc-1","attributes":{"locale":"en-US","name":"My App","privacyPolicyUrl":"https://example.com/privacy"}}]
			}`)

		// Age rating
		case req.Method == http.MethodGet && path == "/v1/appStoreVersions/ver-1/ageRatingDeclaration":
			return submitValidateJSONResponse(http.StatusOK, `{
				"data":{"type":"ageRatingDeclarations","id":"age-1","attributes":{}}
			}`)

		default:
			return nil, fmt.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
		}
	})

	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	stdout, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{"submit", "validate", "--app", "app-1", "--version-id", "ver-1"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		if err := root.Run(context.Background()); err != nil {
			t.Fatalf("run error: %v", err)
		}
	})

	if stderr != "" {
		t.Fatalf("expected empty stderr, got %q", stderr)
	}

	var result submit.SubmitValidateResult
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("failed to parse JSON output: %v (stdout=%q)", err, stdout)
	}
	if !result.Ready {
		t.Fatalf("expected ready=true, got %+v", result)
	}
	if result.ErrorCount != 0 {
		t.Fatalf("expected 0 errors, got %d: %+v", result.ErrorCount, result.Issues)
	}
}

func TestSubmitValidateDetectsMissingBuildAndDescription(t *testing.T) {
	setupSubmitValidateAuth(t)

	originalTransport := http.DefaultTransport
	t.Cleanup(func() {
		http.DefaultTransport = originalTransport
	})

	http.DefaultTransport = submitValidateRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		path := req.URL.Path
		switch {
		// Version — editable state
		case req.Method == http.MethodGet && path == "/v1/appStoreVersions/ver-1":
			return submitValidateJSONResponse(http.StatusOK, `{
				"data":{"type":"appStoreVersions","id":"ver-1","attributes":{"appVersionState":"PREPARE_FOR_SUBMISSION","platform":"IOS","versionString":"1.0.0"}}
			}`)

		// Build NOT attached (404)
		case req.Method == http.MethodGet && path == "/v1/appStoreVersions/ver-1/build":
			return submitValidateJSONResponse(http.StatusNotFound, `{"errors":[{"status":"404","code":"NOT_FOUND","title":"Not Found"}]}`)

		// Version localizations — missing description
		case req.Method == http.MethodGet && path == "/v1/appStoreVersions/ver-1/appStoreVersionLocalizations":
			return submitValidateJSONResponse(http.StatusOK, `{
				"data":[{"type":"appStoreVersionLocalizations","id":"loc-1","attributes":{"locale":"en-US","description":"","keywords":"app"}}]
			}`)

		// Screenshot sets — none
		case req.Method == http.MethodGet && path == "/v1/appStoreVersionLocalizations/loc-1/appScreenshotSets":
			return submitValidateJSONResponse(http.StatusOK, `{"data":[]}`)

		// App infos
		case req.Method == http.MethodGet && path == "/v1/apps/app-1/appInfos":
			return submitValidateJSONResponse(http.StatusOK, `{"data":[{"type":"appInfos","id":"info-1"}]}`)

		// App info localizations — valid
		case req.Method == http.MethodGet && path == "/v1/appInfos/info-1/appInfoLocalizations":
			return submitValidateJSONResponse(http.StatusOK, `{
				"data":[{"type":"appInfoLocalizations","id":"ailoc-1","attributes":{"locale":"en-US","name":"My App","privacyPolicyUrl":"https://example.com/privacy"}}]
			}`)

		// Age rating — exists
		case req.Method == http.MethodGet && path == "/v1/appStoreVersions/ver-1/ageRatingDeclaration":
			return submitValidateJSONResponse(http.StatusOK, `{"data":{"type":"ageRatingDeclarations","id":"age-1","attributes":{}}}`)

		default:
			return nil, fmt.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
		}
	})

	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	stdout, _ := captureOutput(t, func() {
		if err := root.Parse([]string{"submit", "validate", "--app", "app-1", "--version-id", "ver-1"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		// Expect error because there are validation issues
		err := root.Run(context.Background())
		if err == nil {
			t.Fatal("expected error due to validation issues, got nil")
		}
	})

	var result submit.SubmitValidateResult
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("failed to parse JSON output: %v (stdout=%q)", err, stdout)
	}
	if result.Ready {
		t.Fatal("expected ready=false")
	}

	// Should have errors for: build missing, description empty, screenshots missing
	if result.ErrorCount < 2 {
		t.Fatalf("expected at least 2 errors, got %d: %+v", result.ErrorCount, result.Issues)
	}

	// Check specific issues
	issueChecks := make(map[string]bool)
	for _, issue := range result.Issues {
		if issue.Severity == "error" {
			issueChecks[issue.Check] = true
		}
	}
	if !issueChecks["build"] {
		t.Error("expected 'build' error issue")
	}
	if !issueChecks["description"] {
		t.Error("expected 'description' error issue")
	}
	if !issueChecks["screenshots"] {
		t.Error("expected 'screenshots' error issue")
	}
}

func TestSubmitValidateDetectsNonEditableState(t *testing.T) {
	setupSubmitValidateAuth(t)

	originalTransport := http.DefaultTransport
	t.Cleanup(func() {
		http.DefaultTransport = originalTransport
	})

	http.DefaultTransport = submitValidateRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		path := req.URL.Path
		switch {
		// Version in WAITING_FOR_REVIEW (non-editable)
		case req.Method == http.MethodGet && path == "/v1/appStoreVersions/ver-1":
			return submitValidateJSONResponse(http.StatusOK, `{
				"data":{"type":"appStoreVersions","id":"ver-1","attributes":{"appVersionState":"WAITING_FOR_REVIEW","platform":"IOS","versionString":"1.0.0"}}
			}`)

		// Build exists
		case req.Method == http.MethodGet && path == "/v1/appStoreVersions/ver-1/build":
			return submitValidateJSONResponse(http.StatusOK, `{"data":{"type":"builds","id":"build-1","attributes":{"version":"100"}}}`)

		// Version localizations
		case req.Method == http.MethodGet && path == "/v1/appStoreVersions/ver-1/appStoreVersionLocalizations":
			return submitValidateJSONResponse(http.StatusOK, `{
				"data":[{"type":"appStoreVersionLocalizations","id":"loc-1","attributes":{"locale":"en-US","description":"Desc","keywords":"kw"}}]
			}`)

		// Screenshot sets
		case req.Method == http.MethodGet && path == "/v1/appStoreVersionLocalizations/loc-1/appScreenshotSets":
			return submitValidateJSONResponse(http.StatusOK, `{
				"data":[{"type":"appScreenshotSets","id":"set-1","attributes":{"screenshotDisplayType":"APP_IPHONE_67"}}]
			}`)

		// Screenshots
		case req.Method == http.MethodGet && path == "/v1/appScreenshotSets/set-1/appScreenshots":
			return submitValidateJSONResponse(http.StatusOK, `{"data":[{"type":"appScreenshots","id":"ss-1","attributes":{"fileName":"s.png"}}]}`)

		// App infos
		case req.Method == http.MethodGet && path == "/v1/apps/app-1/appInfos":
			return submitValidateJSONResponse(http.StatusOK, `{"data":[{"type":"appInfos","id":"info-1"}]}`)

		// App info localizations
		case req.Method == http.MethodGet && path == "/v1/appInfos/info-1/appInfoLocalizations":
			return submitValidateJSONResponse(http.StatusOK, `{
				"data":[{"type":"appInfoLocalizations","id":"ailoc-1","attributes":{"locale":"en-US","name":"My App","privacyPolicyUrl":"https://example.com/privacy"}}]
			}`)

		// Age rating
		case req.Method == http.MethodGet && path == "/v1/appStoreVersions/ver-1/ageRatingDeclaration":
			return submitValidateJSONResponse(http.StatusOK, `{"data":{"type":"ageRatingDeclarations","id":"age-1","attributes":{}}}`)

		default:
			return nil, fmt.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
		}
	})

	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	stdout, _ := captureOutput(t, func() {
		if err := root.Parse([]string{"submit", "validate", "--app", "app-1", "--version-id", "ver-1"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		err := root.Run(context.Background())
		if err == nil {
			t.Fatal("expected error due to non-editable state, got nil")
		}
	})

	var result submit.SubmitValidateResult
	if err := json.Unmarshal([]byte(stdout), &result); err != nil {
		t.Fatalf("failed to parse JSON output: %v (stdout=%q)", err, stdout)
	}
	if result.Ready {
		t.Fatal("expected ready=false")
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Check == "version_state" && issue.Severity == "error" {
			found = true
			if !strings.Contains(issue.Message, "non-editable") {
				t.Fatalf("expected non-editable message, got %q", issue.Message)
			}
		}
	}
	if !found {
		t.Fatalf("expected version_state error, got issues: %+v", result.Issues)
	}
}

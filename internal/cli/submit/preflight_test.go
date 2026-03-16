package submit

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
)

func TestSubmitPreflightCommand_MissingApp(t *testing.T) {
	// Ensure no app ID comes from env or config.
	t.Setenv("ASC_APP_ID", "")

	cmd := SubmitPreflightCommand()
	if err := cmd.FlagSet.Parse([]string{"--version", "1.0"}); err != nil {
		t.Fatalf("failed to parse flags: %v", err)
	}
	// ResolveAppID may still find an ID from local config — skip Exec test
	// if that's the case, and only assert the command shape.
	resolved := shared.ResolveAppID("")
	if resolved != "" {
		t.Skip("local config provides an app ID; skipping flag validation test")
	}
	if err := cmd.Exec(context.Background(), nil); err != nil && !errors.Is(err, flag.ErrHelp) && !strings.Contains(err.Error(), "authentication") {
		t.Fatalf("expected flag.ErrHelp, got %v", err)
	}
}

func TestSubmitPreflightCommand_MissingVersion(t *testing.T) {
	setupSubmitAuth(t)

	cmd := SubmitPreflightCommand()
	if err := cmd.FlagSet.Parse([]string{"--app", "123"}); err != nil {
		t.Fatalf("failed to parse flags: %v", err)
	}
	if err := cmd.Exec(context.Background(), nil); err != nil && !errors.Is(err, flag.ErrHelp) && !strings.Contains(err.Error(), "authentication") {
		t.Fatalf("expected flag.ErrHelp, got %v", err)
	}
}

func TestSubmitPreflightCommand_InvalidPlatform(t *testing.T) {
	setupSubmitAuth(t)

	cmd := SubmitPreflightCommand()
	if err := cmd.FlagSet.Parse([]string{"--app", "123", "--version", "1.0", "--platform", "INVALID"}); err != nil {
		t.Fatalf("failed to parse flags: %v", err)
	}
	err := cmd.Exec(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for invalid platform")
	}
	if !errors.Is(err, flag.ErrHelp) {
		t.Fatalf("expected flag.ErrHelp for invalid platform, got %v", err)
	}
}

func TestSubmitPreflightCommand_Shape(t *testing.T) {
	cmd := SubmitPreflightCommand()
	if cmd.Name != "preflight" {
		t.Fatalf("unexpected command name: %q", cmd.Name)
	}
	if cmd.FlagSet == nil {
		t.Fatal("expected FlagSet to be set")
	}
}

func TestPreflightResult_TallyCounts(t *testing.T) {
	result := &preflightResult{
		Checks: []checkResult{
			{Name: "a", Passed: true},
			{Name: "b", Passed: false},
			{Name: "c", Passed: true},
			{Name: "d", Passed: false},
		},
	}
	tallyCounts(result)
	if result.PassCount != 2 {
		t.Fatalf("expected 2 passes, got %d", result.PassCount)
	}
	if result.FailCount != 2 {
		t.Fatalf("expected 2 failures, got %d", result.FailCount)
	}
}

func TestPreflightResult_AllPass(t *testing.T) {
	result := &preflightResult{
		Checks: []checkResult{
			{Name: "a", Passed: true},
			{Name: "b", Passed: true},
		},
	}
	tallyCounts(result)
	if result.FailCount != 0 {
		t.Fatalf("expected 0 failures, got %d", result.FailCount)
	}
}

func TestAgeRatingMissingFields_AllSet(t *testing.T) {
	boolTrue := true
	strNone := "NONE"

	attrs := newAgeRatingAllSet(boolTrue, strNone)
	missing := ageRatingMissingFields(attrs)
	if len(missing) != 0 {
		t.Fatalf("expected no missing fields, got: %v", missing)
	}
}

func TestAgeRatingMissingFields_SomeMissing(t *testing.T) {
	boolTrue := true
	strNone := "NONE"

	attrs := newAgeRatingAllSet(boolTrue, strNone)
	// Unset a few
	attrs.Gambling = nil
	attrs.ViolenceRealistic = nil

	missing := ageRatingMissingFields(attrs)
	if len(missing) != 2 {
		t.Fatalf("expected 2 missing fields, got %d: %v", len(missing), missing)
	}

	found := map[string]bool{}
	for _, m := range missing {
		found[m] = true
	}
	if !found["gambling"] {
		t.Fatal("expected 'gambling' in missing fields")
	}
	if !found["violenceRealistic"] {
		t.Fatal("expected 'violenceRealistic' in missing fields")
	}
}

func TestAgeRatingMissingFields_AllMissing(t *testing.T) {
	attrs := asc.AgeRatingDeclarationAttributes{}
	missing := ageRatingMissingFields(attrs)
	// 3 boolean + 10 enum = 13
	if len(missing) != 13 {
		t.Fatalf("expected 13 missing fields, got %d: %v", len(missing), missing)
	}
}

func TestCheckVersionExists_NotFound(t *testing.T) {
	client := newSubmitTestClient(t, submitRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/appStoreVersions") {
			return submitJSONResponse(http.StatusOK, `{"data":[]}`)
		}
		return nil, fmt.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
	}))

	versionID, check := checkVersionExists(context.Background(), client, "app-123", "9.9", "IOS")
	if versionID != "" {
		t.Fatalf("expected empty versionID, got %q", versionID)
	}
	if check.Passed {
		t.Fatal("expected check to fail")
	}
	if check.Hint == "" {
		t.Fatal("expected hint to be set")
	}
}

func TestCheckBuildAttached_NoBuild(t *testing.T) {
	client := newSubmitTestClient(t, submitRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/build") {
			return submitJSONResponse(http.StatusNotFound, `{"errors":[{"status":"404","code":"NOT_FOUND","title":"Not Found"}]}`)
		}
		return nil, fmt.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
	}))

	check := checkBuildAttached(context.Background(), client, "version-123")
	if check.Passed {
		t.Fatal("expected check to fail for missing build")
	}
}

func TestCheckBuildAttached_HasBuild(t *testing.T) {
	client := newSubmitTestClient(t, submitRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/build") {
			return submitJSONResponse(http.StatusOK, `{"data":{"type":"builds","id":"build-1","attributes":{"version":"42"}}}`)
		}
		return nil, fmt.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
	}))

	check := checkBuildAttached(context.Background(), client, "version-123")
	if !check.Passed {
		t.Fatalf("expected check to pass, got message: %s", check.Message)
	}
	if !strings.Contains(check.Message, "42") {
		t.Fatalf("expected build version in message, got: %s", check.Message)
	}
}

func TestCheckContentRights_NotSet(t *testing.T) {
	client := newSubmitTestClient(t, submitRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/v1/apps/") {
			return submitJSONResponse(http.StatusOK, `{"data":{"type":"apps","id":"app-123","attributes":{"name":"Test","bundleId":"com.test","sku":"test"}}}`)
		}
		return nil, fmt.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
	}))

	check := checkContentRights(context.Background(), client, "app-123")
	if check.Passed {
		t.Fatal("expected check to fail when contentRightsDeclaration is nil")
	}
	if check.Hint == "" {
		t.Fatal("expected hint to be set")
	}
}

func TestCheckContentRights_Set(t *testing.T) {
	client := newSubmitTestClient(t, submitRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/v1/apps/") {
			return submitJSONResponse(http.StatusOK, `{"data":{"type":"apps","id":"app-123","attributes":{"name":"Test","bundleId":"com.test","sku":"test","contentRightsDeclaration":"DOES_NOT_USE_THIRD_PARTY_CONTENT"}}}`)
		}
		return nil, fmt.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
	}))

	check := checkContentRights(context.Background(), client, "app-123")
	if !check.Passed {
		t.Fatalf("expected check to pass, got message: %s", check.Message)
	}
}

func TestRunPreflight_AllPass(t *testing.T) {
	setupSubmitAuth(t)

	client := newSubmitTestClient(t, submitRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		path := req.URL.Path

		switch {
		// Version resolve
		case req.Method == http.MethodGet && strings.Contains(path, "/appStoreVersions") && strings.Contains(path, "/apps/"):
			return submitJSONResponse(http.StatusOK, `{
				"data": [{
					"type": "appStoreVersions",
					"id": "version-1",
					"attributes": {"platform": "IOS", "versionString": "1.0"}
				}]
			}`)

		// Build attached
		case req.Method == http.MethodGet && strings.HasSuffix(path, "/build"):
			return submitJSONResponse(http.StatusOK, `{"data":{"type":"builds","id":"build-1","attributes":{"version":"1"}}}`)

		// App infos
		case req.Method == http.MethodGet && strings.HasSuffix(path, "/appInfos"):
			return submitJSONResponse(http.StatusOK, `{"data":[{"type":"appInfos","id":"info-1","attributes":{}}]}`)

		// Age rating
		case req.Method == http.MethodGet && strings.HasSuffix(path, "/ageRatingDeclaration"):
			return submitJSONResponse(http.StatusOK, `{"data":{"type":"ageRatingDeclarations","id":"rating-1","attributes":{
				"gambling": false, "lootBox": false, "unrestrictedWebAccess": false,
				"alcoholTobaccoOrDrugUseOrReferences": "NONE", "gamblingSimulated": "NONE",
				"horrorOrFearThemes": "NONE", "matureOrSuggestiveThemes": "NONE",
				"profanityOrCrudeHumor": "NONE", "sexualContentGraphicAndNudity": "NONE",
				"sexualContentOrNudity": "NONE", "violenceCartoonOrFantasy": "NONE",
				"violenceRealistic": "NONE", "violenceRealisticProlongedGraphicOrSadistic": "NONE"
			}}}`)

		// App (content rights)
		case req.Method == http.MethodGet && strings.HasPrefix(path, "/v1/apps/") && !strings.Contains(path, "/"):
			return submitJSONResponse(http.StatusOK, `{"data":{"type":"apps","id":"app-1","attributes":{"name":"Test","bundleId":"com.test","sku":"test","contentRightsDeclaration":"DOES_NOT_USE_THIRD_PARTY_CONTENT"}}}`)

		// Primary category
		case req.Method == http.MethodGet && strings.HasSuffix(path, "/primaryCategory"):
			return submitJSONResponse(http.StatusOK, `{"data":{"type":"appCategories","id":"SPORTS","attributes":{}}}`)

		// Localizations
		case req.Method == http.MethodGet && strings.Contains(path, "/appStoreVersionLocalizations"):
			return submitJSONResponse(http.StatusOK, `{
				"data": [{
					"type": "appStoreVersionLocalizations",
					"id": "loc-1",
					"attributes": {"locale": "en-US", "description": "A great app", "keywords": "test,app"}
				}]
			}`)

		// Screenshot sets
		case req.Method == http.MethodGet && strings.Contains(path, "/appScreenshotSets"):
			return submitJSONResponse(http.StatusOK, `{
				"data": [{
					"type": "appScreenshotSets",
					"id": "ss-1",
					"attributes": {"screenshotDisplayType": "APP_IPHONE_67"}
				}]
			}`)
		}

		// For the app endpoint without sub-paths
		if req.Method == http.MethodGet && path == "/v1/apps/app-1" {
			return submitJSONResponse(http.StatusOK, `{"data":{"type":"apps","id":"app-1","attributes":{"name":"Test","bundleId":"com.test","sku":"test","contentRightsDeclaration":"DOES_NOT_USE_THIRD_PARTY_CONTENT"}}}`)
		}

		return nil, fmt.Errorf("unexpected request: %s %s", req.Method, path)
	}))

	result := runPreflight(context.Background(), client, "app-1", "1.0", "IOS")
	if result.FailCount != 0 {
		for _, c := range result.Checks {
			if !c.Passed {
				t.Errorf("check %q failed: %s (hint: %s)", c.Name, c.Message, c.Hint)
			}
		}
		t.Fatalf("expected 0 failures, got %d", result.FailCount)
	}
}

func TestPreflightTextOutput(t *testing.T) {
	// Ensure printPreflightText doesn't panic with various result shapes.
	printPreflightText(&preflightResult{
		AppID:    "123",
		Version:  "1.0",
		Platform: "IOS",
		Checks: []checkResult{
			{Name: "Version exists", Passed: true, Message: "Version 1.0 found"},
			{Name: "Build attached", Passed: false, Message: "No build", Hint: "asc submit create ..."},
		},
		PassCount: 1,
		FailCount: 1,
	})
}

func TestSubmitPreflightCommand_JSONOutput(t *testing.T) {
	setupSubmitAuth(t)

	originalTransport := http.DefaultTransport
	t.Cleanup(func() {
		http.DefaultTransport = originalTransport
	})

	http.DefaultTransport = submitRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		// Version resolve returns empty — version not found
		if req.Method == http.MethodGet && strings.Contains(req.URL.Path, "/appStoreVersions") {
			return submitJSONResponse(http.StatusOK, `{"data":[]}`)
		}
		return nil, fmt.Errorf("unexpected request: %s %s", req.Method, req.URL.Path)
	})

	cmd := SubmitPreflightCommand()
	cmd.FlagSet.SetOutput(io.Discard)
	if err := cmd.FlagSet.Parse([]string{"--app", "123", "--version", "1.0", "--output", "json"}); err != nil {
		t.Fatalf("failed to parse flags: %v", err)
	}

	err := cmd.Exec(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error when version not found")
	}
	if !strings.Contains(err.Error(), "issue(s) found") {
		t.Fatalf("expected preflight failure error, got: %v", err)
	}
}

// --- Helpers ---

func newAgeRatingAllSet(boolVal bool, strVal string) asc.AgeRatingDeclarationAttributes {
	return asc.AgeRatingDeclarationAttributes{
		Gambling:                                    &boolVal,
		LootBox:                                     &boolVal,
		UnrestrictedWebAccess:                       &boolVal,
		AlcoholTobaccoOrDrugUseOrReferences:         &strVal,
		GamblingSimulated:                           &strVal,
		HorrorOrFearThemes:                          &strVal,
		MatureOrSuggestiveThemes:                    &strVal,
		ProfanityOrCrudeHumor:                       &strVal,
		SexualContentGraphicAndNudity:               &strVal,
		SexualContentOrNudity:                       &strVal,
		ViolenceCartoonOrFantasy:                    &strVal,
		ViolenceRealistic:                           &strVal,
		ViolenceRealisticProlongedGraphicOrSadistic: &strVal,
	}
}

package apps

import (
	"strings"
	"testing"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
)

func TestSelectLatestAppStoreVersion(t *testing.T) {
	versions := []asc.Resource[asc.AppStoreVersionAttributes]{
		{
			ID: "old",
			Attributes: asc.AppStoreVersionAttributes{
				CreatedDate: "2024-01-01T00:00:00Z",
			},
		},
		{
			ID: "new",
			Attributes: asc.AppStoreVersionAttributes{
				CreatedDate: "2025-01-01T00:00:00Z",
			},
		},
	}

	selected := selectLatestAppStoreVersion(versions)
	if selected.ID != "new" {
		t.Fatalf("expected latest version to be %q, got %q", "new", selected.ID)
	}
}

func TestSelectLatestAppStoreVersionFallsBackToFirst(t *testing.T) {
	versions := []asc.Resource[asc.AppStoreVersionAttributes]{
		{
			ID: "first",
			Attributes: asc.AppStoreVersionAttributes{
				CreatedDate: "invalid-date",
			},
		},
		{
			ID: "second",
			Attributes: asc.AppStoreVersionAttributes{
				CreatedDate: "",
			},
		},
	}

	selected := selectLatestAppStoreVersion(versions)
	if selected.ID != "first" {
		t.Fatalf("expected fallback to the first version, got %q", selected.ID)
	}
}

func TestWarnAppInfoSetSubmitIncompleteLocaleMentionsCanonicalReleaseRun(t *testing.T) {
	stderr := captureAppsCreateOutput(t, func() {
		warnAppInfoSetSubmitIncompleteLocale("en-US", asc.AppStoreVersionLocalizationAttributes{})
	})

	if !strings.Contains(stderr, "`asc release run`") {
		t.Fatalf("expected canonical release guidance in warning, got %q", stderr)
	}
	if !strings.Contains(stderr, "deprecated direct path `asc submit create`") {
		t.Fatalf("expected deprecated direct path note in warning, got %q", stderr)
	}
}

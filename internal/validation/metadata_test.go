package validation

import (
	"strings"
	"testing"
)

func TestMetadataLengthChecks_OverLimit(t *testing.T) {
	loc := VersionLocalization{
		Locale:      "en-US",
		Description: strings.Repeat("a", LimitDescription+1),
		Keywords:    strings.Repeat("b", LimitKeywords+1),
	}
	appInfo := AppInfoLocalization{
		Locale: "en-US",
		Name:   strings.Repeat("n", LimitName+1),
	}

	checks := metadataLengthChecks([]VersionLocalization{loc}, []AppInfoLocalization{appInfo})

	if !hasCheckID(checks, "metadata.length.description") {
		t.Fatalf("expected description length check")
	}
	if !hasCheckID(checks, "metadata.length.keywords") {
		t.Fatalf("expected keywords length check")
	}
	if !hasCheckID(checks, "metadata.length.name") {
		t.Fatalf("expected name length check")
	}
}

func TestMetadataLengthChecks_Valid(t *testing.T) {
	loc := VersionLocalization{
		Locale:      "en-US",
		Description: strings.Repeat("a", LimitDescription),
		Keywords:    strings.Repeat("b", LimitKeywords),
		WhatsNew:    strings.Repeat("c", LimitWhatsNew),
	}
	appInfo := AppInfoLocalization{
		Locale:   "en-US",
		Name:     strings.Repeat("n", LimitName),
		Subtitle: strings.Repeat("s", LimitSubtitle),
	}

	checks := metadataLengthChecks([]VersionLocalization{loc}, []AppInfoLocalization{appInfo})
	if len(checks) != 0 {
		t.Fatalf("expected no checks, got %d", len(checks))
	}
}

package metadata

import (
	"strings"
	"testing"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/validation"
)

func TestVersionLengthIssuesBoundaries(t *testing.T) {
	noIssues := versionLengthIssues("file", "1.2.3", "en-US", VersionLocalization{
		Description:     strings.Repeat("a", validation.LimitDescription),
		Keywords:        strings.Repeat("b", validation.LimitKeywords),
		WhatsNew:        strings.Repeat("c", validation.LimitWhatsNew),
		PromotionalText: strings.Repeat("d", validation.LimitPromotionalText),
	})
	if len(noIssues) != 0 {
		t.Fatalf("expected no issues at limits, got %+v", noIssues)
	}

	withIssues := versionLengthIssues("file", "1.2.3", "en-US", VersionLocalization{
		Description:     strings.Repeat("a", validation.LimitDescription+1),
		Keywords:        strings.Repeat("b", validation.LimitKeywords+1),
		WhatsNew:        strings.Repeat("c", validation.LimitWhatsNew+1),
		PromotionalText: strings.Repeat("d", validation.LimitPromotionalText+1),
	})
	if len(withIssues) != 4 {
		t.Fatalf("expected 4 issues above limits, got %d", len(withIssues))
	}
}

func TestAppInfoLengthIssuesBoundaries(t *testing.T) {
	noIssues := appInfoLengthIssues("file", "en-US", AppInfoLocalization{
		Name:     strings.Repeat("n", validation.LimitName),
		Subtitle: strings.Repeat("s", validation.LimitSubtitle),
	})
	if len(noIssues) != 0 {
		t.Fatalf("expected no issues at limits, got %+v", noIssues)
	}

	withIssues := appInfoLengthIssues("file", "en-US", AppInfoLocalization{
		Name:     strings.Repeat("n", validation.LimitName+1),
		Subtitle: strings.Repeat("s", validation.LimitSubtitle+1),
	})
	if len(withIssues) != 2 {
		t.Fatalf("expected 2 issues above limits, got %d", len(withIssues))
	}
}

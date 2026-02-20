package status

import (
	"testing"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
)

func TestParseInclude_DefaultsToAllSections(t *testing.T) {
	includes, err := parseInclude("")
	if err != nil {
		t.Fatalf("parseInclude error: %v", err)
	}

	if !includes.builds || !includes.testflight || !includes.appstore || !includes.submission || !includes.review || !includes.phasedRelease || !includes.links {
		t.Fatalf("expected all sections enabled by default, got %+v", includes)
	}
}

func TestParseInclude_RejectsUnknownSection(t *testing.T) {
	_, err := parseInclude("builds,unknown")
	if err == nil {
		t.Fatal("expected error for unknown include section")
	}
}

func TestSelectLatestAppStoreVersion_DeterministicTieBreak(t *testing.T) {
	versions := []asc.Resource[asc.AppStoreVersionAttributes]{
		{
			ID: "ver-1",
			Attributes: asc.AppStoreVersionAttributes{
				CreatedDate: "2026-02-20T00:00:00Z",
			},
		},
		{
			ID: "ver-2",
			Attributes: asc.AppStoreVersionAttributes{
				CreatedDate: "2026-02-20T00:00:00Z",
			},
		},
	}

	selected := selectLatestAppStoreVersion(versions)
	if selected == nil {
		t.Fatal("expected selected version, got nil")
	}
	if selected.ID != "ver-2" {
		t.Fatalf("expected deterministic tie-break to choose ver-2, got %q", selected.ID)
	}
}

func TestSelectLatestReviewSubmission_DeterministicTieBreak(t *testing.T) {
	submissions := []asc.ReviewSubmissionResource{
		{
			ID: "sub-1",
			Attributes: asc.ReviewSubmissionAttributes{
				SubmittedDate: "2026-02-20T00:00:00Z",
			},
		},
		{
			ID: "sub-2",
			Attributes: asc.ReviewSubmissionAttributes{
				SubmittedDate: "2026-02-20T00:00:00Z",
			},
		},
	}

	selected := selectLatestReviewSubmission(submissions)
	if selected == nil {
		t.Fatal("expected selected submission, got nil")
	}
	if selected.ID != "sub-2" {
		t.Fatalf("expected deterministic tie-break to choose sub-2, got %q", selected.ID)
	}
}

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

func TestBuildStatusSummary_RedWhenBlockingIssuesExist(t *testing.T) {
	resp := &dashboardResponse{
		Submission: &submissionSection{
			InFlight:       true,
			BlockingIssues: []string{"submission abc has unresolved issues"},
		},
	}

	summary := buildStatusSummary(resp)
	if summary.Health != "red" {
		t.Fatalf("expected health=red, got %q", summary.Health)
	}
	if summary.NextAction == "" {
		t.Fatal("expected next action")
	}
	if len(summary.Blockers) == 0 {
		t.Fatal("expected blockers")
	}
}

func TestBuildStatusSummary_YellowWhenReviewInFlight(t *testing.T) {
	resp := &dashboardResponse{
		Review: &reviewSection{
			State: "WAITING_FOR_REVIEW",
		},
	}

	summary := buildStatusSummary(resp)
	if summary.Health != "yellow" {
		t.Fatalf("expected health=yellow, got %q", summary.Health)
	}
}

func TestBuildStatusSummary_GreenWhenReadyForSale(t *testing.T) {
	resp := &dashboardResponse{
		AppStore: &appStoreSection{
			State: "READY_FOR_SALE",
		},
		Builds: &buildsSection{
			Latest: &latestBuild{ID: "build-1"},
		},
	}

	summary := buildStatusSummary(resp)
	if summary.Health != "green" {
		t.Fatalf("expected health=green, got %q", summary.Health)
	}
	if summary.NextAction != "No action needed." {
		t.Fatalf("expected no action needed, got %q", summary.NextAction)
	}
}

func TestPhasedReleaseProgressBar(t *testing.T) {
	bar := phasedReleaseProgressBar(&phasedReleaseSection{
		Configured:       true,
		CurrentDayNumber: 3,
	})
	if bar == "" {
		t.Fatal("expected progress bar")
	}
	if bar != "[####------] 3/7" {
		t.Fatalf("expected deterministic bar, got %q", bar)
	}
}

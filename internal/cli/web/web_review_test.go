package web

import (
	"testing"

	webcore "github.com/rudrankriyam/App-Store-Connect-CLI/internal/web"
)

func TestNormalizeAttachmentFilenameStripsPathComponents(t *testing.T) {
	attachment := webcore.ReviewAttachment{
		AttachmentID: "attachment-id",
		FileName:     "../../etc/passwd",
	}

	got := normalizeAttachmentFilename(attachment)
	if got != "passwd" {
		t.Fatalf("expected sanitized filename %q, got %q", "passwd", got)
	}
}

func TestNormalizeAttachmentFilenameFallsBackWhenBasenameIsInvalid(t *testing.T) {
	attachment := webcore.ReviewAttachment{
		AttachmentID: "attachment-id",
		FileName:     "../",
	}

	got := normalizeAttachmentFilename(attachment)
	if got != "attachment-id.bin" {
		t.Fatalf("expected fallback filename %q, got %q", "attachment-id.bin", got)
	}
}

func TestBuildReviewListTableRows(t *testing.T) {
	submissions := []webcore.ReviewSubmission{
		{
			ID:            "sub-1",
			State:         "UNRESOLVED_ISSUES",
			SubmittedDate: "2026-02-25T10:00:00Z",
			Platform:      "IOS",
			AppStoreVersionForReview: &webcore.AppStoreVersionForReview{
				ID:      "ver-1",
				Version: "1.2.3",
			},
		},
		{
			ID: "sub-2",
		},
	}

	rows := buildReviewListTableRows(submissions)
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if got := rows[0][0]; got != "sub-1" {
		t.Fatalf("expected first row submission id %q, got %q", "sub-1", got)
	}
	if got := rows[0][3]; got != "1.2.3" {
		t.Fatalf("expected first row version %q, got %q", "1.2.3", got)
	}
	if got := rows[1][1]; got != "n/a" {
		t.Fatalf("expected fallback state %q, got %q", "n/a", got)
	}
	if got := rows[1][3]; got != "n/a" {
		t.Fatalf("expected fallback version %q, got %q", "n/a", got)
	}
}

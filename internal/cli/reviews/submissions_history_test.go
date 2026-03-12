package reviews

import (
	"context"
	"errors"
	"flag"
	"strings"
	"testing"
)

func TestDeriveOutcome(t *testing.T) {
	tests := []struct {
		name            string
		submissionState string
		itemStates      []string
		want            string
	}{
		{
			name:            "all items approved",
			submissionState: "COMPLETE",
			itemStates:      []string{"APPROVED"},
			want:            "approved",
		},
		{
			name:            "any item rejected",
			submissionState: "COMPLETE",
			itemStates:      []string{"APPROVED", "REJECTED"},
			want:            "rejected",
		},
		{
			name:            "unresolved issues no rejected items",
			submissionState: "UNRESOLVED_ISSUES",
			itemStates:      []string{"ACCEPTED"},
			want:            "rejected",
		},
		{
			name:            "rejected item takes priority over unresolved",
			submissionState: "UNRESOLVED_ISSUES",
			itemStates:      []string{"REJECTED"},
			want:            "rejected",
		},
		{
			name:            "mixed non-rejected states falls through to submission state",
			submissionState: "COMPLETE",
			itemStates:      []string{"APPROVED", "ACCEPTED"},
			want:            "complete",
		},
		{
			name:            "no items uses submission state",
			submissionState: "WAITING_FOR_REVIEW",
			itemStates:      nil,
			want:            "waiting_for_review",
		},
		{
			name:            "in review state",
			submissionState: "IN_REVIEW",
			itemStates:      []string{"READY_FOR_REVIEW"},
			want:            "in_review",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deriveOutcome(tt.submissionState, tt.itemStates)
			if got != tt.want {
				t.Errorf("deriveOutcome(%q, %v) = %q, want %q", tt.submissionState, tt.itemStates, got, tt.want)
			}
		})
	}
}

func TestSubmissionsHistoryCommand_MissingApp(t *testing.T) {
	cmd := SubmissionsHistoryCommand()
	if cmd.Name != "submissions-history" {
		t.Fatalf("unexpected command name: %s", cmd.Name)
	}

	// Unset any env that could provide app ID
	t.Setenv("ASC_APP_ID", "")

	err := cmd.ParseAndRun(context.Background(), []string{})
	if err == nil {
		t.Fatal("expected error for missing --app, got nil")
	}
	if !errors.Is(err, flag.ErrHelp) {
		t.Fatalf("expected flag.ErrHelp, got: %v", err)
	}
}

func TestSubmissionsHistoryCommand_InvalidLimit(t *testing.T) {
	cmd := SubmissionsHistoryCommand()
	t.Setenv("ASC_APP_ID", "test-app")
	t.Setenv("ASC_BYPASS_KEYCHAIN", "1")

	err := cmd.ParseAndRun(context.Background(), []string{"--limit", "999"})
	if err == nil {
		t.Fatal("expected error for invalid limit, got nil")
	}
	if !strings.Contains(err.Error(), "--limit must be between 1 and 200") {
		t.Fatalf("unexpected error: %v", err)
	}
}

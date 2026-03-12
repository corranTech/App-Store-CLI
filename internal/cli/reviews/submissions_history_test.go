package reviews

import "testing"

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

package reviews

import (
	"strings"
)

// SubmissionHistoryEntry is the assembled result for one submission.
type SubmissionHistoryEntry struct {
	SubmissionID  string                  `json:"submissionId"`
	VersionString string                  `json:"versionString"`
	Platform      string                  `json:"platform"`
	State         string                  `json:"state"`
	SubmittedDate string                  `json:"submittedDate"`
	Outcome       string                  `json:"outcome"`
	Items         []SubmissionHistoryItem `json:"items"`
}

// SubmissionHistoryItem is a summary of one item in a submission.
type SubmissionHistoryItem struct {
	ID         string `json:"id"`
	State      string `json:"state"`
	Type       string `json:"type"`
	ResourceID string `json:"resourceId"`
}

// deriveOutcome computes a human-readable outcome from submission and item states.
// Priority order:
// 1. Any item REJECTED → "rejected"
// 2. All items APPROVED → "approved"
// 3. Submission state UNRESOLVED_ISSUES → "rejected"
// 4. Fallback → lowercase submission state
func deriveOutcome(submissionState string, itemStates []string) string {
	hasRejected := false
	allApproved := len(itemStates) > 0

	for _, s := range itemStates {
		if s == "REJECTED" {
			hasRejected = true
		}
		if s != "APPROVED" {
			allApproved = false
		}
	}

	if hasRejected {
		return "rejected"
	}
	if allApproved {
		return "approved"
	}
	if submissionState == "UNRESOLVED_ISSUES" {
		return "rejected"
	}
	return strings.ToLower(submissionState)
}

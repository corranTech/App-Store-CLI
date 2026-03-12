package reviews

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
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

// SubmissionsHistoryCommand returns the submissions-history subcommand.
func SubmissionsHistoryCommand() *ffcli.Command {
	fs := flag.NewFlagSet("submissions-history", flag.ExitOnError)

	appID := fs.String("app", "", "App Store Connect app ID (or ASC_APP_ID)")
	platform := fs.String("platform", "", "Filter by platform: IOS, MAC_OS, TV_OS, VISION_OS (comma-separated)")
	state := fs.String("state", "", "Filter by state (comma-separated)")
	version := fs.String("version", "", "Filter by version string (e.g. 1.2.0)")
	limit := fs.Int("limit", 0, "Maximum results per page (1-200)")
	paginate := fs.Bool("paginate", false, "Automatically fetch all pages (aggregate results)")
	output := shared.BindOutputFlags(fs)

	return &ffcli.Command{
		Name:       "submissions-history",
		ShortUsage: "asc review submissions-history [flags]",
		ShortHelp:  "Show enriched review submission history for an app.",
		LongHelp: `Show enriched review submission history for an app.

Each entry includes the submission state, platform, version string, submitted
date, and a derived outcome (approved, rejected, or the raw state).

Examples:
  asc review submissions-history --app "123456789"
  asc review submissions-history --app "123456789" --platform IOS --state COMPLETE
  asc review submissions-history --app "123456789" --version "1.2.0"
  asc review submissions-history --app "123456789" --paginate`,
		FlagSet:   fs,
		UsageFunc: shared.DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			if *limit != 0 && (*limit < 1 || *limit > 200) {
				return fmt.Errorf("review submissions-history: --limit must be between 1 and 200")
			}

			platforms, err := shared.NormalizeAppStoreVersionPlatforms(shared.SplitCSVUpper(*platform))
			if err != nil {
				return fmt.Errorf("review submissions-history: %w", err)
			}
			states := shared.SplitCSVUpper(*state)

			resolvedAppID := shared.ResolveAppID(*appID)
			if resolvedAppID == "" {
				fmt.Fprintln(os.Stderr, "Error: --app is required (or set ASC_APP_ID)")
				return flag.ErrHelp
			}

			client, err := shared.GetASCClient()
			if err != nil {
				return fmt.Errorf("review submissions-history: %w", err)
			}

			requestCtx, cancel := shared.ContextWithTimeout(ctx)
			defer cancel()

			// Suppress unused variable warnings until Task 3 wires up API calls.
			_ = platforms
			_ = states
			_ = paginate

			entries, err := enrichSubmissions(requestCtx, client, nil, strings.TrimSpace(*version))
			if err != nil {
				return fmt.Errorf("review submissions-history: %w", err)
			}

			tableFunc := func() error { return printHistoryTable(entries) }
			markdownFunc := func() error { return printHistoryMarkdown(entries) }
			return shared.PrintOutputWithRenderers(entries, *output.Output, *output.Pretty, tableFunc, markdownFunc)
		},
	}
}

// enrichSubmissions aggregates submission data with items.
// This is a placeholder that will be fully implemented in Task 3.
func enrichSubmissions(_ context.Context, _ *asc.Client, _ []asc.ReviewSubmissionResource, _ string) ([]SubmissionHistoryEntry, error) {
	return nil, nil
}

// printHistoryTable renders submission history as a table.
// This is a placeholder that will be fully implemented in Task 4.
func printHistoryTable(_ []SubmissionHistoryEntry) error {
	return nil
}

// printHistoryMarkdown renders submission history as markdown.
// This is a placeholder that will be fully implemented in Task 4.
func printHistoryMarkdown(_ []SubmissionHistoryEntry) error {
	return nil
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

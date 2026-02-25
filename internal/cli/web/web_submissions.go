package web

import (
	"context"
	"flag"
	"strings"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
	webcore "github.com/rudrankriyam/App-Store-Connect-CLI/internal/web"
)

var allowedReviewSubmissionStates = map[string]struct{}{
	"READY_FOR_REVIEW":   {},
	"WAITING_FOR_REVIEW": {},
	"IN_REVIEW":          {},
	"UNRESOLVED_ISSUES":  {},
	"CANCELING":          {},
	"COMPLETING":         {},
	"COMPLETE":           {},
}

func parseSubmissionStates(stateCSV string) ([]string, error) {
	states := shared.SplitCSVUpper(stateCSV)
	if len(states) == 0 {
		return nil, nil
	}
	invalid := make([]string, 0)
	seen := map[string]struct{}{}
	filtered := make([]string, 0, len(states))
	for _, state := range states {
		if _, exists := allowedReviewSubmissionStates[state]; !exists {
			invalid = append(invalid, state)
			continue
		}
		if _, exists := seen[state]; exists {
			continue
		}
		seen[state] = struct{}{}
		filtered = append(filtered, state)
	}
	if len(invalid) > 0 {
		return nil, shared.UsageErrorf("--state contains unsupported value(s): %s", strings.Join(invalid, ", "))
	}
	return filtered, nil
}

func filterSubmissionsByState(submissions []webcore.ReviewSubmission, states []string) []webcore.ReviewSubmission {
	if len(states) == 0 {
		return submissions
	}
	allowed := make(map[string]struct{}, len(states))
	for _, state := range states {
		allowed[strings.ToUpper(strings.TrimSpace(state))] = struct{}{}
	}
	result := make([]webcore.ReviewSubmission, 0, len(submissions))
	for _, submission := range submissions {
		state := strings.ToUpper(strings.TrimSpace(submission.State))
		if _, ok := allowed[state]; ok {
			result = append(result, submission)
		}
	}
	return result
}

// WebSubmissionsCommand returns the detached web review submissions group.
func WebSubmissionsCommand() *ffcli.Command {
	fs := flag.NewFlagSet("web submissions", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "submissions",
		ShortUsage: "asc web submissions <subcommand> [flags]",
		ShortHelp:  "EXPERIMENTAL: Review submissions via unofficial web APIs.",
		LongHelp: `EXPERIMENTAL / UNOFFICIAL / DISCOURAGED

Inspect App Review submission lifecycle with Apple web-session endpoints.

` + webWarningText,
		FlagSet:   fs,
		UsageFunc: shared.DefaultUsageFunc,
		Subcommands: []*ffcli.Command{
			WebSubmissionsListCommand(),
			WebSubmissionsShowCommand(),
			WebSubmissionsItemsCommand(),
		},
		Exec: func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		},
	}
}

// WebSubmissionsListCommand lists review submissions for an app.
func WebSubmissionsListCommand() *ffcli.Command {
	fs := flag.NewFlagSet("web submissions list", flag.ExitOnError)

	appID := fs.String("app", "", "App ID")
	stateCSV := fs.String("state", "", "Optional comma-separated state filter")
	authFlags := bindWebSessionFlags(fs)
	output := shared.BindOutputFlags(fs)

	return &ffcli.Command{
		Name:       "list",
		ShortUsage: "asc web submissions list --app APP_ID [--state CSV] [flags]",
		ShortHelp:  "EXPERIMENTAL: List review submissions for an app.",
		FlagSet:    fs,
		UsageFunc:  shared.DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			trimmedAppID := strings.TrimSpace(*appID)
			if trimmedAppID == "" {
				return shared.UsageError("--app is required")
			}
			states, err := parseSubmissionStates(*stateCSV)
			if err != nil {
				return err
			}

			requestCtx, cancel := shared.ContextWithTimeout(ctx)
			defer cancel()

			session, err := resolveWebSessionForCommand(requestCtx, authFlags)
			if err != nil {
				return err
			}
			client := webcore.NewClient(session)

			submissions, err := client.ListReviewSubmissions(requestCtx, trimmedAppID)
			if err != nil {
				return withWebAuthHint(err, "web submissions list")
			}
			filtered := filterSubmissionsByState(submissions, states)
			return shared.PrintOutput(filtered, *output.Output, *output.Pretty)
		},
	}
}

// WebSubmissionsShowCommand shows one submission by ID.
func WebSubmissionsShowCommand() *ffcli.Command {
	fs := flag.NewFlagSet("web submissions show", flag.ExitOnError)

	submissionID := fs.String("id", "", "Review submission ID")
	authFlags := bindWebSessionFlags(fs)
	output := shared.BindOutputFlags(fs)

	return &ffcli.Command{
		Name:       "show",
		ShortUsage: "asc web submissions show --id REVIEW_SUBMISSION_ID [flags]",
		ShortHelp:  "EXPERIMENTAL: Show one review submission.",
		FlagSet:    fs,
		UsageFunc:  shared.DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			trimmedID := strings.TrimSpace(*submissionID)
			if trimmedID == "" {
				return shared.UsageError("--id is required")
			}

			requestCtx, cancel := shared.ContextWithTimeout(ctx)
			defer cancel()

			session, err := resolveWebSessionForCommand(requestCtx, authFlags)
			if err != nil {
				return err
			}
			client := webcore.NewClient(session)

			submission, err := client.GetReviewSubmission(requestCtx, trimmedID)
			if err != nil {
				return withWebAuthHint(err, "web submissions show")
			}
			if submission == nil {
				return shared.PrintOutput(map[string]any{}, *output.Output, *output.Pretty)
			}
			return shared.PrintOutput(submission, *output.Output, *output.Pretty)
		},
	}
}

// WebSubmissionsItemsCommand lists items for one submission.
func WebSubmissionsItemsCommand() *ffcli.Command {
	fs := flag.NewFlagSet("web submissions items", flag.ExitOnError)

	submissionID := fs.String("id", "", "Review submission ID")
	authFlags := bindWebSessionFlags(fs)
	output := shared.BindOutputFlags(fs)

	return &ffcli.Command{
		Name:       "items",
		ShortUsage: "asc web submissions items --id REVIEW_SUBMISSION_ID [flags]",
		ShortHelp:  "EXPERIMENTAL: List review submission items.",
		FlagSet:    fs,
		UsageFunc:  shared.DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			trimmedID := strings.TrimSpace(*submissionID)
			if trimmedID == "" {
				return shared.UsageError("--id is required")
			}

			requestCtx, cancel := shared.ContextWithTimeout(ctx)
			defer cancel()

			session, err := resolveWebSessionForCommand(requestCtx, authFlags)
			if err != nil {
				return err
			}
			client := webcore.NewClient(session)

			items, err := client.ListReviewSubmissionItems(requestCtx, trimmedID)
			if err != nil {
				return withWebAuthHint(err, "web submissions items")
			}
			return shared.PrintOutput(items, *output.Output, *output.Pretty)
		},
	}
}

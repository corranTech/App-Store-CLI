package registry

import (
	"context"
	"fmt"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/completion"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
)

// VersionCommand returns a version subcommand.
func VersionCommand(version string) *ffcli.Command {
	return &ffcli.Command{
		Name:       "version",
		ShortUsage: "asc version",
		ShortHelp:  "Print version information and exit.",
		UsageFunc:  shared.DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			fmt.Println(version)
			return nil
		},
	}
}

// Subcommands returns all root subcommands in display order.
// Commands are organized by domain for maintainability.
func Subcommands(version string) []*ffcli.Command {
	var subs []*ffcli.Command

	// Core commands (auth, init, install, docs)
	subs = append(subs, CoreCommands()...)

	// Feedback and reviews
	subs = append(subs, FeedbackCommands()...)
	subs = append(subs, ReviewCommands()...)

	// Analytics and reporting
	subs = append(subs, AnalyticsCommands()...)

	// App management
	subs = append(subs, AppCommands()...)

	// Version and metadata
	subs = append(subs, VersionCommands()...)

	// Build management
	subs = append(subs, BuildCommands()...)

	// TestFlight
	subs = append(subs, TestFlightCommands()...)

	// Assets and localizations
	subs = append(subs, AssetCommands()...)

	// Sandbox
	subs = append(subs, SandboxCommands()...)

	// Signing
	subs = append(subs, SigningCommands()...)

	// In-app purchases
	subs = append(subs, IAPCommands()...)

	// Users and devices
	subs = append(subs, UserCommands()...)

	// Bundle IDs
	subs = append(subs, BundleIDCommands()...)

	// Webhooks and notifications
	subs = append(subs, WebhookCommands()...)

	// Submission
	subs = append(subs, SubmissionCommands()...)

	// Metadata
	subs = append(subs, MetadataCommands()...)

	// App Store metadata
	subs = append(subs, AppStoreCommands()...)

	// Game Center
	subs = append(subs, GameCenterCommands()...)

	// Version command
	subs = append(subs, VersionCommand(version))

	// Shell completion (must be last)
	subs = append(subs, completion.CompletionCommand(subs))
	return subs
}

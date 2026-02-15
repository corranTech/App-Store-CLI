package validate

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/validation"
)

type validateSubscriptionsOptions struct {
	AppID  string
	Strict bool
	Output string
	Pretty bool
}

// ValidateSubscriptionsCommand returns the asc validate subscriptions subcommand.
func ValidateSubscriptionsCommand() *ffcli.Command {
	fs := flag.NewFlagSet("subscriptions", flag.ExitOnError)

	appID := fs.String("app", "", "App Store Connect app ID (or ASC_APP_ID)")
	strict := fs.Bool("strict", false, "Treat warnings as errors (exit non-zero)")
	output := shared.BindOutputFlags(fs)

	return &ffcli.Command{
		Name:       "subscriptions",
		ShortUsage: "asc validate subscriptions --app \"APP_ID\" [flags]",
		ShortHelp:  "Validate subscription review readiness (warning-only by default).",
		LongHelp: `Validate review readiness for auto-renewable subscriptions.

This command is conservative: it emits warnings for subscriptions that look
unsubmitted or need action, but it does not block by default (use --strict for CI).

Examples:
  asc validate subscriptions --app "APP_ID"
  asc validate subscriptions --app "APP_ID" --output table
  asc validate subscriptions --app "APP_ID" --strict`,
		FlagSet:   fs,
		UsageFunc: shared.DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			resolvedAppID := shared.ResolveAppID(*appID)
			if resolvedAppID == "" {
				fmt.Fprintln(os.Stderr, "Error: --app is required (or set ASC_APP_ID)")
				return flag.ErrHelp
			}

			return runValidateSubscriptions(ctx, validateSubscriptionsOptions{
				AppID:  resolvedAppID,
				Strict: *strict,
				Output: *output.Output,
				Pretty: *output.Pretty,
			})
		},
	}
}

func runValidateSubscriptions(ctx context.Context, opts validateSubscriptionsOptions) error {
	client, err := clientFactory()
	if err != nil {
		return fmt.Errorf("validate subscriptions: %w", err)
	}

	requestCtx, cancel := shared.ContextWithTimeout(ctx)
	defer cancel()

	const pageLimit = 200

	nextGroupsURL := ""
	groupIDs := make([]string, 0)
	for {
		var groupsResp *asc.SubscriptionGroupsResponse
		if strings.TrimSpace(nextGroupsURL) != "" {
			groupsResp, err = client.GetSubscriptionGroups(requestCtx, opts.AppID, asc.WithSubscriptionGroupsNextURL(nextGroupsURL))
		} else {
			groupsResp, err = client.GetSubscriptionGroups(requestCtx, opts.AppID, asc.WithSubscriptionGroupsLimit(pageLimit))
		}
		if err != nil {
			return fmt.Errorf("validate subscriptions: failed to fetch subscription groups: %w", err)
		}

		for _, group := range groupsResp.Data {
			if strings.TrimSpace(group.ID) == "" {
				continue
			}
			groupIDs = append(groupIDs, group.ID)
		}

		nextGroupsURL = strings.TrimSpace(groupsResp.Links.Next)
		if nextGroupsURL == "" {
			break
		}
	}

	subs := make([]validation.Subscription, 0)
	for _, groupID := range groupIDs {
		nextSubsURL := ""
		for {
			var subsResp *asc.SubscriptionsResponse
			if strings.TrimSpace(nextSubsURL) != "" {
				subsResp, err = client.GetSubscriptions(requestCtx, groupID, asc.WithSubscriptionsNextURL(nextSubsURL))
			} else {
				subsResp, err = client.GetSubscriptions(requestCtx, groupID, asc.WithSubscriptionsLimit(pageLimit))
			}
			if err != nil {
				return fmt.Errorf("validate subscriptions: failed to fetch subscriptions for group %s: %w", groupID, err)
			}

			for _, sub := range subsResp.Data {
				attrs := sub.Attributes
				subs = append(subs, validation.Subscription{
					ID:        sub.ID,
					Name:      attrs.Name,
					ProductID: attrs.ProductID,
					State:     attrs.State,
					GroupID:   groupID,
				})
			}

			nextSubsURL = strings.TrimSpace(subsResp.Links.Next)
			if nextSubsURL == "" {
				break
			}
		}
	}

	report := validation.ValidateSubscriptions(validation.SubscriptionsInput{
		AppID:         opts.AppID,
		Subscriptions: subs,
	}, opts.Strict)

	if err := shared.PrintOutput(&report, opts.Output, opts.Pretty); err != nil {
		return err
	}

	if report.Summary.Blocking > 0 {
		return shared.NewReportedError(fmt.Errorf("validate subscriptions: found %d blocking issue(s)", report.Summary.Blocking))
	}

	return nil
}

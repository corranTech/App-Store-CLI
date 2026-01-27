package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
)

// AppTagsCommand returns the app-tags command with subcommands.
func AppTagsCommand() *ffcli.Command {
	fs := flag.NewFlagSet("app-tags", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "app-tags",
		ShortUsage: "asc app-tags <subcommand> [flags]",
		ShortHelp:  "Manage app tags for App Store visibility.",
		LongHelp: `Manage app tags for App Store visibility.

Examples:
  asc app-tags list --app "APP_ID"
  asc app-tags get --app "APP_ID" --id "TAG_ID"
  asc app-tags update --id "TAG_ID" --visible-in-app-store=false --confirm`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Subcommands: []*ffcli.Command{
			AppTagsListCommand(),
			AppTagsGetCommand(),
			AppTagsUpdateCommand(),
		},
		Exec: func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		},
	}
}

// AppTagsListCommand returns the list subcommand.
func AppTagsListCommand() *ffcli.Command {
	fs := flag.NewFlagSet("app-tags list", flag.ExitOnError)

	appID := fs.String("app", "", "App Store Connect app ID (or ASC_APP_ID env)")
	visible := fs.String("visible-in-app-store", "", "Filter by visibility (true/false), comma-separated")
	sort := fs.String("sort", "", "Sort by name or -name")
	limit := fs.Int("limit", 0, "Maximum results per page (1-200)")
	next := fs.String("next", "", "Fetch next page using a links.next URL")
	paginate := fs.Bool("paginate", false, "Automatically fetch all pages (aggregate results)")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "list",
		ShortUsage: "asc app-tags list [flags]",
		ShortHelp:  "List app tags for an app.",
		LongHelp: `List app tags for an app.

Examples:
  asc app-tags list --app "APP_ID"
  asc app-tags list --app "APP_ID" --visible-in-app-store true
  asc app-tags list --app "APP_ID" --sort -name --limit 10
  asc app-tags list --next "<links.next>"
  asc app-tags list --app "APP_ID" --paginate`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			if *limit != 0 && (*limit < 1 || *limit > 200) {
				return fmt.Errorf("app-tags list: --limit must be between 1 and 200")
			}
			if err := validateNextURL(*next); err != nil {
				return fmt.Errorf("app-tags list: %w", err)
			}
			if err := validateSort(*sort, "name", "-name"); err != nil {
				return fmt.Errorf("app-tags list: %w", err)
			}

			visibleValues, err := normalizeAppTagVisibilityFilter(*visible)
			if err != nil {
				return fmt.Errorf("app-tags list: %w", err)
			}

			resolvedAppID := resolveAppID(*appID)
			if resolvedAppID == "" && strings.TrimSpace(*next) == "" {
				fmt.Fprintf(os.Stderr, "Error: --app is required (or set ASC_APP_ID)\n\n")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("app-tags list: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			opts := []asc.AppTagsOption{
				asc.WithAppTagsVisibleInAppStore(visibleValues),
				asc.WithAppTagsLimit(*limit),
				asc.WithAppTagsNextURL(*next),
			}
			if strings.TrimSpace(*sort) != "" {
				opts = append(opts, asc.WithAppTagsSort(*sort))
			}

			if *paginate {
				paginateOpts := append(opts, asc.WithAppTagsLimit(200))
				firstPage, err := client.GetAppTags(requestCtx, resolvedAppID, paginateOpts...)
				if err != nil {
					return fmt.Errorf("app-tags list: failed to fetch: %w", err)
				}

				paginated, err := asc.PaginateAll(requestCtx, firstPage, func(ctx context.Context, nextURL string) (asc.PaginatedResponse, error) {
					return client.GetAppTags(ctx, resolvedAppID, asc.WithAppTagsNextURL(nextURL))
				})
				if err != nil {
					return fmt.Errorf("app-tags list: %w", err)
				}

				return printOutput(paginated, *output, *pretty)
			}

			resp, err := client.GetAppTags(requestCtx, resolvedAppID, opts...)
			if err != nil {
				return fmt.Errorf("app-tags list: failed to fetch: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// AppTagsGetCommand returns the get subcommand.
func AppTagsGetCommand() *ffcli.Command {
	fs := flag.NewFlagSet("app-tags get", flag.ExitOnError)

	appID := fs.String("app", "", "App Store Connect app ID (or ASC_APP_ID env)")
	tagID := fs.String("id", "", "App tag ID")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "get",
		ShortUsage: "asc app-tags get [flags]",
		ShortHelp:  "Get an app tag by ID.",
		LongHelp: `Get an app tag by ID.

This command searches the app's tags for the specified ID.

Examples:
  asc app-tags get --app "APP_ID" --id "TAG_ID"`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			trimmedID := strings.TrimSpace(*tagID)
			if trimmedID == "" {
				fmt.Fprintln(os.Stderr, "Error: --id is required")
				return flag.ErrHelp
			}

			resolvedAppID := resolveAppID(*appID)
			if resolvedAppID == "" {
				fmt.Fprintln(os.Stderr, "Error: --app is required (or set ASC_APP_ID)")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("app-tags get: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			resp, err := findAppTagByID(requestCtx, client, resolvedAppID, trimmedID)
			if err != nil {
				return fmt.Errorf("app-tags get: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// AppTagsUpdateCommand returns the update subcommand.
func AppTagsUpdateCommand() *ffcli.Command {
	fs := flag.NewFlagSet("app-tags update", flag.ExitOnError)

	tagID := fs.String("id", "", "App tag ID")
	visibleInAppStore := fs.Bool("visible-in-app-store", false, "Set visibility in the App Store")
	confirm := fs.Bool("confirm", false, "Confirm update")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "update",
		ShortUsage: "asc app-tags update --id TAG_ID --visible-in-app-store [true|false] --confirm",
		ShortHelp:  "Update an app tag.",
		LongHelp: `Update an app tag.

Examples:
  asc app-tags update --id "TAG_ID" --visible-in-app-store --confirm
  asc app-tags update --id "TAG_ID" --visible-in-app-store=false --confirm`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			trimmedID := strings.TrimSpace(*tagID)
			if trimmedID == "" {
				fmt.Fprintln(os.Stderr, "Error: --id is required")
				return flag.ErrHelp
			}

			visited := map[string]bool{}
			fs.Visit(func(f *flag.Flag) {
				visited[f.Name] = true
			})
			if !visited["visible-in-app-store"] {
				fmt.Fprintln(os.Stderr, "Error: --visible-in-app-store is required")
				return flag.ErrHelp
			}
			if !*confirm {
				fmt.Fprintln(os.Stderr, "Error: --confirm is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("app-tags update: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			attrs := asc.AppTagUpdateAttributes{}
			if visited["visible-in-app-store"] {
				attrs.VisibleInAppStore = visibleInAppStore
			}

			resp, err := client.UpdateAppTag(requestCtx, trimmedID, attrs)
			if err != nil {
				return fmt.Errorf("app-tags update: failed to update: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

func normalizeAppTagVisibilityFilter(value string) ([]string, error) {
	values := splitCSV(value)
	if len(values) == 0 {
		return nil, nil
	}

	normalized := make([]string, 0, len(values))
	for _, item := range values {
		lower := strings.ToLower(strings.TrimSpace(item))
		switch lower {
		case "true", "false":
			normalized = append(normalized, lower)
		default:
			return nil, fmt.Errorf("--visible-in-app-store must be true or false")
		}
	}

	return normalized, nil
}

func findAppTagByID(ctx context.Context, client *asc.Client, appID, tagID string) (*asc.AppTagResponse, error) {
	resp, err := client.GetAppTags(ctx, appID, asc.WithAppTagsLimit(200))
	if err != nil {
		return nil, err
	}

	for {
		for _, item := range resp.Data {
			if item.ID == tagID {
				return &asc.AppTagResponse{Data: item}, nil
			}
		}

		if resp.Links.Next == "" {
			break
		}

		resp, err = client.GetAppTags(ctx, appID, asc.WithAppTagsNextURL(resp.Links.Next))
		if err != nil {
			return nil, err
		}
	}

	return nil, fmt.Errorf("tag %q not found", tagID)
}

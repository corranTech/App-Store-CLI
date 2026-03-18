package builds

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

// BuildsCountCommand returns the builds count subcommand.
func BuildsCountCommand() *ffcli.Command {
	fs := flag.NewFlagSet("count", flag.ExitOnError)

	appID := fs.String("app", "", "App Store Connect app ID, bundle ID, or exact app name (required, or ASC_APP_ID env)")
	version := fs.String("version", "", "Filter by marketing version string (CFBundleShortVersionString)")
	buildNumber := fs.String("build-number", "", "Filter by build number (CFBundleVersion)")
	platform := fs.String("platform", "", "Filter by platform: IOS, MAC_OS, TV_OS, VISION_OS")
	processingState := fs.String("processing-state", "", "Filter by processing state: VALID, PROCESSING, FAILED, INVALID, or all")
	output := shared.BindOutputFlags(fs)

	return &ffcli.Command{
		Name:       "count",
		ShortUsage: "asc builds count [flags]",
		ShortHelp:  "Count total builds for an app in App Store Connect.",
		LongHelp: `Count the total number of builds for an app in App Store Connect.

Uses meta.paging.total from the API response — a single request with limit=1
is sufficient, so this is fast regardless of how many builds exist.

Supports the same filters as "asc builds list" so you can count builds
matching specific criteria without fetching them all.

Examples:
  asc builds count --app "123456789"
  asc builds count --app "123456789" --version "2.1.0"
  asc builds count --app "123456789" --build-number "42"
  asc builds count --app "123456789" --platform IOS
  asc builds count --app "123456789" --processing-state VALID
  asc builds count --app "123456789" --processing-state all
  asc builds count --app "123456789" --version "2.1.0" --platform IOS --output json`,
		FlagSet:   fs,
		UsageFunc: shared.DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			resolvedAppID := shared.ResolveAppID(*appID)
			if resolvedAppID == "" {
				fmt.Fprintf(os.Stderr, "Error: --app is required (or set ASC_APP_ID)\n\n")
				return flag.ErrHelp
			}

			platformValue := ""
			if strings.TrimSpace(*platform) != "" {
				normalizedPlatform, err := shared.NormalizePlatform(*platform)
				if err != nil {
					return shared.UsageError(err.Error())
				}
				platformValue = string(normalizedPlatform)
			}

			buildNumberValue := strings.TrimSpace(*buildNumber)
			processingStateValues, err := normalizeBuildProcessingStateFilter(*processingState)
			if err != nil {
				return err
			}

			client, err := shared.GetASCClient()
			if err != nil {
				return fmt.Errorf("builds count: %w", err)
			}

			requestCtx, cancel := shared.ContextWithTimeout(ctx)
			defer cancel()

			resolvedAppID, err = shared.ResolveAppIDWithLookup(requestCtx, client, resolvedAppID)
			if err != nil {
				return fmt.Errorf("builds count: %w", err)
			}

			// Resolve pre-release version IDs when filtering by marketing version,
			// matching the same lookup used in builds list.
			var preReleaseVersionIDs []string
			versionValue := strings.TrimSpace(*version)
			if versionValue != "" {
				preReleaseVersionIDs, err = findPreReleaseVersionIDsForBuildsList(requestCtx, client, resolvedAppID, versionValue)
				if err != nil {
					return fmt.Errorf("builds count: %w", err)
				}
				if len(preReleaseVersionIDs) == 0 {
					result := &asc.BuildsCountResult{AppID: resolvedAppID, Total: 0}
					return shared.PrintOutput(result, *output.Output, *output.Pretty)
				}
			}

			// Request only 1 item — we only need meta.paging.total, not the builds themselves.
			opts := []asc.BuildsOption{
				asc.WithBuildsLimit(1),
			}
			if buildNumberValue != "" {
				opts = append(opts, asc.WithBuildsBuildNumber(buildNumberValue))
			}
			if platformValue != "" {
				opts = append(opts, asc.WithBuildsPreReleaseVersionPlatforms([]string{platformValue}))
			}
			if len(processingStateValues) > 0 {
				opts = append(opts, asc.WithBuildsProcessingStates(processingStateValues))
			}
			if len(preReleaseVersionIDs) > 0 {
				opts = append(opts, asc.WithBuildsPreReleaseVersions(preReleaseVersionIDs))
			}

			resp, err := client.GetBuilds(requestCtx, resolvedAppID, opts...)
			if err != nil {
				return fmt.Errorf("builds count: failed to fetch: %w", err)
			}

			total, ok := asc.ParsePagingTotalOK(resp.Meta)
			if !ok {
				return fmt.Errorf("builds count: App Store Connect did not return a total count in this response; try 'asc builds list --paginate' to count manually")
			}

			result := &asc.BuildsCountResult{AppID: resolvedAppID, Total: total}
			return shared.PrintOutput(result, *output.Output, *output.Pretty)
		},
	}
}

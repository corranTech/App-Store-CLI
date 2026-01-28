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

// GameCenterCommand returns the game-center command group.
func GameCenterCommand() *ffcli.Command {
	fs := flag.NewFlagSet("game-center", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "game-center",
		ShortUsage: "asc game-center <subcommand> [flags]",
		ShortHelp:  "Manage Game Center resources in App Store Connect.",
		LongHelp: `Manage Game Center resources in App Store Connect.

Examples:
  asc game-center achievements list --app "APP_ID"
  asc game-center achievements create --app "APP_ID" --reference-name "First Win" --vendor-id "com.example.firstwin" --points 10
  asc game-center leaderboards list --app "APP_ID"
  asc game-center leaderboards create --app "APP_ID" --reference-name "High Score" --vendor-id "com.example.highscore" --formatter INTEGER --sort DESC --submission-type BEST_SCORE
  asc game-center leaderboard-sets list --app "APP_ID"
  asc game-center leaderboard-sets create --app "APP_ID" --reference-name "Season 1" --vendor-id "com.example.season1"`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Subcommands: []*ffcli.Command{
			GameCenterAchievementsCommand(),
			GameCenterLeaderboardsCommand(),
			GameCenterLeaderboardSetsCommand(),
		},
		Exec: func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		},
	}
}

// GameCenterAchievementsCommand returns the achievements command group.
func GameCenterAchievementsCommand() *ffcli.Command {
	fs := flag.NewFlagSet("achievements", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "achievements",
		ShortUsage: "asc game-center achievements <subcommand> [flags]",
		ShortHelp:  "Manage Game Center achievements.",
		LongHelp: `Manage Game Center achievements.

Examples:
  asc game-center achievements list --app "APP_ID"
  asc game-center achievements get --id "ACHIEVEMENT_ID"
  asc game-center achievements create --app "APP_ID" --reference-name "First Win" --vendor-id "com.example.firstwin" --points 10
  asc game-center achievements update --id "ACHIEVEMENT_ID" --points 20
  asc game-center achievements delete --id "ACHIEVEMENT_ID" --confirm
  asc game-center achievements localizations list --achievement-id "ACHIEVEMENT_ID"
  asc game-center achievements localizations create --achievement-id "ACHIEVEMENT_ID" --locale en-US --name "First Win"
  asc game-center achievements localizations update --id "LOC_ID" --name "New Name"
  asc game-center achievements localizations delete --id "LOC_ID" --confirm
  asc game-center achievements images upload --localization-id "LOC_ID" --file "path/to/image.png"
  asc game-center achievements images delete --id "IMAGE_ID" --confirm`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Subcommands: []*ffcli.Command{
			GameCenterAchievementsListCommand(),
			GameCenterAchievementsGetCommand(),
			GameCenterAchievementsCreateCommand(),
			GameCenterAchievementsUpdateCommand(),
			GameCenterAchievementsDeleteCommand(),
			GameCenterAchievementLocalizationsCommand(),
			GameCenterAchievementImagesCommand(),
			GameCenterAchievementReleasesCommand(),
		},
		Exec: func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		},
	}
}

// GameCenterAchievementsListCommand returns the achievements list subcommand.
func GameCenterAchievementsListCommand() *ffcli.Command {
	fs := flag.NewFlagSet("list", flag.ExitOnError)

	appID := fs.String("app", "", "App Store Connect app ID (or ASC_APP_ID env)")
	limit := fs.Int("limit", 0, "Maximum results per page (1-200)")
	next := fs.String("next", "", "Fetch next page using a links.next URL")
	paginate := fs.Bool("paginate", false, "Automatically fetch all pages (aggregate results)")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "list",
		ShortUsage: "asc game-center achievements list [flags]",
		ShortHelp:  "List Game Center achievements for an app.",
		LongHelp: `List Game Center achievements for an app.

Examples:
  asc game-center achievements list --app "APP_ID"
  asc game-center achievements list --app "APP_ID" --limit 50
  asc game-center achievements list --app "APP_ID" --paginate`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			if *limit != 0 && (*limit < 1 || *limit > 200) {
				return fmt.Errorf("game-center achievements list: --limit must be between 1 and 200")
			}
			if err := validateNextURL(*next); err != nil {
				return fmt.Errorf("game-center achievements list: %w", err)
			}

			resolvedAppID := resolveAppID(*appID)
			if resolvedAppID == "" && strings.TrimSpace(*next) == "" {
				fmt.Fprintln(os.Stderr, "Error: --app is required (or set ASC_APP_ID)")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center achievements list: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			// Get Game Center detail ID first
			gcDetailID, err := client.GetGameCenterDetailID(requestCtx, resolvedAppID)
			if err != nil {
				return fmt.Errorf("game-center achievements list: failed to get Game Center detail: %w", err)
			}

			opts := []asc.GCAchievementsOption{
				asc.WithGCAchievementsLimit(*limit),
				asc.WithGCAchievementsNextURL(*next),
			}

			if *paginate {
				paginateOpts := append(opts, asc.WithGCAchievementsLimit(200))
				firstPage, err := client.GetGameCenterAchievements(requestCtx, gcDetailID, paginateOpts...)
				if err != nil {
					return fmt.Errorf("game-center achievements list: failed to fetch: %w", err)
				}

				resp, err := asc.PaginateAll(requestCtx, firstPage, func(ctx context.Context, nextURL string) (asc.PaginatedResponse, error) {
					return client.GetGameCenterAchievements(ctx, gcDetailID, asc.WithGCAchievementsNextURL(nextURL))
				})
				if err != nil {
					return fmt.Errorf("game-center achievements list: %w", err)
				}

				return printOutput(resp, *output, *pretty)
			}

			resp, err := client.GetGameCenterAchievements(requestCtx, gcDetailID, opts...)
			if err != nil {
				return fmt.Errorf("game-center achievements list: failed to fetch: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterAchievementsGetCommand returns the achievements get subcommand.
func GameCenterAchievementsGetCommand() *ffcli.Command {
	fs := flag.NewFlagSet("get", flag.ExitOnError)

	achievementID := fs.String("id", "", "Game Center achievement ID")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "get",
		ShortUsage: "asc game-center achievements get --id \"ACHIEVEMENT_ID\"",
		ShortHelp:  "Get a Game Center achievement by ID.",
		LongHelp: `Get a Game Center achievement by ID.

Examples:
  asc game-center achievements get --id "ACHIEVEMENT_ID"`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			id := strings.TrimSpace(*achievementID)
			if id == "" {
				fmt.Fprintln(os.Stderr, "Error: --id is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center achievements get: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			resp, err := client.GetGameCenterAchievement(requestCtx, id)
			if err != nil {
				return fmt.Errorf("game-center achievements get: failed to fetch: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterAchievementsCreateCommand returns the achievements create subcommand.
func GameCenterAchievementsCreateCommand() *ffcli.Command {
	fs := flag.NewFlagSet("create", flag.ExitOnError)

	appID := fs.String("app", "", "App Store Connect app ID (or ASC_APP_ID env)")
	referenceName := fs.String("reference-name", "", "Reference name for the achievement")
	vendorID := fs.String("vendor-id", "", "Vendor identifier (e.g., com.example.achievement)")
	points := fs.Int("points", 0, "Points value (1-100)")
	showBeforeEarned := fs.Bool("show-before-earned", true, "Show achievement before it is earned")
	repeatable := fs.Bool("repeatable", false, "Achievement can be earned multiple times")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "create",
		ShortUsage: "asc game-center achievements create [flags]",
		ShortHelp:  "Create a new Game Center achievement.",
		LongHelp: `Create a new Game Center achievement.

Examples:
  asc game-center achievements create --app "APP_ID" --reference-name "First Win" --vendor-id "com.example.firstwin" --points 10
  asc game-center achievements create --app "APP_ID" --reference-name "Master" --vendor-id "com.example.master" --points 100 --repeatable`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			resolvedAppID := resolveAppID(*appID)
			if resolvedAppID == "" {
				fmt.Fprintln(os.Stderr, "Error: --app is required (or set ASC_APP_ID)")
				return flag.ErrHelp
			}

			name := strings.TrimSpace(*referenceName)
			if name == "" {
				fmt.Fprintln(os.Stderr, "Error: --reference-name is required")
				return flag.ErrHelp
			}

			vendor := strings.TrimSpace(*vendorID)
			if vendor == "" {
				fmt.Fprintln(os.Stderr, "Error: --vendor-id is required")
				return flag.ErrHelp
			}

			if *points < 1 || *points > 100 {
				fmt.Fprintln(os.Stderr, "Error: --points must be between 1 and 100")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center achievements create: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			// Get Game Center detail ID first
			gcDetailID, err := client.GetGameCenterDetailID(requestCtx, resolvedAppID)
			if err != nil {
				return fmt.Errorf("game-center achievements create: failed to get Game Center detail: %w", err)
			}

			attrs := asc.GameCenterAchievementCreateAttributes{
				ReferenceName:    name,
				VendorIdentifier: vendor,
				Points:           *points,
				ShowBeforeEarned: *showBeforeEarned,
				Repeatable:       *repeatable,
			}

			resp, err := client.CreateGameCenterAchievement(requestCtx, gcDetailID, attrs)
			if err != nil {
				return fmt.Errorf("game-center achievements create: failed to create: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterAchievementsUpdateCommand returns the achievements update subcommand.
func GameCenterAchievementsUpdateCommand() *ffcli.Command {
	fs := flag.NewFlagSet("update", flag.ExitOnError)

	achievementID := fs.String("id", "", "Game Center achievement ID")
	referenceName := fs.String("reference-name", "", "Reference name for the achievement")
	points := fs.Int("points", 0, "Points value (1-100)")
	showBeforeEarned := fs.String("show-before-earned", "", "Show achievement before it is earned (true/false)")
	repeatable := fs.String("repeatable", "", "Achievement can be earned multiple times (true/false)")
	archived := fs.String("archived", "", "Archive the achievement (true/false)")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "update",
		ShortUsage: "asc game-center achievements update [flags]",
		ShortHelp:  "Update a Game Center achievement.",
		LongHelp: `Update a Game Center achievement.

Examples:
  asc game-center achievements update --id "ACHIEVEMENT_ID" --reference-name "New Name"
  asc game-center achievements update --id "ACHIEVEMENT_ID" --points 20
  asc game-center achievements update --id "ACHIEVEMENT_ID" --archived true`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			id := strings.TrimSpace(*achievementID)
			if id == "" {
				fmt.Fprintln(os.Stderr, "Error: --id is required")
				return flag.ErrHelp
			}

			attrs := asc.GameCenterAchievementUpdateAttributes{}
			hasUpdate := false

			if strings.TrimSpace(*referenceName) != "" {
				name := strings.TrimSpace(*referenceName)
				attrs.ReferenceName = &name
				hasUpdate = true
			}

			if *points != 0 {
				if *points < 1 || *points > 100 {
					fmt.Fprintln(os.Stderr, "Error: --points must be between 1 and 100")
					return flag.ErrHelp
				}
				attrs.Points = points
				hasUpdate = true
			}

			if strings.TrimSpace(*showBeforeEarned) != "" {
				val, err := parseBool(*showBeforeEarned, "--show-before-earned")
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error:", err.Error())
					return flag.ErrHelp
				}
				attrs.ShowBeforeEarned = &val
				hasUpdate = true
			}

			if strings.TrimSpace(*repeatable) != "" {
				val, err := parseBool(*repeatable, "--repeatable")
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error:", err.Error())
					return flag.ErrHelp
				}
				attrs.Repeatable = &val
				hasUpdate = true
			}

			if strings.TrimSpace(*archived) != "" {
				val, err := parseBool(*archived, "--archived")
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error:", err.Error())
					return flag.ErrHelp
				}
				attrs.Archived = &val
				hasUpdate = true
			}

			if !hasUpdate {
				fmt.Fprintln(os.Stderr, "Error: at least one update flag is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center achievements update: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			resp, err := client.UpdateGameCenterAchievement(requestCtx, id, attrs)
			if err != nil {
				return fmt.Errorf("game-center achievements update: failed to update: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterAchievementsDeleteCommand returns the achievements delete subcommand.
func GameCenterAchievementsDeleteCommand() *ffcli.Command {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)

	achievementID := fs.String("id", "", "Game Center achievement ID")
	confirm := fs.Bool("confirm", false, "Confirm deletion")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "delete",
		ShortUsage: "asc game-center achievements delete --id \"ACHIEVEMENT_ID\" --confirm",
		ShortHelp:  "Delete a Game Center achievement.",
		LongHelp: `Delete a Game Center achievement.

Examples:
  asc game-center achievements delete --id "ACHIEVEMENT_ID" --confirm`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			id := strings.TrimSpace(*achievementID)
			if id == "" {
				fmt.Fprintln(os.Stderr, "Error: --id is required")
				return flag.ErrHelp
			}
			if !*confirm {
				fmt.Fprintln(os.Stderr, "Error: --confirm is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center achievements delete: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			if err := client.DeleteGameCenterAchievement(requestCtx, id); err != nil {
				return fmt.Errorf("game-center achievements delete: failed to delete: %w", err)
			}

			result := &asc.GameCenterAchievementDeleteResult{
				ID:      id,
				Deleted: true,
			}

			return printOutput(result, *output, *pretty)
		},
	}
}

func parseBool(value, flagName string) (bool, error) {
	v := strings.ToLower(strings.TrimSpace(value))
	switch v {
	case "true", "1", "yes":
		return true, nil
	case "false", "0", "no":
		return false, nil
	default:
		return false, fmt.Errorf("%s must be true or false", flagName)
	}
}

// GameCenterAchievementLocalizationsCommand returns the achievement localizations command group.
func GameCenterAchievementLocalizationsCommand() *ffcli.Command {
	fs := flag.NewFlagSet("localizations", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "localizations",
		ShortUsage: "asc game-center achievements localizations <subcommand> [flags]",
		ShortHelp:  "Manage Game Center achievement localizations.",
		LongHelp: `Manage Game Center achievement localizations.

Examples:
  asc game-center achievements localizations list --achievement-id "ACHIEVEMENT_ID"
  asc game-center achievements localizations get --id "LOC_ID"
  asc game-center achievements localizations create --achievement-id "ACHIEVEMENT_ID" --locale en-US --name "First Win" --before-earned-description "Win your first game" --after-earned-description "You won!"
  asc game-center achievements localizations update --id "LOC_ID" --name "New Name"
  asc game-center achievements localizations delete --id "LOC_ID" --confirm`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Subcommands: []*ffcli.Command{
			GameCenterAchievementLocalizationsListCommand(),
			GameCenterAchievementLocalizationsGetCommand(),
			GameCenterAchievementLocalizationsCreateCommand(),
			GameCenterAchievementLocalizationsUpdateCommand(),
			GameCenterAchievementLocalizationsDeleteCommand(),
		},
		Exec: func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		},
	}
}

// GameCenterAchievementLocalizationsListCommand returns the localizations list subcommand.
func GameCenterAchievementLocalizationsListCommand() *ffcli.Command {
	fs := flag.NewFlagSet("list", flag.ExitOnError)

	achievementID := fs.String("achievement-id", "", "Game Center achievement ID")
	limit := fs.Int("limit", 0, "Maximum results per page (1-200)")
	next := fs.String("next", "", "Fetch next page using a links.next URL")
	paginate := fs.Bool("paginate", false, "Automatically fetch all pages (aggregate results)")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "list",
		ShortUsage: "asc game-center achievements localizations list --achievement-id \"ACHIEVEMENT_ID\"",
		ShortHelp:  "List localizations for a Game Center achievement.",
		LongHelp: `List localizations for a Game Center achievement.

Examples:
  asc game-center achievements localizations list --achievement-id "ACHIEVEMENT_ID"
  asc game-center achievements localizations list --achievement-id "ACHIEVEMENT_ID" --limit 50
  asc game-center achievements localizations list --achievement-id "ACHIEVEMENT_ID" --paginate`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			if *limit != 0 && (*limit < 1 || *limit > 200) {
				return fmt.Errorf("game-center achievements localizations list: --limit must be between 1 and 200")
			}
			if err := validateNextURL(*next); err != nil {
				return fmt.Errorf("game-center achievements localizations list: %w", err)
			}

			achID := strings.TrimSpace(*achievementID)
			if achID == "" && strings.TrimSpace(*next) == "" {
				fmt.Fprintln(os.Stderr, "Error: --achievement-id is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center achievements localizations list: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			opts := []asc.GCAchievementLocalizationsOption{
				asc.WithGCAchievementLocalizationsLimit(*limit),
				asc.WithGCAchievementLocalizationsNextURL(*next),
			}

			if *paginate {
				paginateOpts := append(opts, asc.WithGCAchievementLocalizationsLimit(200))
				firstPage, err := client.GetGameCenterAchievementLocalizations(requestCtx, achID, paginateOpts...)
				if err != nil {
					return fmt.Errorf("game-center achievements localizations list: failed to fetch: %w", err)
				}

				resp, err := asc.PaginateAll(requestCtx, firstPage, func(ctx context.Context, nextURL string) (asc.PaginatedResponse, error) {
					return client.GetGameCenterAchievementLocalizations(ctx, achID, asc.WithGCAchievementLocalizationsNextURL(nextURL))
				})
				if err != nil {
					return fmt.Errorf("game-center achievements localizations list: %w", err)
				}

				return printOutput(resp, *output, *pretty)
			}

			resp, err := client.GetGameCenterAchievementLocalizations(requestCtx, achID, opts...)
			if err != nil {
				return fmt.Errorf("game-center achievements localizations list: failed to fetch: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterAchievementLocalizationsGetCommand returns the localizations get subcommand.
func GameCenterAchievementLocalizationsGetCommand() *ffcli.Command {
	fs := flag.NewFlagSet("get", flag.ExitOnError)

	localizationID := fs.String("id", "", "Game Center achievement localization ID")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "get",
		ShortUsage: "asc game-center achievements localizations get --id \"LOC_ID\"",
		ShortHelp:  "Get a Game Center achievement localization by ID.",
		LongHelp: `Get a Game Center achievement localization by ID.

Examples:
  asc game-center achievements localizations get --id "LOC_ID"`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			id := strings.TrimSpace(*localizationID)
			if id == "" {
				fmt.Fprintln(os.Stderr, "Error: --id is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center achievements localizations get: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			resp, err := client.GetGameCenterAchievementLocalization(requestCtx, id)
			if err != nil {
				return fmt.Errorf("game-center achievements localizations get: failed to fetch: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterAchievementLocalizationsCreateCommand returns the localizations create subcommand.
func GameCenterAchievementLocalizationsCreateCommand() *ffcli.Command {
	fs := flag.NewFlagSet("create", flag.ExitOnError)

	achievementID := fs.String("achievement-id", "", "Game Center achievement ID")
	locale := fs.String("locale", "", "Locale code (e.g., en-US)")
	name := fs.String("name", "", "Display name")
	beforeEarnedDescription := fs.String("before-earned-description", "", "Description shown before achievement is earned")
	afterEarnedDescription := fs.String("after-earned-description", "", "Description shown after achievement is earned")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "create",
		ShortUsage: "asc game-center achievements localizations create [flags]",
		ShortHelp:  "Create a new Game Center achievement localization.",
		LongHelp: `Create a new Game Center achievement localization.

Examples:
  asc game-center achievements localizations create --achievement-id "ACHIEVEMENT_ID" --locale en-US --name "First Win"
  asc game-center achievements localizations create --achievement-id "ACHIEVEMENT_ID" --locale en-US --name "First Win" --before-earned-description "Win your first game" --after-earned-description "You won!"`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			achID := strings.TrimSpace(*achievementID)
			if achID == "" {
				fmt.Fprintln(os.Stderr, "Error: --achievement-id is required")
				return flag.ErrHelp
			}

			localeVal := strings.TrimSpace(*locale)
			if localeVal == "" {
				fmt.Fprintln(os.Stderr, "Error: --locale is required")
				return flag.ErrHelp
			}

			nameVal := strings.TrimSpace(*name)
			if nameVal == "" {
				fmt.Fprintln(os.Stderr, "Error: --name is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center achievements localizations create: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			attrs := asc.GameCenterAchievementLocalizationCreateAttributes{
				Locale:                  localeVal,
				Name:                    nameVal,
				BeforeEarnedDescription: strings.TrimSpace(*beforeEarnedDescription),
				AfterEarnedDescription:  strings.TrimSpace(*afterEarnedDescription),
			}

			resp, err := client.CreateGameCenterAchievementLocalization(requestCtx, achID, attrs)
			if err != nil {
				return fmt.Errorf("game-center achievements localizations create: failed to create: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterAchievementLocalizationsUpdateCommand returns the localizations update subcommand.
func GameCenterAchievementLocalizationsUpdateCommand() *ffcli.Command {
	fs := flag.NewFlagSet("update", flag.ExitOnError)

	localizationID := fs.String("id", "", "Game Center achievement localization ID")
	name := fs.String("name", "", "Display name")
	beforeEarnedDescription := fs.String("before-earned-description", "", "Description shown before achievement is earned")
	afterEarnedDescription := fs.String("after-earned-description", "", "Description shown after achievement is earned")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "update",
		ShortUsage: "asc game-center achievements localizations update [flags]",
		ShortHelp:  "Update a Game Center achievement localization.",
		LongHelp: `Update a Game Center achievement localization.

Examples:
  asc game-center achievements localizations update --id "LOC_ID" --name "New Name"
  asc game-center achievements localizations update --id "LOC_ID" --before-earned-description "Win a game" --after-earned-description "Winner!"`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			id := strings.TrimSpace(*localizationID)
			if id == "" {
				fmt.Fprintln(os.Stderr, "Error: --id is required")
				return flag.ErrHelp
			}

			attrs := asc.GameCenterAchievementLocalizationUpdateAttributes{}
			hasUpdate := false

			if strings.TrimSpace(*name) != "" {
				nameVal := strings.TrimSpace(*name)
				attrs.Name = &nameVal
				hasUpdate = true
			}

			if strings.TrimSpace(*beforeEarnedDescription) != "" {
				beforeVal := strings.TrimSpace(*beforeEarnedDescription)
				attrs.BeforeEarnedDescription = &beforeVal
				hasUpdate = true
			}

			if strings.TrimSpace(*afterEarnedDescription) != "" {
				afterVal := strings.TrimSpace(*afterEarnedDescription)
				attrs.AfterEarnedDescription = &afterVal
				hasUpdate = true
			}

			if !hasUpdate {
				fmt.Fprintln(os.Stderr, "Error: at least one update flag is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center achievements localizations update: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			resp, err := client.UpdateGameCenterAchievementLocalization(requestCtx, id, attrs)
			if err != nil {
				return fmt.Errorf("game-center achievements localizations update: failed to update: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterAchievementLocalizationsDeleteCommand returns the localizations delete subcommand.
func GameCenterAchievementLocalizationsDeleteCommand() *ffcli.Command {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)

	localizationID := fs.String("id", "", "Game Center achievement localization ID")
	confirm := fs.Bool("confirm", false, "Confirm deletion")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "delete",
		ShortUsage: "asc game-center achievements localizations delete --id \"LOC_ID\" --confirm",
		ShortHelp:  "Delete a Game Center achievement localization.",
		LongHelp: `Delete a Game Center achievement localization.

Examples:
  asc game-center achievements localizations delete --id "LOC_ID" --confirm`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			id := strings.TrimSpace(*localizationID)
			if id == "" {
				fmt.Fprintln(os.Stderr, "Error: --id is required")
				return flag.ErrHelp
			}
			if !*confirm {
				fmt.Fprintln(os.Stderr, "Error: --confirm is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center achievements localizations delete: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			if err := client.DeleteGameCenterAchievementLocalization(requestCtx, id); err != nil {
				return fmt.Errorf("game-center achievements localizations delete: failed to delete: %w", err)
			}

			result := &asc.GameCenterAchievementLocalizationDeleteResult{
				ID:      id,
				Deleted: true,
			}

			return printOutput(result, *output, *pretty)
		},
	}
}

// GameCenterLeaderboardsCommand returns the leaderboards command group.
func GameCenterLeaderboardsCommand() *ffcli.Command {
	fs := flag.NewFlagSet("leaderboards", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "leaderboards",
		ShortUsage: "asc game-center leaderboards <subcommand> [flags]",
		ShortHelp:  "Manage Game Center leaderboards.",
		LongHelp: `Manage Game Center leaderboards.

Examples:
  asc game-center leaderboards list --app "APP_ID"
  asc game-center leaderboards get --id "LEADERBOARD_ID"
  asc game-center leaderboards create --app "APP_ID" --reference-name "High Score" --vendor-id "com.example.highscore" --formatter INTEGER --sort DESC --submission-type BEST_SCORE
  asc game-center leaderboards update --id "LEADERBOARD_ID" --reference-name "New Name"
  asc game-center leaderboards delete --id "LEADERBOARD_ID" --confirm
  asc game-center leaderboards localizations list --leaderboard-id "LEADERBOARD_ID"
  asc game-center leaderboards localizations create --leaderboard-id "LEADERBOARD_ID" --locale en-US --name "High Score"
  asc game-center leaderboards releases list --leaderboard-id "LEADERBOARD_ID"
  asc game-center leaderboards releases create --app "APP_ID" --leaderboard-id "LEADERBOARD_ID"
  asc game-center leaderboards releases delete --id "RELEASE_ID" --confirm`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Subcommands: []*ffcli.Command{
			GameCenterLeaderboardsListCommand(),
			GameCenterLeaderboardsGetCommand(),
			GameCenterLeaderboardsCreateCommand(),
			GameCenterLeaderboardsUpdateCommand(),
			GameCenterLeaderboardsDeleteCommand(),
			GameCenterLeaderboardLocalizationsCommand(),
			GameCenterLeaderboardReleasesCommand(),
			GameCenterLeaderboardImagesCommand(),
		},
		Exec: func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		},
	}
}

// GameCenterLeaderboardImagesCommand returns the leaderboard images command group.
func GameCenterLeaderboardImagesCommand() *ffcli.Command {
	fs := flag.NewFlagSet("images", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "images",
		ShortUsage: "asc game-center leaderboards images <subcommand> [flags]",
		ShortHelp:  "Manage Game Center leaderboard images.",
		LongHelp: `Manage Game Center leaderboard images.

Images are attached to leaderboard localizations. Use the localization ID when uploading.

Examples:
  asc game-center leaderboards images upload --localization-id "LOC_ID" --file path/to/image.png
  asc game-center leaderboards images delete --id "IMAGE_ID" --confirm`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Subcommands: []*ffcli.Command{
			GameCenterLeaderboardImagesUploadCommand(),
			GameCenterLeaderboardImagesDeleteCommand(),
		},
		Exec: func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		},
	}
}

// GameCenterLeaderboardImagesUploadCommand returns the leaderboard images upload subcommand.
func GameCenterLeaderboardImagesUploadCommand() *ffcli.Command {
	fs := flag.NewFlagSet("upload", flag.ExitOnError)

	localizationID := fs.String("localization-id", "", "Game Center leaderboard localization ID")
	filePath := fs.String("file", "", "Path to image file")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "upload",
		ShortUsage: "asc game-center leaderboards images upload --localization-id \"LOC_ID\" --file path/to/image.png",
		ShortHelp:  "Upload an image for a Game Center leaderboard localization.",
		LongHelp: `Upload an image for a Game Center leaderboard localization.

This command performs the full upload flow: reserves the upload, uploads the file, and commits.

Examples:
  asc game-center leaderboards images upload --localization-id "LOC_ID" --file leaderboard.png`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			locID := strings.TrimSpace(*localizationID)
			if locID == "" {
				fmt.Fprintln(os.Stderr, "Error: --localization-id is required")
				return flag.ErrHelp
			}

			file := strings.TrimSpace(*filePath)
			if file == "" {
				fmt.Fprintln(os.Stderr, "Error: --file is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center leaderboards images upload: %w", err)
			}

			requestCtx, cancel := contextWithUploadTimeout(ctx)
			defer cancel()

			result, err := client.UploadGameCenterLeaderboardImage(requestCtx, locID, file)
			if err != nil {
				return fmt.Errorf("game-center leaderboards images upload: %w", err)
			}

			return printOutput(result, *output, *pretty)
		},
	}
}

// GameCenterLeaderboardImagesDeleteCommand returns the leaderboard images delete subcommand.
func GameCenterLeaderboardImagesDeleteCommand() *ffcli.Command {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)

	imageID := fs.String("id", "", "Game Center leaderboard image ID")
	confirm := fs.Bool("confirm", false, "Confirm deletion")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "delete",
		ShortUsage: "asc game-center leaderboards images delete --id \"IMAGE_ID\" --confirm",
		ShortHelp:  "Delete a Game Center leaderboard image.",
		LongHelp: `Delete a Game Center leaderboard image.

Examples:
  asc game-center leaderboards images delete --id "IMAGE_ID" --confirm`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			id := strings.TrimSpace(*imageID)
			if id == "" {
				fmt.Fprintln(os.Stderr, "Error: --id is required")
				return flag.ErrHelp
			}
			if !*confirm {
				fmt.Fprintln(os.Stderr, "Error: --confirm is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center leaderboards images delete: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			if err := client.DeleteGameCenterLeaderboardImage(requestCtx, id); err != nil {
				return fmt.Errorf("game-center leaderboards images delete: failed to delete: %w", err)
			}

			result := &asc.GameCenterLeaderboardImageDeleteResult{
				ID:      id,
				Deleted: true,
			}

			return printOutput(result, *output, *pretty)
		},
	}
}

// GameCenterLeaderboardsListCommand returns the leaderboards list subcommand.
func GameCenterLeaderboardsListCommand() *ffcli.Command {
	fs := flag.NewFlagSet("list", flag.ExitOnError)

	appID := fs.String("app", "", "App Store Connect app ID (or ASC_APP_ID env)")
	limit := fs.Int("limit", 0, "Maximum results per page (1-200)")
	next := fs.String("next", "", "Fetch next page using a links.next URL")
	paginate := fs.Bool("paginate", false, "Automatically fetch all pages (aggregate results)")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "list",
		ShortUsage: "asc game-center leaderboards list [flags]",
		ShortHelp:  "List Game Center leaderboards for an app.",
		LongHelp: `List Game Center leaderboards for an app.

Examples:
  asc game-center leaderboards list --app "APP_ID"
  asc game-center leaderboards list --app "APP_ID" --limit 50
  asc game-center leaderboards list --app "APP_ID" --paginate`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			if *limit != 0 && (*limit < 1 || *limit > 200) {
				return fmt.Errorf("game-center leaderboards list: --limit must be between 1 and 200")
			}
			if err := validateNextURL(*next); err != nil {
				return fmt.Errorf("game-center leaderboards list: %w", err)
			}

			resolvedAppID := resolveAppID(*appID)
			if resolvedAppID == "" && strings.TrimSpace(*next) == "" {
				fmt.Fprintln(os.Stderr, "Error: --app is required (or set ASC_APP_ID)")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center leaderboards list: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			// Get Game Center detail ID first
			gcDetailID, err := client.GetGameCenterDetailID(requestCtx, resolvedAppID)
			if err != nil {
				return fmt.Errorf("game-center leaderboards list: failed to get Game Center detail: %w", err)
			}

			opts := []asc.GCLeaderboardsOption{
				asc.WithGCLeaderboardsLimit(*limit),
				asc.WithGCLeaderboardsNextURL(*next),
			}

			if *paginate {
				paginateOpts := append(opts, asc.WithGCLeaderboardsLimit(200))
				firstPage, err := client.GetGameCenterLeaderboards(requestCtx, gcDetailID, paginateOpts...)
				if err != nil {
					return fmt.Errorf("game-center leaderboards list: failed to fetch: %w", err)
				}

				resp, err := asc.PaginateAll(requestCtx, firstPage, func(ctx context.Context, nextURL string) (asc.PaginatedResponse, error) {
					return client.GetGameCenterLeaderboards(ctx, gcDetailID, asc.WithGCLeaderboardsNextURL(nextURL))
				})
				if err != nil {
					return fmt.Errorf("game-center leaderboards list: %w", err)
				}

				return printOutput(resp, *output, *pretty)
			}

			resp, err := client.GetGameCenterLeaderboards(requestCtx, gcDetailID, opts...)
			if err != nil {
				return fmt.Errorf("game-center leaderboards list: failed to fetch: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterLeaderboardsGetCommand returns the leaderboards get subcommand.
func GameCenterLeaderboardsGetCommand() *ffcli.Command {
	fs := flag.NewFlagSet("get", flag.ExitOnError)

	leaderboardID := fs.String("id", "", "Game Center leaderboard ID")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "get",
		ShortUsage: "asc game-center leaderboards get --id \"LEADERBOARD_ID\"",
		ShortHelp:  "Get a Game Center leaderboard by ID.",
		LongHelp: `Get a Game Center leaderboard by ID.

Examples:
  asc game-center leaderboards get --id "LEADERBOARD_ID"`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			id := strings.TrimSpace(*leaderboardID)
			if id == "" {
				fmt.Fprintln(os.Stderr, "Error: --id is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center leaderboards get: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			resp, err := client.GetGameCenterLeaderboard(requestCtx, id)
			if err != nil {
				return fmt.Errorf("game-center leaderboards get: failed to fetch: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterLeaderboardsCreateCommand returns the leaderboards create subcommand.
func GameCenterLeaderboardsCreateCommand() *ffcli.Command {
	fs := flag.NewFlagSet("create", flag.ExitOnError)

	appID := fs.String("app", "", "App Store Connect app ID (or ASC_APP_ID env)")
	referenceName := fs.String("reference-name", "", "Reference name for the leaderboard")
	vendorID := fs.String("vendor-id", "", "Vendor identifier (e.g., com.example.leaderboard)")
	formatter := fs.String("formatter", "", "Score formatter: INTEGER, DECIMAL_POINT_1_PLACE, DECIMAL_POINT_2_PLACE, DECIMAL_POINT_3_PLACE, ELAPSED_TIME_MILLISECOND, ELAPSED_TIME_SECOND, ELAPSED_TIME_MINUTE, MONEY_WHOLE, MONEY_POINT_2_PLACE")
	sortType := fs.String("sort", "", "Score sort type: ASC, DESC")
	submissionType := fs.String("submission-type", "", "Submission type: BEST_SCORE, MOST_RECENT_SCORE")
	scoreRangeStart := fs.String("score-range-start", "", "Score range start (optional)")
	scoreRangeEnd := fs.String("score-range-end", "", "Score range end (optional)")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "create",
		ShortUsage: "asc game-center leaderboards create [flags]",
		ShortHelp:  "Create a new Game Center leaderboard.",
		LongHelp: `Create a new Game Center leaderboard.

Examples:
  asc game-center leaderboards create --app "APP_ID" --reference-name "High Score" --vendor-id "com.example.highscore" --formatter INTEGER --sort DESC --submission-type BEST_SCORE
  asc game-center leaderboards create --app "APP_ID" --reference-name "Time Trial" --vendor-id "com.example.timetrial" --formatter ELAPSED_TIME_MILLISECOND --sort ASC --submission-type BEST_SCORE`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			resolvedAppID := resolveAppID(*appID)
			if resolvedAppID == "" {
				fmt.Fprintln(os.Stderr, "Error: --app is required (or set ASC_APP_ID)")
				return flag.ErrHelp
			}

			name := strings.TrimSpace(*referenceName)
			if name == "" {
				fmt.Fprintln(os.Stderr, "Error: --reference-name is required")
				return flag.ErrHelp
			}

			vendor := strings.TrimSpace(*vendorID)
			if vendor == "" {
				fmt.Fprintln(os.Stderr, "Error: --vendor-id is required")
				return flag.ErrHelp
			}

			formatterVal := strings.TrimSpace(strings.ToUpper(*formatter))
			if formatterVal == "" {
				fmt.Fprintln(os.Stderr, "Error: --formatter is required")
				return flag.ErrHelp
			}
			if !isValidLeaderboardFormatter(formatterVal) {
				fmt.Fprintf(os.Stderr, "Error: --formatter must be one of: %s\n", strings.Join(asc.ValidLeaderboardFormatters, ", "))
				return flag.ErrHelp
			}

			sortVal := strings.TrimSpace(strings.ToUpper(*sortType))
			if sortVal == "" {
				fmt.Fprintln(os.Stderr, "Error: --sort is required")
				return flag.ErrHelp
			}
			if !isValidScoreSortType(sortVal) {
				fmt.Fprintf(os.Stderr, "Error: --sort must be one of: %s\n", strings.Join(asc.ValidScoreSortTypes, ", "))
				return flag.ErrHelp
			}

			submissionVal := strings.TrimSpace(strings.ToUpper(*submissionType))
			if submissionVal == "" {
				fmt.Fprintln(os.Stderr, "Error: --submission-type is required")
				return flag.ErrHelp
			}
			if !isValidSubmissionType(submissionVal) {
				fmt.Fprintf(os.Stderr, "Error: --submission-type must be one of: %s\n", strings.Join(asc.ValidSubmissionTypes, ", "))
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center leaderboards create: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			// Get Game Center detail ID first
			gcDetailID, err := client.GetGameCenterDetailID(requestCtx, resolvedAppID)
			if err != nil {
				return fmt.Errorf("game-center leaderboards create: failed to get Game Center detail: %w", err)
			}

			attrs := asc.GameCenterLeaderboardCreateAttributes{
				ReferenceName:    name,
				VendorIdentifier: vendor,
				DefaultFormatter: formatterVal,
				ScoreSortType:    sortVal,
				SubmissionType:   submissionVal,
				ScoreRangeStart:  strings.TrimSpace(*scoreRangeStart),
				ScoreRangeEnd:    strings.TrimSpace(*scoreRangeEnd),
			}

			resp, err := client.CreateGameCenterLeaderboard(requestCtx, gcDetailID, attrs)
			if err != nil {
				return fmt.Errorf("game-center leaderboards create: failed to create: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterLeaderboardsUpdateCommand returns the leaderboards update subcommand.
func GameCenterLeaderboardsUpdateCommand() *ffcli.Command {
	fs := flag.NewFlagSet("update", flag.ExitOnError)

	leaderboardID := fs.String("id", "", "Game Center leaderboard ID")
	referenceName := fs.String("reference-name", "", "Reference name for the leaderboard")
	archived := fs.String("archived", "", "Archive the leaderboard (true/false)")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "update",
		ShortUsage: "asc game-center leaderboards update [flags]",
		ShortHelp:  "Update a Game Center leaderboard.",
		LongHelp: `Update a Game Center leaderboard.

Examples:
  asc game-center leaderboards update --id "LEADERBOARD_ID" --reference-name "New Name"
  asc game-center leaderboards update --id "LEADERBOARD_ID" --archived true`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			id := strings.TrimSpace(*leaderboardID)
			if id == "" {
				fmt.Fprintln(os.Stderr, "Error: --id is required")
				return flag.ErrHelp
			}

			attrs := asc.GameCenterLeaderboardUpdateAttributes{}
			hasUpdate := false

			if strings.TrimSpace(*referenceName) != "" {
				name := strings.TrimSpace(*referenceName)
				attrs.ReferenceName = &name
				hasUpdate = true
			}

			if strings.TrimSpace(*archived) != "" {
				val, err := parseBool(*archived, "--archived")
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error:", err.Error())
					return flag.ErrHelp
				}
				attrs.Archived = &val
				hasUpdate = true
			}

			if !hasUpdate {
				fmt.Fprintln(os.Stderr, "Error: at least one update flag is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center leaderboards update: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			resp, err := client.UpdateGameCenterLeaderboard(requestCtx, id, attrs)
			if err != nil {
				return fmt.Errorf("game-center leaderboards update: failed to update: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterLeaderboardsDeleteCommand returns the leaderboards delete subcommand.
func GameCenterLeaderboardsDeleteCommand() *ffcli.Command {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)

	leaderboardID := fs.String("id", "", "Game Center leaderboard ID")
	confirm := fs.Bool("confirm", false, "Confirm deletion")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "delete",
		ShortUsage: "asc game-center leaderboards delete --id \"LEADERBOARD_ID\" --confirm",
		ShortHelp:  "Delete a Game Center leaderboard.",
		LongHelp: `Delete a Game Center leaderboard.

Examples:
  asc game-center leaderboards delete --id "LEADERBOARD_ID" --confirm`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			id := strings.TrimSpace(*leaderboardID)
			if id == "" {
				fmt.Fprintln(os.Stderr, "Error: --id is required")
				return flag.ErrHelp
			}
			if !*confirm {
				fmt.Fprintln(os.Stderr, "Error: --confirm is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center leaderboards delete: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			if err := client.DeleteGameCenterLeaderboard(requestCtx, id); err != nil {
				return fmt.Errorf("game-center leaderboards delete: failed to delete: %w", err)
			}

			result := &asc.GameCenterLeaderboardDeleteResult{
				ID:      id,
				Deleted: true,
			}

			return printOutput(result, *output, *pretty)
		},
	}
}

func isValidLeaderboardFormatter(value string) bool {
	for _, v := range asc.ValidLeaderboardFormatters {
		if value == v {
			return true
		}
	}
	return false
}

func isValidScoreSortType(value string) bool {
	for _, v := range asc.ValidScoreSortTypes {
		if value == v {
			return true
		}
	}
	return false
}

func isValidSubmissionType(value string) bool {
	for _, v := range asc.ValidSubmissionTypes {
		if value == v {
			return true
		}
	}
	return false
}

// GameCenterLeaderboardReleasesCommand returns the releases command group.
func GameCenterLeaderboardReleasesCommand() *ffcli.Command {
	fs := flag.NewFlagSet("releases", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "releases",
		ShortUsage: "asc game-center leaderboards releases <subcommand> [flags]",
		ShortHelp:  "Manage Game Center leaderboard releases.",
		LongHelp: `Manage Game Center leaderboard releases.

Examples:
  asc game-center leaderboards releases list --leaderboard-id "LEADERBOARD_ID"
  asc game-center leaderboards releases create --app "APP_ID" --leaderboard-id "LEADERBOARD_ID"
  asc game-center leaderboards releases delete --id "RELEASE_ID" --confirm`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Subcommands: []*ffcli.Command{
			GameCenterLeaderboardReleasesListCommand(),
			GameCenterLeaderboardReleasesCreateCommand(),
			GameCenterLeaderboardReleasesDeleteCommand(),
		},
		Exec: func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		},
	}
}

// GameCenterLeaderboardReleasesListCommand returns the releases list subcommand.
func GameCenterLeaderboardReleasesListCommand() *ffcli.Command {
	fs := flag.NewFlagSet("list", flag.ExitOnError)

	leaderboardID := fs.String("leaderboard-id", "", "Game Center leaderboard ID")
	limit := fs.Int("limit", 0, "Maximum results per page (1-200)")
	next := fs.String("next", "", "Fetch next page using a links.next URL")
	paginate := fs.Bool("paginate", false, "Automatically fetch all pages (aggregate results)")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "list",
		ShortUsage: "asc game-center leaderboards releases list --leaderboard-id \"LEADERBOARD_ID\"",
		ShortHelp:  "List releases for a Game Center leaderboard.",
		LongHelp: `List releases for a Game Center leaderboard.

Examples:
  asc game-center leaderboards releases list --leaderboard-id "LEADERBOARD_ID"
  asc game-center leaderboards releases list --leaderboard-id "LEADERBOARD_ID" --limit 50
  asc game-center leaderboards releases list --leaderboard-id "LEADERBOARD_ID" --paginate`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			if *limit != 0 && (*limit < 1 || *limit > 200) {
				return fmt.Errorf("game-center leaderboards releases list: --limit must be between 1 and 200")
			}
			if err := validateNextURL(*next); err != nil {
				return fmt.Errorf("game-center leaderboards releases list: %w", err)
			}

			lbID := strings.TrimSpace(*leaderboardID)
			if lbID == "" && strings.TrimSpace(*next) == "" {
				fmt.Fprintln(os.Stderr, "Error: --leaderboard-id is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center leaderboards releases list: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			opts := []asc.GCLeaderboardReleasesOption{
				asc.WithGCLeaderboardReleasesLimit(*limit),
				asc.WithGCLeaderboardReleasesNextURL(*next),
			}

			if *paginate {
				paginateOpts := append(opts, asc.WithGCLeaderboardReleasesLimit(200))
				firstPage, err := client.GetGameCenterLeaderboardReleases(requestCtx, lbID, paginateOpts...)
				if err != nil {
					return fmt.Errorf("game-center leaderboards releases list: failed to fetch: %w", err)
				}

				resp, err := asc.PaginateAll(requestCtx, firstPage, func(ctx context.Context, nextURL string) (asc.PaginatedResponse, error) {
					return client.GetGameCenterLeaderboardReleases(ctx, lbID, asc.WithGCLeaderboardReleasesNextURL(nextURL))
				})
				if err != nil {
					return fmt.Errorf("game-center leaderboards releases list: %w", err)
				}

				return printOutput(resp, *output, *pretty)
			}

			resp, err := client.GetGameCenterLeaderboardReleases(requestCtx, lbID, opts...)
			if err != nil {
				return fmt.Errorf("game-center leaderboards releases list: failed to fetch: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterLeaderboardReleasesCreateCommand returns the releases create subcommand.
func GameCenterLeaderboardReleasesCreateCommand() *ffcli.Command {
	fs := flag.NewFlagSet("create", flag.ExitOnError)

	appID := fs.String("app", "", "App Store Connect app ID (or ASC_APP_ID env)")
	leaderboardID := fs.String("leaderboard-id", "", "Game Center leaderboard ID")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "create",
		ShortUsage: "asc game-center leaderboards releases create --app \"APP_ID\" --leaderboard-id \"LEADERBOARD_ID\"",
		ShortHelp:  "Create a release for a Game Center leaderboard.",
		LongHelp: `Create a release for a Game Center leaderboard.

A release associates a leaderboard with a Game Center detail, making it live.

Examples:
  asc game-center leaderboards releases create --app "APP_ID" --leaderboard-id "LEADERBOARD_ID"`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			resolvedAppID := resolveAppID(*appID)
			if resolvedAppID == "" {
				fmt.Fprintln(os.Stderr, "Error: --app is required (or set ASC_APP_ID)")
				return flag.ErrHelp
			}

			lbID := strings.TrimSpace(*leaderboardID)
			if lbID == "" {
				fmt.Fprintln(os.Stderr, "Error: --leaderboard-id is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center leaderboards releases create: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			// Get Game Center detail ID first
			gcDetailID, err := client.GetGameCenterDetailID(requestCtx, resolvedAppID)
			if err != nil {
				return fmt.Errorf("game-center leaderboards releases create: failed to get Game Center detail: %w", err)
			}

			resp, err := client.CreateGameCenterLeaderboardRelease(requestCtx, gcDetailID, lbID)
			if err != nil {
				return fmt.Errorf("game-center leaderboards releases create: failed to create: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterLeaderboardReleasesDeleteCommand returns the releases delete subcommand.
func GameCenterLeaderboardReleasesDeleteCommand() *ffcli.Command {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)

	releaseID := fs.String("id", "", "Game Center leaderboard release ID")
	confirm := fs.Bool("confirm", false, "Confirm deletion")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "delete",
		ShortUsage: "asc game-center leaderboards releases delete --id \"RELEASE_ID\" --confirm",
		ShortHelp:  "Delete a Game Center leaderboard release.",
		LongHelp: `Delete a Game Center leaderboard release.

Examples:
  asc game-center leaderboards releases delete --id "RELEASE_ID" --confirm`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			id := strings.TrimSpace(*releaseID)
			if id == "" {
				fmt.Fprintln(os.Stderr, "Error: --id is required")
				return flag.ErrHelp
			}
			if !*confirm {
				fmt.Fprintln(os.Stderr, "Error: --confirm is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center leaderboards releases delete: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			if err := client.DeleteGameCenterLeaderboardRelease(requestCtx, id); err != nil {
				return fmt.Errorf("game-center leaderboards releases delete: failed to delete: %w", err)
			}

			result := &asc.GameCenterLeaderboardReleaseDeleteResult{
				ID:      id,
				Deleted: true,
			}

			return printOutput(result, *output, *pretty)
		},
	}
}

// GameCenterLeaderboardSetsCommand returns the leaderboard-sets command group.
func GameCenterLeaderboardSetsCommand() *ffcli.Command {
	fs := flag.NewFlagSet("leaderboard-sets", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "leaderboard-sets",
		ShortUsage: "asc game-center leaderboard-sets <subcommand> [flags]",
		ShortHelp:  "Manage Game Center leaderboard sets.",
		LongHelp: `Manage Game Center leaderboard sets.

Examples:
  asc game-center leaderboard-sets list --app "APP_ID"
  asc game-center leaderboard-sets get --id "SET_ID"
  asc game-center leaderboard-sets create --app "APP_ID" --reference-name "Season 1" --vendor-id "com.example.season1"
  asc game-center leaderboard-sets update --id "SET_ID" --reference-name "Season 1 - Updated"
  asc game-center leaderboard-sets delete --id "SET_ID" --confirm
  asc game-center leaderboard-sets members list --set-id "SET_ID"
  asc game-center leaderboard-sets members set --set-id "SET_ID" --leaderboard-ids "id1,id2,id3"
  asc game-center leaderboard-sets releases list --set-id "SET_ID"
  asc game-center leaderboard-sets releases create --app "APP_ID" --set-id "SET_ID"
  asc game-center leaderboard-sets releases delete --id "RELEASE_ID" --confirm
  asc game-center leaderboard-sets localizations list --set-id "SET_ID"
  asc game-center leaderboard-sets localizations create --set-id "SET_ID" --locale en-US --name "Season 1"
  asc game-center leaderboard-sets images upload --localization-id "LOC_ID" --file path/to/image.png
  asc game-center leaderboard-sets images delete --id "IMAGE_ID" --confirm`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Subcommands: []*ffcli.Command{
			GameCenterLeaderboardSetsListCommand(),
			GameCenterLeaderboardSetsGetCommand(),
			GameCenterLeaderboardSetsCreateCommand(),
			GameCenterLeaderboardSetsUpdateCommand(),
			GameCenterLeaderboardSetsDeleteCommand(),
			GameCenterLeaderboardSetMembersCommand(),
			GameCenterLeaderboardSetReleasesCommand(),
			GameCenterLeaderboardSetImagesCommand(),
			GameCenterLeaderboardSetLocalizationsCommand(),
		},
		Exec: func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		},
	}
}

// GameCenterLeaderboardSetLocalizationsCommand returns the localizations command group.
func GameCenterLeaderboardSetLocalizationsCommand() *ffcli.Command {
	fs := flag.NewFlagSet("localizations", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "localizations",
		ShortUsage: "asc game-center leaderboard-sets localizations <subcommand> [flags]",
		ShortHelp:  "Manage Game Center leaderboard set localizations.",
		LongHelp: `Manage Game Center leaderboard set localizations.

Examples:
  asc game-center leaderboard-sets localizations list --set-id "SET_ID"
  asc game-center leaderboard-sets localizations get --id "LOC_ID"
  asc game-center leaderboard-sets localizations create --set-id "SET_ID" --locale en-US --name "Season 1"
  asc game-center leaderboard-sets localizations update --id "LOC_ID" --name "New Name"
  asc game-center leaderboard-sets localizations delete --id "LOC_ID" --confirm`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Subcommands: []*ffcli.Command{
			GameCenterLeaderboardSetLocalizationsListCommand(),
			GameCenterLeaderboardSetLocalizationsGetCommand(),
			GameCenterLeaderboardSetLocalizationsCreateCommand(),
			GameCenterLeaderboardSetLocalizationsUpdateCommand(),
			GameCenterLeaderboardSetLocalizationsDeleteCommand(),
		},
		Exec: func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		},
	}
}

// GameCenterLeaderboardSetLocalizationsListCommand returns the localizations list subcommand.
func GameCenterLeaderboardSetLocalizationsListCommand() *ffcli.Command {
	fs := flag.NewFlagSet("list", flag.ExitOnError)

	setID := fs.String("set-id", "", "Game Center leaderboard set ID")
	limit := fs.Int("limit", 0, "Maximum results per page (1-200)")
	next := fs.String("next", "", "Fetch next page using a links.next URL")
	paginate := fs.Bool("paginate", false, "Automatically fetch all pages (aggregate results)")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "list",
		ShortUsage: "asc game-center leaderboard-sets localizations list --set-id \"SET_ID\"",
		ShortHelp:  "List localizations for a Game Center leaderboard set.",
		LongHelp: `List localizations for a Game Center leaderboard set.

Examples:
  asc game-center leaderboard-sets localizations list --set-id "SET_ID"
  asc game-center leaderboard-sets localizations list --set-id "SET_ID" --limit 50
  asc game-center leaderboard-sets localizations list --set-id "SET_ID" --paginate`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			if *limit != 0 && (*limit < 1 || *limit > 200) {
				return fmt.Errorf("game-center leaderboard-sets localizations list: --limit must be between 1 and 200")
			}
			if err := validateNextURL(*next); err != nil {
				return fmt.Errorf("game-center leaderboard-sets localizations list: %w", err)
			}

			id := strings.TrimSpace(*setID)
			if id == "" && strings.TrimSpace(*next) == "" {
				fmt.Fprintln(os.Stderr, "Error: --set-id is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center leaderboard-sets localizations list: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			opts := []asc.GCLeaderboardSetLocalizationsOption{
				asc.WithGCLeaderboardSetLocalizationsLimit(*limit),
				asc.WithGCLeaderboardSetLocalizationsNextURL(*next),
			}

			if *paginate {
				paginateOpts := append(opts, asc.WithGCLeaderboardSetLocalizationsLimit(200))
				firstPage, err := client.GetGameCenterLeaderboardSetLocalizations(requestCtx, id, paginateOpts...)
				if err != nil {
					return fmt.Errorf("game-center leaderboard-sets localizations list: failed to fetch: %w", err)
				}

				resp, err := asc.PaginateAll(requestCtx, firstPage, func(ctx context.Context, nextURL string) (asc.PaginatedResponse, error) {
					return client.GetGameCenterLeaderboardSetLocalizations(ctx, id, asc.WithGCLeaderboardSetLocalizationsNextURL(nextURL))
				})
				if err != nil {
					return fmt.Errorf("game-center leaderboard-sets localizations list: %w", err)
				}

				return printOutput(resp, *output, *pretty)
			}

			resp, err := client.GetGameCenterLeaderboardSetLocalizations(requestCtx, id, opts...)
			if err != nil {
				return fmt.Errorf("game-center leaderboard-sets localizations list: failed to fetch: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterLeaderboardSetLocalizationsGetCommand returns the localizations get subcommand.
func GameCenterLeaderboardSetLocalizationsGetCommand() *ffcli.Command {
	fs := flag.NewFlagSet("get", flag.ExitOnError)

	localizationID := fs.String("id", "", "Game Center leaderboard set localization ID")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "get",
		ShortUsage: "asc game-center leaderboard-sets localizations get --id \"LOC_ID\"",
		ShortHelp:  "Get a Game Center leaderboard set localization by ID.",
		LongHelp: `Get a Game Center leaderboard set localization by ID.

Examples:
  asc game-center leaderboard-sets localizations get --id "LOC_ID"`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			id := strings.TrimSpace(*localizationID)
			if id == "" {
				fmt.Fprintln(os.Stderr, "Error: --id is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center leaderboard-sets localizations get: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			resp, err := client.GetGameCenterLeaderboardSetLocalization(requestCtx, id)
			if err != nil {
				return fmt.Errorf("game-center leaderboard-sets localizations get: failed to fetch: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterLeaderboardSetLocalizationsCreateCommand returns the localizations create subcommand.
func GameCenterLeaderboardSetLocalizationsCreateCommand() *ffcli.Command {
	fs := flag.NewFlagSet("create", flag.ExitOnError)

	setID := fs.String("set-id", "", "Game Center leaderboard set ID")
	locale := fs.String("locale", "", "Locale code (e.g., en-US, de-DE)")
	name := fs.String("name", "", "Display name for the leaderboard set")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "create",
		ShortUsage: "asc game-center leaderboard-sets localizations create --set-id \"SET_ID\" --locale \"LOCALE\" --name \"NAME\"",
		ShortHelp:  "Create a localization for a Game Center leaderboard set.",
		LongHelp: `Create a localization for a Game Center leaderboard set.

Examples:
  asc game-center leaderboard-sets localizations create --set-id "SET_ID" --locale en-US --name "Season 1"
  asc game-center leaderboard-sets localizations create --set-id "SET_ID" --locale de-DE --name "Staffel 1"`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			id := strings.TrimSpace(*setID)
			if id == "" {
				fmt.Fprintln(os.Stderr, "Error: --set-id is required")
				return flag.ErrHelp
			}

			localeVal := strings.TrimSpace(*locale)
			if localeVal == "" {
				fmt.Fprintln(os.Stderr, "Error: --locale is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center leaderboard-sets localizations create: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			attrs := asc.GameCenterLeaderboardSetLocalizationCreateAttributes{
				Locale: localeVal,
				Name:   strings.TrimSpace(*name),
			}

			resp, err := client.CreateGameCenterLeaderboardSetLocalization(requestCtx, id, attrs)
			if err != nil {
				return fmt.Errorf("game-center leaderboard-sets localizations create: failed to create: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterLeaderboardSetLocalizationsUpdateCommand returns the localizations update subcommand.
func GameCenterLeaderboardSetLocalizationsUpdateCommand() *ffcli.Command {
	fs := flag.NewFlagSet("update", flag.ExitOnError)

	localizationID := fs.String("id", "", "Game Center leaderboard set localization ID")
	name := fs.String("name", "", "Display name for the leaderboard set")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "update",
		ShortUsage: "asc game-center leaderboard-sets localizations update --id \"LOC_ID\" --name \"NAME\"",
		ShortHelp:  "Update a Game Center leaderboard set localization.",
		LongHelp: `Update a Game Center leaderboard set localization.

Examples:
  asc game-center leaderboard-sets localizations update --id "LOC_ID" --name "New Name"`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			id := strings.TrimSpace(*localizationID)
			if id == "" {
				fmt.Fprintln(os.Stderr, "Error: --id is required")
				return flag.ErrHelp
			}

			attrs := asc.GameCenterLeaderboardSetLocalizationUpdateAttributes{}
			hasUpdate := false

			if strings.TrimSpace(*name) != "" {
				nameVal := strings.TrimSpace(*name)
				attrs.Name = &nameVal
				hasUpdate = true
			}

			if !hasUpdate {
				fmt.Fprintln(os.Stderr, "Error: at least one update flag is required (--name)")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center leaderboard-sets localizations update: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			resp, err := client.UpdateGameCenterLeaderboardSetLocalization(requestCtx, id, attrs)
			if err != nil {
				return fmt.Errorf("game-center leaderboard-sets localizations update: failed to update: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterLeaderboardSetLocalizationsDeleteCommand returns the localizations delete subcommand.
func GameCenterLeaderboardSetLocalizationsDeleteCommand() *ffcli.Command {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)

	localizationID := fs.String("id", "", "Game Center leaderboard set localization ID")
	confirm := fs.Bool("confirm", false, "Confirm deletion")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "delete",
		ShortUsage: "asc game-center leaderboard-sets localizations delete --id \"LOC_ID\" --confirm",
		ShortHelp:  "Delete a Game Center leaderboard set localization.",
		LongHelp: `Delete a Game Center leaderboard set localization.

Examples:
  asc game-center leaderboard-sets localizations delete --id "LOC_ID" --confirm`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			id := strings.TrimSpace(*localizationID)
			if id == "" {
				fmt.Fprintln(os.Stderr, "Error: --id is required")
				return flag.ErrHelp
			}
			if !*confirm {
				fmt.Fprintln(os.Stderr, "Error: --confirm is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center leaderboard-sets localizations delete: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			if err := client.DeleteGameCenterLeaderboardSetLocalization(requestCtx, id); err != nil {
				return fmt.Errorf("game-center leaderboard-sets localizations delete: failed to delete: %w", err)
			}

			result := &asc.GameCenterLeaderboardSetLocalizationDeleteResult{
				ID:      id,
				Deleted: true,
			}

			return printOutput(result, *output, *pretty)
		},
	}
}

// GameCenterLeaderboardSetsListCommand returns the leaderboard-sets list subcommand.
func GameCenterLeaderboardSetsListCommand() *ffcli.Command {
	fs := flag.NewFlagSet("list", flag.ExitOnError)

	appID := fs.String("app", "", "App Store Connect app ID (or ASC_APP_ID env)")
	limit := fs.Int("limit", 0, "Maximum results per page (1-200)")
	next := fs.String("next", "", "Fetch next page using a links.next URL")
	paginate := fs.Bool("paginate", false, "Automatically fetch all pages (aggregate results)")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "list",
		ShortUsage: "asc game-center leaderboard-sets list [flags]",
		ShortHelp:  "List Game Center leaderboard sets for an app.",
		LongHelp: `List Game Center leaderboard sets for an app.

Examples:
  asc game-center leaderboard-sets list --app "APP_ID"
  asc game-center leaderboard-sets list --app "APP_ID" --limit 50
  asc game-center leaderboard-sets list --app "APP_ID" --paginate`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			if *limit != 0 && (*limit < 1 || *limit > 200) {
				return fmt.Errorf("game-center leaderboard-sets list: --limit must be between 1 and 200")
			}
			if err := validateNextURL(*next); err != nil {
				return fmt.Errorf("game-center leaderboard-sets list: %w", err)
			}

			resolvedAppID := resolveAppID(*appID)
			if resolvedAppID == "" && strings.TrimSpace(*next) == "" {
				fmt.Fprintln(os.Stderr, "Error: --app is required (or set ASC_APP_ID)")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center leaderboard-sets list: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			// Get Game Center detail ID first
			gcDetailID, err := client.GetGameCenterDetailID(requestCtx, resolvedAppID)
			if err != nil {
				return fmt.Errorf("game-center leaderboard-sets list: failed to get Game Center detail: %w", err)
			}

			opts := []asc.GCLeaderboardSetsOption{
				asc.WithGCLeaderboardSetsLimit(*limit),
				asc.WithGCLeaderboardSetsNextURL(*next),
			}

			if *paginate {
				paginateOpts := append(opts, asc.WithGCLeaderboardSetsLimit(200))
				firstPage, err := client.GetGameCenterLeaderboardSets(requestCtx, gcDetailID, paginateOpts...)
				if err != nil {
					return fmt.Errorf("game-center leaderboard-sets list: failed to fetch: %w", err)
				}

				resp, err := asc.PaginateAll(requestCtx, firstPage, func(ctx context.Context, nextURL string) (asc.PaginatedResponse, error) {
					return client.GetGameCenterLeaderboardSets(ctx, gcDetailID, asc.WithGCLeaderboardSetsNextURL(nextURL))
				})
				if err != nil {
					return fmt.Errorf("game-center leaderboard-sets list: %w", err)
				}

				return printOutput(resp, *output, *pretty)
			}

			resp, err := client.GetGameCenterLeaderboardSets(requestCtx, gcDetailID, opts...)
			if err != nil {
				return fmt.Errorf("game-center leaderboard-sets list: failed to fetch: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterLeaderboardSetsGetCommand returns the leaderboard-sets get subcommand.
func GameCenterLeaderboardSetsGetCommand() *ffcli.Command {
	fs := flag.NewFlagSet("get", flag.ExitOnError)

	setID := fs.String("id", "", "Game Center leaderboard set ID")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "get",
		ShortUsage: "asc game-center leaderboard-sets get --id \"SET_ID\"",
		ShortHelp:  "Get a Game Center leaderboard set by ID.",
		LongHelp: `Get a Game Center leaderboard set by ID.

Examples:
  asc game-center leaderboard-sets get --id "SET_ID"`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			id := strings.TrimSpace(*setID)
			if id == "" {
				fmt.Fprintln(os.Stderr, "Error: --id is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center leaderboard-sets get: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			resp, err := client.GetGameCenterLeaderboardSet(requestCtx, id)
			if err != nil {
				return fmt.Errorf("game-center leaderboard-sets get: failed to fetch: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterLeaderboardSetsCreateCommand returns the leaderboard-sets create subcommand.
func GameCenterLeaderboardSetsCreateCommand() *ffcli.Command {
	fs := flag.NewFlagSet("create", flag.ExitOnError)

	appID := fs.String("app", "", "App Store Connect app ID (or ASC_APP_ID env)")
	referenceName := fs.String("reference-name", "", "Reference name for the leaderboard set")
	vendorID := fs.String("vendor-id", "", "Vendor identifier (e.g., com.example.set)")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "create",
		ShortUsage: "asc game-center leaderboard-sets create [flags]",
		ShortHelp:  "Create a new Game Center leaderboard set.",
		LongHelp: `Create a new Game Center leaderboard set.

Examples:
  asc game-center leaderboard-sets create --app "APP_ID" --reference-name "Season 1" --vendor-id "com.example.season1"
  asc game-center leaderboard-sets create --app "APP_ID" --reference-name "Weekly Challenge" --vendor-id "com.example.weekly"`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			resolvedAppID := resolveAppID(*appID)
			if resolvedAppID == "" {
				fmt.Fprintln(os.Stderr, "Error: --app is required (or set ASC_APP_ID)")
				return flag.ErrHelp
			}

			name := strings.TrimSpace(*referenceName)
			if name == "" {
				fmt.Fprintln(os.Stderr, "Error: --reference-name is required")
				return flag.ErrHelp
			}

			vendor := strings.TrimSpace(*vendorID)
			if vendor == "" {
				fmt.Fprintln(os.Stderr, "Error: --vendor-id is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center leaderboard-sets create: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			// Get Game Center detail ID first
			gcDetailID, err := client.GetGameCenterDetailID(requestCtx, resolvedAppID)
			if err != nil {
				return fmt.Errorf("game-center leaderboard-sets create: failed to get Game Center detail: %w", err)
			}

			attrs := asc.GameCenterLeaderboardSetCreateAttributes{
				ReferenceName:    name,
				VendorIdentifier: vendor,
			}

			resp, err := client.CreateGameCenterLeaderboardSet(requestCtx, gcDetailID, attrs)
			if err != nil {
				return fmt.Errorf("game-center leaderboard-sets create: failed to create: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterLeaderboardSetsUpdateCommand returns the leaderboard-sets update subcommand.
func GameCenterLeaderboardSetsUpdateCommand() *ffcli.Command {
	fs := flag.NewFlagSet("update", flag.ExitOnError)

	setID := fs.String("id", "", "Game Center leaderboard set ID")
	referenceName := fs.String("reference-name", "", "Reference name for the leaderboard set")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "update",
		ShortUsage: "asc game-center leaderboard-sets update [flags]",
		ShortHelp:  "Update a Game Center leaderboard set.",
		LongHelp: `Update a Game Center leaderboard set.

Examples:
  asc game-center leaderboard-sets update --id "SET_ID" --reference-name "Season 1 - Updated"`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			id := strings.TrimSpace(*setID)
			if id == "" {
				fmt.Fprintln(os.Stderr, "Error: --id is required")
				return flag.ErrHelp
			}

			attrs := asc.GameCenterLeaderboardSetUpdateAttributes{}
			hasUpdate := false

			if strings.TrimSpace(*referenceName) != "" {
				name := strings.TrimSpace(*referenceName)
				attrs.ReferenceName = &name
				hasUpdate = true
			}

			if !hasUpdate {
				fmt.Fprintln(os.Stderr, "Error: at least one update flag is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center leaderboard-sets update: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			resp, err := client.UpdateGameCenterLeaderboardSet(requestCtx, id, attrs)
			if err != nil {
				return fmt.Errorf("game-center leaderboard-sets update: failed to update: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterLeaderboardSetsDeleteCommand returns the leaderboard-sets delete subcommand.
func GameCenterLeaderboardSetsDeleteCommand() *ffcli.Command {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)

	setID := fs.String("id", "", "Game Center leaderboard set ID")
	confirm := fs.Bool("confirm", false, "Confirm deletion")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "delete",
		ShortUsage: "asc game-center leaderboard-sets delete --id \"SET_ID\" --confirm",
		ShortHelp:  "Delete a Game Center leaderboard set.",
		LongHelp: `Delete a Game Center leaderboard set.

Examples:
  asc game-center leaderboard-sets delete --id "SET_ID" --confirm`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			id := strings.TrimSpace(*setID)
			if id == "" {
				fmt.Fprintln(os.Stderr, "Error: --id is required")
				return flag.ErrHelp
			}
			if !*confirm {
				fmt.Fprintln(os.Stderr, "Error: --confirm is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center leaderboard-sets delete: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			if err := client.DeleteGameCenterLeaderboardSet(requestCtx, id); err != nil {
				return fmt.Errorf("game-center leaderboard-sets delete: failed to delete: %w", err)
			}

			result := &asc.GameCenterLeaderboardSetDeleteResult{
				ID:      id,
				Deleted: true,
			}

			return printOutput(result, *output, *pretty)
		},
	}
}

// GameCenterAchievementReleasesCommand returns the achievement releases command group.
func GameCenterAchievementReleasesCommand() *ffcli.Command {
	fs := flag.NewFlagSet("releases", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "releases",
		ShortUsage: "asc game-center achievements releases <subcommand> [flags]",
		ShortHelp:  "Manage Game Center achievement releases.",
		LongHelp: `Manage Game Center achievement releases. Releases are used to publish achievements to live.

Examples:
  asc game-center achievements releases list --achievement-id "ACHIEVEMENT_ID"
  asc game-center achievements releases create --app "APP_ID" --achievement-id "ACHIEVEMENT_ID"
  asc game-center achievements releases delete --id "RELEASE_ID" --confirm`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Subcommands: []*ffcli.Command{
			GameCenterAchievementReleasesListCommand(),
			GameCenterAchievementReleasesCreateCommand(),
			GameCenterAchievementReleasesDeleteCommand(),
		},
		Exec: func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		},
	}
}

// GameCenterAchievementReleasesListCommand returns the achievement releases list subcommand.
func GameCenterAchievementReleasesListCommand() *ffcli.Command {
	fs := flag.NewFlagSet("list", flag.ExitOnError)

	achievementID := fs.String("achievement-id", "", "Game Center achievement ID")
	limit := fs.Int("limit", 0, "Maximum results per page (1-200)")
	next := fs.String("next", "", "Fetch next page using a links.next URL")
	paginate := fs.Bool("paginate", false, "Automatically fetch all pages (aggregate results)")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "list",
		ShortUsage: "asc game-center achievements releases list --achievement-id \"ACHIEVEMENT_ID\"",
		ShortHelp:  "List releases for a Game Center achievement.",
		LongHelp: `List releases for a Game Center achievement.

Examples:
  asc game-center achievements releases list --achievement-id "ACHIEVEMENT_ID"
  asc game-center achievements releases list --achievement-id "ACHIEVEMENT_ID" --limit 50
  asc game-center achievements releases list --achievement-id "ACHIEVEMENT_ID" --paginate`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			if *limit != 0 && (*limit < 1 || *limit > 200) {
				return fmt.Errorf("game-center achievements releases list: --limit must be between 1 and 200")
			}
			if err := validateNextURL(*next); err != nil {
				return fmt.Errorf("game-center achievements releases list: %w", err)
			}

			id := strings.TrimSpace(*achievementID)
			if id == "" && strings.TrimSpace(*next) == "" {
				fmt.Fprintln(os.Stderr, "Error: --achievement-id is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center achievements releases list: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			opts := []asc.GCAchievementReleasesOption{
				asc.WithGCAchievementReleasesLimit(*limit),
				asc.WithGCAchievementReleasesNextURL(*next),
			}

			if *paginate {
				paginateOpts := append(opts, asc.WithGCAchievementReleasesLimit(200))
				firstPage, err := client.GetGameCenterAchievementReleases(requestCtx, id, paginateOpts...)
				if err != nil {
					return fmt.Errorf("game-center achievements releases list: failed to fetch: %w", err)
				}

				resp, err := asc.PaginateAll(requestCtx, firstPage, func(ctx context.Context, nextURL string) (asc.PaginatedResponse, error) {
					return client.GetGameCenterAchievementReleases(ctx, id, asc.WithGCAchievementReleasesNextURL(nextURL))
				})
				if err != nil {
					return fmt.Errorf("game-center achievements releases list: %w", err)
				}

				return printOutput(resp, *output, *pretty)
			}

			resp, err := client.GetGameCenterAchievementReleases(requestCtx, id, opts...)
			if err != nil {
				return fmt.Errorf("game-center achievements releases list: failed to fetch: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterAchievementReleasesCreateCommand returns the achievement releases create subcommand.
func GameCenterAchievementReleasesCreateCommand() *ffcli.Command {
	fs := flag.NewFlagSet("create", flag.ExitOnError)

	appID := fs.String("app", "", "App Store Connect app ID (or ASC_APP_ID env)")
	achievementID := fs.String("achievement-id", "", "Game Center achievement ID")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "create",
		ShortUsage: "asc game-center achievements releases create --app \"APP_ID\" --achievement-id \"ACHIEVEMENT_ID\"",
		ShortHelp:  "Create a new Game Center achievement release.",
		LongHelp: `Create a new Game Center achievement release. This publishes the achievement to live.

Examples:
  asc game-center achievements releases create --app "APP_ID" --achievement-id "ACHIEVEMENT_ID"`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			resolvedAppID := resolveAppID(*appID)
			if resolvedAppID == "" {
				fmt.Fprintln(os.Stderr, "Error: --app is required (or set ASC_APP_ID)")
				return flag.ErrHelp
			}

			id := strings.TrimSpace(*achievementID)
			if id == "" {
				fmt.Fprintln(os.Stderr, "Error: --achievement-id is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center achievements releases create: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			// Get Game Center detail ID first
			gcDetailID, err := client.GetGameCenterDetailID(requestCtx, resolvedAppID)
			if err != nil {
				return fmt.Errorf("game-center achievements releases create: failed to get Game Center detail: %w", err)
			}

			resp, err := client.CreateGameCenterAchievementRelease(requestCtx, gcDetailID, id)
			if err != nil {
				return fmt.Errorf("game-center achievements releases create: failed to create: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterAchievementReleasesDeleteCommand returns the achievement releases delete subcommand.
func GameCenterAchievementReleasesDeleteCommand() *ffcli.Command {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)

	releaseID := fs.String("id", "", "Game Center achievement release ID")
	confirm := fs.Bool("confirm", false, "Confirm deletion")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "delete",
		ShortUsage: "asc game-center achievements releases delete --id \"RELEASE_ID\" --confirm",
		ShortHelp:  "Delete a Game Center achievement release.",
		LongHelp: `Delete a Game Center achievement release.

Examples:
  asc game-center achievements releases delete --id "RELEASE_ID" --confirm`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			id := strings.TrimSpace(*releaseID)
			if id == "" {
				fmt.Fprintln(os.Stderr, "Error: --id is required")
				return flag.ErrHelp
			}
			if !*confirm {
				fmt.Fprintln(os.Stderr, "Error: --confirm is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center achievements releases delete: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			if err := client.DeleteGameCenterAchievementRelease(requestCtx, id); err != nil {
				return fmt.Errorf("game-center achievements releases delete: failed to delete: %w", err)
			}

			result := &asc.GameCenterAchievementReleaseDeleteResult{
				ID:      id,
				Deleted: true,
			}

			return printOutput(result, *output, *pretty)
		},
	}
}

// GameCenterAchievementImagesCommand returns the achievement images command group.
func GameCenterAchievementImagesCommand() *ffcli.Command {
	fs := flag.NewFlagSet("images", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "images",
		ShortUsage: "asc game-center achievements images <subcommand> [flags]",
		ShortHelp:  "Manage Game Center achievement images.",
		LongHelp: `Manage Game Center achievement images. Images are attached to achievement localizations.

Examples:
  asc game-center achievements images upload --localization-id "LOC_ID" --file "path/to/image.png"
  asc game-center achievements images get --id "IMAGE_ID"
  asc game-center achievements images delete --id "IMAGE_ID" --confirm`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Subcommands: []*ffcli.Command{
			GameCenterAchievementImagesUploadCommand(),
			GameCenterAchievementImagesGetCommand(),
			GameCenterAchievementImagesDeleteCommand(),
		},
		Exec: func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		},
	}
}

// GameCenterAchievementImagesUploadCommand returns the achievement images upload subcommand.
func GameCenterAchievementImagesUploadCommand() *ffcli.Command {
	fs := flag.NewFlagSet("upload", flag.ExitOnError)

	localizationID := fs.String("localization-id", "", "Game Center achievement localization ID")
	filePath := fs.String("file", "", "Path to the image file to upload")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "upload",
		ShortUsage: "asc game-center achievements images upload --localization-id \"LOC_ID\" --file \"path/to/image.png\"",
		ShortHelp:  "Upload an image for a Game Center achievement localization.",
		LongHelp: `Upload an image for a Game Center achievement localization.

The image file will be validated, reserved, uploaded in chunks, and committed.

Examples:
  asc game-center achievements images upload --localization-id "LOC_ID" --file "path/to/image.png"`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			locID := strings.TrimSpace(*localizationID)
			if locID == "" {
				fmt.Fprintln(os.Stderr, "Error: --localization-id is required")
				return flag.ErrHelp
			}

			path := strings.TrimSpace(*filePath)
			if path == "" {
				fmt.Fprintln(os.Stderr, "Error: --file is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center achievements images upload: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			result, err := client.UploadGameCenterAchievementImage(requestCtx, locID, path)
			if err != nil {
				return fmt.Errorf("game-center achievements images upload: %w", err)
			}

			return printOutput(result, *output, *pretty)
		},
	}
}

// GameCenterAchievementImagesGetCommand returns the achievement images get subcommand.
func GameCenterAchievementImagesGetCommand() *ffcli.Command {
	fs := flag.NewFlagSet("get", flag.ExitOnError)

	imageID := fs.String("id", "", "Game Center achievement image ID")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "get",
		ShortUsage: "asc game-center achievements images get --id \"IMAGE_ID\"",
		ShortHelp:  "Get a Game Center achievement image by ID.",
		LongHelp: `Get a Game Center achievement image by ID.

Examples:
  asc game-center achievements images get --id "IMAGE_ID"`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			id := strings.TrimSpace(*imageID)
			if id == "" {
				fmt.Fprintln(os.Stderr, "Error: --id is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center achievements images get: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			resp, err := client.GetGameCenterAchievementImage(requestCtx, id)
			if err != nil {
				return fmt.Errorf("game-center achievements images get: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterAchievementImagesDeleteCommand returns the achievement images delete subcommand.
func GameCenterAchievementImagesDeleteCommand() *ffcli.Command {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)

	imageID := fs.String("id", "", "Game Center achievement image ID")
	confirm := fs.Bool("confirm", false, "Confirm deletion")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "delete",
		ShortUsage: "asc game-center achievements images delete --id \"IMAGE_ID\" --confirm",
		ShortHelp:  "Delete a Game Center achievement image.",
		LongHelp: `Delete a Game Center achievement image.

Examples:
  asc game-center achievements images delete --id "IMAGE_ID" --confirm`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			id := strings.TrimSpace(*imageID)
			if id == "" {
				fmt.Fprintln(os.Stderr, "Error: --id is required")
				return flag.ErrHelp
			}
			if !*confirm {
				fmt.Fprintln(os.Stderr, "Error: --confirm is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center achievements images delete: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			if err := client.DeleteGameCenterAchievementImage(requestCtx, id); err != nil {
				return fmt.Errorf("game-center achievements images delete: failed to delete: %w", err)
			}

			result := &asc.GameCenterAchievementImageDeleteResult{
				ID:      id,
				Deleted: true,
			}

			return printOutput(result, *output, *pretty)
		},
	}
}

// GameCenterLeaderboardSetReleasesCommand returns the leaderboard-sets releases command group.
func GameCenterLeaderboardSetReleasesCommand() *ffcli.Command {
	fs := flag.NewFlagSet("releases", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "releases",
		ShortUsage: "asc game-center leaderboard-sets releases <subcommand> [flags]",
		ShortHelp:  "Manage Game Center leaderboard set releases.",
		LongHelp: `Manage Game Center leaderboard set releases.

Releases control which Game Center details (apps) a leaderboard set is associated with.

Examples:
  asc game-center leaderboard-sets releases list --set-id "SET_ID"
  asc game-center leaderboard-sets releases create --app "APP_ID" --set-id "SET_ID"
  asc game-center leaderboard-sets releases delete --id "RELEASE_ID" --confirm`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Subcommands: []*ffcli.Command{
			GameCenterLeaderboardSetReleasesListCommand(),
			GameCenterLeaderboardSetReleasesCreateCommand(),
			GameCenterLeaderboardSetReleasesDeleteCommand(),
		},
		Exec: func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		},
	}
}

// GameCenterLeaderboardSetReleasesListCommand returns the leaderboard-sets releases list subcommand.
func GameCenterLeaderboardSetReleasesListCommand() *ffcli.Command {
	fs := flag.NewFlagSet("list", flag.ExitOnError)

	setID := fs.String("set-id", "", "Game Center leaderboard set ID")
	limit := fs.Int("limit", 0, "Maximum results per page (1-200)")
	next := fs.String("next", "", "Fetch next page using a links.next URL")
	paginate := fs.Bool("paginate", false, "Automatically fetch all pages (aggregate results)")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "list",
		ShortUsage: "asc game-center leaderboard-sets releases list --set-id \"SET_ID\"",
		ShortHelp:  "List releases for a Game Center leaderboard set.",
		LongHelp: `List releases for a Game Center leaderboard set.

Examples:
  asc game-center leaderboard-sets releases list --set-id "SET_ID"
  asc game-center leaderboard-sets releases list --set-id "SET_ID" --limit 50
  asc game-center leaderboard-sets releases list --set-id "SET_ID" --paginate`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			if *limit != 0 && (*limit < 1 || *limit > 200) {
				return fmt.Errorf("game-center leaderboard-sets releases list: --limit must be between 1 and 200")
			}
			if err := validateNextURL(*next); err != nil {
				return fmt.Errorf("game-center leaderboard-sets releases list: %w", err)
			}

			id := strings.TrimSpace(*setID)
			if id == "" && strings.TrimSpace(*next) == "" {
				fmt.Fprintln(os.Stderr, "Error: --set-id is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center leaderboard-sets releases list: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			opts := []asc.GCLeaderboardSetReleasesOption{
				asc.WithGCLeaderboardSetReleasesLimit(*limit),
				asc.WithGCLeaderboardSetReleasesNextURL(*next),
			}

			if *paginate {
				paginateOpts := append(opts, asc.WithGCLeaderboardSetReleasesLimit(200))
				firstPage, err := client.GetGameCenterLeaderboardSetReleases(requestCtx, id, paginateOpts...)
				if err != nil {
					return fmt.Errorf("game-center leaderboard-sets releases list: failed to fetch: %w", err)
				}

				resp, err := asc.PaginateAll(requestCtx, firstPage, func(ctx context.Context, nextURL string) (asc.PaginatedResponse, error) {
					return client.GetGameCenterLeaderboardSetReleases(ctx, id, asc.WithGCLeaderboardSetReleasesNextURL(nextURL))
				})
				if err != nil {
					return fmt.Errorf("game-center leaderboard-sets releases list: %w", err)
				}

				return printOutput(resp, *output, *pretty)
			}

			resp, err := client.GetGameCenterLeaderboardSetReleases(requestCtx, id, opts...)
			if err != nil {
				return fmt.Errorf("game-center leaderboard-sets releases list: failed to fetch: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterLeaderboardSetReleasesCreateCommand returns the leaderboard-sets releases create subcommand.
func GameCenterLeaderboardSetReleasesCreateCommand() *ffcli.Command {
	fs := flag.NewFlagSet("create", flag.ExitOnError)

	appID := fs.String("app", "", "App Store Connect app ID (or ASC_APP_ID env)")
	setID := fs.String("set-id", "", "Game Center leaderboard set ID")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "create",
		ShortUsage: "asc game-center leaderboard-sets releases create --app \"APP_ID\" --set-id \"SET_ID\"",
		ShortHelp:  "Create a release for a Game Center leaderboard set.",
		LongHelp: `Create a release for a Game Center leaderboard set.

This associates the leaderboard set with the app's Game Center detail.

Examples:
  asc game-center leaderboard-sets releases create --app "APP_ID" --set-id "SET_ID"`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			resolvedAppID := resolveAppID(*appID)
			if resolvedAppID == "" {
				fmt.Fprintln(os.Stderr, "Error: --app is required (or set ASC_APP_ID)")
				return flag.ErrHelp
			}

			id := strings.TrimSpace(*setID)
			if id == "" {
				fmt.Fprintln(os.Stderr, "Error: --set-id is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center leaderboard-sets releases create: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			// Get Game Center detail ID first
			gcDetailID, err := client.GetGameCenterDetailID(requestCtx, resolvedAppID)
			if err != nil {
				return fmt.Errorf("game-center leaderboard-sets releases create: failed to get Game Center detail: %w", err)
			}

			resp, err := client.CreateGameCenterLeaderboardSetRelease(requestCtx, gcDetailID, id)
			if err != nil {
				return fmt.Errorf("game-center leaderboard-sets releases create: failed to create: %w", err)
			}

			return printOutput(resp, *output, *pretty)
		},
	}
}

// GameCenterLeaderboardSetReleasesDeleteCommand returns the leaderboard-sets releases delete subcommand.
func GameCenterLeaderboardSetReleasesDeleteCommand() *ffcli.Command {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)

	releaseID := fs.String("id", "", "Game Center leaderboard set release ID")
	confirm := fs.Bool("confirm", false, "Confirm deletion")
	output := fs.String("output", "json", "Output format: json (default), table, markdown")
	pretty := fs.Bool("pretty", false, "Pretty-print JSON output")

	return &ffcli.Command{
		Name:       "delete",
		ShortUsage: "asc game-center leaderboard-sets releases delete --id \"RELEASE_ID\" --confirm",
		ShortHelp:  "Delete a Game Center leaderboard set release.",
		LongHelp: `Delete a Game Center leaderboard set release.

Examples:
  asc game-center leaderboard-sets releases delete --id "RELEASE_ID" --confirm`,
		FlagSet:   fs,
		UsageFunc: DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			id := strings.TrimSpace(*releaseID)
			if id == "" {
				fmt.Fprintln(os.Stderr, "Error: --id is required")
				return flag.ErrHelp
			}
			if !*confirm {
				fmt.Fprintln(os.Stderr, "Error: --confirm is required")
				return flag.ErrHelp
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("game-center leaderboard-sets releases delete: %w", err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			if err := client.DeleteGameCenterLeaderboardSetRelease(requestCtx, id); err != nil {
				return fmt.Errorf("game-center leaderboard-sets releases delete: failed to delete: %w", err)
			}

			result := &asc.GameCenterLeaderboardSetReleaseDeleteResult{
				ID:      id,
				Deleted: true,
			}

			return printOutput(result, *output, *pretty)
		},
	}
}

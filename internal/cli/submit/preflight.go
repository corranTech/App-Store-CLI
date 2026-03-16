package submit

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

// checkResult represents the outcome of a single preflight check.
type checkResult struct {
	Name    string `json:"name"`
	Passed  bool   `json:"passed"`
	Message string `json:"message,omitempty"`
	Hint    string `json:"hint,omitempty"`
}

// preflightResult aggregates all preflight check outcomes.
type preflightResult struct {
	AppID     string        `json:"app_id"`
	Version   string        `json:"version"`
	Platform  string        `json:"platform"`
	Checks    []checkResult `json:"checks"`
	PassCount int           `json:"pass_count"`
	FailCount int           `json:"fail_count"`
}

// SubmitPreflightCommand returns the "submit preflight" subcommand.
func SubmitPreflightCommand() *ffcli.Command {
	fs := flag.NewFlagSet("submit preflight", flag.ExitOnError)

	appID := fs.String("app", "", "App Store Connect app ID (or ASC_APP_ID)")
	version := fs.String("version", "", "App Store version string")
	platform := fs.String("platform", "IOS", "Platform: IOS, MAC_OS, TV_OS, VISION_OS")
	output := shared.BindOutputFlags(fs)

	return &ffcli.Command{
		Name:       "preflight",
		ShortUsage: "asc submit preflight [flags]",
		ShortHelp:  "Check submission readiness without submitting.",
		LongHelp: `Check all submission requirements upfront and report issues with fix commands.

Examples:
  asc submit preflight --app "123456789" --version "1.0"
  asc submit preflight --app "123456789" --version "1.0" --platform TV_OS
  asc submit preflight --app "123456789" --version "2.0" --output json`,
		FlagSet:   fs,
		UsageFunc: shared.DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			resolvedAppID := shared.ResolveAppID(*appID)
			if resolvedAppID == "" {
				fmt.Fprintln(os.Stderr, "Error: --app is required (or set ASC_APP_ID)")
				return flag.ErrHelp
			}
			if strings.TrimSpace(*version) == "" {
				fmt.Fprintln(os.Stderr, "Error: --version is required")
				return flag.ErrHelp
			}

			normalizedPlatform, err := shared.NormalizeAppStoreVersionPlatform(*platform)
			if err != nil {
				return shared.UsageError(err.Error())
			}

			client, err := shared.GetASCClient()
			if err != nil {
				return fmt.Errorf("submit preflight: %w", err)
			}

			result := runPreflight(ctx, client, resolvedAppID, strings.TrimSpace(*version), normalizedPlatform)

			outputFormat := *output.Output
			if outputFormat == "" || outputFormat == "text" {
				printPreflightText(result)
				if result.FailCount > 0 {
					return fmt.Errorf("submit preflight: %d issue(s) found", result.FailCount)
				}
				return nil
			}

			if err := shared.PrintOutput(result, outputFormat, *output.Pretty); err != nil {
				return err
			}
			if result.FailCount > 0 {
				return fmt.Errorf("submit preflight: %d issue(s) found", result.FailCount)
			}
			return nil
		},
	}
}

// runPreflight executes all checks and collects results. Individual check
// failures do not short-circuit — every check runs independently.
func runPreflight(ctx context.Context, client *asc.Client, appID, version, platform string) *preflightResult {
	result := &preflightResult{
		AppID:    appID,
		Version:  version,
		Platform: platform,
	}

	// 1. Version exists — resolve version ID
	versionID, versionCheck := checkVersionExists(ctx, client, appID, version, platform)
	result.Checks = append(result.Checks, versionCheck)

	if versionID == "" {
		// Cannot continue without a version ID.
		tallyCounts(result)
		return result
	}

	// 2. Build attached
	result.Checks = append(result.Checks, checkBuildAttached(ctx, client, versionID))

	// 3. Age rating (requires appInfoID)
	appInfoID, appInfoCheck := resolveAppInfoID(ctx, client, appID)
	if appInfoID != "" {
		result.Checks = append(result.Checks, checkAgeRating(ctx, client, appInfoID, appID))
	} else {
		result.Checks = append(result.Checks, appInfoCheck)
	}

	// 4. Content rights
	result.Checks = append(result.Checks, checkContentRights(ctx, client, appID))

	// 5. Primary category (requires appInfoID)
	if appInfoID != "" {
		result.Checks = append(result.Checks, checkPrimaryCategory(ctx, client, appInfoID, appID))
	} else {
		result.Checks = append(result.Checks, checkResult{
			Name:    "Primary category",
			Passed:  false,
			Message: "Could not resolve app info to check primary category",
		})
	}

	// 6 & 7. Localizations + screenshots
	locChecks := checkLocalizations(ctx, client, versionID, appID)
	result.Checks = append(result.Checks, locChecks...)

	tallyCounts(result)
	return result
}

func tallyCounts(result *preflightResult) {
	result.PassCount = 0
	result.FailCount = 0
	for _, c := range result.Checks {
		if c.Passed {
			result.PassCount++
		} else {
			result.FailCount++
		}
	}
}

// --- Individual checks ---

func checkVersionExists(ctx context.Context, client *asc.Client, appID, version, platform string) (string, checkResult) {
	resolveCtx, cancel := shared.ContextWithTimeout(ctx)
	defer cancel()

	versionID, err := shared.ResolveAppStoreVersionID(resolveCtx, client, appID, version, platform)
	if err != nil {
		return "", checkResult{
			Name:    "Version exists",
			Passed:  false,
			Message: fmt.Sprintf("Version %s not found for platform %s: %v", version, platform, err),
			Hint:    fmt.Sprintf("asc versions create --app %s --version %s --platform %s", appID, version, platform),
		}
	}
	return versionID, checkResult{
		Name:    "Version exists",
		Passed:  true,
		Message: fmt.Sprintf("Version %s found", version),
	}
}

func checkBuildAttached(ctx context.Context, client *asc.Client, versionID string) checkResult {
	buildCtx, cancel := shared.ContextWithTimeout(ctx)
	defer cancel()

	buildResp, err := client.GetAppStoreVersionBuild(buildCtx, versionID)
	if err != nil {
		if asc.IsNotFound(err) {
			return checkResult{
				Name:    "Build attached",
				Passed:  false,
				Message: "No build attached to this version",
				Hint:    fmt.Sprintf("asc submit create --version-id %s --build BUILD_ID --confirm", versionID),
			}
		}
		return checkResult{
			Name:    "Build attached",
			Passed:  false,
			Message: fmt.Sprintf("Failed to check build: %v", err),
		}
	}

	buildVersion := buildResp.Data.Attributes.Version
	if strings.TrimSpace(buildResp.Data.ID) == "" {
		return checkResult{
			Name:    "Build attached",
			Passed:  false,
			Message: "No build attached to this version",
			Hint:    fmt.Sprintf("asc submit create --version-id %s --build BUILD_ID --confirm", versionID),
		}
	}

	msg := "Build attached"
	if buildVersion != "" {
		msg = fmt.Sprintf("Build attached (build %s)", buildVersion)
	}
	return checkResult{
		Name:    "Build attached",
		Passed:  true,
		Message: msg,
	}
}

func resolveAppInfoID(ctx context.Context, client *asc.Client, appID string) (string, checkResult) {
	infoCtx, cancel := shared.ContextWithTimeout(ctx)
	defer cancel()

	infos, err := client.GetAppInfos(infoCtx, appID)
	if err != nil {
		return "", checkResult{
			Name:    "App info",
			Passed:  false,
			Message: fmt.Sprintf("Failed to fetch app info: %v", err),
		}
	}
	if len(infos.Data) == 0 {
		return "", checkResult{
			Name:    "App info",
			Passed:  false,
			Message: "No app info records found",
		}
	}

	// Prefer the app info in PREPARE_FOR_SUBMISSION state — that is the
	// editable draft the submission will use.
	candidates := asc.AppInfoCandidates(infos.Data)
	for _, c := range candidates {
		if strings.EqualFold(c.State, "PREPARE_FOR_SUBMISSION") {
			return c.ID, checkResult{}
		}
	}

	// Fall back to the first app info if none is in the expected state.
	return infos.Data[0].ID, checkResult{}
}

func checkAgeRating(ctx context.Context, client *asc.Client, appInfoID, appID string) checkResult {
	ratingCtx, cancel := shared.ContextWithTimeout(ctx)
	defer cancel()

	resp, err := client.GetAgeRatingDeclarationForAppInfo(ratingCtx, appInfoID)
	if err != nil {
		if asc.IsNotFound(err) {
			return checkResult{
				Name:    "Age rating",
				Passed:  false,
				Message: "Age rating declaration not found",
				Hint:    fmt.Sprintf("asc age-rating set --app %s --gambling false --violence-realistic NONE ...", appID),
			}
		}
		return checkResult{
			Name:    "Age rating",
			Passed:  false,
			Message: fmt.Sprintf("Failed to check age rating: %v", err),
		}
	}

	attrs := resp.Data.Attributes
	missing := ageRatingMissingFields(attrs)
	if len(missing) > 0 {
		return checkResult{
			Name:    "Age rating",
			Passed:  false,
			Message: fmt.Sprintf("Age rating incomplete (missing: %s)", strings.Join(missing, ", ")),
			Hint:    fmt.Sprintf("asc age-rating set --app %s --gambling false --violence-realistic NONE ...", appID),
		}
	}

	return checkResult{
		Name:    "Age rating",
		Passed:  true,
		Message: "Age rating complete",
	}
}

// ageRatingMissingFields checks which required age rating fields are nil.
// Apple requires all boolean and enum content descriptors to be explicitly set.
func ageRatingMissingFields(a asc.AgeRatingDeclarationAttributes) []string {
	var missing []string

	// Boolean descriptors
	if a.Gambling == nil {
		missing = append(missing, "gambling")
	}
	if a.LootBox == nil {
		missing = append(missing, "lootBox")
	}
	if a.UnrestrictedWebAccess == nil {
		missing = append(missing, "unrestrictedWebAccess")
	}

	// Enum descriptors
	if a.AlcoholTobaccoOrDrugUseOrReferences == nil {
		missing = append(missing, "alcoholTobaccoOrDrugUseOrReferences")
	}
	if a.GamblingSimulated == nil {
		missing = append(missing, "gamblingSimulated")
	}
	if a.HorrorOrFearThemes == nil {
		missing = append(missing, "horrorOrFearThemes")
	}
	if a.MatureOrSuggestiveThemes == nil {
		missing = append(missing, "matureOrSuggestiveThemes")
	}
	if a.ProfanityOrCrudeHumor == nil {
		missing = append(missing, "profanityOrCrudeHumor")
	}
	if a.SexualContentGraphicAndNudity == nil {
		missing = append(missing, "sexualContentGraphicAndNudity")
	}
	if a.SexualContentOrNudity == nil {
		missing = append(missing, "sexualContentOrNudity")
	}
	if a.ViolenceCartoonOrFantasy == nil {
		missing = append(missing, "violenceCartoonOrFantasy")
	}
	if a.ViolenceRealistic == nil {
		missing = append(missing, "violenceRealistic")
	}
	if a.ViolenceRealisticProlongedGraphicOrSadistic == nil {
		missing = append(missing, "violenceRealisticProlongedGraphicOrSadistic")
	}

	return missing
}

func checkContentRights(ctx context.Context, client *asc.Client, appID string) checkResult {
	appCtx, cancel := shared.ContextWithTimeout(ctx)
	defer cancel()

	appResp, err := client.GetApp(appCtx, appID)
	if err != nil {
		return checkResult{
			Name:    "Content rights",
			Passed:  false,
			Message: fmt.Sprintf("Failed to fetch app: %v", err),
		}
	}

	if appResp.Data.Attributes.ContentRightsDeclaration == nil {
		return checkResult{
			Name:    "Content rights",
			Passed:  false,
			Message: "Content rights declaration not set",
			Hint:    fmt.Sprintf("asc apps update --app %s --content-rights DOES_NOT_USE_THIRD_PARTY_CONTENT", appID),
		}
	}

	return checkResult{
		Name:    "Content rights",
		Passed:  true,
		Message: "Content rights set",
	}
}

func checkPrimaryCategory(ctx context.Context, client *asc.Client, appInfoID, appID string) checkResult {
	catCtx, cancel := shared.ContextWithTimeout(ctx)
	defer cancel()

	catResp, err := client.GetAppInfoPrimaryCategory(catCtx, appInfoID)
	if err != nil {
		if asc.IsNotFound(err) {
			return checkResult{
				Name:    "Primary category",
				Passed:  false,
				Message: "Primary category not set",
				Hint:    fmt.Sprintf("asc app-setup categories set --app %s --primary SPORTS", appID),
			}
		}
		return checkResult{
			Name:    "Primary category",
			Passed:  false,
			Message: fmt.Sprintf("Failed to check primary category: %v", err),
		}
	}

	if strings.TrimSpace(catResp.Data.ID) == "" {
		return checkResult{
			Name:    "Primary category",
			Passed:  false,
			Message: "Primary category not set",
			Hint:    fmt.Sprintf("asc app-setup categories set --app %s --primary SPORTS", appID),
		}
	}

	return checkResult{
		Name:    "Primary category",
		Passed:  true,
		Message: fmt.Sprintf("Primary category set (%s)", catResp.Data.ID),
	}
}

func checkLocalizations(ctx context.Context, client *asc.Client, versionID, appID string) []checkResult {
	locCtx, cancel := shared.ContextWithTimeout(ctx)
	defer cancel()

	localizations, err := client.GetAppStoreVersionLocalizations(locCtx, versionID, asc.WithAppStoreVersionLocalizationsLimit(200))
	if err != nil {
		return []checkResult{{
			Name:    "Localizations",
			Passed:  false,
			Message: fmt.Sprintf("Failed to fetch localizations: %v", err),
		}}
	}

	if len(localizations.Data) == 0 {
		return []checkResult{{
			Name:    "Localizations",
			Passed:  false,
			Message: "No localizations found for this version",
		}}
	}

	var checks []checkResult

	// Check description
	hasDescription := false
	for _, loc := range localizations.Data {
		if strings.TrimSpace(loc.Attributes.Description) != "" {
			hasDescription = true
			break
		}
	}
	if hasDescription {
		checks = append(checks, checkResult{
			Name:    "Description",
			Passed:  true,
			Message: "Description present",
		})
	} else {
		checks = append(checks, checkResult{
			Name:    "Description",
			Passed:  false,
			Message: "No localization has a description",
			Hint:    "asc metadata push --version-id " + versionID,
		})
	}

	// Check keywords
	hasKeywords := false
	for _, loc := range localizations.Data {
		if strings.TrimSpace(loc.Attributes.Keywords) != "" {
			hasKeywords = true
			break
		}
	}
	if hasKeywords {
		checks = append(checks, checkResult{
			Name:    "Keywords",
			Passed:  true,
			Message: "Keywords present",
		})
	} else {
		checks = append(checks, checkResult{
			Name:    "Keywords",
			Passed:  false,
			Message: "No localization has keywords",
			Hint:    "asc metadata push --version-id " + versionID,
		})
	}

	// Check screenshots — look at the first localization that has screenshot sets.
	screenshotCheck := checkScreenshots(ctx, client, localizations.Data, appID)
	checks = append(checks, screenshotCheck)

	return checks
}

func checkScreenshots(ctx context.Context, client *asc.Client, localizations []asc.Resource[asc.AppStoreVersionLocalizationAttributes], appID string) checkResult {
	var fetchErrors []string
	for _, loc := range localizations {
		locID := loc.ID

		ssCtx, cancel := shared.ContextWithTimeout(ctx)
		sets, err := client.GetAppStoreVersionLocalizationScreenshotSets(ssCtx, locID, asc.WithAppStoreVersionLocalizationScreenshotSetsLimit(50))
		cancel()

		if err != nil {
			fetchErrors = append(fetchErrors, fmt.Sprintf("%s: %v", loc.Attributes.Locale, err))
			continue
		}

		if len(sets.Data) > 0 {
			return checkResult{
				Name:    "Screenshots",
				Passed:  true,
				Message: fmt.Sprintf("Screenshots uploaded (%d set(s) in %s)", len(sets.Data), loc.Attributes.Locale),
			}
		}
	}

	// If every localization failed to fetch, report the errors.
	if len(fetchErrors) > 0 && len(fetchErrors) == len(localizations) {
		return checkResult{
			Name:    "Screenshots",
			Passed:  false,
			Message: fmt.Sprintf("Failed to fetch screenshots for all localizations: %s", strings.Join(fetchErrors, "; ")),
		}
	}

	locID := ""
	if len(localizations) > 0 {
		locID = localizations[0].ID
	}
	hint := "asc screenshots upload --version-localization LOC_ID --path ./screenshots --device-type APPLE_TV"
	if locID != "" {
		hint = fmt.Sprintf("asc screenshots upload --version-localization %s --path ./screenshots --device-type APPLE_TV", locID)
	}

	msg := "No screenshots uploaded"
	if len(fetchErrors) > 0 {
		msg = fmt.Sprintf("No screenshots uploaded (some locales failed: %s)", strings.Join(fetchErrors, "; "))
	}

	return checkResult{
		Name:    "Screenshots",
		Passed:  false,
		Message: msg,
		Hint:    hint,
	}
}

// --- Text output ---

func printPreflightText(result *preflightResult) {
	header := fmt.Sprintf("Preflight check for app %s v%s (%s)", result.AppID, result.Version, result.Platform)
	fmt.Fprintln(os.Stderr, header)
	fmt.Fprintln(os.Stderr, strings.Repeat("\u2500", len(header)))

	for _, c := range result.Checks {
		if c.Passed {
			fmt.Fprintf(os.Stderr, "\u2713 %s\n", c.Message)
		} else {
			fmt.Fprintf(os.Stderr, "\u2717 %s\n", c.Message)
			if c.Hint != "" {
				fmt.Fprintf(os.Stderr, "  Hint: %s\n", c.Hint)
			}
		}
	}

	fmt.Fprintln(os.Stderr)
	if result.FailCount == 0 {
		fmt.Fprintln(os.Stderr, "Result: All checks passed. Ready to submit.")
	} else {
		fmt.Fprintf(os.Stderr, "Result: %d issue(s) found. Fix them before submitting.\n", result.FailCount)
	}
}

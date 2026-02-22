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

// SubmitValidateIssue represents one pre-submission validation issue.
type SubmitValidateIssue struct {
	Check    string `json:"check"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

// SubmitValidateResult is the structured result for submit validate.
type SubmitValidateResult struct {
	AppID      string                `json:"appId"`
	VersionID  string                `json:"versionId"`
	Platform   string                `json:"platform"`
	Issues     []SubmitValidateIssue `json:"issues"`
	ErrorCount int                   `json:"errorCount"`
	WarnCount  int                   `json:"warningCount"`
	Ready      bool                  `json:"ready"`
}

func (r *SubmitValidateResult) addError(check, message string) {
	r.Issues = append(r.Issues, SubmitValidateIssue{Check: check, Severity: "error", Message: message})
	r.ErrorCount++
}

func (r *SubmitValidateResult) addWarning(check, message string) {
	r.Issues = append(r.Issues, SubmitValidateIssue{Check: check, Severity: "warning", Message: message})
	r.WarnCount++
}

// SubmitValidateCommand returns the submit validate subcommand.
func SubmitValidateCommand() *ffcli.Command {
	fs := flag.NewFlagSet("submit validate", flag.ExitOnError)

	appID := fs.String("app", "", "App Store Connect app ID (or ASC_APP_ID)")
	version := fs.String("version", "", "App Store version string")
	versionID := fs.String("version-id", "", "App Store version ID")
	platform := fs.String("platform", "IOS", "Platform: IOS, MAC_OS, TV_OS, VISION_OS")
	output := shared.BindOutputFlags(fs)

	return &ffcli.Command{
		Name:       "validate",
		ShortUsage: "asc submit validate [flags]",
		ShortHelp:  "Check submission readiness without submitting.",
		LongHelp: `Check submission readiness without submitting.

Performs live API checks to detect common submission blockers:
  - Version exists and is in an editable state
  - Build is attached to the version
  - App info localizations have name set
  - Version localizations have description and keywords
  - Screenshots exist for each localization
  - Privacy policy URL is set
  - Age rating declaration exists

Examples:
  asc submit validate --app "123456789" --version "1.0.0"
  asc submit validate --app "123456789" --version-id "VERSION_ID"
  asc submit validate --app "123456789" --version "1.0.0" --output table`,
		FlagSet:   fs,
		UsageFunc: shared.DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			if strings.TrimSpace(*version) == "" && strings.TrimSpace(*versionID) == "" {
				fmt.Fprintln(os.Stderr, "Error: --version or --version-id is required")
				return flag.ErrHelp
			}
			if strings.TrimSpace(*version) != "" && strings.TrimSpace(*versionID) != "" {
				return shared.UsageError("--version and --version-id are mutually exclusive")
			}

			resolvedAppID := shared.ResolveAppID(*appID)
			if resolvedAppID == "" {
				fmt.Fprintln(os.Stderr, "Error: --app is required (or set ASC_APP_ID)")
				return flag.ErrHelp
			}

			normalizedPlatform, err := shared.NormalizeAppStoreVersionPlatform(*platform)
			if err != nil {
				return shared.UsageError(err.Error())
			}

			client, err := shared.GetASCClient()
			if err != nil {
				return fmt.Errorf("submit validate: %w", err)
			}

			requestCtx, cancel := shared.ContextWithTimeout(ctx)
			defer cancel()

			resolvedVersionID := strings.TrimSpace(*versionID)
			if resolvedVersionID == "" {
				resolvedVersionID, err = shared.ResolveAppStoreVersionID(requestCtx, client, resolvedAppID, strings.TrimSpace(*version), normalizedPlatform)
				if err != nil {
					return fmt.Errorf("submit validate: %w", err)
				}
			}

			result := runValidation(requestCtx, client, resolvedAppID, resolvedVersionID, normalizedPlatform)

			if err := shared.PrintOutputWithRenderers(
				result,
				*output.Output,
				*output.Pretty,
				func() error { return printValidateTable(result) },
				func() error { return printValidateMarkdown(result) },
			); err != nil {
				return err
			}

			if result.ErrorCount > 0 {
				return shared.NewReportedError(fmt.Errorf("submit validate: %d error(s) found", result.ErrorCount))
			}
			return nil
		},
	}
}

func runValidation(ctx context.Context, client *asc.Client, appID, versionID, platform string) *SubmitValidateResult {
	result := &SubmitValidateResult{
		AppID:     appID,
		VersionID: versionID,
		Platform:  platform,
		Issues:    make([]SubmitValidateIssue, 0),
	}

	// 1. Check version exists and state
	versionResp, err := client.GetAppStoreVersion(ctx, versionID)
	if err != nil {
		result.addError("version", fmt.Sprintf("failed to fetch version: %v", err))
		return result
	}
	state := shared.ResolveAppStoreVersionState(versionResp.Data.Attributes)
	if !isEditableState(state) {
		result.addError("version_state", fmt.Sprintf("version is in non-editable state: %s", state))
	}

	// 2. Check build attached
	checkBuildAttached(ctx, client, versionID, result)

	// 3. Check version localizations (description, keywords)
	checkVersionLocalizations(ctx, client, versionID, result)

	// 4. Check app info localizations (name, privacy policy URL)
	checkAppInfoLocalizations(ctx, client, appID, result)

	// 5. Check age rating
	checkAgeRating(ctx, client, versionID, result)

	result.Ready = result.ErrorCount == 0
	return result
}

func isEditableState(state string) bool {
	switch strings.ToUpper(state) {
	case "PREPARE_FOR_SUBMISSION", "DEVELOPER_REJECTED", "REJECTED",
		"METADATA_REJECTED", "INVALID_BINARY", "DEVELOPER_REMOVED_FROM_SALE":
		return true
	default:
		return false
	}
}

func checkBuildAttached(ctx context.Context, client *asc.Client, versionID string, result *SubmitValidateResult) {
	_, err := client.GetAppStoreVersionBuild(ctx, versionID)
	if err != nil {
		if asc.IsNotFound(err) {
			result.addError("build", "no build attached to this version")
		} else {
			result.addWarning("build", fmt.Sprintf("unable to check build: %v", err))
		}
	}
}

func checkVersionLocalizations(ctx context.Context, client *asc.Client, versionID string, result *SubmitValidateResult) {
	resp, err := client.GetAppStoreVersionLocalizations(ctx, versionID, asc.WithAppStoreVersionLocalizationsLimit(200))
	if err != nil {
		result.addWarning("version_localizations", fmt.Sprintf("unable to fetch: %v", err))
		return
	}

	if len(resp.Data) == 0 {
		result.addError("version_localizations", "no version localizations found")
		return
	}

	for _, loc := range resp.Data {
		locale := loc.Attributes.Locale
		if strings.TrimSpace(loc.Attributes.Description) == "" {
			result.addError("description", fmt.Sprintf("locale %s: description is empty", locale))
		}
		if strings.TrimSpace(loc.Attributes.Keywords) == "" {
			result.addWarning("keywords", fmt.Sprintf("locale %s: keywords are empty", locale))
		}

		// Check screenshots for this localization
		checkScreenshots(ctx, client, loc.ID, locale, result)
	}
}

func checkScreenshots(ctx context.Context, client *asc.Client, localizationID, locale string, result *SubmitValidateResult) {
	sets, err := client.GetAppScreenshotSets(ctx, localizationID)
	if err != nil {
		result.addWarning("screenshots", fmt.Sprintf("locale %s: unable to check screenshots: %v", locale, err))
		return
	}

	if len(sets.Data) == 0 {
		result.addError("screenshots", fmt.Sprintf("locale %s: no screenshot sets found", locale))
		return
	}

	for _, set := range sets.Data {
		screenshots, err := client.GetAppScreenshots(ctx, set.ID)
		if err != nil {
			result.addWarning("screenshots", fmt.Sprintf("locale %s (%s): unable to check: %v", locale, set.Attributes.ScreenshotDisplayType, err))
			continue
		}
		if len(screenshots.Data) == 0 {
			result.addWarning("screenshots", fmt.Sprintf("locale %s (%s): empty screenshot set", locale, set.Attributes.ScreenshotDisplayType))
		}
	}
}

func checkAppInfoLocalizations(ctx context.Context, client *asc.Client, appID string, result *SubmitValidateResult) {
	appInfoResp, err := client.GetAppInfos(ctx, appID)
	if err != nil {
		result.addWarning("app_info", fmt.Sprintf("unable to fetch app info: %v", err))
		return
	}
	if len(appInfoResp.Data) == 0 {
		result.addError("app_info", "no app info records found")
		return
	}

	appInfoID := appInfoResp.Data[0].ID

	locs, err := client.GetAppInfoLocalizations(ctx, appInfoID, asc.WithAppInfoLocalizationsLimit(200))
	if err != nil {
		result.addWarning("app_info_localizations", fmt.Sprintf("unable to fetch: %v", err))
		return
	}

	if len(locs.Data) == 0 {
		result.addError("app_info_localizations", "no app info localizations found")
		return
	}

	for _, loc := range locs.Data {
		locale := loc.Attributes.Locale
		if strings.TrimSpace(loc.Attributes.Name) == "" {
			result.addError("name", fmt.Sprintf("locale %s: app name is empty", locale))
		}
		if strings.TrimSpace(loc.Attributes.PrivacyPolicyURL) == "" {
			result.addWarning("privacy_policy_url", fmt.Sprintf("locale %s: privacy policy URL is empty", locale))
		}
	}
}

func checkAgeRating(ctx context.Context, client *asc.Client, versionID string, result *SubmitValidateResult) {
	_, err := client.GetAgeRatingDeclarationForAppStoreVersion(ctx, versionID)
	if err != nil {
		if asc.IsNotFound(err) {
			result.addError("age_rating", "no age rating declaration found")
		} else {
			result.addWarning("age_rating", fmt.Sprintf("unable to check: %v", err))
		}
	}
}

func printValidateTable(result *SubmitValidateResult) error {
	if result.Ready {
		fmt.Println("Status: READY")
	} else {
		fmt.Println("Status: NOT READY")
	}
	fmt.Printf("App: %s  Version: %s  Platform: %s\n", result.AppID, result.VersionID, result.Platform)
	fmt.Printf("Errors: %d  Warnings: %d\n\n", result.ErrorCount, result.WarnCount)

	rows := make([][]string, 0, len(result.Issues))
	for _, issue := range result.Issues {
		rows = append(rows, []string{issue.Check, issue.Severity, issue.Message})
	}
	if len(rows) == 0 {
		rows = append(rows, []string{"all", "info", "no issues found"})
	}
	asc.RenderTable([]string{"check", "severity", "message"}, rows)
	return nil
}

func printValidateMarkdown(result *SubmitValidateResult) error {
	status := "READY"
	if !result.Ready {
		status = "NOT READY"
	}
	fmt.Printf("**Status:** %s\n\n", status)
	fmt.Printf("**App:** %s  **Version:** %s  **Platform:** %s\n\n", result.AppID, result.VersionID, result.Platform)
	fmt.Printf("**Errors:** %d  **Warnings:** %d\n\n", result.ErrorCount, result.WarnCount)

	rows := make([][]string, 0, len(result.Issues))
	for _, issue := range result.Issues {
		rows = append(rows, []string{issue.Check, issue.Severity, issue.Message})
	}
	if len(rows) == 0 {
		rows = append(rows, []string{"all", "info", "no issues found"})
	}
	asc.RenderMarkdown([]string{"check", "severity", "message"}, rows)
	return nil
}

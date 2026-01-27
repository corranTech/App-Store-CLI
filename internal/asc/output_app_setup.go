package asc

import (
	"fmt"
	"os"
	"text/tabwriter"
)

// AppSetupInfoResult represents CLI output for app-setup info updates.
type AppSetupInfoResult struct {
	AppID               string                       `json:"appId"`
	App                 *AppResponse                 `json:"app,omitempty"`
	AppInfoLocalization *AppInfoLocalizationResponse `json:"appInfoLocalization,omitempty"`
}

func printAppSetupInfoResultTable(result *AppSetupInfoResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Resource\tID\tLocale\tName\tSubtitle\tBundle ID\tPrimary Locale\tPrivacy Policy URL")
	if result.App != nil {
		attrs := result.App.Data.Attributes
		fmt.Fprintf(w, "app\t%s\t\t\t\t%s\t%s\t\n", result.App.Data.ID, attrs.BundleID, attrs.PrimaryLocale)
	}
	if result.AppInfoLocalization != nil {
		attrs := result.AppInfoLocalization.Data.Attributes
		fmt.Fprintf(
			w,
			"appInfoLocalization\t%s\t%s\t%s\t%s\t\t\t%s\n",
			result.AppInfoLocalization.Data.ID,
			attrs.Locale,
			compactWhitespace(attrs.Name),
			compactWhitespace(attrs.Subtitle),
			attrs.PrivacyPolicyURL,
		)
	}
	return w.Flush()
}

func printAppSetupInfoResultMarkdown(result *AppSetupInfoResult) error {
	fmt.Fprintln(os.Stdout, "| Resource | ID | Locale | Name | Subtitle | Bundle ID | Primary Locale | Privacy Policy URL |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- | --- | --- | --- | --- |")
	if result.App != nil {
		attrs := result.App.Data.Attributes
		fmt.Fprintf(
			os.Stdout,
			"| app | %s |  |  |  | %s | %s |  |\n",
			escapeMarkdown(result.App.Data.ID),
			escapeMarkdown(attrs.BundleID),
			escapeMarkdown(attrs.PrimaryLocale),
		)
	}
	if result.AppInfoLocalization != nil {
		attrs := result.AppInfoLocalization.Data.Attributes
		fmt.Fprintf(
			os.Stdout,
			"| appInfoLocalization | %s | %s | %s | %s |  |  | %s |\n",
			escapeMarkdown(result.AppInfoLocalization.Data.ID),
			escapeMarkdown(attrs.Locale),
			escapeMarkdown(attrs.Name),
			escapeMarkdown(attrs.Subtitle),
			escapeMarkdown(attrs.PrivacyPolicyURL),
		)
	}
	return nil
}

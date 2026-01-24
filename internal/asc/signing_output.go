package asc

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
)

// BundleIDDeleteResult represents CLI output for bundle ID deletions.
type BundleIDDeleteResult struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// BundleIDCapabilityDeleteResult represents CLI output for capability deletions.
type BundleIDCapabilityDeleteResult struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

func printBundleIDsTable(resp *BundleIDsResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tName\tIdentifier\tPlatform\tSeed ID")
	for _, item := range resp.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			item.ID,
			compactWhitespace(item.Attributes.Name),
			item.Attributes.Identifier,
			item.Attributes.Platform,
			item.Attributes.SeedID,
		)
	}
	return w.Flush()
}

func printBundleIDsMarkdown(resp *BundleIDsResponse) error {
	fmt.Fprintln(os.Stdout, "| ID | Name | Identifier | Platform | Seed ID |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- | --- |")
	for _, item := range resp.Data {
		fmt.Fprintf(os.Stdout, "| %s | %s | %s | %s | %s |\n",
			item.ID,
			escapeMarkdown(item.Attributes.Name),
			escapeMarkdown(item.Attributes.Identifier),
			escapeMarkdown(string(item.Attributes.Platform)),
			escapeMarkdown(item.Attributes.SeedID),
		)
	}
	return nil
}

func printBundleIDCapabilitiesTable(resp *BundleIDCapabilitiesResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tCapability\tSettings")
	for _, item := range resp.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			item.ID,
			item.Attributes.CapabilityType,
			formatCapabilitySettings(item.Attributes.Settings),
		)
	}
	return w.Flush()
}

func printBundleIDCapabilitiesMarkdown(resp *BundleIDCapabilitiesResponse) error {
	fmt.Fprintln(os.Stdout, "| ID | Capability | Settings |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- |")
	for _, item := range resp.Data {
		fmt.Fprintf(os.Stdout, "| %s | %s | %s |\n",
			item.ID,
			escapeMarkdown(item.Attributes.CapabilityType),
			escapeMarkdown(formatCapabilitySettings(item.Attributes.Settings)),
		)
	}
	return nil
}

func printBundleIDDeleteResultTable(result *BundleIDDeleteResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tDeleted")
	fmt.Fprintf(w, "%s\t%t\n", result.ID, result.Deleted)
	return w.Flush()
}

func printBundleIDDeleteResultMarkdown(result *BundleIDDeleteResult) error {
	fmt.Fprintln(os.Stdout, "| ID | Deleted |")
	fmt.Fprintln(os.Stdout, "| --- | --- |")
	fmt.Fprintf(os.Stdout, "| %s | %t |\n",
		escapeMarkdown(result.ID),
		result.Deleted,
	)
	return nil
}

func printBundleIDCapabilityDeleteResultTable(result *BundleIDCapabilityDeleteResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tDeleted")
	fmt.Fprintf(w, "%s\t%t\n", result.ID, result.Deleted)
	return w.Flush()
}

func printBundleIDCapabilityDeleteResultMarkdown(result *BundleIDCapabilityDeleteResult) error {
	fmt.Fprintln(os.Stdout, "| ID | Deleted |")
	fmt.Fprintln(os.Stdout, "| --- | --- |")
	fmt.Fprintf(os.Stdout, "| %s | %t |\n",
		escapeMarkdown(result.ID),
		result.Deleted,
	)
	return nil
}

func formatCapabilitySettings(settings []CapabilitySetting) string {
	if len(settings) == 0 {
		return ""
	}
	payload, err := json.Marshal(settings)
	if err != nil {
		return ""
	}
	return sanitizeTerminal(string(payload))
}

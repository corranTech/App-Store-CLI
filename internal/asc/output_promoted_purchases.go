package asc

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

// PromotedPurchaseDeleteResult represents CLI output for promoted purchase deletions.
type PromotedPurchaseDeleteResult struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// AppPromotedPurchasesLinkResult represents CLI output for linking promoted purchases.
type AppPromotedPurchasesLinkResult struct {
	AppID               string   `json:"appId"`
	PromotedPurchaseIDs []string `json:"promotedPurchaseIds"`
	Action              string   `json:"action"`
}

func promotedPurchaseBool(value *bool) string {
	if value == nil {
		return ""
	}
	return strconv.FormatBool(*value)
}

func printPromotedPurchasesTable(resp *PromotedPurchasesResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tVisible For All Users\tEnabled\tState")
	for _, item := range resp.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			item.ID,
			promotedPurchaseBool(item.Attributes.VisibleForAllUsers),
			promotedPurchaseBool(item.Attributes.Enabled),
			item.Attributes.State,
		)
	}
	return w.Flush()
}

func printPromotedPurchasesMarkdown(resp *PromotedPurchasesResponse) error {
	fmt.Fprintln(os.Stdout, "| ID | Visible For All Users | Enabled | State |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- |")
	for _, item := range resp.Data {
		fmt.Fprintf(os.Stdout, "| %s | %s | %s | %s |\n",
			escapeMarkdown(item.ID),
			escapeMarkdown(promotedPurchaseBool(item.Attributes.VisibleForAllUsers)),
			escapeMarkdown(promotedPurchaseBool(item.Attributes.Enabled)),
			escapeMarkdown(item.Attributes.State),
		)
	}
	return nil
}

func printPromotedPurchaseDeleteResultTable(result *PromotedPurchaseDeleteResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tDeleted")
	fmt.Fprintf(w, "%s\t%t\n", result.ID, result.Deleted)
	return w.Flush()
}

func printPromotedPurchaseDeleteResultMarkdown(result *PromotedPurchaseDeleteResult) error {
	fmt.Fprintln(os.Stdout, "| ID | Deleted |")
	fmt.Fprintln(os.Stdout, "| --- | --- |")
	fmt.Fprintf(os.Stdout, "| %s | %t |\n", escapeMarkdown(result.ID), result.Deleted)
	return nil
}

func printAppPromotedPurchasesLinkResultTable(result *AppPromotedPurchasesLinkResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "App ID\tPromoted Purchase IDs\tAction")
	fmt.Fprintf(w, "%s\t%s\t%s\n",
		result.AppID,
		strings.Join(result.PromotedPurchaseIDs, ", "),
		result.Action,
	)
	return w.Flush()
}

func printAppPromotedPurchasesLinkResultMarkdown(result *AppPromotedPurchasesLinkResult) error {
	fmt.Fprintln(os.Stdout, "| App ID | Promoted Purchase IDs | Action |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- |")
	fmt.Fprintf(os.Stdout, "| %s | %s | %s |\n",
		escapeMarkdown(result.AppID),
		escapeMarkdown(strings.Join(result.PromotedPurchaseIDs, ", ")),
		escapeMarkdown(result.Action),
	)
	return nil
}

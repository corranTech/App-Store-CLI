package asc

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func printAppsTable(resp *AppsResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tName\tBundle ID\tSKU")
	for _, item := range resp.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			item.ID,
			compactWhitespace(item.Attributes.Name),
			item.Attributes.BundleID,
			item.Attributes.SKU,
		)
	}
	return w.Flush()
}

func printAppsMarkdown(resp *AppsResponse) error {
	fmt.Fprintln(os.Stdout, "| ID | Name | Bundle ID | SKU |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- |")
	for _, item := range resp.Data {
		fmt.Fprintf(os.Stdout, "| %s | %s | %s | %s |\n",
			item.ID,
			escapeMarkdown(item.Attributes.Name),
			escapeMarkdown(item.Attributes.BundleID),
			escapeMarkdown(item.Attributes.SKU),
		)
	}
	return nil
}

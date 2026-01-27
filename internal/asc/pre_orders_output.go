package asc

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func printEndAppAvailabilityPreOrderTable(resp *EndAppAvailabilityPreOrderResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID")
	fmt.Fprintf(w, "%s\n", resp.Data.ID)
	return w.Flush()
}

func printEndAppAvailabilityPreOrderMarkdown(resp *EndAppAvailabilityPreOrderResponse) error {
	fmt.Fprintln(os.Stdout, "| ID |")
	fmt.Fprintln(os.Stdout, "| --- |")
	fmt.Fprintf(os.Stdout, "| %s |\n", escapeMarkdown(resp.Data.ID))
	return nil
}

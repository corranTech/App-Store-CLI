package asc

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func printActorsTable(resp *ActorsResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tType\tName\tEmail\tAPI Key ID")
	for _, item := range resp.Data {
		attr := item.Attributes
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			item.ID,
			compactWhitespace(attr.ActorType),
			compactWhitespace(formatPersonName(attr.UserFirstName, attr.UserLastName)),
			compactWhitespace(attr.UserEmail),
			compactWhitespace(attr.APIKeyID),
		)
	}
	return w.Flush()
}

func printActorsMarkdown(resp *ActorsResponse) error {
	fmt.Fprintln(os.Stdout, "| ID | Type | Name | Email | API Key ID |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- | --- |")
	for _, item := range resp.Data {
		attr := item.Attributes
		fmt.Fprintf(os.Stdout, "| %s | %s | %s | %s | %s |\n",
			escapeMarkdown(item.ID),
			escapeMarkdown(attr.ActorType),
			escapeMarkdown(formatPersonName(attr.UserFirstName, attr.UserLastName)),
			escapeMarkdown(attr.UserEmail),
			escapeMarkdown(attr.APIKeyID),
		)
	}
	return nil
}

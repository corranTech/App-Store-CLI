package asc

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func printNotarySubmissionStatusTable(resp *NotarySubmissionStatusResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tSTATUS\tNAME\tCREATED")
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
		resp.Data.ID,
		resp.Data.Attributes.Status,
		compactWhitespace(resp.Data.Attributes.Name),
		resp.Data.Attributes.CreatedDate,
	)
	return w.Flush()
}

func printNotarySubmissionsListTable(resp *NotarySubmissionsListResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tSTATUS\tNAME\tCREATED")
	for _, item := range resp.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			item.ID,
			item.Attributes.Status,
			compactWhitespace(item.Attributes.Name),
			item.Attributes.CreatedDate,
		)
	}
	return w.Flush()
}

func printNotarySubmissionLogsTable(resp *NotarySubmissionLogsResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tDEVELOPER LOG URL")
	fmt.Fprintf(w, "%s\t%s\n",
		resp.Data.ID,
		resp.Data.Attributes.DeveloperLogURL,
	)
	return w.Flush()
}

func printNotarySubmissionStatusMarkdown(resp *NotarySubmissionStatusResponse) error {
	fmt.Println("| ID | Status | Name | Created |")
	fmt.Println("|----|--------|------|---------|")
	fmt.Printf("| %s | %s | %s | %s |\n",
		escapeMarkdown(resp.Data.ID),
		escapeMarkdown(string(resp.Data.Attributes.Status)),
		escapeMarkdown(resp.Data.Attributes.Name),
		escapeMarkdown(resp.Data.Attributes.CreatedDate),
	)
	return nil
}

func printNotarySubmissionsListMarkdown(resp *NotarySubmissionsListResponse) error {
	fmt.Println("| ID | Status | Name | Created |")
	fmt.Println("|----|--------|------|---------|")
	for _, item := range resp.Data {
		fmt.Printf("| %s | %s | %s | %s |\n",
			escapeMarkdown(item.ID),
			escapeMarkdown(string(item.Attributes.Status)),
			escapeMarkdown(item.Attributes.Name),
			escapeMarkdown(item.Attributes.CreatedDate),
		)
	}
	return nil
}

func printNotarySubmissionLogsMarkdown(resp *NotarySubmissionLogsResponse) error {
	fmt.Println("| ID | Developer Log URL |")
	fmt.Println("|----|-------------------|")
	fmt.Printf("| %s | %s |\n",
		escapeMarkdown(resp.Data.ID),
		escapeMarkdown(resp.Data.Attributes.DeveloperLogURL),
	)
	return nil
}

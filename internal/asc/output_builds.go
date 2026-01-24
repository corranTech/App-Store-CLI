package asc

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

// BuildUploadResult represents CLI output for build upload preparation.
type BuildUploadResult struct {
	UploadID   string            `json:"uploadId"`
	FileID     string            `json:"fileId"`
	FileName   string            `json:"fileName"`
	FileSize   int64             `json:"fileSize"`
	Operations []UploadOperation `json:"operations,omitempty"`
}

// BuildBetaGroupsUpdateResult represents CLI output for build beta group updates.
type BuildBetaGroupsUpdateResult struct {
	BuildID  string   `json:"buildId"`
	GroupIDs []string `json:"groupIds"`
	Action   string   `json:"action"`
}

func printBuildsTable(resp *BuildsResponse) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Version\tUploaded\tProcessing\tExpired")
	for _, item := range resp.Data {
		fmt.Fprintf(w, "%s\t%s\t%s\t%t\n",
			item.Attributes.Version,
			item.Attributes.UploadedDate,
			item.Attributes.ProcessingState,
			item.Attributes.Expired,
		)
	}
	return w.Flush()
}

func printBuildsMarkdown(resp *BuildsResponse) error {
	fmt.Fprintln(os.Stdout, "| Version | Uploaded | Processing | Expired |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- |")
	for _, item := range resp.Data {
		fmt.Fprintf(os.Stdout, "| %s | %s | %s | %t |\n",
			escapeMarkdown(item.Attributes.Version),
			escapeMarkdown(item.Attributes.UploadedDate),
			escapeMarkdown(item.Attributes.ProcessingState),
			item.Attributes.Expired,
		)
	}
	return nil
}

func printBuildUploadResultTable(result *BuildUploadResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Upload ID\tFile ID\tFile Name\tFile Size")
	fmt.Fprintf(w, "%s\t%s\t%s\t%d\n",
		result.UploadID,
		result.FileID,
		result.FileName,
		result.FileSize,
	)
	if err := w.Flush(); err != nil {
		return err
	}
	if len(result.Operations) == 0 {
		return nil
	}
	fmt.Fprintln(os.Stdout, "\nUpload Operations")
	opsWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(opsWriter, "Method\tURL\tLength\tOffset")
	for _, op := range result.Operations {
		fmt.Fprintf(opsWriter, "%s\t%s\t%d\t%d\n",
			op.Method,
			op.URL,
			op.Length,
			op.Offset,
		)
	}
	return opsWriter.Flush()
}

func printBuildUploadResultMarkdown(result *BuildUploadResult) error {
	fmt.Fprintln(os.Stdout, "| Upload ID | File ID | File Name | File Size |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- |")
	fmt.Fprintf(os.Stdout, "| %s | %s | %s | %d |\n",
		escapeMarkdown(result.UploadID),
		escapeMarkdown(result.FileID),
		escapeMarkdown(result.FileName),
		result.FileSize,
	)
	if len(result.Operations) == 0 {
		return nil
	}
	fmt.Fprintln(os.Stdout, "\n| Method | URL | Length | Offset |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- | --- |")
	for _, op := range result.Operations {
		fmt.Fprintf(os.Stdout, "| %s | %s | %d | %d |\n",
			escapeMarkdown(op.Method),
			escapeMarkdown(op.URL),
			op.Length,
			op.Offset,
		)
	}
	return nil
}

func printBuildBetaGroupsUpdateTable(result *BuildBetaGroupsUpdateResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Build ID\tGroup IDs\tAction")
	fmt.Fprintf(w, "%s\t%s\t%s\n",
		result.BuildID,
		strings.Join(result.GroupIDs, ", "),
		result.Action,
	)
	return w.Flush()
}

func printBuildBetaGroupsUpdateMarkdown(result *BuildBetaGroupsUpdateResult) error {
	fmt.Fprintln(os.Stdout, "| Build ID | Group IDs | Action |")
	fmt.Fprintln(os.Stdout, "| --- | --- | --- |")
	fmt.Fprintf(os.Stdout, "| %s | %s | %s |\n",
		escapeMarkdown(result.BuildID),
		escapeMarkdown(strings.Join(result.GroupIDs, ", ")),
		escapeMarkdown(result.Action),
	)
	return nil
}

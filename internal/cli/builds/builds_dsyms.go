package builds

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
)

// dsymHTTPClient is the HTTP client used for dSYM downloads.
// Tests can replace this.
var dsymHTTPClient = &http.Client{
	Timeout: 5 * time.Minute,
}

// DSYMDownloadResult is the structured output for dSYM downloads.
type DSYMDownloadResult struct {
	BuildID string             `json:"buildId"`
	Dir     string             `json:"dir"`
	Files   []DSYMDownloadFile `json:"files"`
}

// DSYMDownloadFile describes one downloaded dSYM file.
type DSYMDownloadFile struct {
	BundleID string `json:"bundleId,omitempty"`
	FileName string `json:"fileName"`
	FilePath string `json:"filePath"`
	FileSize int64  `json:"fileSize"`
}

type dsymBundleInfo struct {
	BundleID string
	DSYMURL  *string
}

// BuildsDsymsCommand returns the builds dsyms subcommand.
func BuildsDsymsCommand() *ffcli.Command {
	fs := flag.NewFlagSet("dsyms", flag.ExitOnError)

	buildID := fs.String("build", "", "Build ID (required)")
	outputDir := fs.String("output-dir", ".", "Output directory for dSYM files")
	output := shared.BindOutputFlags(fs)

	return &ffcli.Command{
		Name:       "dsyms",
		ShortUsage: "asc builds dsyms --build \"BUILD_ID\" [--output-dir \"./dsyms\"]",
		ShortHelp:  "Download dSYM files for a build.",
		LongHelp: `Download dSYM debug symbol files for a build.

dSYM files are used for crash symbolication with tools like Crashlytics
and Sentry. Each build bundle that includes symbols will have a dSYM
download URL.

Examples:
  asc builds dsyms --build "BUILD_ID"
  asc builds dsyms --build "BUILD_ID" --output-dir "./dsyms"
  asc builds dsyms --build "BUILD_ID" --output json`,
		FlagSet:   fs,
		UsageFunc: shared.DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			trimmedBuildID := strings.TrimSpace(*buildID)
			if trimmedBuildID == "" {
				return shared.UsageError("--build is required")
			}

			dirValue := strings.TrimSpace(*outputDir)
			if dirValue == "" {
				dirValue = "."
			}

			client, err := shared.GetASCClient()
			if err != nil {
				return fmt.Errorf("builds dsyms: %w", err)
			}

			requestCtx, cancel := shared.ContextWithTimeout(ctx)
			defer cancel()

			bundlesResp, err := client.GetBuildBundlesForBuild(requestCtx, trimmedBuildID)
			if err != nil {
				return fmt.Errorf("builds dsyms: %w", err)
			}

			bundles := make([]dsymBundleInfo, 0, len(bundlesResp.Data))
			for _, b := range bundlesResp.Data {
				bundleID := ""
				if b.Attributes.BundleID != nil {
					bundleID = *b.Attributes.BundleID
				}
				bundles = append(bundles, dsymBundleInfo{
					BundleID: bundleID,
					DSYMURL:  b.Attributes.DSYMURL,
				})
			}

			downloadable := filterBundlesWithDSYM(bundles)
			if len(downloadable) == 0 {
				fmt.Fprintln(os.Stderr, "No dSYM files available for this build")
				result := DSYMDownloadResult{
					BuildID: trimmedBuildID,
					Dir:     dirValue,
					Files:   []DSYMDownloadFile{},
				}
				return shared.PrintOutput(result, *output.Output, *output.Pretty)
			}

			if err := os.MkdirAll(dirValue, 0o755); err != nil {
				return fmt.Errorf("builds dsyms: failed to create output directory: %w", err)
			}

			files := make([]DSYMDownloadFile, 0, len(downloadable))
			for i, bundle := range downloadable {
				fileName := dsymFileName(bundle.BundleID, trimmedBuildID, i)
				filePath := filepath.Join(dirValue, fileName)

				fmt.Fprintf(os.Stderr, "Downloading dSYM for %s...\n", displayBundleID(bundle.BundleID, i))

				size, err := downloadDSYM(*bundle.DSYMURL, filePath)
				if err != nil {
					return fmt.Errorf("builds dsyms: failed to download %s: %w", fileName, err)
				}

				fmt.Fprintf(os.Stderr, "  Saved %s (%d bytes)\n", filePath, size)

				files = append(files, DSYMDownloadFile{
					BundleID: bundle.BundleID,
					FileName: fileName,
					FilePath: filePath,
					FileSize: size,
				})
			}

			result := DSYMDownloadResult{
				BuildID: trimmedBuildID,
				Dir:     dirValue,
				Files:   files,
			}

			return shared.PrintOutput(result, *output.Output, *output.Pretty)
		},
	}
}

func filterBundlesWithDSYM(bundles []dsymBundleInfo) []dsymBundleInfo {
	result := make([]dsymBundleInfo, 0, len(bundles))
	for _, b := range bundles {
		if b.DSYMURL != nil && strings.TrimSpace(*b.DSYMURL) != "" {
			result = append(result, b)
		}
	}
	return result
}

func dsymFileName(bundleID, buildID string, index int) string {
	if bundleID != "" {
		return bundleID + ".dsym.zip"
	}
	return fmt.Sprintf("%s_%d.dsym.zip", buildID, index)
}

func displayBundleID(bundleID string, index int) string {
	if bundleID != "" {
		return bundleID
	}
	return fmt.Sprintf("bundle %d", index)
}

func downloadDSYM(url, destPath string) (int64, error) {
	resp, err := dsymHTTPClient.Get(url)
	if err != nil {
		return 0, fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("download returned HTTP %d", resp.StatusCode)
	}

	return shared.WriteStreamToFile(destPath, resp.Body)
}

// SetDSYMHTTPClient replaces the HTTP client for tests.
func SetDSYMHTTPClient(c *http.Client) func() {
	prev := dsymHTTPClient
	dsymHTTPClient = c
	return func() { dsymHTTPClient = prev }
}

// dsymTableRows returns table headers and rows for DSYMDownloadResult.
func dsymTableRows(result DSYMDownloadResult) ([]string, [][]string) {
	headers := []string{"bundleId", "fileName", "filePath", "fileSize"}
	rows := make([][]string, 0, len(result.Files))
	for _, f := range result.Files {
		rows = append(rows, []string{
			f.BundleID,
			f.FileName,
			f.FilePath,
			fmt.Sprintf("%d", f.FileSize),
		})
	}
	return headers, rows
}

func init() {
	// Register table/markdown renderers if needed in the future.
	_ = dsymTableRows
}

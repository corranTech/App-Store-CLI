package assets

import (
	"os"
	"testing"
)

func TestValidateScreenshotAssetsSortsEntriesAndKeepsHiddenWarningsNonBlocking(t *testing.T) {
	dir := t.TempDir()
	writeAssetsTestPNGWithSize(t, dir, "02-details.png", 1242, 2688)
	writeAssetsTestPNGWithSize(t, dir, "01-home.png", 1242, 2688)
	writeAssetsTestPNGWithSize(t, dir, ".hidden.png", 1242, 2688)

	result, err := validateScreenshotAssets(dir, "APP_IPHONE_65")
	if err != nil {
		t.Fatalf("validateScreenshotAssets() error: %v", err)
	}

	if result.ErrorCount != 0 {
		t.Fatalf("expected 0 errors, got %d", result.ErrorCount)
	}
	if result.WarningCount != 1 {
		t.Fatalf("expected 1 warning, got %d", result.WarningCount)
	}
	if result.ReadyFiles != 3 {
		t.Fatalf("expected 3 ready files, got %d", result.ReadyFiles)
	}

	wantOrder := []string{".hidden.png", "01-home.png", "02-details.png"}
	for i, want := range wantOrder {
		if result.Files[i].FileName != want {
			t.Fatalf("expected file %q at index %d, got %q", want, i, result.Files[i].FileName)
		}
		if result.Files[i].Order != i+1 {
			t.Fatalf("expected order %d at index %d, got %d", i+1, i, result.Files[i].Order)
		}
	}

	if !hasScreenshotValidateIssueWithSeverity(result.Issues, "hidden_file", screenshotValidateSeverityWarning, ".hidden.png") {
		t.Fatalf("expected hidden-file warning, got %+v", result.Issues)
	}
}

func TestValidateScreenshotAssetsReportsUnreadableDotfilesAndDimensionMismatch(t *testing.T) {
	dir := t.TempDir()
	writeAssetsTestPNGWithSize(t, dir, "01-home.png", 1242, 2688)
	writeAssetsTestPNGWithSize(t, dir, "03-bad.png", 100, 100)
	if err := os.WriteFile(dir+"/.DS_Store", []byte("not an image"), 0o644); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	result, err := validateScreenshotAssets(dir, "APP_IPHONE_65")
	if err != nil {
		t.Fatalf("validateScreenshotAssets() error: %v", err)
	}

	if result.ErrorCount != 2 {
		t.Fatalf("expected 2 errors, got %d", result.ErrorCount)
	}
	if result.WarningCount != 1 {
		t.Fatalf("expected 1 warning, got %d", result.WarningCount)
	}
	if result.ReadyFiles != 1 {
		t.Fatalf("expected 1 ready file, got %d", result.ReadyFiles)
	}

	if !hasScreenshotValidateIssueWithSeverity(result.Issues, "hidden_file", screenshotValidateSeverityWarning, ".DS_Store") {
		t.Fatalf("expected hidden-file warning, got %+v", result.Issues)
	}
	if !hasScreenshotValidateIssueWithSeverity(result.Issues, "read_failure", screenshotValidateSeverityError, ".DS_Store") {
		t.Fatalf("expected read-failure error, got %+v", result.Issues)
	}
	if !hasScreenshotValidateIssueWithSeverity(result.Issues, "dimension_mismatch", screenshotValidateSeverityError, "03-bad.png") {
		t.Fatalf("expected dimension mismatch error, got %+v", result.Issues)
	}
}

func hasScreenshotValidateIssueWithSeverity(issues []screenshotValidateIssue, code, severity, fileName string) bool {
	for _, issue := range issues {
		if issue.Code == code && issue.Severity == severity && issue.FileName == fileName {
			return true
		}
	}
	return false
}

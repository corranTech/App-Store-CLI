package cmdtest

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScreenshotsUploadAppScopedModeRequiresVersionSelector(t *testing.T) {
	stdout, stderr, runErr := runRootCommand(t, []string{
		"screenshots", "upload",
		"--app", "123456789",
		"--path", "./screenshots",
		"--device-type", "IPHONE_65",
	})

	if stdout != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
	if !errors.Is(runErr, flag.ErrHelp) {
		t.Fatalf("expected flag.ErrHelp, got %v", runErr)
	}
	if !strings.Contains(stderr, "Error: --version or --version-id is required with --app") {
		t.Fatalf("expected missing app-scoped version selector error, got %q", stderr)
	}
}

func TestScreenshotsUploadRejectsMixingDirectAndAppScopedSelectors(t *testing.T) {
	stdout, stderr, runErr := runRootCommand(t, []string{
		"screenshots", "upload",
		"--version-localization", "LOC_ID",
		"--app", "123456789",
		"--version", "1.2.3",
		"--path", "./screenshots",
		"--device-type", "IPHONE_65",
	})

	if stdout != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
	if !errors.Is(runErr, flag.ErrHelp) {
		t.Fatalf("expected flag.ErrHelp, got %v", runErr)
	}
	if !strings.Contains(stderr, "Error: --version-localization cannot be combined with --app, --version, --version-id, or --platform") {
		t.Fatalf("expected direct/app-scoped selector conflict error, got %q", stderr)
	}
}

func TestScreenshotsUploadIgnoresASCAppIDUntilAppScopedModeIsRequested(t *testing.T) {
	t.Setenv("ASC_APP_ID", "123456789")

	stdout, stderr, runErr := runRootCommand(t, []string{
		"screenshots", "upload",
		"--path", "./screenshots",
		"--device-type", "IPHONE_65",
	})

	if stdout != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
	if !errors.Is(runErr, flag.ErrHelp) {
		t.Fatalf("expected flag.ErrHelp, got %v", runErr)
	}
	if !strings.Contains(stderr, "Error: --version-localization is required") {
		t.Fatalf("expected direct-mode selector error, got %q", stderr)
	}
}

func TestScreenshotsUploadAppScopedModeRejectsInvalidPlatformBeforeAuth(t *testing.T) {
	t.Setenv("ASC_KEY_ID", "")
	t.Setenv("ASC_ISSUER_ID", "")
	t.Setenv("ASC_PRIVATE_KEY_PATH", "")
	t.Setenv("ASC_PRIVATE_KEY", "")
	t.Setenv("ASC_PRIVATE_KEY_B64", "")
	t.Setenv("ASC_APP_ID", "")
	t.Setenv("ASC_PROFILE", "")
	t.Setenv("ASC_STRICT_AUTH", "")
	t.Setenv("ASC_CONFIG_PATH", filepath.Join(t.TempDir(), "nonexistent.json"))

	rootDir := t.TempDir()
	localeDir := filepath.Join(rootDir, "en-US", "iphone")
	if err := os.MkdirAll(localeDir, 0o755); err != nil {
		t.Fatalf("mkdir locale dir: %v", err)
	}
	writePNG(t, filepath.Join(localeDir, "01-home.png"), 1284, 2778)

	stdout, stderr, runErr := runRootCommand(t, []string{
		"screenshots", "upload",
		"--app", "123456789",
		"--version", "1.2.3",
		"--platform", "ANDROID",
		"--path", rootDir,
		"--device-type", "IPHONE_65",
	})

	if stdout != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
	if !errors.Is(runErr, flag.ErrHelp) {
		t.Fatalf("expected flag.ErrHelp, got %v", runErr)
	}
	if !strings.Contains(stderr, "Error: --platform must be one of: IOS, MAC_OS, TV_OS, VISION_OS") {
		t.Fatalf("expected invalid platform usage error, got %q", stderr)
	}
	if strings.Contains(stderr, "screenshots upload:") {
		t.Fatalf("expected raw usage error without command prefix, got %q", stderr)
	}
}

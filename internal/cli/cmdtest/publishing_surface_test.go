package cmdtest

import (
	"context"
	"errors"
	"flag"
	"io"
	"strings"
	"testing"
)

func TestPublishHelpPointsToCanonicalPublishCommands(t *testing.T) {
	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	_, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{"publish"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		err := root.Run(context.Background())
		if !errors.Is(err, flag.ErrHelp) {
			t.Fatalf("expected ErrHelp, got %v", err)
		}
	})

	if !strings.Contains(stderr, "asc publish testflight") {
		t.Fatalf("expected publish help to mention canonical TestFlight publish command, got %q", stderr)
	}
	if !strings.Contains(stderr, "asc release run") {
		t.Fatalf("expected publish help to mention canonical App Store publish command, got %q", stderr)
	}
	if strings.Contains(stderr, "\n  appstore") {
		t.Fatalf("expected deprecated publish appstore subcommand to be hidden from publish help, got %q", stderr)
	}
}

func TestPublishAppStoreWarnsDeprecatedPath(t *testing.T) {
	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	_, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{"publish", "appstore"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		err := root.Run(context.Background())
		if !errors.Is(err, flag.ErrHelp) {
			t.Fatalf("expected ErrHelp, got %v", err)
		}
	})

	requireStderrContainsWarning(t, stderr, "Warning: `asc publish appstore` is deprecated. Use `asc release run` for the canonical App Store publish flow.")
	if !strings.Contains(stderr, "Error: --ipa is required") {
		t.Fatalf("expected validation error after deprecation warning, got %q", stderr)
	}
}

func TestSubmitHelpPointsToReleaseRunAndHidesCreate(t *testing.T) {
	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	_, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{"submit"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		err := root.Run(context.Background())
		if !errors.Is(err, flag.ErrHelp) {
			t.Fatalf("expected ErrHelp, got %v", err)
		}
	})

	if !strings.Contains(stderr, "asc release run") {
		t.Fatalf("expected submit help to point to asc release run, got %q", stderr)
	}
	if !strings.Contains(stderr, "asc submit preflight") {
		t.Fatalf("expected submit help to keep preflight guidance, got %q", stderr)
	}
	if strings.Contains(stderr, "\n  create") {
		t.Fatalf("expected deprecated submit create subcommand to be hidden from submit help, got %q", stderr)
	}
}

func TestSubmitCreateWarnsDeprecatedPath(t *testing.T) {
	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	_, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{
			"submit", "create",
			"--app", "app-1",
			"--version", "1.0.0",
			"--build", "build-1",
		}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		err := root.Run(context.Background())
		if !errors.Is(err, flag.ErrHelp) {
			t.Fatalf("expected ErrHelp, got %v", err)
		}
	})

	requireStderrContainsWarning(t, stderr, "Warning: `asc submit create` is deprecated. Use `asc release run` for the canonical App Store publish flow.")
	if !strings.Contains(stderr, "--confirm is required") {
		t.Fatalf("expected normal validation error after deprecation warning, got %q", stderr)
	}
}

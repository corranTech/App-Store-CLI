package cmdtest

import (
	"context"
	"errors"
	"flag"
	"io"
	"strings"
	"testing"
)

func TestSubscriptionsPricesAdd_TierAndPricePointMutualExclusion(t *testing.T) {
	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	_, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{
			"subscriptions", "prices", "add",
			"--id", "SUB_ID",
			"--price-point", "PP",
			"--tier", "5",
			"--territory", "USA",
		}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		err := root.Run(context.Background())
		if !errors.Is(err, flag.ErrHelp) {
			t.Fatalf("expected ErrHelp, got %v", err)
		}
	})

	if !strings.Contains(stderr, "mutually exclusive") {
		t.Fatalf("expected mutually exclusive error, got %q", stderr)
	}
}

func TestSubscriptionsPricesAdd_TierRequiresApp(t *testing.T) {
	t.Setenv("ASC_APP_ID", "")

	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	_, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{
			"subscriptions", "prices", "add",
			"--id", "SUB_ID",
			"--tier", "5",
			"--territory", "USA",
		}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		err := root.Run(context.Background())
		if !errors.Is(err, flag.ErrHelp) {
			t.Fatalf("expected ErrHelp, got %v", err)
		}
	})

	if !strings.Contains(stderr, "--app is required") {
		t.Fatalf("expected --app required error, got %q", stderr)
	}
}

func TestSubscriptionsPricesAdd_TierRequiresTerritory(t *testing.T) {
	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	_, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{
			"subscriptions", "prices", "add",
			"--id", "SUB_ID",
			"--tier", "5",
			"--app", "APP",
		}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		err := root.Run(context.Background())
		if !errors.Is(err, flag.ErrHelp) {
			t.Fatalf("expected ErrHelp, got %v", err)
		}
	})

	if !strings.Contains(stderr, "--territory is required") {
		t.Fatalf("expected --territory required error, got %q", stderr)
	}
}

package cmdtest

import (
	"context"
	"errors"
	"flag"
	"io"
	"strings"
	"testing"
)

func TestIAPPriceSchedulesCreate_TierAndPricesMutualExclusion(t *testing.T) {
	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	_, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{
			"iap", "price-schedules", "create",
			"--iap-id", "IAP_ID",
			"--base-territory", "USA",
			"--tier", "5",
			"--prices", "PP:2026-03-01",
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

func TestIAPPriceSchedulesCreate_TierAndPriceMutualExclusion(t *testing.T) {
	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	_, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{
			"iap", "price-schedules", "create",
			"--iap-id", "IAP_ID",
			"--base-territory", "USA",
			"--tier", "5",
			"--price", "4.99",
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

func TestIAPPriceSchedulesCreate_TierRequiresApp(t *testing.T) {
	t.Setenv("ASC_APP_ID", "")

	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	_, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{
			"iap", "price-schedules", "create",
			"--iap-id", "IAP_ID",
			"--base-territory", "USA",
			"--tier", "5",
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

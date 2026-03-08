package cmdtest

import (
	"context"
	"errors"
	"flag"
	"io"
	"strings"
	"testing"
)

func TestSubscriptionsHelpShowsCanonicalCommerceSubcommands(t *testing.T) {
	root := RootCommand("1.2.3")

	subscriptionsCmd := findSubcommand(root, "subscriptions")
	if subscriptionsCmd == nil {
		t.Fatal("expected subscriptions command")
	}
	subscriptionsUsage := subscriptionsCmd.UsageFunc(subscriptionsCmd)
	if !strings.Contains(subscriptionsUsage, "  win-back-offers") {
		t.Fatalf("expected subscriptions help to list win-back-offers, got %q", subscriptionsUsage)
	}
	if !strings.Contains(subscriptionsUsage, "  promoted-purchases") {
		t.Fatalf("expected subscriptions help to list promoted-purchases, got %q", subscriptionsUsage)
	}
	if usageListsSubcommand(subscriptionsUsage, "promoted-purchase") {
		t.Fatalf("expected subscriptions help to hide deprecated singular promoted-purchase shim, got %q", subscriptionsUsage)
	}

	offerCodesCmd := findSubcommand(root, "subscriptions", "offer-codes")
	if offerCodesCmd == nil {
		t.Fatal("expected subscriptions offer-codes command")
	}
	offerCodesUsage := offerCodesCmd.UsageFunc(offerCodesCmd)
	if !strings.Contains(offerCodesUsage, "  generate") {
		t.Fatalf("expected subscriptions offer-codes help to list generate, got %q", offerCodesUsage)
	}
	if !strings.Contains(offerCodesUsage, "  values") {
		t.Fatalf("expected subscriptions offer-codes help to list values, got %q", offerCodesUsage)
	}

	iapCmd := findSubcommand(root, "iap")
	if iapCmd == nil {
		t.Fatal("expected iap command")
	}
	iapUsage := iapCmd.UsageFunc(iapCmd)
	if !strings.Contains(iapUsage, "  promoted-purchases") {
		t.Fatalf("expected iap help to list promoted-purchases, got %q", iapUsage)
	}
	if usageListsSubcommand(iapUsage, "promoted-purchase") {
		t.Fatalf("expected iap help to hide deprecated singular promoted-purchase shim, got %q", iapUsage)
	}
}

func TestLegacyCommerceRootCommandsAreDeprecatedInHelp(t *testing.T) {
	root := RootCommand("1.2.3")

	tests := []struct {
		name      string
		command   string
		wantUsage string
	}{
		{
			name:      "offer-codes",
			command:   "offer-codes",
			wantUsage: `asc subscriptions offer-codes`,
		},
		{
			name:      "win-back-offers",
			command:   "win-back-offers",
			wantUsage: `asc subscriptions win-back-offers`,
		},
		{
			name:      "promoted-purchases",
			command:   "promoted-purchases",
			wantUsage: `asc subscriptions promoted-purchases`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cmd := findSubcommand(root, test.command)
			if cmd == nil {
				t.Fatalf("expected %s command", test.command)
			}
			usage := cmd.UsageFunc(cmd)
			if !strings.Contains(usage, "DEPRECATED:") {
				t.Fatalf("expected deprecated help text, got %q", usage)
			}
			if !strings.Contains(usage, test.wantUsage) {
				t.Fatalf("expected canonical replacement %q in usage, got %q", test.wantUsage, usage)
			}
		})
	}
}

func TestLegacyOfferCodesExecutionWarnsToStderr(t *testing.T) {
	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	stdout, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{"offer-codes", "values"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		err := root.Run(context.Background())
		if !errors.Is(err, flag.ErrHelp) {
			t.Fatalf("expected ErrHelp, got %v", err)
		}
	})

	if stdout != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
	if !strings.Contains(stderr, `Warning: "asc offer-codes values" is deprecated. Use "asc subscriptions offer-codes values" instead.`) {
		t.Fatalf("expected deprecation warning in stderr, got %q", stderr)
	}
	if !strings.Contains(stderr, "--id is required") {
		t.Fatalf("expected original validation error in stderr, got %q", stderr)
	}
}

func TestCanonicalWrapperErrorsUseCanonicalPaths(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "subscriptions win-back-offers next validation",
			args:    []string{"subscriptions", "win-back-offers", "list", "--next", "http://api.appstoreconnect.apple.com/v1/subscriptions/sub-1/winBackOffers?cursor=AQ"},
			wantErr: "subscriptions win-back-offers list: --next must be an App Store Connect URL",
		},
		{
			name:    "subscriptions promoted-purchases next validation",
			args:    []string{"subscriptions", "promoted-purchases", "list", "--next", "http://api.appstoreconnect.apple.com/v1/apps/app-1/promotedPurchases?cursor=AQ"},
			wantErr: "subscriptions promoted-purchases list: --next must be an App Store Connect URL",
		},
		{
			name:    "iap promoted-purchases next validation",
			args:    []string{"iap", "promoted-purchases", "list", "--next", "http://api.appstoreconnect.apple.com/v1/apps/app-1/promotedPurchases?cursor=AQ"},
			wantErr: "iap promoted-purchases list: --next must be an App Store Connect URL",
		},
		{
			name:    "subscriptions offer-codes values auth error",
			args:    []string{"subscriptions", "offer-codes", "values", "--id", "batch-1"},
			wantErr: "subscriptions offer-codes values:",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := RootCommand("1.2.3")
			root.FlagSet.SetOutput(io.Discard)

			var runErr error
			stdout, stderr := captureOutput(t, func() {
				if err := root.Parse(test.args); err != nil {
					t.Fatalf("parse error: %v", err)
				}
				runErr = root.Run(context.Background())
			})

			if runErr == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(runErr.Error(), test.wantErr) {
				t.Fatalf("expected error %q, got %v", test.wantErr, runErr)
			}
			if stdout != "" {
				t.Fatalf("expected empty stdout, got %q", stdout)
			}
			if stderr != "" {
				t.Fatalf("expected empty stderr, got %q", stderr)
			}
		})
	}
}

func usageListsSubcommand(usage string, name string) bool {
	for _, line := range strings.Split(usage, "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		if fields[0] == name {
			return true
		}
	}
	return false
}

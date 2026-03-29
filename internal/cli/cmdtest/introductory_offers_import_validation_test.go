package cmdtest

import (
	"context"
	"errors"
	"flag"
	"io"
	"path/filepath"
	"strings"
	"testing"
)

func TestSubscriptionsIntroductoryOffersImport_MissingRequiredFlagsReturnUsage(t *testing.T) {
	t.Setenv("ASC_APP_ID", "")
	t.Setenv("ASC_CONFIG_PATH", filepath.Join(t.TempDir(), "config.json"))

	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "missing subscription id",
			args:    []string{"subscriptions", "offers", "introductory", "import", "--input", "offers.csv"},
			wantErr: "Error: --subscription-id is required",
		},
		{
			name:    "missing input",
			args:    []string{"subscriptions", "offers", "introductory", "import", "--subscription-id", "SUB_ID"},
			wantErr: "Error: --input is required",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := RootCommand("1.2.3")
			root.FlagSet.SetOutput(io.Discard)

			stdout, stderr := captureOutput(t, func() {
				if err := root.Parse(test.args); err != nil {
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
			if !strings.Contains(stderr, test.wantErr) {
				t.Fatalf("expected %q in stderr, got %q", test.wantErr, stderr)
			}
		})
	}
}

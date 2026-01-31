package cmdtest

import (
	"context"
	"errors"
	"flag"
	"io"
	"strings"
	"testing"
)

func TestRelationshipCommandValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "bundle-ids app get missing id",
			args:    []string{"bundle-ids", "app", "get"},
			wantErr: "--id is required",
		},
		{
			name:    "bundle-ids profiles list missing id",
			args:    []string{"bundle-ids", "profiles", "list"},
			wantErr: "--id is required",
		},
		{
			name:    "users invites visible-apps list missing id",
			args:    []string{"users", "invites", "visible-apps", "list"},
			wantErr: "--id is required",
		},
		{
			name:    "agreements territories list missing id",
			args:    []string{"agreements", "territories", "list"},
			wantErr: "--id is required",
		},
		{
			name:    "background-assets app-store-releases get missing id",
			args:    []string{"background-assets", "app-store-releases", "get"},
			wantErr: "--id is required",
		},
		{
			name:    "background-assets external-beta-releases get missing id",
			args:    []string{"background-assets", "external-beta-releases", "get"},
			wantErr: "--id is required",
		},
		{
			name:    "background-assets internal-beta-releases get missing id",
			args:    []string{"background-assets", "internal-beta-releases", "get"},
			wantErr: "--id is required",
		},
		{
			name:    "apps ci-product get missing id",
			args:    []string{"apps", "ci-product", "get"},
			wantErr: "--id is required",
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
				t.Fatalf("expected error %q, got %q", test.wantErr, stderr)
			}
		})
	}
}

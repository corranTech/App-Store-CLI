package cmdtest

import (
	"context"
	"errors"
	"flag"
	"io"
	"strings"
	"testing"
)

func TestBuildsLegacySelectorMigrationErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "builds wait legacy build",
			args:    []string{"builds", "wait", "--build", "BUILD_123"},
			wantErr: `--build was removed; use --build-id "BUILD_123"`,
		},
		{
			name:    "builds wait legacy newest",
			args:    []string{"builds", "wait", "--app", "APP_123", "--newest"},
			wantErr: "--newest was removed; use --latest",
		},
		{
			name:    "builds dsyms legacy build",
			args:    []string{"builds", "dsyms", "--build", "BUILD_123"},
			wantErr: `--build was removed; use --build-id "BUILD_123"`,
		},
		{
			name:    "builds info legacy build",
			args:    []string{"builds", "info", "--build", "BUILD_123"},
			wantErr: `--build was removed; use --build-id "BUILD_123"`,
		},
		{
			name:    "builds app get legacy build",
			args:    []string{"builds", "app", "get", "--build", "BUILD_123"},
			wantErr: `--build was removed; use --build-id "BUILD_123"`,
		},
		{
			name:    "builds app get legacy id",
			args:    []string{"builds", "app", "get", "--id", "BUILD_123"},
			wantErr: `--id was removed as a build selector; use --build-id "BUILD_123"`,
		},
		{
			name:    "builds pre-release-version get legacy id",
			args:    []string{"builds", "pre-release-version", "get", "--id", "BUILD_123"},
			wantErr: `--id was removed as a build selector; use --build-id "BUILD_123"`,
		},
		{
			name:    "builds icons list legacy build",
			args:    []string{"builds", "icons", "list", "--build", "BUILD_123"},
			wantErr: `--build was removed; use --build-id "BUILD_123"`,
		},
		{
			name:    "builds beta-app-review-submission get legacy id",
			args:    []string{"builds", "beta-app-review-submission", "get", "--id", "BUILD_123"},
			wantErr: `--id was removed as a build selector; use --build-id "BUILD_123"`,
		},
		{
			name:    "builds build-beta-detail get legacy build",
			args:    []string{"builds", "build-beta-detail", "get", "--build", "BUILD_123"},
			wantErr: `--build was removed; use --build-id "BUILD_123"`,
		},
		{
			name:    "builds links view legacy build",
			args:    []string{"builds", "links", "view", "--build", "BUILD_123", "--type", "app"},
			wantErr: `--build was removed; use --build-id "BUILD_123"`,
		},
		{
			name:    "builds metrics beta-usages legacy build",
			args:    []string{"builds", "metrics", "beta-usages", "--build", "BUILD_123"},
			wantErr: `--build was removed; use --build-id "BUILD_123"`,
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
				t.Fatalf("expected migration error %q, got %q", test.wantErr, stderr)
			}
			if strings.Contains(stderr, "flag provided but not defined") {
				t.Fatalf("expected friendly migration guidance instead of parse error, got %q", stderr)
			}
		})
	}
}

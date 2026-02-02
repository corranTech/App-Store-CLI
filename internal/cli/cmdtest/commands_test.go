package cmdtest

import (
	"context"
	"errors"
	"flag"
	"io"
	"strings"
	"testing"
)

func TestAppsSearchKeywordsValidationErrors(t *testing.T) {
	t.Setenv("ASC_APP_ID", "")

	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "apps search-keywords list missing app",
			args:    []string{"apps", "search-keywords", "list"},
			wantErr: "--app is required",
		},
		{
			name:    "apps search-keywords set missing app",
			args:    []string{"apps", "search-keywords", "set", "--keywords", "kw1", "--confirm"},
			wantErr: "--app is required",
		},
		{
			name:    "apps search-keywords set missing confirm",
			args:    []string{"apps", "search-keywords", "set", "--app", "123", "--keywords", "kw1"},
			wantErr: "--confirm is required",
		},
		{
			name:    "apps search-keywords set missing keywords",
			args:    []string{"apps", "search-keywords", "set", "--app", "123", "--confirm"},
			wantErr: "--keywords is required",
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

func TestLocalizationsSearchKeywordsValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "localizations search-keywords list missing localization",
			args:    []string{"localizations", "search-keywords", "list"},
			wantErr: "--localization-id is required",
		},
		{
			name:    "localizations search-keywords add missing keywords",
			args:    []string{"localizations", "search-keywords", "add", "--localization-id", "loc-1"},
			wantErr: "--keywords is required",
		},
		{
			name:    "localizations search-keywords delete missing confirm",
			args:    []string{"localizations", "search-keywords", "delete", "--localization-id", "loc-1", "--keywords", "kw1"},
			wantErr: "--confirm is required",
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

func TestLocalizationsMediaSetsValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "preview sets list missing localization",
			args:    []string{"localizations", "preview-sets", "list"},
			wantErr: "--localization-id is required",
		},
		{
			name:    "preview sets relationships missing localization",
			args:    []string{"localizations", "preview-sets", "relationships"},
			wantErr: "--localization-id is required",
		},
		{
			name:    "screenshot sets list missing localization",
			args:    []string{"localizations", "screenshot-sets", "list"},
			wantErr: "--localization-id is required",
		},
		{
			name:    "screenshot sets relationships missing localization",
			args:    []string{"localizations", "screenshot-sets", "relationships"},
			wantErr: "--localization-id is required",
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

func TestVersionsRelationshipsValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "versions relationships missing type",
			args:    []string{"versions", "relationships", "--version-id", "id-1"},
			wantErr: "--type is required",
		},
		{
			name:    "versions relationships missing version id",
			args:    []string{"versions", "relationships", "--type", "appStoreReviewDetail"},
			wantErr: "--version-id is required",
		},
		{
			name:    "versions relationships invalid type",
			args:    []string{"versions", "relationships", "--version-id", "id-1", "--type", "nope"},
			wantErr: "--type must be one of",
		},
		{
			name:    "versions relationships invalid limit for single",
			args:    []string{"versions", "relationships", "--version-id", "id-1", "--type", "appStoreReviewDetail", "--limit", "10"},
			wantErr: "--limit, --next, and --paginate are only valid for to-many relationships",
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

func TestAppsUpdateValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "apps update missing id",
			args:    []string{"apps", "update"},
			wantErr: "Error: --id is required",
		},
		{
			name:    "apps update missing fields",
			args:    []string{"apps", "update", "--id", "APP_ID"},
			wantErr: "Error: --bundle-id, --primary-locale, or --content-rights is required",
		},
		{
			name:    "apps update invalid content rights",
			args:    []string{"apps", "update", "--id", "APP_ID", "--content-rights", "INVALID"},
			wantErr: "Error: --content-rights must be DOES_NOT_USE_THIRD_PARTY_CONTENT or USES_THIRD_PARTY_CONTENT",
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

func TestAppSetupInfoSetValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "app-setup info set missing app",
			args:    []string{"app-setup", "info", "set", "--content-rights", "DOES_NOT_USE_THIRD_PARTY_CONTENT"},
			wantErr: "Error: --app is required",
		},
		{
			name:    "app-setup info set missing updates",
			args:    []string{"app-setup", "info", "set", "--app", "APP_ID"},
			wantErr: "Error: provide at least one update flag",
		},
		{
			name:    "app-setup info set invalid content rights",
			args:    []string{"app-setup", "info", "set", "--app", "APP_ID", "--content-rights", "INVALID"},
			wantErr: "Error: --content-rights must be DOES_NOT_USE_THIRD_PARTY_CONTENT or USES_THIRD_PARTY_CONTENT",
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

func TestAppsAppEncryptionDeclarationsValidationErrors(t *testing.T) {
	t.Setenv("ASC_APP_ID", "")

	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "apps app-encryption-declarations list missing id",
			args:    []string{"apps", "app-encryption-declarations", "list"},
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

func TestAppInfoMutualExclusiveFlags(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "app-info get invalid version and version-id",
			args:    []string{"app-info", "get", "--app", "APP_ID", "--version", "1.0", "--version-id", "VERSION_ID"},
			wantErr: "app-info get: --version and --version-id are mutually exclusive",
		},
		{
			name:    "app-info get invalid version without platform",
			args:    []string{"app-info", "get", "--app", "APP_ID", "--version", "1.0"},
			wantErr: "Error: --platform is required with --version",
		},
		{
			name:    "app-info get invalid app-info with version flags",
			args:    []string{"app-info", "get", "--app-info", "APP_INFO_ID", "--include", "ageRatingDeclaration", "--version-id", "VERSION_ID"},
			wantErr: "Error: --include cannot be used with version localization flags",
		},
		{
			name:    "app-info get invalid app-info without include",
			args:    []string{"app-info", "get", "--app-info", "APP_INFO_ID"},
			wantErr: "Error: --app-info requires --include",
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

func TestAppInfoGetIncludeValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "app-info get invalid include",
			args:    []string{"app-info", "get", "--app", "APP_ID", "--include", "invalid"},
			wantErr: "app-info get: invalid include value(s)",
		},
		{
			name:    "app-info get invalid include with version flags",
			args:    []string{"app-info", "get", "--app", "APP_ID", "--include", "ageRatingDeclaration", "--version-id", "VERSION_ID"},
			wantErr: "Error: --include cannot be used with version localization flags",
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

func TestAppInfoGetValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "app-info get invalid limit",
			args:    []string{"app-info", "get", "--app", "APP_ID", "--limit", "201"},
			wantErr: "app-info get: --limit must be between 1 and 200",
		},
		{
			name:    "app-info get invalid next",
			args:    []string{"app-info", "get", "--app", "APP_ID", "--next", "https://bad.example.com"},
			wantErr: "app-info get: invalid next URL",
		},
		{
			name:    "app-info get missing app and app-info",
			args:    []string{"app-info", "get", "--version-id", "VERSION_ID"},
			wantErr: "Error: --app or --app-info is required (or set ASC_APP_ID)",
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

func TestAppInfoAppStoreVersionResolutionErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "app-info get invalid platform",
			args:    []string{"app-info", "get", "--app", "APP_ID", "--version", "1.0", "--platform", "INVALID"},
			wantErr: "app-info get: invalid platform",
		},
		{
			name:    "app-info get invalid state",
			args:    []string{"app-info", "get", "--app", "APP_ID", "--state", "INVALID"},
			wantErr: "app-info get: invalid app store state",
		},
		{
			name:    "app-info get invalid locale",
			args:    []string{"app-info", "get", "--app", "APP_ID", "--locale", "invalid"},
			wantErr: "app-info get: invalid locale",
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

func TestAppInfoGetPaginationValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "app-info get invalid paginate with include",
			args:    []string{"app-info", "get", "--app", "APP_ID", "--include", "ageRatingDeclaration", "--paginate"},
			wantErr: "Error: --include cannot be used with version localization flags",
		},
		{
			name:    "app-info get invalid paginate with include and version",
			args:    []string{"app-info", "get", "--app", "APP_ID", "--version", "1.0", "--platform", "IOS", "--include", "ageRatingDeclaration", "--paginate"},
			wantErr: "Error: --include cannot be used with version localization flags",
		},
		{
			name:    "app-info get invalid paginate without version flags",
			args:    []string{"app-info", "get", "--app", "APP_ID", "--include", "ageRatingDeclaration", "--paginate"},
			wantErr: "Error: --include cannot be used with version localization flags",
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

func TestAppInfoGetPaginateBehavior(t *testing.T) {
	server := newFakeServer(t)

	server.route("GET", "/v1/apps/123456789/appStoreVersions", withResponse(t, "app_store_versions.json"))
	server.route("GET", "/v1/appStoreVersions/111/relationships/appStoreVersionLocalizations", withResponse(t, "app_store_version_localizations.json"))
	server.route("GET", "/v1/appStoreVersionLocalizations/111", withResponse(t, "app_store_version_localization.json"))
	server.route("GET", "/v1/appStoreVersionLocalizations/111-2", withResponse(t, "app_store_version_localization.json"))

	root := RootCommand("1.2.3")

	stdout, _ := captureOutput(t, func() {
		if err := root.Parse([]string{"app-info", "get", "--app", "123456789", "--paginate"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		if err := root.Run(context.Background()); err != nil {
			t.Fatalf("run error: %v", err)
		}
	})

	if !strings.Contains(stdout, "\"id\":\"111-2\"") {
		t.Fatalf("expected output to contain all pages, got: %s", stdout)
	}
}

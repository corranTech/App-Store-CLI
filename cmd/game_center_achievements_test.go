package cmd

import (
	"context"
	"errors"
	"flag"
	"io"
	"strings"
	"testing"
)

func TestGameCenterAchievementsListValidationErrors(t *testing.T) {
	t.Setenv("ASC_APP_ID", "")

	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	stdout, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{"game-center", "achievements", "list"}); err != nil {
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
	if !strings.Contains(stderr, "--app is required") {
		t.Fatalf("expected missing app error, got %q", stderr)
	}
}

func TestGameCenterAchievementsGetValidationErrors(t *testing.T) {
	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	stdout, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{"game-center", "achievements", "get"}); err != nil {
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
	if !strings.Contains(stderr, "--id is required") {
		t.Fatalf("expected missing id error, got %q", stderr)
	}
}

func TestGameCenterAchievementsCreateValidationErrors(t *testing.T) {
	t.Setenv("ASC_APP_ID", "")

	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "missing app",
			args:    []string{"game-center", "achievements", "create", "--reference-name", "Test", "--vendor-id", "com.test", "--points", "10"},
			wantErr: "--app is required",
		},
		{
			name:    "missing reference-name",
			args:    []string{"game-center", "achievements", "create", "--app", "APP_ID", "--vendor-id", "com.test", "--points", "10"},
			wantErr: "--reference-name is required",
		},
		{
			name:    "missing vendor-id",
			args:    []string{"game-center", "achievements", "create", "--app", "APP_ID", "--reference-name", "Test", "--points", "10"},
			wantErr: "--vendor-id is required",
		},
		{
			name:    "missing points",
			args:    []string{"game-center", "achievements", "create", "--app", "APP_ID", "--reference-name", "Test", "--vendor-id", "com.test"},
			wantErr: "--points must be between 1 and 100",
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

func TestGameCenterAchievementsUpdateValidationErrors(t *testing.T) {
	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	stdout, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{"game-center", "achievements", "update"}); err != nil {
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
	if !strings.Contains(stderr, "--id is required") {
		t.Fatalf("expected missing id error, got %q", stderr)
	}
}

func TestGameCenterAchievementsDeleteValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "missing id",
			args:    []string{"game-center", "achievements", "delete", "--confirm"},
			wantErr: "--id is required",
		},
		{
			name:    "missing confirm",
			args:    []string{"game-center", "achievements", "delete", "--id", "ACH_ID"},
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

func TestGameCenterAchievementLocalizationsListValidationErrors(t *testing.T) {
	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	stdout, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{"game-center", "achievements", "localizations", "list"}); err != nil {
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
	if !strings.Contains(stderr, "--achievement-id is required") {
		t.Fatalf("expected missing achievement-id error, got %q", stderr)
	}
}

func TestGameCenterAchievementLocalizationsCreateValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "missing achievement-id",
			args:    []string{"game-center", "achievements", "localizations", "create", "--locale", "en-US", "--name", "Test", "--before-earned-description", "Before", "--after-earned-description", "After"},
			wantErr: "--achievement-id is required",
		},
		{
			name:    "missing locale",
			args:    []string{"game-center", "achievements", "localizations", "create", "--achievement-id", "ACH_ID", "--name", "Test", "--before-earned-description", "Before", "--after-earned-description", "After"},
			wantErr: "--locale is required",
		},
		{
			name:    "missing name",
			args:    []string{"game-center", "achievements", "localizations", "create", "--achievement-id", "ACH_ID", "--locale", "en-US", "--before-earned-description", "Before", "--after-earned-description", "After"},
			wantErr: "--name is required",
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

func TestGameCenterAchievementImagesUploadValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "missing localization-id",
			args:    []string{"game-center", "achievements", "images", "upload", "--file", "test.png"},
			wantErr: "--localization-id is required",
		},
		{
			name:    "missing file",
			args:    []string{"game-center", "achievements", "images", "upload", "--localization-id", "LOC_ID"},
			wantErr: "--file is required",
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

func TestGameCenterAchievementReleasesListValidationErrors(t *testing.T) {
	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	stdout, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{"game-center", "achievements", "releases", "list"}); err != nil {
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
	if !strings.Contains(stderr, "--achievement-id is required") {
		t.Fatalf("expected missing achievement-id error, got %q", stderr)
	}
}

func TestGameCenterAchievementsListLimitValidation(t *testing.T) {
	t.Setenv("ASC_APP_ID", "APP_ID")

	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	stdout, _ := captureOutput(t, func() {
		if err := root.Parse([]string{"game-center", "achievements", "list", "--app", "APP_ID", "--limit", "201"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		err := root.Run(context.Background())
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "--limit must be between 1 and 200") {
			t.Fatalf("expected error containing %q, got %v", "--limit must be between 1 and 200", err)
		}
	})

	if stdout != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
}

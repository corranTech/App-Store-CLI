package cmd

import (
	"context"
	"errors"
	"flag"
	"io"
	"strings"
	"testing"
)

func TestGameCenterLeaderboardSetsListValidationErrors(t *testing.T) {
	t.Setenv("ASC_APP_ID", "")

	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	stdout, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{"game-center", "leaderboard-sets", "list"}); err != nil {
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

func TestGameCenterLeaderboardSetsCreateValidationErrors(t *testing.T) {
	t.Setenv("ASC_APP_ID", "")

	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "missing app",
			args:    []string{"game-center", "leaderboard-sets", "create", "--reference-name", "Test", "--vendor-id", "com.test"},
			wantErr: "--app is required",
		},
		{
			name:    "missing reference-name",
			args:    []string{"game-center", "leaderboard-sets", "create", "--app", "APP_ID", "--vendor-id", "com.test"},
			wantErr: "--reference-name is required",
		},
		{
			name:    "missing vendor-id",
			args:    []string{"game-center", "leaderboard-sets", "create", "--app", "APP_ID", "--reference-name", "Test"},
			wantErr: "--vendor-id is required",
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

func TestGameCenterLeaderboardSetMembersListValidationErrors(t *testing.T) {
	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	stdout, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{"game-center", "leaderboard-sets", "members", "list"}); err != nil {
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
	if !strings.Contains(stderr, "--set-id is required") {
		t.Fatalf("expected missing set-id error, got %q", stderr)
	}
}

func TestGameCenterLeaderboardSetsListLimitValidation(t *testing.T) {
	t.Setenv("ASC_APP_ID", "APP_ID")

	root := RootCommand("1.2.3")
	root.FlagSet.SetOutput(io.Discard)

	stdout, _ := captureOutput(t, func() {
		if err := root.Parse([]string{"game-center", "leaderboard-sets", "list", "--app", "APP_ID", "--limit", "500"}); err != nil {
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

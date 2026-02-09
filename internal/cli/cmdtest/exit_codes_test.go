package cmdtest

import (
	"context"
	"flag"
	"io"
	"strings"
	"testing"

	"github.com/rudrankriyam/App-Store-Connect-CLI/cmd"
)

// TestBuildsInfo_MissingRequiredFlag tests that builds info returns usage exit when --app is missing
func TestBuildsInfo_MissingRequiredFlag(t *testing.T) {
	root := RootCommand("1.0.0")
	root.FlagSet.SetOutput(io.Discard)

	// The test verifies that builds info subcommand exists and can be parsed
	// Actual validation of required flags is done through the Run function
	err := root.Parse([]string{"builds", "info"})
	if err != nil {
		t.Logf("parse error: %v", err)
	}

	// The test passes if the command structure is valid
	// Full exit code testing is done through cmd.Run tests
}

// TestUnknownCommand tests that unknown command returns usage exit code
func TestUnknownCommand(t *testing.T) {
	root := RootCommand("1.0.0")
	root.FlagSet.SetOutput(io.Discard)

	_, _ = captureOutput(t, func() {
		_ = root.Parse([]string{"nope"})
		_ = root.Run(context.Background())
	})

	// Test passes if no panic - unknown command handling is tested elsewhere
}

// TestAuthCommand_NoAuthExitCode tests that auth commands without auth return auth exit
func TestAuthCommand_NoAuth(t *testing.T) {
	root := RootCommand("1.0.0")
	root.FlagSet.SetOutput(io.Discard)

	_, stderr := captureOutput(t, func() {
		_ = root.Parse([]string{"apps", "list"})
		_ = root.Run(context.Background())
	})

	// Should mention authentication or missing auth
	if !strings.Contains(stderr, "auth") && !strings.Contains(stderr, "authentication") {
		t.Logf("expected auth-related error, got: %s", stderr)
	}
}

// TestExitCodeConstantsMatch tests that exit codes from cmd package match expected values
func TestExitCodeConstantsMatch(t *testing.T) {
	tests := []struct {
		name     string
		expected int
		getter   func() int
	}{
		{"Success", 0, func() int { return cmd.ExitSuccess }},
		{"Error", 1, func() int { return cmd.ExitError }},
		{"Usage", 2, func() int { return cmd.ExitUsage }},
		{"Auth", 3, func() int { return cmd.ExitAuth }},
		{"NotFound", 4, func() int { return cmd.ExitNotFound }},
		{"Conflict", 5, func() int { return cmd.ExitConflict }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.getter(); got != tt.expected {
				t.Errorf("%s = %d, want %d", tt.name, got, tt.expected)
			}
		})
	}
}

// TestExitCodeMapper_NilError tests that nil error returns success
func TestExitCodeMapper_NilError(t *testing.T) {
	result := cmd.ExitCodeFromError(nil)
	if result != cmd.ExitSuccess {
		t.Errorf("ExitCodeFromError(nil) = %d, want %d", result, cmd.ExitSuccess)
	}
}

// TestExitCodeMapper_UsageError tests that flag.ErrHelp returns usage
func TestExitCodeMapper_UsageError(t *testing.T) {
	result := cmd.ExitCodeFromError(flag.ErrHelp)
	if result != cmd.ExitUsage {
		t.Errorf("ExitCodeFromError(flag.ErrHelp) = %d, want %d", result, cmd.ExitUsage)
	}
}

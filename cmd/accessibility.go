package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	accessibilitycli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/accessibility"
)

// AccessibilityCommand returns the accessibility command group.
func AccessibilityCommand() *ffcli.Command {
	return accessibilitycli.AccessibilityCommand()
}

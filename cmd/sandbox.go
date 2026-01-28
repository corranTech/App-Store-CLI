package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	sandboxcli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/sandbox"
)

// SandboxCommand returns the sandbox command group.
func SandboxCommand() *ffcli.Command {
	return sandboxcli.SandboxCommand()
}

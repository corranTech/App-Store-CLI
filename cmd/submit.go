package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	submitcli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/submit"
)

// SubmitCommand returns the submit command group.
func SubmitCommand() *ffcli.Command {
	return submitcli.SubmitCommand()
}

package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	nominationscli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/nominations"
)

// NominationsCommand returns the nominations command group.
func NominationsCommand() *ffcli.Command {
	return nominationscli.NominationsCommand()
}

package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	eulacli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/eula"
)

// EULACommand returns the end user license agreement command group.
func EULACommand() *ffcli.Command {
	return eulacli.EULACommand()
}

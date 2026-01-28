package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	financecli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/finance"
)

// FinanceCommand returns the finance command group.
func FinanceCommand() *ffcli.Command {
	return financecli.FinanceCommand()
}

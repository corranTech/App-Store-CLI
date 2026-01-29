package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	migratecli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/migrate"
)

// MigrateCommand returns the migrate command group.
func MigrateCommand() *ffcli.Command {
	return migratecli.MigrateCommand()
}

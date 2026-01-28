package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	appscli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/apps"
)

// AppSetupCommand returns the app-setup command group.
func AppSetupCommand() *ffcli.Command {
	return appscli.AppSetupCommand()
}

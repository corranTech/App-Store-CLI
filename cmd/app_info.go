package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	appscli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/apps"
)

// AppInfoCommand returns the app-info command group.
func AppInfoCommand() *ffcli.Command {
	return appscli.AppInfoCommand()
}

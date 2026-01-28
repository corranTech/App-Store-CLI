package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	appscli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/apps"
)

// AppsCommand returns the apps command group.
func AppsCommand() *ffcli.Command {
	return appscli.AppsCommand()
}

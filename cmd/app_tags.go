package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	appscli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/apps"
)

// AppTagsCommand returns the app-tags command group.
func AppTagsCommand() *ffcli.Command {
	return appscli.AppTagsCommand()
}

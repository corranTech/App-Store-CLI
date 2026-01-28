package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	buildscli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/builds"
)

// BuildsCommand returns the builds command group.
func BuildsCommand() *ffcli.Command {
	return buildscli.BuildsCommand()
}

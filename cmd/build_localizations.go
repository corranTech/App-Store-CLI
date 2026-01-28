package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	buildlocalizationscli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/buildlocalizations"
)

// BuildLocalizationsCommand returns the build-localizations command group.
func BuildLocalizationsCommand() *ffcli.Command {
	return buildlocalizationscli.BuildLocalizationsCommand()
}

package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	localizationscli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/localizations"
)

// LocalizationsCommand returns the localizations command group.
func LocalizationsCommand() *ffcli.Command {
	return localizationscli.LocalizationsCommand()
}

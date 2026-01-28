package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	assetscli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/assets"
)

// AssetsCommand returns the assets command group.
func AssetsCommand() *ffcli.Command {
	return assetscli.AssetsCommand()
}

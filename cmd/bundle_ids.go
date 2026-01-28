package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	bundleidscli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/bundleids"
)

// BundleIDsCommand returns the bundle-ids command group.
func BundleIDsCommand() *ffcli.Command {
	return bundleidscli.BundleIDsCommand()
}

package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	xcodecloudcli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/xcodecloud"
)

// XcodeCloudCommand returns the xcode-cloud command group.
func XcodeCloudCommand() *ffcli.Command {
	return xcodecloudcli.XcodeCloudCommand()
}

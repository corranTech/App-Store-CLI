package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	devicescli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/devices"
)

// DevicesCommand returns the devices command group.
func DevicesCommand() *ffcli.Command {
	return devicescli.DevicesCommand()
}

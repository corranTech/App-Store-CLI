package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	performancecli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/performance"
)

// PerformanceCommand returns the performance command group.
func PerformanceCommand() *ffcli.Command {
	return performancecli.PerformanceCommand()
}

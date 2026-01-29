package backgroundassets

import "github.com/peterbourgon/ff/v3/ffcli"

// Command returns the background assets command group.
func Command() *ffcli.Command {
	return BackgroundAssetsCommand()
}

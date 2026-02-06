package notarization

import "github.com/peterbourgon/ff/v3/ffcli"

// NotarizationCommand returns the notarization command group.
func NotarizationCommand() *ffcli.Command {
	return notarizationCommand()
}

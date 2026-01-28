package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	authcli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/auth"
)

// AuthCommand returns the auth subcommand.
func AuthCommand() *ffcli.Command {
	return authcli.AuthCommand()
}

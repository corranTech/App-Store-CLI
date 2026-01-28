package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	profilescli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/profiles"
)

// ProfilesCommand returns the profiles command group.
func ProfilesCommand() *ffcli.Command {
	return profilescli.ProfilesCommand()
}

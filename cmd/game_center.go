package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	gamecentercli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/gamecenter"
)

// GameCenterCommand returns the game-center command group.
func GameCenterCommand() *ffcli.Command {
	return gamecentercli.GameCenterCommand()
}

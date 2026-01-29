package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	ageratingcli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/agerating"
)

// AgeRatingCommand returns the age-rating command group.
func AgeRatingCommand() *ffcli.Command {
	return ageratingcli.AgeRatingCommand()
}

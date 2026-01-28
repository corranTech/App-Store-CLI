package cmd

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	offercodescli "github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/offercodes"
)

// OfferCodesCommand returns the offer-codes command group.
func OfferCodesCommand() *ffcli.Command {
	return offercodescli.OfferCodesCommand()
}

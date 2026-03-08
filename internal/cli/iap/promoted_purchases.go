package iap

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/promotedpurchases"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
)

// IAPPromotedPurchasesCommand returns the canonical nested promoted purchases tree.
func IAPPromotedPurchasesCommand() *ffcli.Command {
	cmd := shared.RewriteCommandTreePath(
		promotedpurchases.PromotedPurchasesCommand(),
		"asc promoted-purchases",
		"asc iap promoted-purchases",
	)
	if cmd != nil {
		cmd.ShortHelp = "Manage promoted purchases for in-app purchases."
		configureIAPPromotedPurchasesCreate(cmd)
	}
	return cmd
}

func configureIAPPromotedPurchasesCreate(cmd *ffcli.Command) {
	promotedpurchases.ConfigureFixedProductTypeCreateCommand(cmd, promotedpurchases.FixedProductTypeCreateConfig{
		ShortUsage: "asc iap promoted-purchases create --app APP_ID --product-id PRODUCT_ID --visible-for-all-users",
		ShortHelp:  "Create a promoted purchase for an in-app purchase.",
		LongHelp: `Create a promoted purchase for an in-app purchase.

Examples:
  asc iap promoted-purchases create --app "APP_ID" --product-id "IAP_ID" --visible-for-all-users true
  asc iap promoted-purchases create --app "APP_ID" --product-id "IAP_ID" --visible-for-all-users true --enabled true`,
		ProductType:    "IN_APP_PURCHASE",
		ProductIDUsage: "In-app purchase ID",
	})
}

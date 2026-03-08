package subscriptions

import (
	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/promotedpurchases"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
)

// SubscriptionsPromotedPurchasesCommand returns the canonical nested promoted purchases tree.
func SubscriptionsPromotedPurchasesCommand() *ffcli.Command {
	cmd := shared.RewriteCommandTreePath(
		promotedpurchases.PromotedPurchasesCommand(),
		"asc promoted-purchases",
		"asc subscriptions promoted-purchases",
	)
	if cmd != nil {
		cmd.ShortHelp = "Manage promoted purchases for subscriptions."
		configureSubscriptionsPromotedPurchasesCreate(cmd)
	}
	return cmd
}

func configureSubscriptionsPromotedPurchasesCreate(cmd *ffcli.Command) {
	promotedpurchases.ConfigureFixedProductTypeCreateCommand(cmd, promotedpurchases.FixedProductTypeCreateConfig{
		ShortUsage: "asc subscriptions promoted-purchases create --app APP_ID --product-id PRODUCT_ID --visible-for-all-users",
		ShortHelp:  "Create a promoted purchase for a subscription.",
		LongHelp: `Create a promoted purchase for a subscription.

Examples:
  asc subscriptions promoted-purchases create --app "APP_ID" --product-id "SUB_ID" --visible-for-all-users true
  asc subscriptions promoted-purchases create --app "APP_ID" --product-id "SUB_ID" --visible-for-all-users true --enabled true`,
		ProductType:    "SUBSCRIPTION",
		ProductIDUsage: "Subscription ID",
	})
}

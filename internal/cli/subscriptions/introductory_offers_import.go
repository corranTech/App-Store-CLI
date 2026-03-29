package subscriptions

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
)

// SubscriptionsIntroductoryOffersImportCommand returns the introductory offers import subcommand.
func SubscriptionsIntroductoryOffersImportCommand() *ffcli.Command {
	fs := flag.NewFlagSet("introductory-offers import", flag.ExitOnError)

	subscriptionID := fs.String("subscription-id", "", "Subscription ID")
	inputPath := fs.String("input", "", "Input CSV file path (required)")
	_ = fs.String("offer-duration", "", "Default offer duration")
	_ = fs.String("offer-mode", "", "Default offer mode")
	_ = fs.Int("number-of-periods", 0, "Default number of periods")
	_ = fs.String("start-date", "", "Default start date (YYYY-MM-DD)")
	_ = fs.String("end-date", "", "Default end date (YYYY-MM-DD)")
	_ = fs.Bool("dry-run", false, "Validate input and print summary without creating offers")
	_ = fs.Bool("continue-on-error", true, "Continue processing rows after runtime failures (default true)")
	output := shared.BindOutputFlags(fs)

	return &ffcli.Command{
		Name:       "import",
		ShortUsage: "asc subscriptions introductory-offers import --subscription-id \"SUB_ID\" --input \"./offers.csv\" [flags]",
		ShortHelp:  "Import introductory offers from a CSV file.",
		LongHelp: `Import introductory offers from a CSV file.

Examples:
  asc subscriptions introductory-offers import --subscription-id "SUB_ID" --input "./offers.csv"
  asc subscriptions introductory-offers import --subscription-id "SUB_ID" --input "./offers.csv" --offer-duration ONE_WEEK --offer-mode FREE_TRIAL --number-of-periods 1`,
		FlagSet:   fs,
		UsageFunc: shared.DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			_ = ctx
			_ = output

			if strings.TrimSpace(*subscriptionID) == "" {
				fmt.Fprintln(os.Stderr, "Error: --subscription-id is required")
				return flag.ErrHelp
			}
			if strings.TrimSpace(*inputPath) == "" {
				fmt.Fprintln(os.Stderr, "Error: --input is required")
				return flag.ErrHelp
			}

			return shared.UsageError("introductory-offers import is not implemented yet")
		},
	}
}

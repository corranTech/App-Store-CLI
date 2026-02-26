package shared

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
)

// PricingSetCommandConfig configures pricing set commands.
type PricingSetCommandConfig struct {
	FlagSetName           string
	CommandName           string
	ShortUsage            string
	ShortHelp             string
	LongHelp              string
	ErrorPrefix           string
	StartDateHelp         string
	StartDateDefaultToday bool
	RequireBaseTerritory  bool
	ResolveBaseTerritory  bool
}

// NewPricingSetCommand builds a pricing set command with shared behavior.
func NewPricingSetCommand(config PricingSetCommandConfig) *ffcli.Command {
	fs := flag.NewFlagSet(config.FlagSetName, flag.ExitOnError)

	appID := fs.String("app", "", "App Store Connect app ID (or ASC_APP_ID)")
	pricePointID := fs.String("price-point", "", "App price point ID")
	price := fs.String("price", "", "Customer price (e.g., 0.99) to select price point")
	baseTerritory := fs.String("base-territory", "", "Base territory ID (e.g., USA)")
	startDate := fs.String("start-date", "", config.StartDateHelp)
	output := BindOutputFlags(fs)

	return &ffcli.Command{
		Name:       config.CommandName,
		ShortUsage: config.ShortUsage,
		ShortHelp:  config.ShortHelp,
		LongHelp:   config.LongHelp,
		FlagSet:    fs,
		UsageFunc:  DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			resolvedAppID := resolveAppID(*appID)
			if resolvedAppID == "" {
				fmt.Fprintln(os.Stderr, "Error: --app is required (or set ASC_APP_ID)")
				return flag.ErrHelp
			}
			pricePointValue := strings.TrimSpace(*pricePointID)
			priceValue := strings.TrimSpace(*price)

			if pricePointValue == "" && priceValue == "" {
				fmt.Fprintln(os.Stderr, "Error: --price-point or --price is required")
				return flag.ErrHelp
			}
			if pricePointValue != "" && priceValue != "" {
				fmt.Fprintln(os.Stderr, "Error: --price-point and --price are mutually exclusive")
				return flag.ErrHelp
			}
			if priceValue != "" {
				if _, err := parseFinitePrice(priceValue); err != nil {
					if errors.Is(err, errNonFinitePrice) {
						fmt.Fprintln(os.Stderr, "Error: --price must be a finite number")
					} else {
						fmt.Fprintln(os.Stderr, "Error: --price must be a number")
					}
					return flag.ErrHelp
				}
			}

			baseTerritoryValue := strings.TrimSpace(*baseTerritory)
			if (config.RequireBaseTerritory || priceValue != "") && baseTerritoryValue == "" {
				fmt.Fprintln(os.Stderr, "Error: --base-territory is required")
				return flag.ErrHelp
			}

			startDateValue := strings.TrimSpace(*startDate)
			if startDateValue == "" {
				if config.StartDateDefaultToday {
					startDateValue = time.Now().Format("2006-01-02")
				} else {
					fmt.Fprintln(os.Stderr, "Error: --start-date is required")
					return flag.ErrHelp
				}
			}

			normalizedStartDate, err := normalizePricingStartDate(startDateValue)
			if err != nil {
				return fmt.Errorf("%s: %w", config.ErrorPrefix, err)
			}

			client, err := getASCClient()
			if err != nil {
				return fmt.Errorf("%s: %w", config.ErrorPrefix, err)
			}

			requestCtx, cancel := contextWithTimeout(ctx)
			defer cancel()

			baseTerritoryID := baseTerritoryValue
			if config.ResolveBaseTerritory {
				baseTerritoryID, err = resolveBaseTerritoryID(requestCtx, client, resolvedAppID, baseTerritoryValue)
				if err != nil {
					return fmt.Errorf("%s: %w", config.ErrorPrefix, err)
				}
			}

			if priceValue != "" {
				priceFilter := PriceFilter{Price: priceValue}
				foundID := ""
				fetch := func(nextURL string) (*asc.AppPricePointsV3Response, error) {
					opts := []asc.PricePointsOption{
						asc.WithPricePointsLimit(200),
						asc.WithPricePointsTerritory(baseTerritoryID),
					}
					if nextURL != "" {
						opts = append(opts, asc.WithPricePointsNextURL(nextURL))
					}
					return client.GetAppPricePoints(requestCtx, resolvedAppID, opts...)
				}

				resp, err := fetch("")
				if err != nil {
					return fmt.Errorf("resolve price: %w", err)
				}

				for {
					for _, pp := range resp.Data {
						if priceFilter.MatchesPrice(pp.Attributes.CustomerPrice) {
							foundID = pp.ID
							break
						}
					}
					if foundID != "" {
						break
					}
					if resp.Links.Next == "" {
						break
					}
					resp, err = fetch(resp.Links.Next)
					if err != nil {
						return fmt.Errorf("resolve price (next page): %w", err)
					}
				}

				if foundID == "" {
					return fmt.Errorf("price point not found for price %q in territory %q", priceValue, baseTerritoryID)
				}
				pricePointValue = foundID
			}

			resp, err := client.CreateAppPriceSchedule(requestCtx, resolvedAppID, asc.AppPriceScheduleCreateAttributes{
				PricePointID:    pricePointValue,
				StartDate:       normalizedStartDate,
				BaseTerritoryID: baseTerritoryID,
			})
			if err != nil {
				return fmt.Errorf("%s: %w", config.ErrorPrefix, err)
			}

			return printOutput(resp, *output.Output, *output.Pretty)
		},
	}
}

func normalizePricingStartDate(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", fmt.Errorf("--start-date is required")
	}
	parsed, err := time.Parse("2006-01-02", trimmed)
	if err != nil {
		return "", fmt.Errorf("--start-date must be in YYYY-MM-DD format")
	}
	return parsed.Format("2006-01-02"), nil
}

func resolveBaseTerritoryID(ctx context.Context, client *asc.Client, appID string, baseTerritory string) (string, error) {
	trimmed := strings.ToUpper(strings.TrimSpace(baseTerritory))
	if trimmed != "" {
		return trimmed, nil
	}

	schedule, err := client.GetAppPriceSchedule(ctx, appID)
	if err != nil {
		if asc.IsNotFound(err) {
			return "", fmt.Errorf("--base-territory is required when app price schedule is missing")
		}
		return "", fmt.Errorf("get app price schedule: %w", err)
	}

	scheduleID := strings.TrimSpace(schedule.Data.ID)
	if scheduleID == "" {
		return "", fmt.Errorf("app price schedule ID missing")
	}

	territoryResp, err := client.GetAppPriceScheduleBaseTerritory(ctx, scheduleID)
	if err != nil {
		return "", fmt.Errorf("get base territory: %w", err)
	}

	territoryID := strings.ToUpper(strings.TrimSpace(territoryResp.Data.ID))
	if territoryID == "" {
		return "", fmt.Errorf("base territory ID missing from response")
	}

	return territoryID, nil
}

package shared

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
)

// TierEntry represents a single tier in a territory's price point list.
type TierEntry struct {
	Tier          int    `json:"tier"`
	PricePointID  string `json:"pricePointId"`
	CustomerPrice string `json:"customerPrice"`
	Proceeds      string `json:"proceeds"`
}

// ResolveTiers fetches all price points for a territory, sorts by customerPrice ascending,
// and assigns tier numbers starting at 1. Free (0.00) price points are excluded.
func ResolveTiers(ctx context.Context, client *asc.Client, appID, territory string, refresh bool) ([]TierEntry, error) {
	appID = strings.TrimSpace(appID)
	territory = strings.ToUpper(strings.TrimSpace(territory))
	if appID == "" {
		return nil, fmt.Errorf("app ID is required for tier resolution")
	}
	if territory == "" {
		return nil, fmt.Errorf("territory is required for tier resolution")
	}

	if !refresh {
		cached, err := LoadTierCache(appID, territory)
		if err == nil && len(cached) > 0 {
			return cached, nil
		}
	}

	type pricePointEntry struct {
		id            string
		customerPrice float64
		rawPrice      string
		proceeds      string
	}

	var entries []pricePointEntry
	var nextURL string

	for {
		opts := []asc.PricePointsOption{
			asc.WithPricePointsLimit(200),
			asc.WithPricePointsTerritory(territory),
		}
		if nextURL != "" {
			opts = []asc.PricePointsOption{asc.WithPricePointsNextURL(nextURL)}
		}

		resp, err := client.GetAppPricePoints(ctx, appID, opts...)
		if err != nil {
			return nil, fmt.Errorf("fetch price points: %w", err)
		}

		for _, pp := range resp.Data {
			price, err := strconv.ParseFloat(strings.TrimSpace(pp.Attributes.CustomerPrice), 64)
			if err != nil || price <= 0 {
				continue
			}
			entries = append(entries, pricePointEntry{
				id:            pp.ID,
				customerPrice: price,
				rawPrice:      pp.Attributes.CustomerPrice,
				proceeds:      pp.Attributes.Proceeds,
			})
		}

		if resp.Links.Next == "" {
			break
		}
		nextURL = resp.Links.Next
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].customerPrice < entries[j].customerPrice
	})

	tiers := make([]TierEntry, 0, len(entries))
	for i, e := range entries {
		tiers = append(tiers, TierEntry{
			Tier:          i + 1,
			PricePointID:  e.id,
			CustomerPrice: e.rawPrice,
			Proceeds:      e.proceeds,
		})
	}

	if len(tiers) > 0 {
		_ = SaveTierCache(appID, territory, tiers)
	}

	return tiers, nil
}

// ResolvePricePointByTier finds the price point ID for a given tier number.
func ResolvePricePointByTier(tiers []TierEntry, tier int) (string, error) {
	for _, t := range tiers {
		if t.Tier == tier {
			return t.PricePointID, nil
		}
	}
	return "", fmt.Errorf("tier %d not found (valid range: 1-%d)", tier, len(tiers))
}

// ResolvePricePointByPrice finds the price point ID for a given customer price.
func ResolvePricePointByPrice(tiers []TierEntry, price string) (string, error) {
	target, err := strconv.ParseFloat(strings.TrimSpace(price), 64)
	if err != nil {
		return "", fmt.Errorf("invalid price %q: %w", price, err)
	}
	if math.IsNaN(target) || math.IsInf(target, 0) {
		return "", fmt.Errorf("price must be a finite number")
	}

	for _, t := range tiers {
		cp, err := strconv.ParseFloat(strings.TrimSpace(t.CustomerPrice), 64)
		if err != nil {
			continue
		}
		if math.Abs(cp-target) < 0.005 {
			return t.PricePointID, nil
		}
	}
	return "", fmt.Errorf("no price point found matching price %s in this territory", price)
}

// ValidatePriceSelectionFlags checks that --price-point, --tier, and --price are mutually exclusive.
// Returns a usage-style error if more than one is set.
func ValidatePriceSelectionFlags(pricePoint string, tier int, price string) error {
	count := 0
	if pricePoint != "" {
		count++
	}
	if tier > 0 {
		count++
	}
	if price != "" {
		count++
	}
	if count == 0 {
		return fmt.Errorf("one of --price-point, --tier, or --price is required")
	}
	if count > 1 {
		return fmt.Errorf("--price-point, --tier, and --price are mutually exclusive")
	}
	return nil
}

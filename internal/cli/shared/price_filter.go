package shared

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// PriceFilter holds parsed price filter criteria.
type PriceFilter struct {
	Price    string
	MinPrice string
	MaxPrice string
}

// HasFilter returns true if any price filter is set.
func (pf PriceFilter) HasFilter() bool {
	return pf.Price != "" || pf.MinPrice != "" || pf.MaxPrice != ""
}

// Validate checks that filter values are valid numbers and not contradictory.
func (pf PriceFilter) Validate() error {
	if pf.Price != "" && (pf.MinPrice != "" || pf.MaxPrice != "") {
		return fmt.Errorf("--price and --min-price/--max-price are mutually exclusive")
	}
	if pf.Price != "" {
		if _, err := strconv.ParseFloat(pf.Price, 64); err != nil {
			return fmt.Errorf("--price must be a number: %w", err)
		}
	}
	if pf.MinPrice != "" {
		if _, err := strconv.ParseFloat(pf.MinPrice, 64); err != nil {
			return fmt.Errorf("--min-price must be a number: %w", err)
		}
	}
	if pf.MaxPrice != "" {
		if _, err := strconv.ParseFloat(pf.MaxPrice, 64); err != nil {
			return fmt.Errorf("--max-price must be a number: %w", err)
		}
	}
	if pf.MinPrice != "" && pf.MaxPrice != "" {
		min, _ := strconv.ParseFloat(pf.MinPrice, 64)
		max, _ := strconv.ParseFloat(pf.MaxPrice, 64)
		if min > max {
			return fmt.Errorf("--min-price (%s) cannot exceed --max-price (%s)", pf.MinPrice, pf.MaxPrice)
		}
	}
	return nil
}

// MatchesPrice returns true if the given customerPrice string passes the filter.
func (pf PriceFilter) MatchesPrice(customerPrice string) bool {
	if !pf.HasFilter() {
		return true
	}
	price, err := strconv.ParseFloat(strings.TrimSpace(customerPrice), 64)
	if err != nil {
		return false
	}
	if pf.Price != "" {
		target, _ := strconv.ParseFloat(pf.Price, 64)
		return math.Abs(price-target) < 0.005
	}
	if pf.MinPrice != "" {
		min, _ := strconv.ParseFloat(pf.MinPrice, 64)
		if price < min-0.005 {
			return false
		}
	}
	if pf.MaxPrice != "" {
		max, _ := strconv.ParseFloat(pf.MaxPrice, 64)
		if price > max+0.005 {
			return false
		}
	}
	return true
}

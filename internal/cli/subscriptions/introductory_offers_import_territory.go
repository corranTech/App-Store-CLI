package subscriptions

import (
	"fmt"
	"strings"
	"sync"
	"unicode"

	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

type introductoryOfferTerritoryNameMapResult struct {
	id        string
	ambiguous bool
}

var (
	introductoryOfferTerritoryNamesOnce sync.Once
	introductoryOfferTerritoryNames     map[string]introductoryOfferTerritoryNameMapResult
)

func normalizeSubscriptionIntroductoryOfferImportTerritoryID(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", fmt.Errorf("territory is required")
	}

	upper := strings.ToUpper(trimmed)
	if isThreeLetterCode(upper) {
		return upper, nil
	}
	if len(upper) == 2 {
		if region, err := language.ParseRegion(upper); err == nil {
			if iso3 := strings.ToUpper(strings.TrimSpace(region.ISO3())); isThreeLetterCode(iso3) {
				return iso3, nil
			}
		}
	}

	key := normalizeSubscriptionIntroductoryOfferImportTerritoryName(trimmed)
	entry, ok := subscriptionIntroductoryOfferImportTerritoryNameMap()[key]
	if !ok {
		return "", fmt.Errorf("territory %q could not be mapped to an App Store Connect territory ID", trimmed)
	}
	if entry.ambiguous || entry.id == "" {
		return "", fmt.Errorf("territory %q is ambiguous; use a 3-letter territory ID like USA", trimmed)
	}
	return entry.id, nil
}

func subscriptionIntroductoryOfferImportTerritoryNameMap() map[string]introductoryOfferTerritoryNameMapResult {
	introductoryOfferTerritoryNamesOnce.Do(func() {
		m := make(map[string]introductoryOfferTerritoryNameMapResult)
		regionNamer := display.English.Regions()

		for a := 'A'; a <= 'Z'; a++ {
			for b := 'A'; b <= 'Z'; b++ {
				for c := 'A'; c <= 'Z'; c++ {
					code := string([]rune{a, b, c})
					region, err := language.ParseRegion(code)
					if err != nil {
						continue
					}
					iso3 := strings.ToUpper(strings.TrimSpace(region.ISO3()))
					if iso3 != code {
						continue
					}
					name := strings.TrimSpace(regionNamer.Name(region))
					if name == "" || strings.EqualFold(name, code) {
						continue
					}
					key := normalizeSubscriptionIntroductoryOfferImportTerritoryName(name)
					if key == "" {
						continue
					}
					existing, exists := m[key]
					switch {
					case !exists:
						m[key] = introductoryOfferTerritoryNameMapResult{id: iso3}
					case existing.id != iso3:
						m[key] = introductoryOfferTerritoryNameMapResult{ambiguous: true}
					}
				}
			}
		}

		introductoryOfferTerritoryNames = m
	})

	return introductoryOfferTerritoryNames
}

func normalizeSubscriptionIntroductoryOfferImportTerritoryName(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return ""
	}

	var b strings.Builder
	b.Grow(len(value))
	lastUnderscore := false
	for _, r := range value {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			b.WriteRune(r)
			lastUnderscore = false
		case lastUnderscore:
			continue
		default:
			b.WriteByte('_')
			lastUnderscore = true
		}
	}

	return strings.Trim(b.String(), "_")
}

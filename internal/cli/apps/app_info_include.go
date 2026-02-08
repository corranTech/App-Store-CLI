package apps

import (
	"fmt"
	"strings"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
)

func normalizeAppInfoInclude(value string) ([]string, error) {
	return normalizeInclude(value, appInfoIncludeList(), "--include")
}

func normalizeInclude(value string, allowed []string, flagName string) ([]string, error) {
	include := shared.SplitCSV(value)
	if len(include) == 0 {
		return nil, nil
	}
	allowedMap := map[string]struct{}{}
	for _, option := range allowed {
		allowedMap[option] = struct{}{}
	}
	for _, option := range include {
		if _, ok := allowedMap[option]; !ok {
			return nil, fmt.Errorf("%s must be one of: %s", flagName, strings.Join(allowed, ", "))
		}
	}
	return include, nil
}

func appInfoIncludeList() []string {
	return []string{
		"ageRatingDeclaration",
		"territoryAgeRatings",
		"primaryCategory",
		"primarySubcategoryOne",
		"primarySubcategoryTwo",
		"secondaryCategory",
		"secondarySubcategoryOne",
		"secondarySubcategoryTwo",
	}
}

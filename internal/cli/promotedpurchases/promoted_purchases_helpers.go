package promotedpurchases

import (
	"fmt"
	"strconv"
	"strings"
)

type optionalBool struct {
	set   bool
	value bool
}

func (b *optionalBool) Set(value string) error {
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fmt.Errorf("must be true or false")
	}
	b.value = parsed
	b.set = true
	return nil
}

func (b *optionalBool) String() string {
	if !b.set {
		return ""
	}
	return strconv.FormatBool(b.value)
}

func (b *optionalBool) IsBoolFlag() bool {
	return true
}

type promotedPurchaseProductType string

const (
	promotedPurchaseProductTypeSubscription  promotedPurchaseProductType = "SUBSCRIPTION"
	promotedPurchaseProductTypeInAppPurchase promotedPurchaseProductType = "IN_APP_PURCHASE"
)

func normalizePromotedPurchaseProductType(value string) (promotedPurchaseProductType, error) {
	normalized := strings.ToUpper(strings.TrimSpace(value))
	normalized = strings.ReplaceAll(normalized, "-", "_")
	normalized = strings.ReplaceAll(normalized, " ", "_")
	switch normalized {
	case string(promotedPurchaseProductTypeSubscription):
		return promotedPurchaseProductTypeSubscription, nil
	case string(promotedPurchaseProductTypeInAppPurchase):
		return promotedPurchaseProductTypeInAppPurchase, nil
	default:
		return "", fmt.Errorf("--product-type must be one of: SUBSCRIPTION, IN_APP_PURCHASE")
	}
}

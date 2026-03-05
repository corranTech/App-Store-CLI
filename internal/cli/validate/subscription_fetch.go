package validate

import (
	"context"
	"fmt"
	"strings"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/validation"
)

func fetchSubscriptions(ctx context.Context, client *asc.Client, appID string) ([]validation.Subscription, error) {
	groupsCtx, groupsCancel := shared.ContextWithTimeout(ctx)
	groupsResp, err := client.GetSubscriptionGroups(groupsCtx, appID, asc.WithSubscriptionGroupsLimit(200))
	groupsCancel()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch subscription groups: %w", err)
	}

	paginatedGroups, err := asc.PaginateAll(ctx, groupsResp, func(_ context.Context, nextURL string) (asc.PaginatedResponse, error) {
		pageCtx, pageCancel := shared.ContextWithTimeout(ctx)
		defer pageCancel()
		return client.GetSubscriptionGroups(pageCtx, appID, asc.WithSubscriptionGroupsNextURL(nextURL))
	})
	if err != nil {
		return nil, fmt.Errorf("paginate subscription groups: %w", err)
	}

	groups, ok := paginatedGroups.(*asc.SubscriptionGroupsResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected subscription groups response type %T", paginatedGroups)
	}

	subscriptions := make([]validation.Subscription, 0)
	for _, group := range groups.Data {
		groupID := strings.TrimSpace(group.ID)
		if groupID == "" {
			continue
		}

		subsCtx, subsCancel := shared.ContextWithTimeout(ctx)
		subsResp, err := client.GetSubscriptions(subsCtx, groupID, asc.WithSubscriptionsLimit(200))
		subsCancel()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch subscriptions for group %s: %w", groupID, err)
		}

		paginatedSubs, err := asc.PaginateAll(ctx, subsResp, func(_ context.Context, nextURL string) (asc.PaginatedResponse, error) {
			pageCtx, pageCancel := shared.ContextWithTimeout(ctx)
			defer pageCancel()
			return client.GetSubscriptions(pageCtx, groupID, asc.WithSubscriptionsNextURL(nextURL))
		})
		if err != nil {
			return nil, fmt.Errorf("paginate subscriptions: %w", err)
		}

		subsResult, ok := paginatedSubs.(*asc.SubscriptionsResponse)
		if !ok {
			return nil, fmt.Errorf("unexpected subscriptions response type %T", paginatedSubs)
		}

		for _, sub := range subsResult.Data {
			hasImage, err := subscriptionHasImage(ctx, client, sub.ID)
			if err != nil {
				return nil, fmt.Errorf("fetch subscription images for %s: %w", strings.TrimSpace(sub.ID), err)
			}

			attrs := sub.Attributes
			subscriptions = append(subscriptions, validation.Subscription{
				ID:        sub.ID,
				Name:      attrs.Name,
				ProductID: attrs.ProductID,
				State:     attrs.State,
				GroupID:   groupID,
				HasImage:  hasImage,
			})
		}
	}

	return subscriptions, nil
}

func subscriptionHasImage(ctx context.Context, client *asc.Client, subscriptionID string) (bool, error) {
	requestCtx, cancel := shared.ContextWithTimeout(ctx)
	defer cancel()

	resp, err := client.GetSubscriptionImages(requestCtx, strings.TrimSpace(subscriptionID), asc.WithSubscriptionImagesLimit(1))
	if err != nil {
		if asc.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	return resp != nil && len(resp.Data) > 0, nil
}

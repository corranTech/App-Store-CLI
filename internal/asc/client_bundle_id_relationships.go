package asc

import (
	"context"
	"fmt"
)

// GetBundleIDApp retrieves the app for a bundle ID.
func (c *Client) GetBundleIDApp(ctx context.Context, bundleID string) (*AppResponse, error) {
	_ = ctx
	_ = bundleID
	return nil, fmt.Errorf("GetBundleIDApp not implemented")
}

// GetBundleIDProfiles retrieves profiles for a bundle ID.
func (c *Client) GetBundleIDProfiles(ctx context.Context, bundleID string, opts ...BundleIDProfilesOption) (*ProfilesResponse, error) {
	_ = ctx
	_ = bundleID
	_ = opts
	return nil, fmt.Errorf("GetBundleIDProfiles not implemented")
}

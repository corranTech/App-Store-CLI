package asc

import (
	"context"
	"fmt"
)

// GetUserInvitationVisibleApps retrieves visible apps for a user invitation.
func (c *Client) GetUserInvitationVisibleApps(ctx context.Context, invitationID string, opts ...UserInvitationVisibleAppsOption) (*AppsResponse, error) {
	_ = ctx
	_ = invitationID
	_ = opts
	return nil, fmt.Errorf("GetUserInvitationVisibleApps not implemented")
}

package web

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
	webcore "github.com/rudrankriyam/App-Store-Connect-CLI/internal/web"
	webref "github.com/rudrankriyam/App-Store-Connect-CLI/internal/web/reference"
)

var (
	resolveWebAuthCredentialsFn = shared.ResolveAuthCredentials
	newWebAuthClientFn          = webcore.NewClient
	lookupWebAuthKeyFn          = func(ctx context.Context, client *webcore.Client, keyID string) (*webcore.APIKeyRoleLookup, error) {
		return client.LookupAPIKeyRoles(ctx, keyID)
	}
	resolveWebAuthRefFn = webref.Resolve
)

type webAuthCapabilitiesResult struct {
	KeyID        string                    `json:"keyId"`
	Name         string                    `json:"name,omitempty"`
	Kind         string                    `json:"kind"`
	Roles        []string                  `json:"roles"`
	RoleSource   string                    `json:"roleSource"`
	Active       bool                      `json:"active"`
	KeyType      string                    `json:"keyType,omitempty"`
	LastUsed     string                    `json:"lastUsed,omitempty"`
	Lookup       string                    `json:"lookup"`
	ResolvedFrom string                    `json:"resolvedFrom"`
	Profile      string                    `json:"profile,omitempty"`
	GeneratedBy  *webcoreKeyActorResult    `json:"generatedBy,omitempty"`
	RevokedBy    *webcoreKeyActorResult    `json:"revokedBy,omitempty"`
	RoleDetails  []webAuthRoleDetailResult `json:"roleDetails,omitempty"`
	Capabilities []webAuthCapabilityResult `json:"capabilities,omitempty"`
	UnknownRoles []string                  `json:"unknownRoles,omitempty"`
	LastVerified string                    `json:"referenceLastVerified,omitempty"`
}

type webcoreKeyActorResult struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type webAuthRoleDetailResult struct {
	Code                  string   `json:"code"`
	Label                 string   `json:"label"`
	UIAliases             []string `json:"uiAliases,omitempty"`
	TeamKeySelectable     bool     `json:"teamKeySelectable"`
	IndividualKeyEligible bool     `json:"individualKeyEligible"`
	Summary               string   `json:"summary,omitempty"`
	Capabilities          []string `json:"capabilities,omitempty"`
	Notes                 []string `json:"notes,omitempty"`
	ExampleTasks          []string `json:"exampleTasks,omitempty"`
	NotableActions        []string `json:"notableExclusiveActions,omitempty"`
}

type webAuthCapabilityResult struct {
	ID      string `json:"id"`
	Label   string `json:"label"`
	Summary string `json:"summary,omitempty"`
}

// WebAuthCapabilitiesCommand returns exact key-role lookup via App Store Connect web-session endpoints.
func WebAuthCapabilitiesCommand() *ffcli.Command {
	fs := flag.NewFlagSet("web auth capabilities", flag.ExitOnError)

	authFlags := bindWebSessionFlags(fs)
	keyID := fs.String("key-id", "", "API key ID to inspect (optional; defaults to the currently selected CLI API key)")
	output := shared.BindOutputFlags(fs)

	return &ffcli.Command{
		Name:       "capabilities",
		ShortUsage: "asc web auth capabilities [--key-id ID] [flags]",
		ShortHelp:  "[experimental] Show exact web-visible API key roles and documented capabilities.",
		LongHelp: `EXPERIMENTAL / UNOFFICIAL / DISCOURAGED

Return exact role metadata for an App Store Connect API key using Apple web-session endpoints, then map those roles to the bundled Apple capability reference.
Unlike "asc auth capabilities", which probes effective public-API access, this command reads the web-visible key role assignment directly and expands it with documented role capabilities.

If --key-id is omitted, the command resolves the current API key ID from local asc API auth and uses the active web session only for the exact web lookup.

` + webWarningText + `

Examples:
  asc web auth capabilities
  asc web auth capabilities --output json
  asc web auth capabilities --key-id "39MX87M9Y4" --output table`,
		FlagSet:   fs,
		UsageFunc: shared.DefaultUsageFunc,
		Exec: func(ctx context.Context, args []string) error {
			if len(args) > 0 {
				return shared.UsageError("web auth capabilities does not accept positional arguments")
			}
			if _, err := shared.ValidateOutputFormat(*output.Output, *output.Pretty); err != nil {
				return shared.UsageError(err.Error())
			}

			resolvedKeyID := strings.TrimSpace(*keyID)
			resolvedFrom := "flag"
			profile := ""
			if resolvedKeyID == "" {
				resolved, err := resolveWebAuthCredentialsFn("")
				if err != nil {
					return shared.UsageErrorf("unable to resolve current API key ID; run 'asc auth login' or provide --key-id (%v)", err)
				}
				resolvedKeyID = strings.TrimSpace(resolved.KeyID)
				profile = strings.TrimSpace(resolved.Profile)
				resolvedFrom = "auth"
			}
			if resolvedKeyID == "" {
				return shared.UsageError("unable to resolve current API key ID; run 'asc auth login' or provide --key-id")
			}

			requestCtx, cancel := shared.ContextWithTimeout(ctx)
			defer cancel()

			session, err := resolveWebSessionForCommand(requestCtx, authFlags)
			if err != nil {
				return err
			}

			client := newWebAuthClientFn(session)
			var lookup *webcore.APIKeyRoleLookup
			err = withWebSpinner("Loading exact API key roles", func() error {
				var innerErr error
				lookup, innerErr = lookupWebAuthKeyFn(requestCtx, client, resolvedKeyID)
				return innerErr
			})
			if err != nil {
				return wrapWebAuthCapabilitiesError(resolvedKeyID, err)
			}

			ref, err := resolveWebAuthRefFn(lookup.Roles)
			if err != nil {
				return fmt.Errorf("web auth capabilities failed: unable to load bundled role reference: %w", err)
			}

			result := webAuthCapabilitiesResult{
				KeyID:        lookup.KeyID,
				Name:         lookup.Name,
				Kind:         lookup.Kind,
				Roles:        append([]string(nil), lookup.Roles...),
				RoleSource:   lookup.RoleSource,
				Active:       lookup.Active,
				KeyType:      lookup.KeyType,
				LastUsed:     lookup.LastUsed,
				Lookup:       lookup.Lookup,
				ResolvedFrom: resolvedFrom,
				Profile:      profile,
				GeneratedBy:  convertKeyActor(lookup.GeneratedBy),
				RevokedBy:    convertKeyActor(lookup.RevokedBy),
				RoleDetails:  convertWebAuthRoleDetails(ref.RoleDetails),
				Capabilities: convertWebAuthCapabilities(ref.Capabilities),
				UnknownRoles: append([]string(nil), ref.UnknownRoles...),
				LastVerified: ref.LastVerified,
			}

			return shared.PrintOutputWithRenderers(
				result,
				*output.Output,
				*output.Pretty,
				func() error { return renderWebAuthCapabilitiesTable(result) },
				func() error { return renderWebAuthCapabilitiesMarkdown(result) },
			)
		},
	}
}

func convertKeyActor(actor *webcore.KeyActor) *webcoreKeyActorResult {
	if actor == nil {
		return nil
	}
	return &webcoreKeyActorResult{
		ID:   actor.ID,
		Name: actor.Name,
	}
}

func wrapWebAuthCapabilitiesError(keyID string, err error) error {
	if errors.Is(err, webcore.ErrAPIKeyNotFound) {
		return fmt.Errorf("web auth capabilities failed: key %q not found in App Store Connect web key lists", keyID)
	}
	if errors.Is(err, webcore.ErrAPIKeyNotVisible) {
		return fmt.Errorf("web auth capabilities failed: key %q is not visible in the accessible App Store Connect web key lists (team key list may be unavailable to this account)", keyID)
	}
	if errors.Is(err, webcore.ErrAPIKeyRolesUnresolved) {
		return fmt.Errorf("web auth capabilities failed: exact roles could not be resolved for key %q", keyID)
	}
	return withWebAuthHint(err, "web auth capabilities")
}

func renderWebAuthCapabilitiesTable(result webAuthCapabilitiesResult) error {
	asc.RenderTable(webAuthCapabilitiesHeaders(), webAuthCapabilitiesRows(result))
	return nil
}

func renderWebAuthCapabilitiesMarkdown(result webAuthCapabilitiesResult) error {
	asc.RenderMarkdown(webAuthCapabilitiesHeaders(), webAuthCapabilitiesRows(result))
	return nil
}

func webAuthCapabilitiesHeaders() []string {
	return []string{"KEY ID", "KIND", "ACTIVE", "ROLES", "CAPABILITIES", "NAME", "LOOKUP", "RESOLVED FROM", "PROFILE"}
}

func webAuthCapabilitiesRows(result webAuthCapabilitiesResult) [][]string {
	return [][]string{{
		result.KeyID,
		result.Kind,
		fmt.Sprintf("%t", result.Active),
		strings.Join(webAuthCapabilityRoleLabels(result), ", "),
		strings.Join(webAuthCapabilityLabels(result), ", "),
		result.Name,
		result.Lookup,
		result.ResolvedFrom,
		result.Profile,
	}}
}

func webAuthCapabilityRoleLabels(result webAuthCapabilitiesResult) []string {
	if len(result.RoleDetails) == 0 {
		return append([]string(nil), result.Roles...)
	}
	labels := make([]string, 0, len(result.RoleDetails)+len(result.UnknownRoles))
	for _, role := range result.RoleDetails {
		label := strings.TrimSpace(role.Label)
		if label == "" {
			label = role.Code
		}
		labels = append(labels, label)
	}
	return append(labels, result.UnknownRoles...)
}

func webAuthCapabilityLabels(result webAuthCapabilitiesResult) []string {
	if len(result.Capabilities) == 0 {
		return nil
	}
	labels := make([]string, 0, len(result.Capabilities))
	for _, capability := range result.Capabilities {
		label := strings.TrimSpace(capability.Label)
		if label == "" {
			label = capability.ID
		}
		labels = append(labels, label)
	}
	return labels
}

func convertWebAuthRoleDetails(src []webref.Role) []webAuthRoleDetailResult {
	if len(src) == 0 {
		return nil
	}
	dst := make([]webAuthRoleDetailResult, 0, len(src))
	for _, role := range src {
		dst = append(dst, webAuthRoleDetailResult{
			Code:                  role.Code,
			Label:                 role.Label,
			UIAliases:             append([]string(nil), role.UIAliases...),
			TeamKeySelectable:     role.TeamKeySelectable,
			IndividualKeyEligible: role.IndividualKeyEligible,
			Summary:               role.Summary,
			Capabilities:          append([]string(nil), role.Capabilities...),
			Notes:                 append([]string(nil), role.Notes...),
			ExampleTasks:          append([]string(nil), role.ExampleTasks...),
			NotableActions:        append([]string(nil), role.NotableActions...),
		})
	}
	return dst
}

func convertWebAuthCapabilities(src []webref.CapabilityGroup) []webAuthCapabilityResult {
	if len(src) == 0 {
		return nil
	}
	dst := make([]webAuthCapabilityResult, 0, len(src))
	for _, capability := range src {
		dst = append(dst, webAuthCapabilityResult{
			ID:      capability.ID,
			Label:   capability.Label,
			Summary: capability.Summary,
		})
	}
	return dst
}

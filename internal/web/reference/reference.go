package reference

import (
	_ "embed"
	"encoding/json"
	"strings"
	"sync"
)

//go:embed apple_role_capabilities.json
var raw []byte

var (
	load sync.Once
	snap *Snapshot
	err  error
)

type Snapshot struct {
	LastVerified string            `json:"last_verified"`
	Purpose      string            `json:"purpose"`
	Sources      []Source          `json:"sources"`
	APIKeyNotes  APIKeyNotes       `json:"api_key_notes"`
	Groups       []CapabilityGroup `json:"capability_groups"`
	Roles        []Role            `json:"roles"`
}

type Source struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type APIKeyNotes struct {
	Team       TeamKeyNotes       `json:"team_keys"`
	Individual IndividualKeyNotes `json:"individual_keys"`
}

type TeamKeyNotes struct {
	RequiredCreatorRoles []string `json:"required_creator_roles"`
	TeamScope            string   `json:"team_scope"`
	EditableAfterCreate  bool     `json:"editable_after_creation"`
	SelectableRoles      []string `json:"selectable_roles"`
}

type IndividualKeyNotes struct {
	EligibleUserRoles           []string `json:"eligible_user_roles"`
	OneActiveKeyPerUser         bool     `json:"one_active_key_per_user"`
	DefaultGenerationPermission string   `json:"default_generation_permission"`
	ManageableBy                []string `json:"manageable_by"`
}

type CapabilityGroup struct {
	ID      string `json:"id"`
	Label   string `json:"label"`
	Summary string `json:"summary,omitempty"`
}

type Role struct {
	Code                  string   `json:"code"`
	Label                 string   `json:"label"`
	UIAliases             []string `json:"ui_aliases,omitempty"`
	TeamKeySelectable     bool     `json:"team_key_selectable"`
	IndividualKeyEligible bool     `json:"individual_key_eligible"`
	Summary               string   `json:"summary,omitempty"`
	Capabilities          []string `json:"capabilities,omitempty"`
	Notes                 []string `json:"notes,omitempty"`
	ExampleTasks          []string `json:"example_tasks,omitempty"`
	NotableActions        []string `json:"notable_exclusive_actions,omitempty"`
}

type View struct {
	LastVerified string            `json:"lastVerified"`
	RoleDetails  []Role            `json:"roleDetails,omitempty"`
	Capabilities []CapabilityGroup `json:"capabilities,omitempty"`
	UnknownRoles []string          `json:"unknownRoles,omitempty"`
}

func Load() (*Snapshot, error) {
	load.Do(func() {
		var v Snapshot
		if e := json.Unmarshal(raw, &v); e != nil {
			err = e
			return
		}
		snap = &v
	})
	return snap, err
}

func Resolve(codes []string) (*View, error) {
	v, err := Load()
	if err != nil {
		return nil, err
	}

	roleByCode := make(map[string]Role, len(v.Roles))
	for _, role := range v.Roles {
		roleByCode[role.Code] = role
	}

	groupByID := make(map[string]CapabilityGroup, len(v.Groups))
	for _, group := range v.Groups {
		groupByID[group.ID] = group
	}

	seen := make(map[string]struct{}, len(codes))
	groups := make(map[string]struct{})
	view := &View{LastVerified: v.LastVerified}

	for _, rawCode := range codes {
		code := strings.TrimSpace(rawCode)
		if code == "" {
			continue
		}
		if _, ok := seen[code]; ok {
			continue
		}
		seen[code] = struct{}{}

		role, ok := roleByCode[code]
		if !ok {
			view.UnknownRoles = append(view.UnknownRoles, code)
			continue
		}
		view.RoleDetails = append(view.RoleDetails, role)
		for _, id := range role.Capabilities {
			groups[id] = struct{}{}
		}
	}

	for _, group := range v.Groups {
		if _, ok := groups[group.ID]; ok {
			view.Capabilities = append(view.Capabilities, group)
		}
	}

	for id := range groups {
		if _, ok := groupByID[id]; ok {
			continue
		}
		view.Capabilities = append(view.Capabilities, CapabilityGroup{
			ID:    id,
			Label: id,
		})
	}

	return view, nil
}

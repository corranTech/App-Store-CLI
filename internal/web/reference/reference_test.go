package reference

import "testing"

func TestLoad(t *testing.T) {
	v, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if v.LastVerified == "" {
		t.Fatal("expected last_verified")
	}
	if len(v.Roles) == 0 {
		t.Fatal("expected roles")
	}
	if len(v.Groups) == 0 {
		t.Fatal("expected capability groups")
	}
}

func TestResolve(t *testing.T) {
	view, err := Resolve([]string{"APP_MANAGER", "FINANCE"})
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if len(view.RoleDetails) != 2 {
		t.Fatalf("expected two role details, got %d", len(view.RoleDetails))
	}
	if len(view.Capabilities) == 0 {
		t.Fatal("expected capabilities")
	}
	if view.Capabilities[0].ID != "all_apps_access" && view.Capabilities[0].ID != "app_pricing_and_store_info" {
		t.Fatalf("unexpected first capability: %#v", view.Capabilities[0])
	}
}

func TestResolveUnknownRole(t *testing.T) {
	view, err := Resolve([]string{"NOPE", "APP_MANAGER"})
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if len(view.UnknownRoles) != 1 || view.UnknownRoles[0] != "NOPE" {
		t.Fatalf("unexpected unknown roles: %#v", view.UnknownRoles)
	}
	if len(view.RoleDetails) != 1 || view.RoleDetails[0].Code != "APP_MANAGER" {
		t.Fatalf("unexpected role details: %#v", view.RoleDetails)
	}
}

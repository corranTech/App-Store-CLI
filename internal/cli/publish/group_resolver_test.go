package publish

import (
	"strings"
	"testing"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
)

func TestResolvePublishBetaGroupIDsFromList_ResolvesByNameAndID(t *testing.T) {
	groups := &asc.BetaGroupsResponse{
		Data: []asc.Resource[asc.BetaGroupAttributes]{
			{ID: "GROUP_A", Attributes: asc.BetaGroupAttributes{Name: "External Testers"}},
			{ID: "GROUP_B", Attributes: asc.BetaGroupAttributes{Name: "Internal Team"}},
		},
	}

	got, err := resolvePublishBetaGroupIDsFromList(
		[]string{" external testers ", "GROUP_B", "GROUP_A", "EXTERNAL TESTERS"},
		groups,
	)
	if err != nil {
		t.Fatalf("resolvePublishBetaGroupIDsFromList() error = %v", err)
	}

	want := []string{"GROUP_A", "GROUP_B"}
	if len(got) != len(want) {
		t.Fatalf("expected %d groups, got %d (%v)", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected group %d to be %q, got %q", i, want[i], got[i])
		}
	}
}

func TestResolvePublishBetaGroupIDsFromList_MissingGroup(t *testing.T) {
	groups := &asc.BetaGroupsResponse{
		Data: []asc.Resource[asc.BetaGroupAttributes]{
			{ID: "GROUP_A", Attributes: asc.BetaGroupAttributes{Name: "External Testers"}},
		},
	}

	_, err := resolvePublishBetaGroupIDsFromList([]string{"does-not-exist"}, groups)
	if err == nil {
		t.Fatal("expected error for missing beta group")
	}
	if !strings.Contains(err.Error(), `beta group "does-not-exist" not found`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolvePublishBetaGroupIDsFromList_AmbiguousName(t *testing.T) {
	groups := &asc.BetaGroupsResponse{
		Data: []asc.Resource[asc.BetaGroupAttributes]{
			{ID: "GROUP_A", Attributes: asc.BetaGroupAttributes{Name: "QA"}},
			{ID: "GROUP_B", Attributes: asc.BetaGroupAttributes{Name: "QA"}},
		},
	}

	_, err := resolvePublishBetaGroupIDsFromList([]string{"qa"}, groups)
	if err == nil {
		t.Fatal("expected error for ambiguous beta group name")
	}
	if !strings.Contains(err.Error(), `multiple beta groups named "qa"; use group ID`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolvePublishBetaGroupIDsFromList_NoGroupsReturned(t *testing.T) {
	_, err := resolvePublishBetaGroupIDsFromList([]string{"group"}, nil)
	if err == nil {
		t.Fatal("expected error when no group list is available")
	}
	if !strings.Contains(err.Error(), "no beta groups returned for app") {
		t.Fatalf("unexpected error: %v", err)
	}
}

package shared

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/peterbourgon/ff/v3/ffcli"
)

func TestVisibleHelpSubcommandsFiltersHiddenCommands(t *testing.T) {
	visible := &ffcli.Command{Name: "visible"}
	hidden := HideCommandFromParentHelp(&ffcli.Command{Name: "hidden"})

	filtered := VisibleHelpSubcommands([]*ffcli.Command{visible, hidden, nil})
	if len(filtered) != 1 {
		t.Fatalf("expected 1 visible subcommand, got %d", len(filtered))
	}
	if filtered[0].Name != "visible" {
		t.Fatalf("expected visible subcommand to remain, got %q", filtered[0].Name)
	}
}

func TestRewriteCommandTreePathRewritesRuntimeErrorPrefix(t *testing.T) {
	cmd := &ffcli.Command{
		Name:       "values",
		ShortUsage: "asc offer-codes values [flags]",
		Exec: func(ctx context.Context, args []string) error {
			return fmt.Errorf("offer-codes values: %w", errors.New("boom"))
		},
	}

	rewritten := RewriteCommandTreePath(cmd, "asc offer-codes", "asc subscriptions offer-codes")
	if rewritten == nil {
		t.Fatal("expected rewritten command")
	}

	err := rewritten.Exec(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if got := err.Error(); got != "subscriptions offer-codes values: boom" {
		t.Fatalf("unexpected rewritten error: %q", got)
	}
}

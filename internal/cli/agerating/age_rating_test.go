package agerating

import (
	"context"
	"errors"
	"flag"
	"testing"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/asc"
)

func TestAgeRatingCommandShape(t *testing.T) {
	cmd := AgeRatingCommand()
	if cmd == nil {
		t.Fatal("expected age-rating command")
	}
	if cmd.Name != "age-rating" {
		t.Fatalf("unexpected command name: %q", cmd.Name)
	}
	if len(cmd.Subcommands) != 2 {
		t.Fatalf("expected 2 subcommands, got %d", len(cmd.Subcommands))
	}
	if got := Command(); got == nil {
		t.Fatal("expected Command wrapper to return command")
	}
}

func TestAgeRatingValidationErrors(t *testing.T) {
	t.Setenv("ASC_APP_ID", "")

	t.Run("get missing app or ids", func(t *testing.T) {
		cmd := AgeRatingGetCommand()
		if err := cmd.FlagSet.Parse([]string{}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		if err := cmd.Exec(context.Background(), nil); !errors.Is(err, flag.ErrHelp) {
			t.Fatalf("expected ErrHelp, got %v", err)
		}
	})

	t.Run("get conflicting app-info and version", func(t *testing.T) {
		cmd := AgeRatingGetCommand()
		if err := cmd.FlagSet.Parse([]string{"--app-info-id", "A", "--version-id", "V"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		err := cmd.Exec(context.Background(), nil)
		if err == nil || errors.Is(err, flag.ErrHelp) {
			t.Fatalf("expected non-ErrHelp error, got %v", err)
		}
	})

	t.Run("set missing id and app", func(t *testing.T) {
		cmd := AgeRatingSetCommand()
		if err := cmd.FlagSet.Parse([]string{}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		if err := cmd.Exec(context.Background(), nil); !errors.Is(err, flag.ErrHelp) {
			t.Fatalf("expected ErrHelp, got %v", err)
		}
	})
}

func TestAgeRatingHelpers(t *testing.T) {
	if _, err := buildAgeRatingAttributes(map[string]string{
		"seventeen-plus": "true",
	}); err == nil {
		t.Fatal("expected unsupported flag error")
	}

	if got, err := parseOptionalEnumFlag("--kids-age-band", "five_and_under", kidsAgeBandValues); err != nil || got == nil || *got != "FIVE_AND_UNDER" {
		t.Fatalf("expected normalized enum value, got %v err=%v", got, err)
	}
	if _, err := parseOptionalEnumFlag("--kids-age-band", "bad", kidsAgeBandValues); err == nil {
		t.Fatal("expected enum validation error")
	}

	if hasAgeRatingUpdates(asc.AgeRatingDeclarationAttributes{}) {
		t.Fatal("expected no updates for zero-value attrs")
	}
	value := "NONE"
	if !hasAgeRatingUpdates(asc.AgeRatingDeclarationAttributes{GamblingSimulated: &value}) {
		t.Fatal("expected updates when one pointer attribute is set")
	}
}

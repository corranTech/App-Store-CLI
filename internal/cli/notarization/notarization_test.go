package notarization

import (
	"context"
	"errors"
	"flag"
	"testing"
)

func TestNotarizationCommandConstructors(t *testing.T) {
	top := NotarizationCommand()
	if top == nil {
		t.Fatal("expected notarization command")
	}
	if top.Name == "" {
		t.Fatal("expected command name")
	}
	if len(top.Subcommands) == 0 {
		t.Fatal("expected subcommands")
	}

	constructors := []func() interface{}{
		func() interface{} { return submitCommand() },
		func() interface{} { return statusCommand() },
		func() interface{} { return logCommand() },
		func() interface{} { return listCommand() },
	}
	for _, ctor := range constructors {
		if got := ctor(); got == nil {
			t.Fatal("expected constructor to return command")
		}
	}
}

func TestNotarizationSubmitValidation(t *testing.T) {
	cmd := submitCommand()
	if err := cmd.FlagSet.Parse([]string{}); err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if err := cmd.Exec(context.Background(), nil); !errors.Is(err, flag.ErrHelp) {
		t.Fatalf("expected ErrHelp, got %v", err)
	}
}

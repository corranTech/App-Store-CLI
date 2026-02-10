package appclips

import "testing"

func TestAppClipsCommandConstructors(t *testing.T) {
	top := AppClipsCommand()
	if top == nil {
		t.Fatal("expected app-clips command")
	}
	if top.Name == "" {
		t.Fatal("expected top-level command name")
	}
	if len(top.Subcommands) == 0 {
		t.Fatal("expected app-clips subcommands")
	}

	if got := Command(); got == nil {
		t.Fatal("expected Command wrapper to return command")
	}

	constructors := []func() interface{}{
		func() interface{} { return AppClipDefaultExperiencesCommand() },
		func() interface{} { return AppClipAdvancedExperiencesCommand() },
		func() interface{} { return AppClipHeaderImagesCommand() },
		func() interface{} { return AppClipInvocationsCommand() },
		func() interface{} { return AppClipDomainStatusCommand() },
	}
	for _, ctor := range constructors {
		if got := ctor(); got == nil {
			t.Fatal("expected constructor to return command")
		}
	}
}

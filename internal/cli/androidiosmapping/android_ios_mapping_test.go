package androidiosmapping

import "testing"

func TestAndroidIosMappingCommandConstructors(t *testing.T) {
	cmd := AndroidIosMappingCommand()
	if cmd == nil {
		t.Fatal("expected android-ios-mapping command")
	}
	if cmd.Name == "" {
		t.Fatal("expected command name")
	}
	if len(cmd.Subcommands) != 5 {
		t.Fatalf("expected 5 subcommands, got %d", len(cmd.Subcommands))
	}

	constructors := []func() interface{}{
		func() interface{} { return AndroidIosMappingListCommand() },
		func() interface{} { return AndroidIosMappingGetCommand() },
		func() interface{} { return AndroidIosMappingCreateCommand() },
		func() interface{} { return AndroidIosMappingUpdateCommand() },
		func() interface{} { return AndroidIosMappingDeleteCommand() },
	}
	for _, ctor := range constructors {
		if got := ctor(); got == nil {
			t.Fatal("expected non-nil subcommand constructor")
		}
	}
}

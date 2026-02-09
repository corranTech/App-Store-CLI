package app_events

import "testing"

func TestAppEventsCommandConstructors(t *testing.T) {
	top := AppEventsCommand()
	if top == nil {
		t.Fatal("expected app-events command")
	}
	if top.Name == "" {
		t.Fatal("expected top-level command name")
	}
	if len(top.Subcommands) == 0 {
		t.Fatal("expected app-events subcommands")
	}

	if got := Command(); got == nil {
		t.Fatal("expected Command wrapper to return command")
	}

	constructors := []func() interface{}{
		func() interface{} { return AppEventLocalizationsCommand() },
		func() interface{} { return AppEventScreenshotsCommand() },
		func() interface{} { return AppEventVideoClipsCommand() },
		func() interface{} { return AppEventsRelationshipsCommand() },
		func() interface{} { return AppEventsSubmitCommand() },
		func() interface{} { return AppEventLocalizationScreenshotsCommand() },
	}
	for _, ctor := range constructors {
		if got := ctor(); got == nil {
			t.Fatal("expected constructor to return command")
		}
	}
}

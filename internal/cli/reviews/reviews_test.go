package reviews

import "testing"

func TestReviewsCommandConstructors(t *testing.T) {
	top := ReviewsCommand()
	if top == nil {
		t.Fatal("expected reviews command")
	}
	if top.Name == "" {
		t.Fatal("expected command name")
	}
	if len(top.Subcommands) == 0 {
		t.Fatal("expected reviews subcommands")
	}

	if got := Command(); got == nil {
		t.Fatal("expected Command wrapper to return command")
	}

	constructors := []func() interface{}{
		func() interface{} { return ReviewCommand() },
		func() interface{} { return ReviewsGetCommand() },
		func() interface{} { return ReviewsRatingsCommand() },
		func() interface{} { return ReviewsResponseCommand() },
		func() interface{} { return ReviewDetailsAttachmentsListCommand() },
	}
	for _, ctor := range constructors {
		if got := ctor(); got == nil {
			t.Fatal("expected constructor to return command")
		}
	}
}

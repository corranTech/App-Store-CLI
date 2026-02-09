package productpages

import "testing"

func TestProductPagesCommandConstructors(t *testing.T) {
	top := ProductPagesCommand()
	if top == nil {
		t.Fatal("expected product-pages command")
	}
	if top.Name == "" {
		t.Fatal("expected command name")
	}
	if len(top.Subcommands) == 0 {
		t.Fatal("expected subcommands")
	}

	if got := Command(); got == nil {
		t.Fatal("expected Command wrapper to return command")
	}

	constructors := []func() interface{}{
		func() interface{} { return CustomPagesCommand() },
		func() interface{} { return CustomPageLocalizationsCommand() },
		func() interface{} { return CustomPageVersionsCommand() },
		func() interface{} { return ExperimentsCommand() },
		func() interface{} { return ExperimentTreatmentsCommand() },
	}
	for _, ctor := range constructors {
		if got := ctor(); got == nil {
			t.Fatal("expected constructor to return command")
		}
	}
}

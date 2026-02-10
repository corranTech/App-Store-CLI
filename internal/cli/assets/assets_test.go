package assets

import "testing"

func TestAssetsCommandConstructors(t *testing.T) {
	top := AssetsCommand()
	if top == nil {
		t.Fatal("expected assets command")
	}
	if top.Name == "" {
		t.Fatal("expected command name")
	}
	if len(top.Subcommands) == 0 {
		t.Fatal("expected assets subcommands")
	}

	if got := Command(); got == nil {
		t.Fatal("expected Command wrapper to return command")
	}

	if got := AssetsScreenshotsCommand(); got == nil {
		t.Fatal("expected screenshots command")
	}
	if got := AssetsPreviewsCommand(); got == nil {
		t.Fatal("expected previews command")
	}
}

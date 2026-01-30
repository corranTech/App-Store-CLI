package cmdtest

import (
	"context"
	"errors"
	"flag"
	"strings"
	"testing"
)

func TestProductPagesCustomPagesListRequiresApp(t *testing.T) {
	t.Setenv("ASC_APP_ID", "")
	root := RootCommand("1.2.3")

	stdout, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{"product-pages", "custom-pages", "list"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		err := root.Run(context.Background())
		if !errors.Is(err, flag.ErrHelp) {
			t.Fatalf("expected ErrHelp, got %v", err)
		}
	})

	if stdout != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
	if !strings.Contains(stderr, "--app is required") {
		t.Fatalf("expected missing app error, got %q", stderr)
	}
}

func TestProductPagesCustomPagesDeleteRequiresConfirm(t *testing.T) {
	root := RootCommand("1.2.3")

	stdout, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{"product-pages", "custom-pages", "delete", "--custom-page-id", "page-1"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		err := root.Run(context.Background())
		if !errors.Is(err, flag.ErrHelp) {
			t.Fatalf("expected ErrHelp, got %v", err)
		}
	})

	if stdout != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
	if !strings.Contains(stderr, "--confirm is required") {
		t.Fatalf("expected missing confirm error, got %q", stderr)
	}
}

func TestProductPagesExperimentsCreateRequiresVersionID(t *testing.T) {
	root := RootCommand("1.2.3")

	stdout, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{"product-pages", "experiments", "create", "--name", "Icon Test", "--traffic-proportion", "50"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		err := root.Run(context.Background())
		if !errors.Is(err, flag.ErrHelp) {
			t.Fatalf("expected ErrHelp, got %v", err)
		}
	})

	if stdout != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
	if !strings.Contains(stderr, "--version-id is required") {
		t.Fatalf("expected missing version-id error, got %q", stderr)
	}
}

func TestProductPagesExperimentsCreateRequiresAppForV2(t *testing.T) {
	t.Setenv("ASC_APP_ID", "")
	root := RootCommand("1.2.3")

	stdout, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{"product-pages", "experiments", "create", "--v2", "--platform", "IOS", "--name", "Icon Test", "--traffic-proportion", "50"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		err := root.Run(context.Background())
		if !errors.Is(err, flag.ErrHelp) {
			t.Fatalf("expected ErrHelp, got %v", err)
		}
	})

	if stdout != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
	if !strings.Contains(stderr, "--app is required") {
		t.Fatalf("expected missing app error, got %q", stderr)
	}
}

func TestProductPagesExperimentsDeleteRequiresConfirm(t *testing.T) {
	root := RootCommand("1.2.3")

	stdout, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{"product-pages", "experiments", "delete", "--experiment-id", "exp-1"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		err := root.Run(context.Background())
		if !errors.Is(err, flag.ErrHelp) {
			t.Fatalf("expected ErrHelp, got %v", err)
		}
	})

	if stdout != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
	if !strings.Contains(stderr, "--confirm is required") {
		t.Fatalf("expected missing confirm error, got %q", stderr)
	}
}

func TestProductPagesCustomPagesListRejectsInvalidLimit(t *testing.T) {
	root := RootCommand("1.2.3")

	tests := []struct {
		name  string
		limit string
	}{
		{
			name:  "limit below range",
			limit: "-1",
		},
		{
			name:  "limit above range",
			limit: "201",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stdout, stderr := captureOutput(t, func() {
				if err := root.Parse([]string{"product-pages", "custom-pages", "list", "--app", "APP_ID", "--limit", test.limit}); err != nil {
					t.Fatalf("parse error: %v", err)
				}
				err := root.Run(context.Background())
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if errors.Is(err, flag.ErrHelp) {
					t.Fatalf("unexpected ErrHelp, got %v", err)
				}
			})

			if stdout != "" {
				t.Fatalf("expected empty stdout, got %q", stdout)
			}
			if stderr != "" {
				t.Fatalf("expected empty stderr, got %q", stderr)
			}
		})
	}
}

func TestProductPagesCustomPagesListRejectsInvalidNextURL(t *testing.T) {
	root := RootCommand("1.2.3")

	stdout, stderr := captureOutput(t, func() {
		if err := root.Parse([]string{"product-pages", "custom-pages", "list", "--app", "APP_ID", "--next", "not-a-url"}); err != nil {
			t.Fatalf("parse error: %v", err)
		}
		err := root.Run(context.Background())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if errors.Is(err, flag.ErrHelp) {
			t.Fatalf("unexpected ErrHelp, got %v", err)
		}
	})

	if stdout != "" {
		t.Fatalf("expected empty stdout, got %q", stdout)
	}
	if stderr != "" {
		t.Fatalf("expected empty stderr, got %q", stderr)
	}
}

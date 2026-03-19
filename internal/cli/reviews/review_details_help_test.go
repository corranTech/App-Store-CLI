package reviews

import (
	"strings"
	"testing"
)

func TestReviewDetailsCreateCommandClarifiesReviewerAccessGuidance(t *testing.T) {
	cmd := ReviewDetailsCreateCommand()

	if !strings.Contains(cmd.LongHelp, "Leave `--demo-account-required` false when `--notes` are enough") {
		t.Fatalf("expected create help to explain notes-only guidance, got %q", cmd.LongHelp)
	}
	if !strings.Contains(cmd.LongHelp, "Use `--demo-account-required=true` only when App Review needs demo credentials") {
		t.Fatalf("expected create help to explain demo credential opt-in, got %q", cmd.LongHelp)
	}

	if got := cmd.FlagSet.Lookup("demo-account-required").Usage; !strings.Contains(got, "Set true only when App Review needs demo credentials") {
		t.Fatalf("expected --demo-account-required usage to clarify semantics, got %q", got)
	}
	if got := cmd.FlagSet.Lookup("notes").Usage; !strings.Contains(got, "reviewer instructions") {
		t.Fatalf("expected --notes usage to mention reviewer instructions, got %q", got)
	}
}

func TestReviewDetailsUpdateCommandClarifiesReviewerAccessGuidance(t *testing.T) {
	cmd := ReviewDetailsUpdateCommand()

	if !strings.Contains(cmd.LongHelp, "Leave `--demo-account-required` false when `--notes` are enough") {
		t.Fatalf("expected update help to explain notes-only guidance, got %q", cmd.LongHelp)
	}
	if !strings.Contains(cmd.LongHelp, "Do not use placeholder demo credentials") {
		t.Fatalf("expected update help to discourage placeholder credentials, got %q", cmd.LongHelp)
	}

	if got := cmd.FlagSet.Lookup("demo-account-name").Usage; !strings.Contains(got, "when demo credentials are required") {
		t.Fatalf("expected --demo-account-name usage to clarify when it is needed, got %q", got)
	}
	if got := cmd.FlagSet.Lookup("demo-account-password").Usage; !strings.Contains(got, "when demo credentials are required") {
		t.Fatalf("expected --demo-account-password usage to clarify when it is needed, got %q", got)
	}
}

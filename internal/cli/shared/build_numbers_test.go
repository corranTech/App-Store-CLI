package shared

import (
	"strings"
	"testing"
)

func TestParseBuildNumberRejectsNonNumeric(t *testing.T) {
	_, err := parseBuildNumber("1a", "processed build")
	if err == nil {
		t.Fatal("expected error for non-numeric build number")
	}
	if !strings.Contains(err.Error(), "processed build") {
		t.Fatalf("expected error to mention source, got %v", err)
	}
}

func TestParseBuildNumberRejectsEmpty(t *testing.T) {
	_, err := parseBuildNumber(" ", "build upload")
	if err == nil {
		t.Fatal("expected error for empty build number")
	}
}

func TestParseBuildNumberAllowsNumeric(t *testing.T) {
	got, err := parseBuildNumber("42", "processed build")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.String() != "42" {
		t.Fatalf("expected 42, got %q", got.String())
	}
}

func TestParseBuildNumberAllowsDotSeparatedNumeric(t *testing.T) {
	got, err := parseBuildNumber("1.2.3", "build upload")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.String() != "1.2.3" {
		t.Fatalf("expected 1.2.3, got %q", got.String())
	}
}

func TestBuildNumberNextIncrementsLastSegment(t *testing.T) {
	parsed, err := parseBuildNumber("1.2.3", "processed build")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	next, err := parsed.Next()
	if err != nil {
		t.Fatalf("unexpected error incrementing build number: %v", err)
	}
	if next.String() != "1.2.4" {
		t.Fatalf("expected next build number 1.2.4, got %q", next.String())
	}
}

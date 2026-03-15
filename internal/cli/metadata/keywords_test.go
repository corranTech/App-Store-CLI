package metadata

import "testing"

func TestSplitMetadataKeywordTokensSupportsMixedSeparators(t *testing.T) {
	got := splitMetadataKeywordTokens(" habit tracker，mood journal、sleep log;\nenergy tracker； focus timer ")
	want := []string{
		"habit tracker",
		"mood journal",
		"sleep log",
		"energy tracker",
		"focus timer",
	}
	if len(got) != len(want) {
		t.Fatalf("expected %d tokens, got %d: %v", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected token %d to be %q, got %q (%v)", i, want[i], got[i], got)
		}
	}
}

func TestNormalizeMetadataKeywordListDeduplicatesCaseAndWhitespace(t *testing.T) {
	got, err := normalizeMetadataKeywordList([]string{
		"  habit   tracker ",
		"mood journal",
		"Habit Tracker",
		"mood   journal",
		"  ",
	})
	if err != nil {
		t.Fatalf("normalizeMetadataKeywordList() error: %v", err)
	}
	want := []string{"habit tracker", "mood journal"}
	if len(got) != len(want) {
		t.Fatalf("expected %d normalized keywords, got %d: %v", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected keyword %d to be %q, got %q (%v)", i, want[i], got[i], got)
		}
	}
}

func TestDecodeMetadataKeywordValueArrayExpandsEmbeddedSeparators(t *testing.T) {
	got, err := decodeMetadataKeywordValue([]any{
		"habit tracker，mood journal",
		"sleep log",
		" Focus tracker ",
	})
	if err != nil {
		t.Fatalf("decodeMetadataKeywordValue() error: %v", err)
	}
	want := []string{"habit tracker", "mood journal", "sleep log", "Focus tracker"}
	if len(got) != len(want) {
		t.Fatalf("expected %d tokens, got %d: %v", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected token %d to be %q, got %q (%v)", i, want[i], got[i], got)
		}
	}
}

func TestBuildMetadataKeywordFieldNormalizesMixedInput(t *testing.T) {
	field, count, err := buildMetadataKeywordField([]string{
		"habit tracker",
		" mood journal ",
		"Habit Tracker",
		"sleep log",
	})
	if err != nil {
		t.Fatalf("buildMetadataKeywordField() error: %v", err)
	}
	if field != "habit tracker,mood journal,sleep log" {
		t.Fatalf("expected canonical keyword field, got %q", field)
	}
	if count != 3 {
		t.Fatalf("expected count 3, got %d", count)
	}
}

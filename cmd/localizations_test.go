package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseStringsContent(t *testing.T) {
	input := `
// Comment
"description" = "Hello\nWorld";
/* block comment */
"keywords" = "one, two";
`
	values, err := parseStringsContent(input)
	if err != nil {
		t.Fatalf("parseStringsContent() error: %v", err)
	}
	if values["description"] != "Hello\nWorld" {
		t.Fatalf("expected description to be parsed, got %q", values["description"])
	}
	if values["keywords"] != "one, two" {
		t.Fatalf("expected keywords to be parsed, got %q", values["keywords"])
	}
}

func TestWriteStringsFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "en-US.strings")
	values := map[string]string{
		"description": "Hello",
		"keywords":    "one, two",
	}

	if err := writeStringsFile(path, values, []string{"description", "keywords"}); err != nil {
		t.Fatalf("writeStringsFile() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file error: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "\"description\" = \"Hello\";") {
		t.Fatalf("expected description line, got: %s", content)
	}
	if !strings.Contains(content, "\"keywords\" = \"one, two\";") {
		t.Fatalf("expected keywords line, got: %s", content)
	}
}

func TestReadLocalizationStrings_FileLocale(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "en-US.strings")
	if err := os.WriteFile(path, []byte("\"description\" = \"Hello\";\n"), 0o644); err != nil {
		t.Fatalf("write file error: %v", err)
	}

	values, err := readLocalizationStrings(path, nil)
	if err != nil {
		t.Fatalf("readLocalizationStrings() error: %v", err)
	}
	if values["en-US"]["description"] != "Hello" {
		t.Fatalf("expected description Hello, got %q", values["en-US"]["description"])
	}
}

func TestReadLocalizationStrings_RejectsSymlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target.strings")
	if err := os.WriteFile(target, []byte("\"description\" = \"Hello\";\n"), 0o644); err != nil {
		t.Fatalf("write file error: %v", err)
	}
	link := filepath.Join(dir, "en-US.strings")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlink not supported: %v", err)
	}

	_, err := readLocalizationStrings(dir, nil)
	if err == nil {
		t.Fatal("expected error for symlinked strings file")
	}
	if !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("expected symlink error, got %v", err)
	}
}

package scripts_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func repoRootForScriptTests(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve current file path")
	}
	return filepath.Dir(filepath.Dir(thisFile))
}

func copyGeneratorScript(t *testing.T, dstRepoRoot string) {
	t.Helper()
	src := filepath.Join(repoRootForScriptTests(t), "scripts", "update-wall-of-apps.py")
	content, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read generator script: %v", err)
	}

	dst := filepath.Join(dstRepoRoot, "scripts", "update-wall-of-apps.py")
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		t.Fatalf("mkdir scripts dir: %v", err)
	}
	if err := os.WriteFile(dst, content, 0o755); err != nil {
		t.Fatalf("write generator script: %v", err)
	}
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir parent dir for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file %s: %v", path, err)
	}
}

func runGenerator(repoRoot string) (string, error) {
	cmd := exec.Command("python3", filepath.Join(repoRoot, "scripts", "update-wall-of-apps.py"))
	cmd.Dir = repoRoot
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func TestUpdateWallOfAppsScriptGeneratesDocsAndReadme(t *testing.T) {
	tmpRepo := t.TempDir()
	copyGeneratorScript(t, tmpRepo)

	writeFile(t, filepath.Join(tmpRepo, "docs", "wall-of-apps.json"), `[
  {
    "app": "Zulu App",
    "link": "https://example.com/zulu",
    "creator": "Zulu Creator",
    "platform": ["iOS"]
  },
  {
    "app": "Alpha App",
    "link": "https://example.com/alpha",
    "creator": "Alpha Creator",
    "platform": ["macOS", "iOS"]
  }
]`)

	writeFile(t, filepath.Join(tmpRepo, "README.md"), `# Demo

Before.
<!-- WALL-OF-APPS:START -->
Old content.
<!-- WALL-OF-APPS:END -->
After.
`)

	output, err := runGenerator(tmpRepo)
	if err != nil {
		t.Fatalf("generator failed: %v\noutput:\n%s", err, output)
	}

	generatedPath := filepath.Join(tmpRepo, "docs", "generated", "app-wall.md")
	generatedContentBytes, err := os.ReadFile(generatedPath)
	if err != nil {
		t.Fatalf("read generated app wall: %v", err)
	}
	generatedContent := string(generatedContentBytes)

	if !strings.Contains(generatedContent, "Generated from docs/wall-of-apps.json") {
		t.Fatalf("expected generated header, got:\n%s", generatedContent)
	}
	if !strings.Contains(generatedContent, "| App | Link | Creator | Platform |") {
		t.Fatalf("expected markdown table header, got:\n%s", generatedContent)
	}
	alphaRow := "| Alpha App | [Open](https://example.com/alpha) | Alpha Creator | macOS, iOS |"
	zuluRow := "| Zulu App | [Open](https://example.com/zulu) | Zulu Creator | iOS |"
	alphaIdx := strings.Index(generatedContent, alphaRow)
	zuluIdx := strings.Index(generatedContent, zuluRow)
	if alphaIdx == -1 || zuluIdx == -1 {
		t.Fatalf("expected both generated rows, got:\n%s", generatedContent)
	}
	if alphaIdx > zuluIdx {
		t.Fatalf("expected deterministic app sorting, got:\n%s", generatedContent)
	}

	readmeBytes, err := os.ReadFile(filepath.Join(tmpRepo, "README.md"))
	if err != nil {
		t.Fatalf("read generated README: %v", err)
	}
	readme := string(readmeBytes)
	if strings.Contains(readme, "Old content.") {
		t.Fatalf("expected README snippet to be replaced, got:\n%s", readme)
	}
	if !strings.Contains(readme, alphaRow) || !strings.Contains(readme, zuluRow) {
		t.Fatalf("expected README to include generated rows, got:\n%s", readme)
	}
}

func TestUpdateWallOfAppsScriptValidatesRequiredCreatorField(t *testing.T) {
	tmpRepo := t.TempDir()
	copyGeneratorScript(t, tmpRepo)

	writeFile(t, filepath.Join(tmpRepo, "docs", "wall-of-apps.json"), `[
  {
    "app": "No Creator App",
    "link": "https://example.com/no-creator",
    "platform": ["iOS"]
  }
]`)

	writeFile(t, filepath.Join(tmpRepo, "README.md"), `# Demo
<!-- WALL-OF-APPS:START -->
Old content.
<!-- WALL-OF-APPS:END -->
`)

	output, err := runGenerator(tmpRepo)
	if err == nil {
		t.Fatalf("expected generator to fail for missing creator, output:\n%s", output)
	}
	if !strings.Contains(output, "'creator' is required") {
		t.Fatalf("expected missing creator error, got:\n%s", output)
	}
}

package update

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestRestart_ValidatesInputs(t *testing.T) {
	t.Run("missing executable", func(t *testing.T) {
		code, err := Restart("", []string{"asc"}, nil)
		if err == nil {
			t.Fatal("expected error for missing executable")
		}
		if code != 1 {
			t.Fatalf("exit code = %d, want 1", code)
		}
	})

	t.Run("missing args", func(t *testing.T) {
		code, err := Restart("/tmp/asc", nil, nil)
		if err == nil {
			t.Fatal("expected error for missing args")
		}
		if code != 1 {
			t.Fatalf("exit code = %d, want 1", code)
		}
	})
}

func TestRestart_ReturnsExitCodeFromChildProcess(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell helper script test is unix-only")
	}

	scriptPath := writeRestartTestScript(t, `#!/bin/sh
exit 7
`)

	code, err := Restart(scriptPath, []string{"asc"}, nil)
	if err != nil {
		t.Fatalf("Restart() error = %v", err)
	}
	if code != 7 {
		t.Fatalf("Restart() exit code = %d, want 7", code)
	}
}

func TestRestart_AppendsSkipUpdateEnvVar(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell helper script test is unix-only")
	}

	scriptPath := writeRestartTestScript(t, `#!/bin/sh
if [ "$ASC_SKIP_UPDATE" != "1" ]; then
  exit 23
fi
if [ "$TEST_CUSTOM_ENV" != "present" ]; then
  exit 24
fi
exit 0
`)

	code, err := Restart(scriptPath, []string{"asc", "apps", "list"}, []string{"TEST_CUSTOM_ENV=present"})
	if err != nil {
		t.Fatalf("Restart() error = %v", err)
	}
	if code != 0 {
		t.Fatalf("Restart() exit code = %d, want 0", code)
	}
}

func TestRestart_ReturnsErrorOnCommandStartFailure(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing-asc-executable")

	code, err := Restart(missing, []string{"asc"}, nil)
	if err == nil {
		t.Fatal("expected command start error")
	}
	if code != 1 {
		t.Fatalf("Restart() exit code = %d, want 1", code)
	}
}

func writeRestartTestScript(t *testing.T, body string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "restart-helper.sh")
	if err := os.WriteFile(path, []byte(body), 0o755); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}
	return path
}

func TestRestart_ExecutesWithForwardedArgs(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell helper script test is unix-only")
	}

	scriptPath := writeRestartTestScript(t, `#!/bin/sh
if [ "$1" != "apps" ] || [ "$2" != "list" ]; then
  exit 31
fi
exit 0
`)

	code, err := Restart(scriptPath, []string{"asc", "apps", "list"}, nil)
	if err != nil {
		t.Fatalf("Restart() error = %v", err)
	}
	if code != 0 {
		t.Fatalf("Restart() exit code = %d, want 0", code)
	}
}

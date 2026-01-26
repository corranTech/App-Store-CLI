package cmd

import (
	"encoding/base64"
	"os"
	"strings"
	"testing"
)

func TestResolvePrivateKeyPathPrefersPath(t *testing.T) {
	resetPrivateKeyTemp(t)
	t.Setenv("ASC_PRIVATE_KEY_PATH", "/tmp/AuthKey.p8")
	t.Setenv("ASC_PRIVATE_KEY_B64", base64.StdEncoding.EncodeToString([]byte("ignored")))
	t.Setenv("ASC_PRIVATE_KEY", "ignored")

	path, err := resolvePrivateKeyPath()
	if err != nil {
		t.Fatalf("resolvePrivateKeyPath() error: %v", err)
	}
	if path != "/tmp/AuthKey.p8" {
		t.Fatalf("expected path /tmp/AuthKey.p8, got %q", path)
	}
}

func TestResolvePrivateKeyPathFromBase64(t *testing.T) {
	resetPrivateKeyTemp(t)
	t.Setenv("ASC_PRIVATE_KEY_PATH", "")
	t.Setenv("ASC_PRIVATE_KEY", "")

	encoded := base64.StdEncoding.EncodeToString([]byte("key-data"))
	t.Setenv("ASC_PRIVATE_KEY_B64", encoded)

	path, err := resolvePrivateKeyPath()
	if err != nil {
		t.Fatalf("resolvePrivateKeyPath() error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	if string(data) != "key-data" {
		t.Fatalf("expected key data %q, got %q", "key-data", string(data))
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error: %v", err)
	}
	if info.Mode().Perm()&0o077 != 0 {
		t.Fatalf("expected 0600 permissions, got %v", info.Mode().Perm())
	}
}

func TestResolvePrivateKeyPathFromRawValue(t *testing.T) {
	resetPrivateKeyTemp(t)
	t.Setenv("ASC_PRIVATE_KEY_PATH", "")
	t.Setenv("ASC_PRIVATE_KEY_B64", "")

	t.Setenv("ASC_PRIVATE_KEY", "line1\\nline2")
	path, err := resolvePrivateKeyPath()
	if err != nil {
		t.Fatalf("resolvePrivateKeyPath() error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	if string(data) != "line1\nline2" {
		t.Fatalf("expected newline expansion, got %q", string(data))
	}
}

func TestResolvePrivateKeyPathInvalidBase64(t *testing.T) {
	resetPrivateKeyTemp(t)
	t.Setenv("ASC_PRIVATE_KEY_PATH", "")
	t.Setenv("ASC_PRIVATE_KEY", "")
	t.Setenv("ASC_PRIVATE_KEY_B64", "not-base64")

	if _, err := resolvePrivateKeyPath(); err == nil {
		t.Fatal("expected error for invalid base64")
	}
}

func resetPrivateKeyTemp(t *testing.T) {
	t.Helper()
	if privateKeyTempPath != "" {
		_ = os.Remove(privateKeyTempPath)
		privateKeyTempPath = ""
	}
	t.Cleanup(func() {
		if privateKeyTempPath != "" {
			_ = os.Remove(privateKeyTempPath)
			privateKeyTempPath = ""
		}
	})
	t.Setenv("ASC_PRIVATE_KEY_PATH", "")
	t.Setenv("ASC_PRIVATE_KEY_B64", "")
	t.Setenv("ASC_PRIVATE_KEY", "")
	t.Setenv("ASC_BYPASS_KEYCHAIN", "1")
	t.Setenv("ASC_CONFIG_PATH", os.TempDir()+string(os.PathSeparator)+strings.ReplaceAll(t.Name(), " ", "_"))
}

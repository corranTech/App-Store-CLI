package config

import (
	"path/filepath"
	"testing"
)

func TestConfigSaveLoadRemove(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	cfg := &Config{
		KeyID:          "KEY123",
		IssuerID:       "ISSUER456",
		PrivateKeyPath: "/tmp/AuthKey.p8",
		DefaultKeyName: "default",
	}

	if err := Save(cfg); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if loaded.KeyID != cfg.KeyID {
		t.Fatalf("KeyID mismatch: got %q want %q", loaded.KeyID, cfg.KeyID)
	}
	if loaded.IssuerID != cfg.IssuerID {
		t.Fatalf("IssuerID mismatch: got %q want %q", loaded.IssuerID, cfg.IssuerID)
	}
	if loaded.PrivateKeyPath != cfg.PrivateKeyPath {
		t.Fatalf("PrivateKeyPath mismatch: got %q want %q", loaded.PrivateKeyPath, cfg.PrivateKeyPath)
	}
	if loaded.DefaultKeyName != cfg.DefaultKeyName {
		t.Fatalf("DefaultKeyName mismatch: got %q want %q", loaded.DefaultKeyName, cfg.DefaultKeyName)
	}

	if err := Remove(); err != nil {
		t.Fatalf("Remove() error: %v", err)
	}

	if _, err := Load(); err != ErrNotFound {
		t.Fatalf("expected ErrNotFound after Remove(), got %v", err)
	}
}

func TestLoadMissingConfig(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	if _, err := Load(); err != ErrNotFound {
		t.Fatalf("expected ErrNotFound for missing config, got %v", err)
	}
}

func TestConfigPath(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("HOME", tempDir)

	path, err := Path()
	if err != nil {
		t.Fatalf("Path() error: %v", err)
	}

	expected := filepath.Join(tempDir, ".asc", "config.json")
	if path != expected {
		t.Fatalf("Path() mismatch: got %q want %q", path, expected)
	}
}

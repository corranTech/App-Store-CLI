package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	configDirName  = ".asc"
	configFileName = "config.json"
)

// Config holds the application configuration
type Config struct {
	KeyID          string `json:"key_id"`
	IssuerID       string `json:"issuer_id"`
	PrivateKeyPath string `json:"private_key_path"`
	DefaultKeyName string `json:"default_key_name"`
}

// ErrNotFound is returned when the config file doesn't exist
var ErrNotFound = fmt.Errorf("configuration not found")

// configDir returns the path to the configuration directory
func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, configDirName), nil
}

// configPath returns the path to the config file
func configPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, configFileName), nil
}

// Path returns the configuration file path.
func Path() (string, error) {
	return configPath()
}

// Load loads the configuration from the config file
func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

// Save saves the configuration to the config file
func Save(cfg *Config) error {
	dir, err := configDir()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	path, err := configPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// Remove removes the config file
func Remove() error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to remove config: %w", err)
	}

	return nil
}

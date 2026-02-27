package shared

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const tierCacheTTL = 24 * time.Hour

type tierCacheFile struct {
	AsOf    time.Time   `json:"asOf"`
	AppID   string      `json:"appId"`
	Territory string    `json:"territory"`
	Tiers   []TierEntry `json:"tiers"`
}

func tierCacheDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	dir := filepath.Join(home, ".asc", "cache")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create cache dir: %w", err)
	}
	return dir, nil
}

func tierCachePath(appID, territory string) (string, error) {
	dir, err := tierCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, fmt.Sprintf("tiers-%s-%s.json", appID, territory)), nil
}

// LoadTierCache loads cached tier data. Returns an error if the cache is missing or expired.
func LoadTierCache(appID, territory string) ([]TierEntry, error) {
	path, err := tierCachePath(appID, territory)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cache tierCacheFile
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, fmt.Errorf("parse cache: %w", err)
	}

	if time.Since(cache.AsOf) > tierCacheTTL {
		return nil, fmt.Errorf("cache expired")
	}

	return cache.Tiers, nil
}

// SaveTierCache writes tier data to the cache file.
func SaveTierCache(appID, territory string, tiers []TierEntry) error {
	path, err := tierCachePath(appID, territory)
	if err != nil {
		return err
	}

	cache := tierCacheFile{
		AsOf:      time.Now(),
		AppID:     appID,
		Territory: territory,
		Tiers:     tiers,
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal cache: %w", err)
	}

	return os.WriteFile(path, data, 0o644)
}

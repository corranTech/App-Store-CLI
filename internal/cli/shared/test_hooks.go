package shared

import (
	"os"
	"path/filepath"
)

// ResetTierCacheForTest clears the local tier cache for tests.
func ResetTierCacheForTest() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	_ = os.RemoveAll(filepath.Join(home, ".asc", "cache"))
}

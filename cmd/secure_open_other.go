//go:build !darwin && !linux && !freebsd && !netbsd && !openbsd && !dragonfly

package cmd

import "os"

// openNewFileNoFollow creates a new file without following symlinks.
// Uses O_EXCL to prevent overwriting existing files.
// Note: O_NOFOLLOW is not available on this platform.
func openNewFileNoFollow(path string, perm os.FileMode) (*os.File, error) {
	flags := os.O_WRONLY | os.O_CREATE | os.O_EXCL
	return os.OpenFile(path, flags, perm)
}

// openExistingNoFollow opens an existing file.
func openExistingNoFollow(path string) (*os.File, error) {
	return os.Open(path)
}

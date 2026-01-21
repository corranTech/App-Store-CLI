//go:build darwin || linux || freebsd || netbsd || openbsd || dragonfly

package cmd

import (
	"os"

	"golang.org/x/sys/unix"
)

// openNewFileNoFollow creates a new file without following symlinks.
// Uses O_EXCL to prevent overwriting existing files and O_NOFOLLOW to prevent symlink attacks.
func openNewFileNoFollow(path string, perm os.FileMode) (*os.File, error) {
	flags := os.O_WRONLY | os.O_CREATE | os.O_EXCL | unix.O_NOFOLLOW
	return os.OpenFile(path, flags, perm)
}

// openExistingNoFollow opens an existing file without following symlinks.
func openExistingNoFollow(path string) (*os.File, error) {
	flags := os.O_RDONLY | unix.O_NOFOLLOW
	return os.OpenFile(path, flags, 0)
}

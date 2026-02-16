package profiles

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rudrankriyam/App-Store-Connect-CLI/internal/cli/shared"
)

func writeProfileFile(path string, content []byte, force bool) error {
	if !force {
		return shared.WriteProfileFile(path, content)
	}
	return writeFileBytesNoSymlink(path, bytes.NewReader(content), 0o644)
}

func writeFileBytesNoSymlink(path string, reader io.Reader, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	// Do not remove/replace a symlink.
	hadExisting := false
	if info, err := os.Lstat(path); err == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("refusing to overwrite symlink %q", path)
		}
		if info.IsDir() {
			return fmt.Errorf("output path %q is a directory", path)
		}
		hadExisting = true
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	tempFile, err := os.CreateTemp(filepath.Dir(path), ".asc-profile-*")
	if err != nil {
		return err
	}
	defer tempFile.Close()

	tempPath := tempFile.Name()
	success := false
	defer func() {
		if !success {
			_ = os.Remove(tempPath)
		}
	}()

	if err := tempFile.Chmod(perm); err != nil {
		return err
	}
	if _, err := io.Copy(tempFile, reader); err != nil {
		return err
	}
	if err := tempFile.Sync(); err != nil {
		return err
	}
	if err := tempFile.Close(); err != nil {
		return err
	}

	if err := os.Rename(tempPath, path); err != nil {
		if !hadExisting {
			return err
		}

		backupFile, backupErr := os.CreateTemp(filepath.Dir(path), ".asc-profile-backup-*")
		if backupErr != nil {
			return err
		}
		backupPath := backupFile.Name()
		if closeErr := backupFile.Close(); closeErr != nil {
			return closeErr
		}
		if removeErr := os.Remove(backupPath); removeErr != nil {
			return removeErr
		}

		if moveErr := os.Rename(path, backupPath); moveErr != nil {
			return moveErr
		}
		if moveErr := os.Rename(tempPath, path); moveErr != nil {
			_ = os.Rename(backupPath, path)
			return moveErr
		}
		_ = os.Remove(backupPath)
	}

	success = true
	return nil
}

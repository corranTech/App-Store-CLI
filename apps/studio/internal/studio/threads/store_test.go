package threads

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSaveThreadRoundTrip(t *testing.T) {
	store := NewStore(t.TempDir())
	now := time.Now().UTC()
	thread := Thread{
		ID:        "thread-1",
		Title:     "Release Prep",
		CreatedAt: now,
		UpdatedAt: now,
		Messages: []Message{
			{ID: "msg-1", Role: RoleUser, Kind: KindMessage, Content: "Validate 2.3.0", CreatedAt: now},
		},
	}

	if err := store.SaveThread(thread); err != nil {
		t.Fatalf("SaveThread() error = %v", err)
	}

	got, err := store.Get("thread-1")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.Title != thread.Title {
		t.Fatalf("Title = %q, want %q", got.Title, thread.Title)
	}
	if len(got.Messages) != 1 {
		t.Fatalf("len(Messages) = %d, want 1", len(got.Messages))
	}
}

func TestSaveThreadUsesOwnerOnlyPermissions(t *testing.T) {
	root := t.TempDir()
	store := NewStore(root)
	now := time.Now().UTC()

	if err := store.SaveThread(Thread{
		ID:        "thread-1",
		Title:     "Release Prep",
		CreatedAt: now,
		UpdatedAt: now,
	}); err != nil {
		t.Fatalf("SaveThread() error = %v", err)
	}

	info, err := os.Stat(filepath.Join(root, "threads.json"))
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("permissions = %#o, want 0o600", got)
	}
}

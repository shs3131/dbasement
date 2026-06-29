package watcher

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupTestRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "src"), 0755)
	os.WriteFile(filepath.Join(dir, "src", "main.go"), []byte("package main\n"), 0644)
	os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Test\n"), 0644)
	return dir
}

func TestNew(t *testing.T) {
	dir := setupTestRepo(t)
	changes := make(chan ChangeEvent, 10)

	w := New(dir, func(ev ChangeEvent) {
		changes <- ev
	})
	defer w.Stop()

	if w == nil {
		t.Fatal("New() returned nil")
	}
}

func TestDetectFileCreation(t *testing.T) {
	dir := setupTestRepo(t)

	w := New(dir, nil)

	w.scanOnce()

	newFile := filepath.Join(dir, "src", "new.go")
	os.WriteFile(newFile, []byte("package main\n"), 0644)

	ev := w.detectChanges()

	found := false
	for _, c := range ev.FileChanges {
		if c.Type == ChangeCreated && filepath.ToSlash(c.Path) == "src/new.go" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Should detect creation of src/new.go, got %+v", ev.FileChanges)
	}
}

func TestDetectFileModification(t *testing.T) {
	dir := setupTestRepo(t)

	w := New(dir, nil)

	w.scanOnce()

	mainFile := filepath.Join(dir, "src", "main.go")
	os.WriteFile(mainFile, []byte("package main\n\nfunc main() {}\n"), 0644)

	ev := w.detectChanges()

	found := false
	for _, c := range ev.FileChanges {
		if c.Type == ChangeModified && filepath.ToSlash(c.Path) == "src/main.go" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Should detect modification of src/main.go, got %+v", ev.FileChanges)
	}
}

func TestDetectFileDeletion(t *testing.T) {
	dir := setupTestRepo(t)

	w := New(dir, nil)

	w.scanOnce()

	os.Remove(filepath.Join(dir, "README.md"))

	ev := w.detectChanges()

	found := false
	for _, c := range ev.FileChanges {
		if c.Type == ChangeDeleted && c.Path == "README.md" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Should detect deletion of README.md, got %+v", ev.FileChanges)
	}
}

func TestNoChange(t *testing.T) {
	dir := setupTestRepo(t)

	w := New(dir, nil)

	w.scanOnce()

	ev := w.detectChanges()

	if len(ev.FileChanges) != 0 {
		t.Errorf("Should detect no changes, got %+v", ev.FileChanges)
	}
}

func TestShouldIgnoreGit(t *testing.T) {
	dir := setupTestRepo(t)
	os.MkdirAll(filepath.Join(dir, ".git"), 0755)

	w := New(dir, nil)

	if !w.shouldIgnore(filepath.Join(dir, ".git")) {
		t.Error("Should ignore .git directory")
	}

	if !w.shouldIgnore(filepath.Join(dir, ".git", "objects")) {
		t.Error("Should ignore .git subdirectories")
	}
}

func TestShouldIgnoreNodeModules(t *testing.T) {
	dir := setupTestRepo(t)

	w := New(dir, nil)

	if !w.shouldIgnore(filepath.Join(dir, "node_modules")) {
		t.Error("Should ignore node_modules directory")
	}
}

func TestShouldIgnoreDbasement(t *testing.T) {
	dir := setupTestRepo(t)

	w := New(dir, nil)

	if !w.shouldIgnore(filepath.Join(dir, ".dbasement")) {
		t.Error("Should ignore .dbasement directory")
	}
}

func TestStartStop(t *testing.T) {
	dir := setupTestRepo(t)
	w := New(dir, nil)

	w.Start()
	w.Stop()

	// Should not block or panic
}

func TestSetInterval(t *testing.T) {
	dir := setupTestRepo(t)
	w := New(dir, nil)

	w.SetInterval(time.Second)
}

func TestIgnoreNestedDirectories(t *testing.T) {
	dir := setupTestRepo(t)
	os.MkdirAll(filepath.Join(dir, "vendor", "github.com", "foo"), 0755)

	w := New(dir, nil)

	if !w.shouldIgnore(filepath.Join(dir, "vendor")) {
		t.Error("Should ignore vendor directory")
	}
}

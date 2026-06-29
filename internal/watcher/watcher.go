package watcher

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type ChangeType int

const (
	ChangeUnknown ChangeType = iota
	ChangeModified
	ChangeCreated
	ChangeDeleted
	ChangeRenamed
)

type FileChange struct {
	Type   ChangeType
	Path   string
	OldPath string
}

type ChangeEvent struct {
	FileChanges []FileChange
	HasGitDiff  bool
	Timestamp   time.Time
}

type Watcher struct {
	repoPath  string
	ignore    []string
	interval  time.Duration
	fileHashes map[string]string
	mu        sync.RWMutex
	onChange  func(ChangeEvent)
	stopCh    chan struct{}
	running   bool
}

func New(repoPath string, onChange func(ChangeEvent)) *Watcher {
	return &Watcher{
		repoPath:    repoPath,
		interval:    5 * time.Second,
		fileHashes:  make(map[string]string),
		onChange:    onChange,
		stopCh:      make(chan struct{}),
		ignore: []string{
			".git",
			"node_modules",
			"vendor",
			".dbasement",
			".DS_Store",
			"*.pyc",
			"__pycache__",
			"*.exe",
			"*.dll",
			"*.so",
			"*.dylib",
			"*.bin",
			"*.log",
		},
	}
}

func (w *Watcher) SetInterval(d time.Duration) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.interval = d
}

func (w *Watcher) Start() {
	w.mu.Lock()
	if w.running {
		w.mu.Unlock()
		return
	}
	w.running = true
	w.mu.Unlock()

	go w.loop()
}

func (w *Watcher) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.running {
		close(w.stopCh)
		w.running = false
	}
}

func (w *Watcher) loop() {
	w.scanOnce()

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			changes := w.detectChanges()
			if len(changes.FileChanges) > 0 && w.onChange != nil {
				w.onChange(changes)
			}
		case <-w.stopCh:
			return
		}
	}
}

func (w *Watcher) scanOnce() {
	hashes := make(map[string]string)
	filepath.Walk(w.repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if w.shouldIgnore(path) {
				return filepath.SkipDir
			}
			return nil
		}
		if w.shouldIgnore(path) {
			return nil
		}
		if !info.Mode().IsRegular() {
			return nil
		}

		relPath, _ := filepath.Rel(w.repoPath, path)
		hash, err := hashFile(path)
		if err == nil {
			hashes[relPath] = hash
		}
		return nil
	})

	w.mu.Lock()
	w.fileHashes = hashes
	w.mu.Unlock()
}

func (w *Watcher) detectChanges() ChangeEvent {
	w.mu.Lock()
	oldHashes := w.fileHashes
	w.fileHashes = make(map[string]string)
	w.mu.Unlock()

	newHashes := make(map[string]string)
	var changes []FileChange

	filepath.Walk(w.repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if w.shouldIgnore(path) {
				return filepath.SkipDir
			}
			return nil
		}
		if !info.Mode().IsRegular() {
			return nil
		}

		relPath, _ := filepath.Rel(w.repoPath, path)
		hash, err := hashFile(path)
		if err != nil {
			return nil
		}

		newHashes[relPath] = hash

		oldHash, exists := oldHashes[relPath]
		if !exists {
			changes = append(changes, FileChange{
				Type: ChangeCreated,
				Path: relPath,
			})
		} else if oldHash != hash {
			changes = append(changes, FileChange{
				Type: ChangeModified,
				Path: relPath,
			})
		}

		return nil
	})

	for path := range oldHashes {
		if _, exists := newHashes[path]; !exists {
			changes = append(changes, FileChange{
				Type: ChangeDeleted,
				Path: path,
			})
		}
	}

	w.mu.Lock()
	w.fileHashes = newHashes
	w.mu.Unlock()

	return ChangeEvent{
		FileChanges: changes,
		Timestamp:   time.Now(),
	}
}

func (w *Watcher) shouldIgnore(path string) bool {
	rel, err := filepath.Rel(w.repoPath, path)
	if err != nil {
		return true
	}

	parts := strings.Split(rel, string(filepath.Separator))
	for _, part := range parts {
		for _, pattern := range w.ignore {
			if matched, _ := filepath.Match(pattern, part); matched {
				return true
			}
			if part == pattern {
				return true
			}
		}
	}
	return false
}

func hashFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	h := sha256.Sum256(data)
	return fmt.Sprintf("%x", h[:8]), nil
}

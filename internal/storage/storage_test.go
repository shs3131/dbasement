package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestDB(t *testing.T) *DB {
	t.Helper()
	dir := t.TempDir()
	db, err := New(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestNew(t *testing.T) {
	dir := t.TempDir()
	db, err := New(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer db.Close()

	// Verify tables exist by writing and reading
	if err := db.SetMeta("test", "value"); err != nil {
		t.Fatalf("SetMeta() error = %v", err)
	}
	v, err := db.GetMeta("test")
	if err != nil {
		t.Fatalf("GetMeta() error = %v", err)
	}
	if v != "value" {
		t.Fatalf("GetMeta() = %q, want %q", v, "value")
	}
}

func TestMeta(t *testing.T) {
	db := setupTestDB(t)

	tests := []struct {
		key   string
		value string
	}{
		{"key1", "value1"},
		{"key2", "value2 with spaces"},
		{"empty", ""},
	}

	for _, tt := range tests {
		if err := db.SetMeta(tt.key, tt.value); err != nil {
			t.Errorf("SetMeta(%q, %q) error = %v", tt.key, tt.value, err)
		}
		v, err := db.GetMeta(tt.key)
		if err != nil {
			t.Errorf("GetMeta(%q) error = %v", tt.key, err)
		}
		if v != tt.value {
			t.Errorf("GetMeta(%q) = %q, want %q", tt.key, v, tt.value)
		}
	}

	v, err := db.GetMeta("nonexistent")
	if err != nil {
		t.Errorf("GetMeta(nonexistent) error = %v", err)
	}
	if v != "" {
		t.Errorf("GetMeta(nonexistent) = %q, want empty", v)
	}
}

func TestMetaOverwrite(t *testing.T) {
	db := setupTestDB(t)

	db.SetMeta("key", "original")
	db.SetMeta("key", "updated")

	v, _ := db.GetMeta("key")
	if v != "updated" {
		t.Errorf("GetMeta() after overwrite = %q, want %q", v, "updated")
	}
}

func TestSections(t *testing.T) {
	db := setupTestDB(t)

	if err := db.SetSection("architecture", "_main", "frontend and backend"); err != nil {
		t.Fatalf("SetSection() error = %v", err)
	}

	content, err := db.GetSection("architecture", "_main")
	if err != nil {
		t.Fatalf("GetSection() error = %v", err)
	}
	if content != "frontend and backend" {
		t.Errorf("GetSection() = %q, want %q", content, "frontend and backend")
	}

	// Get nonexistent section
	content, err = db.GetSection("nonexistent", "")
	if err != nil {
		t.Errorf("GetSection(nonexistent) error = %v", err)
	}
	if content != "" {
		t.Errorf("GetSection(nonexistent) = %q, want empty", content)
	}
}

func TestSectionsWithSubkeys(t *testing.T) {
	db := setupTestDB(t)

	db.SetSection("api", "/users", "GET /users")
	db.SetSection("api", "/posts", "GET /posts")

	all, err := db.GetSectionAll("api")
	if err != nil {
		t.Fatalf("GetSectionAll() error = %v", err)
	}

	if len(all) != 2 {
		t.Errorf("GetSectionAll() returned %d items, want 2", len(all))
	}
}

func TestChangelog(t *testing.T) {
	db := setupTestDB(t)

	if err := db.AddChangelog("first change"); err != nil {
		t.Fatalf("AddChangelog() error = %v", err)
	}
	if err := db.AddChangelog("second change"); err != nil {
		t.Fatalf("AddChangelog() error = %v", err)
	}

	entries, err := db.GetChangelog(5)
	if err != nil {
		t.Fatalf("GetChangelog() error = %v", err)
	}

	if len(entries) != 2 {
		t.Errorf("GetChangelog() returned %d entries, want 2", len(entries))
	}

	if entries[0].Entry != "second change" {
		t.Errorf("Latest entry = %q, want %q", entries[0].Entry, "second change")
	}
}

func TestChangelogLimit(t *testing.T) {
	db := setupTestDB(t)

	for i := 0; i < 10; i++ {
		db.AddChangelog("change")
	}

	entries, _ := db.GetChangelog(3)
	if len(entries) != 3 {
		t.Errorf("GetChangelog(3) returned %d entries, want 3", len(entries))
	}

	entries, _ = db.GetChangelog(0)
	if len(entries) > 50 {
		t.Errorf("GetChangelog(0) returned %d entries, want <=50", len(entries))
	}
}

func TestDesignDecisions(t *testing.T) {
	db := setupTestDB(t)

	if err := db.AddDesignDecision("Use Go", "Performance and simplicity"); err != nil {
		t.Fatalf("AddDesignDecision() error = %v", err)
	}

	decisions, err := db.GetDesignDecisions()
	if err != nil {
		t.Fatalf("GetDesignDecisions() error = %v", err)
	}

	if len(decisions) != 1 {
		t.Errorf("GetDesignDecisions() returned %d items, want 1", len(decisions))
	}

	if decisions[0].Decision != "Use Go" {
		t.Errorf("Decision = %q, want %q", decisions[0].Decision, "Use Go")
	}
	if decisions[0].Reason != "Performance and simplicity" {
		t.Errorf("Reason = %q, want %q", decisions[0].Reason, "Performance and simplicity")
	}
}

func TestTodo(t *testing.T) {
	db := setupTestDB(t)

	if err := db.AddTodo("Write tests", "ai"); err != nil {
		t.Fatalf("AddTodo() error = %v", err)
	}
	if err := db.AddTodo("Add CI pipeline", "user"); err != nil {
		t.Fatalf("AddTodo() error = %v", err)
	}

	items, err := db.GetTodos(false)
	if err != nil {
		t.Fatalf("GetTodos() error = %v", err)
	}
	if len(items) != 2 {
		t.Errorf("GetTodos() returned %d items, want 2", len(items))
	}

	// Mark one as done
	if err := db.MarkTodoDone(items[0].ID); err != nil {
		t.Fatalf("MarkTodoDone() error = %v", err)
	}

	items, _ = db.GetTodos(false)
	if len(items) != 1 {
		t.Errorf("GetTodos(false) after done = %d, want 1", len(items))
	}

	items, _ = db.GetTodos(true)
	if len(items) != 2 {
		t.Errorf("GetTodos(true) after done = %d, want 2", len(items))
	}
}

func TestKnownIssues(t *testing.T) {
	db := setupTestDB(t)

	if err := db.AddKnownIssue("Memory leak in cache", 90); err != nil {
		t.Fatalf("AddKnownIssue() error = %v", err)
	}

	issues, err := db.GetKnownIssues(false)
	if err != nil {
		t.Fatalf("GetKnownIssues() error = %v", err)
	}
	if len(issues) != 1 {
		t.Errorf("GetKnownIssues() returned %d items, want 1", len(issues))
	}
	if issues[0].Confidence != 90 {
		t.Errorf("Confidence = %d, want 90", issues[0].Confidence)
	}

	if err := db.ResolveKnownIssue(issues[0].ID); err != nil {
		t.Fatalf("ResolveKnownIssue() error = %v", err)
	}

	open, _ := db.GetKnownIssues(false)
	if len(open) != 0 {
		t.Errorf("Unresolved issues after resolve = %d, want 0", len(open))
	}

	all, _ := db.GetKnownIssues(true)
	if len(all) != 1 {
		t.Errorf("All issues after resolve = %d, want 1", len(all))
	}
}

func TestSearchMemory(t *testing.T) {
	db := setupTestDB(t)

	db.SetSection("architecture", "_main", "React frontend, Go backend")
	db.SetSection("api", "_all", "REST API with JWT authentication")
	db.SetSection("features", "_all", "User authentication, File upload")

	results, err := db.SearchMemory("authentication")
	if err != nil {
		t.Fatalf("SearchMemory() error = %v", err)
	}

	if len(results) == 0 {
		t.Errorf("SearchMemory('authentication') returned no results")
	}

	// Check that api section was found
	apiResults, ok := results["api/_all"]
	if !ok {
		t.Errorf("SearchMemory() should find api/_all")
	} else {
		found := false
		for _, r := range apiResults {
			if contains(r, "authentication") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("api/_all should contain 'authentication'")
		}
	}

	// Search for something that doesn't exist
	noResults, _ := db.SearchMemory("zzzznotfound")
	if len(noResults) != 0 {
		t.Errorf("SearchMemory('zzzznotfound') should return no results")
	}
}

func TestDeleteSection(t *testing.T) {
	db := setupTestDB(t)

	db.SetSection("test", "item", "content")
	content, _ := db.GetSection("test", "item")
	if content != "content" {
		t.Fatal("Failed to set section")
	}

	if err := db.DeleteSection("test", "item"); err != nil {
		t.Fatalf("DeleteSection() error = %v", err)
	}

	content, _ = db.GetSection("test", "item")
	if content != "" {
		t.Errorf("Section should be empty after delete, got %q", content)
	}
}

func TestConcurrentAccess(t *testing.T) {
	db := setupTestDB(t)

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(n int) {
			db.SetSection("concurrent", "test", "value")
			db.GetSection("concurrent", "test")
			db.AddChangelog("test")
			db.GetChangelog(10)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestDatabaseFileCreated(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "memory.db")

	db, err := New(path)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	db.Close()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Database file was not created at %s", path)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && search(s, substr)
}

func search(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

package analyzer

import (
	"testing"

	"github.com/shs3131/dbasement/internal/watcher"
)

func TestIsSmallChangeEmpty(t *testing.T) {
	a := New("/test")
	if !a.isSmallChange("") {
		t.Error("Empty diff should be small")
	}
}

func TestIsSmallChangeWhitespace(t *testing.T) {
	a := New("/test")
	diff := `diff --git a/file.go b/file.go
index abc..def 100644
--- a/file.go
+++ b/file.go
@@ -1,3 +1,3 @@
 
 // comment
-// old comment
+// new comment`

	if !a.isSmallChange(diff) {
		t.Error("Comment-only change should be small")
	}
}

func TestIsSmallChangeReal(t *testing.T) {
	a := New("/test")
	diff := `diff --git a/main.go b/main.go
index abc..def 100644
--- a/main.go
+++ b/main.go
@@ -1,5 +1,8 @@
 package main
 
+import "fmt"
+
 func main() {
+	fmt.Println("hello")
 }`

	if a.isSmallChange(diff) {
		t.Error("Real code change should not be small")
	}
}

func TestJudgeRelevanceLow(t *testing.T) {
	a := New("/test")
	relevance := a.judgeRelevance("", "fix typo in comment", map[string]string{})
	if relevance >= RelevanceMedium {
		t.Error("Typo fix should be low relevance")
	}
}

func TestJudgeRelevanceHigh(t *testing.T) {
	a := New("/test")
	relevance := a.judgeRelevance(
		"added new authentication middleware",
		"feat: add JWT auth middleware",
		map[string]string{"auth/middleware.go": "M"},
	)
	if relevance < RelevanceHigh {
		t.Error("Auth middleware should be high relevance")
	}
}

func TestJudgeRelevanceDatabaseMigration(t *testing.T) {
	a := New("/test")
	relevance := a.judgeRelevance(
		"",
		"add users table",
		map[string]string{"db/migrations/001_users.sql": "A"},
	)
	if relevance < RelevanceHigh {
		t.Error("Database migration should be high relevance")
	}
}

func TestIsMeaningfulFile(t *testing.T) {
	a := New("/test")

	tests := []struct {
		path    string
		want    bool
	}{
		{"main.go", true},
		{"api/routes.ts", true},
		{"Dockerfile", true},
		{"package.json", true},
		{"README.md", false},
		{".gitignore", false},
		{"assets/logo.png", false},
	}

	for _, tt := range tests {
		got := a.isMeaningfulFile(tt.path)
		if got != tt.want {
			t.Errorf("isMeaningfulFile(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestAnalyzeFileChangesEmpty(t *testing.T) {
	a := New("/test")
	result, err := a.AnalyzeFileChanges(nil)
	if err != nil {
		t.Fatalf("AnalyzeFileChanges() error = %v", err)
	}
	if result.Relevant {
		t.Error("Empty changes should not be relevant")
	}
}

func TestAnalyzeFileChangesMeaningful(t *testing.T) {
	a := New("/test")
	changes := []watcher.FileChange{
		{Type: watcher.ChangeCreated, Path: "api/users.go"},
		{Type: watcher.ChangeModified, Path: "internal/db/schema.go"},
	}

	result, err := a.AnalyzeFileChanges(changes)
	if err != nil {
		t.Fatalf("AnalyzeFileChanges() error = %v", err)
	}
	if !result.Relevant {
		t.Error("Creating meaningful files should be relevant")
	}
}

func TestAnalyzeFileChangesIgnore(t *testing.T) {
	a := New("/test")
	changes := []watcher.FileChange{
		{Type: watcher.ChangeModified, Path: "README.md"},
		{Type: watcher.ChangeCreated, Path: "assets/logo.png"},
	}

	result, err := a.AnalyzeFileChanges(changes)
	if err != nil {
		t.Fatalf("AnalyzeFileChanges() error = %v", err)
	}
	if result.Relevant {
		t.Error("Non-meaningful file changes should not be relevant")
	}
}

func TestCalculateConfidence(t *testing.T) {
	a := New("/test")

	tests := []struct {
		name string
		diff string
		msg  string
		min  int
	}{
		{"small diff", "changelog entry", "", 50},
		{"with commit msg", "several\nlines\nof\nchanges\nhere\nand\nthere\nfor\nthe\ntest", "feat: add new feature", 70},
		{"large diff", generateLargeDiff(60), "", 60},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := a.calculateConfidence(tt.diff, tt.msg)
			if conf < tt.min {
				t.Errorf("calculateConfidence() = %d, want >= %d", conf, tt.min)
			}
			if conf > 98 {
				t.Errorf("calculateConfidence() = %d, want <= 98", conf)
			}
		})
	}
}

func generateLargeDiff(n int) string {
	diff := ""
	for i := 0; i < n; i++ {
		diff += "line of code\n"
	}
	return diff
}

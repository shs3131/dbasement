package summarizer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzeProject(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "cmd", "app"), 0755)
	os.MkdirAll(filepath.Join(dir, "internal", "db"), 0755)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n"), 0644)
	os.WriteFile(filepath.Join(dir, "internal", "db", "db.go"), []byte("package db\n"), 0644)
	os.WriteFile(filepath.Join(dir, "Dockerfile"), []byte("FROM golang\n"), 0644)

	info := AnalyzeProject(dir)

	if info.Name != filepath.Base(dir) {
		t.Errorf("Name = %q, want %q", info.Name, filepath.Base(dir))
	}
	if info.FileCount < 3 {
		t.Errorf("FileCount = %d, want >= 3", info.FileCount)
	}
	if !info.HasDocker {
		t.Error("Should detect Dockerfile")
	}
	if info.Language != "Go" {
		t.Errorf("Language = %q, want %q", info.Language, "Go")
	}
}

func TestAnalyzeProjectWithTests(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "tests"), 0755)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n"), 0644)
	os.WriteFile(filepath.Join(dir, "main_test.go"), []byte("package main\n"), 0644)

	info := AnalyzeProject(dir)

	if !info.HasTests {
		t.Error("Should detect tests")
	}
}

func TestGenerateSummary(t *testing.T) {
	info := &ProjectInfo{
		Name:      "test-project",
		Language:  "Go",
		FileCount: 10,
		DirCount:  3,
		HasDocker: true,
		HasTests: true,
	}

	summary := GenerateSummary(info)
	if summary == "" {
		t.Error("Summary should not be empty")
	}
}

func TestGenerateArchitecture(t *testing.T) {
	info := &ProjectInfo{
		HasFrontend: true,
		HasBackend:  true,
		HasDatabase: true,
		HasDocker:   true,
	}

	arch := GenerateArchitecture(info)
	if arch == "" {
		t.Error("Architecture should not be empty")
	}
}

func TestInitializeFromGit(t *testing.T) {
	info := &ProjectInfo{
		Name:      "test",
		FileCount: 5,
		DirCount:  2,
	}

	summary, arch, err := InitializeFromGit(info, nil)
	if err != nil {
		t.Fatalf("InitializeFromGit() error = %v", err)
	}
	if summary == "" {
		t.Error("Summary should not be empty")
	}
	if arch == "" {
		t.Error("Architecture should not be empty")
	}
}

func TestAnalyzeProjectEmpty(t *testing.T) {
	dir := t.TempDir()
	info := AnalyzeProject(dir)

	if info.Name != filepath.Base(dir) {
		t.Errorf("Name = %q, want %q", info.Name, filepath.Base(dir))
	}
}

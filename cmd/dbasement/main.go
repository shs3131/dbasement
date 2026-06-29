package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/shs3131/dbasement/internal/analyzer"
	"github.com/shs3131/dbasement/internal/git"
	"github.com/shs3131/dbasement/internal/mcp"
	"github.com/shs3131/dbasement/internal/memory"
	"github.com/shs3131/dbasement/internal/storage"
	"github.com/shs3131/dbasement/internal/watcher"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("dbasement: ")

	var projectPath string
	flag.StringVar(&projectPath, "project", "", "Path to the project root (default: current directory)")
	flag.Parse()

	if projectPath == "" {
		var err error
		projectPath, err = os.Getwd()
		if err != nil {
			log.Fatalf("Cannot determine current directory: %v", err)
		}
	}

	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		log.Fatalf("Invalid project path: %v", err)
	}

	if fi, err := os.Stat(absPath); err != nil || !fi.IsDir() {
		log.Fatalf("Project path must be a valid directory: %s", absPath)
	}

	dbasementDir := filepath.Join(absPath, ".dbasement")
	if err := os.MkdirAll(dbasementDir, 0755); err != nil {
		log.Fatalf("Cannot create .dbasement directory: %v", err)
	}

	dbPath := filepath.Join(dbasementDir, "memory.db")
	db, err := storage.New(dbPath)
	if err != nil {
		log.Fatalf("Cannot open database: %v", err)
	}
	defer db.Close()

	mem := memory.New(db, dbasementDir)
	g := git.New(absPath)
	an := analyzer.New(absPath)

	watcher := watcher.New(absPath, func(ev watcher.ChangeEvent) {
		if len(ev.FileChanges) == 0 {
			return
		}
		analysis, err := an.AnalyzeFileChanges(ev.FileChanges)
		if err != nil || !analysis.Relevant {
			return
		}
		if analysis.AddChangelog {
			mem.AddChangelog(analysis.Summary)
		}
	})
	watcher.Start()
	defer watcher.Stop()

	if g.IsRepo() {
		hash, err := g.LatestCommitHash()
		if err == nil && hash != "" {
			mem.SetLastAnalysis(hash)
		}
	}

	srv := mcp.New(mem, g, an)
	log.Printf("Dbasement started for project: %s", absPath)

	if err := srv.Run(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

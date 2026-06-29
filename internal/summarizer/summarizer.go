package summarizer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shs3131/dbasement/internal/git"
)

type DepInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type RouteInfo struct {
	Method string
	Path   string
	File   string
}

type ProjectInfo struct {
	Name        string
	Language    string
	Framework   string
	FileCount   int
	DirCount    int
	HasDocker   bool
	HasCI       bool
	HasDatabase bool
	HasAPI      bool
	HasFrontend bool
	HasBackend  bool
	HasTests    bool
	Files       []string
}

func AnalyzeProject(projectPath string) *ProjectInfo {
	info := &ProjectInfo{
		Name: filepath.Base(projectPath),
	}

	filepath.Walk(projectPath, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if fi.IsDir() {
			name := fi.Name()
			if name == ".git" || name == "node_modules" || name == "vendor" || name == ".dbasement" {
				return filepath.SkipDir
			}
			info.DirCount++
			return nil
		}

		info.FileCount++
		rel, _ := filepath.Rel(projectPath, path)
		info.Files = append(info.Files, rel)

		ext := filepath.Ext(path)
		switch ext {
		case ".go":
			info.Language = "Go"
		case ".py":
			info.Language = "Python"
		case ".js":
			info.Language = "JavaScript"
		case ".ts", ".tsx":
			info.Language = "TypeScript"
		case ".rs":
			info.Language = "Rust"
		case ".java":
			info.Language = "Java"
		case ".rb":
			info.Language = "Ruby"
		case ".cs":
			info.Language = "C#"
		case ".php":
			info.Language = "PHP"
		case ".swift":
			info.Language = "Swift"
		case ".kt":
			info.Language = "Kotlin"
		}

		name := strings.ToLower(fi.Name())
		if strings.Contains(name, "docker") || strings.Contains(name, "dockerfile") {
			info.HasDocker = true
		}
		if strings.Contains(path, ".github") || strings.Contains(name, "jenkins") || strings.Contains(name, ".gitlab") {
			info.HasCI = true
		}
		if strings.Contains(path, "migration") || strings.Contains(path, "schema") || ext == ".sql" {
			info.HasDatabase = true
		}
		if strings.Contains(path, "api") || strings.Contains(path, "route") || strings.Contains(path, "endpoint") || strings.Contains(path, "controller") {
			info.HasAPI = true
		}
		if strings.Contains(path, "frontend") || strings.Contains(path, "ui") || strings.Contains(path, "component") || strings.Contains(path, "page") || ext == ".css" || ext == ".html" {
			info.HasFrontend = true
		}
		if strings.Contains(path, "backend") || strings.Contains(path, "server") || strings.Contains(path, "service") || strings.Contains(path, "handler") {
			info.HasBackend = true
		}
		if strings.Contains(path, "test") || strings.Contains(path, "spec") || strings.Contains(name, "test") {
			info.HasTests = true
		}

		return nil
	})

	if strings.Contains(info.Name, "-") {
		parts := strings.Split(info.Name, "-")
		for i, p := range parts {
			parts[i] = strings.Title(p)
		}
		info.Name = strings.Join(parts, " ")
	}

	return info
}

func GenerateSummary(info *ProjectInfo) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# %s\n\n", info.Name))

	if info.Language != "" {
		b.WriteString(fmt.Sprintf("Language: %s\n", info.Language))
	}

	b.WriteString(fmt.Sprintf("\n## Overview\n\n"))
	b.WriteString(fmt.Sprintf("This project contains %d files across %d directories.", info.FileCount, info.DirCount))

	var traits []string
	if info.HasFrontend {
		traits = append(traits, "frontend")
	}
	if info.HasBackend {
		traits = append(traits, "backend")
	}
	if info.HasAPI {
		traits = append(traits, "API")
	}
	if info.HasDatabase {
		traits = append(traits, "database")
	}
	if info.HasDocker {
		traits = append(traits, "Docker")
	}
	if info.HasCI {
		traits = append(traits, "CI/CD")
	}
	if info.HasTests {
		traits = append(traits, "tests")
	}

	if len(traits) > 0 {
		b.WriteString(fmt.Sprintf("\nIncludes: %s.", strings.Join(traits, ", ")))
	}

	return b.String()
}

func GenerateArchitecture(info *ProjectInfo) string {
	var b strings.Builder

	b.WriteString("## Architecture\n\n")

	if info.HasFrontend && info.HasBackend {
		b.WriteString("- Architecture: Full-stack application\n")
		b.WriteString("- Communication: API-based\n")
	} else if info.HasFrontend {
		b.WriteString("- Architecture: Frontend application\n")
	} else if info.HasBackend {
		b.WriteString("- Architecture: Backend service\n")
	}

	if info.HasDocker {
		b.WriteString("- Containerization: Docker\n")
	}

	if info.HasDatabase {
		b.WriteString("- Database: Present (see database section)\n")
	}

	return b.String()
}

type DeepInfo struct {
	Features    string
	API         string
	Database    string
	Deps        string
	KnownIssues string
	Todos       string
	Decisions   string
}

func DeepAnalyze(projectPath string) *DeepInfo {
	di := &DeepInfo{}
	pkgJSON := filepath.Join(projectPath, "package.json")
	goMod := filepath.Join(projectPath, "go.mod")

	var deps []string
	var apiRoutes []string
	var dbModels []string

	if data, err := os.ReadFile(pkgJSON); err == nil {
		var pkg struct {
			Name         string            `json:"name"`
			Description  string            `json:"description"`
			Dependencies map[string]string `json:"dependencies"`
			DevDeps      map[string]string `json:"devDependencies"`
			Scripts      map[string]string `json:"scripts"`
		}
		if json.Unmarshal(data, &pkg) == nil {
			if pkg.Description != "" {
				deps = append(deps, fmt.Sprintf("Description: %s", pkg.Description))
			}
			for name, ver := range pkg.Dependencies {
				deps = append(deps, fmt.Sprintf("- %s: %s", name, ver))
			}
			for name, ver := range pkg.DevDeps {
				deps = append(deps, fmt.Sprintf("- %s: %s (dev)", name, ver))
			}
		}
	}

	if data, err := os.ReadFile(goMod); err == nil {
		di.Deps = fmt.Sprintf("## Dependencies\n\nGo module detected.\n```\n%s\n```\n", string(data))
	}

	if len(deps) > 0 {
		di.Deps = fmt.Sprintf("## Dependencies\n\n%s\n", strings.Join(deps, "\n"))
	}

	featureSet := make(map[string]bool)
	filepath.Walk(projectPath, func(path string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() {
			name := strings.ToLower(fi.Name())
			if name == ".git" || name == "node_modules" || name == "vendor" || name == ".dbasement" || name == "dist" {
				return filepath.SkipDir
			}
			return nil
		}

		rel, _ := filepath.Rel(projectPath, path)
		lower := strings.ToLower(rel)

		if strings.HasPrefix(lower, "src") || strings.HasPrefix(lower, "server") {
			if strings.Contains(lower, "route") || strings.Contains(lower, "api") || strings.Contains(lower, "controller") || strings.Contains(lower, "handler") || strings.Contains(lower, "endpoint") {
				if data, err := os.ReadFile(path); err == nil {
					content := string(data)
					lines := strings.Split(content, "\n")
					for _, line := range lines {
						trimmed := strings.TrimSpace(line)
						for _, method := range []string{"GET", "POST", "PUT", "PATCH", "DELETE"} {
							if strings.Contains(trimmed, method+" ") && (strings.Contains(trimmed, "/api/") || strings.Contains(trimmed, "'/api") || strings.Contains(trimmed, `"/api`)) {
								apiRoutes = append(apiRoutes, fmt.Sprintf("- %s (from %s)", trimmed, rel))
							}
						}
					}
				}
			}

			if strings.Contains(lower, "model") || strings.Contains(lower, "schema") || strings.Contains(lower, "migration") {
				if strings.HasSuffix(lower, ".go") || strings.HasSuffix(lower, ".js") || strings.HasSuffix(lower, ".ts") || strings.HasSuffix(lower, ".sql") {
					dbModels = append(dbModels, fmt.Sprintf("- %s", rel))
				}
			}
		}

		if strings.Contains(lower, "test") || strings.Contains(lower, "spec") {
			featureSet["Testing"] = true
		}
		if strings.Contains(lower, "auth") {
			featureSet["Authentication"] = true
		}
		if strings.Contains(lower, "docker") || strings.Contains(lower, "dockerfile") {
			featureSet["Containerization (Docker)"] = true
		}
		if strings.Contains(lower, "admin") {
			featureSet["Admin Panel"] = true
		}
		if strings.Contains(lower, "websocket") || strings.Contains(lower, "socket") || strings.Contains(lower, "ws") {
			featureSet["WebSocket"] = true
		}
		if strings.Contains(lower, "monitor") || strings.Contains(lower, "heartbeat") || strings.Contains(lower, "health") {
			featureSet["Monitoring / Health Checks"] = true
		}
		if strings.Contains(lower, "cron") || strings.Contains(lower, "scheduler") || strings.Contains(lower, "scheduled") {
			featureSet["Scheduled Tasks / Cron"] = true
		}
		if strings.Contains(lower, "email") || strings.Contains(lower, "mail") || strings.Contains(lower, "notify") {
			featureSet["Email / Notifications"] = true
		}
		if strings.Contains(lower, "webhook") {
			featureSet["Webhooks"] = true
		}
		if strings.Contains(lower, "cache") {
			featureSet["Caching"] = true
		}
		if strings.Contains(lower, "rate") && strings.Contains(lower, "limit") {
			featureSet["Rate Limiting"] = true
		}
		if strings.Contains(lower, "log") {
			featureSet["Logging / Audit"] = true
		}
		if strings.Contains(lower, "search") {
			featureSet["Search"] = true
		}
		if strings.Contains(lower, "analytics") || strings.Contains(lower, "stats") || strings.Contains(lower, "statistics") {
			featureSet["Analytics / Statistics"] = true
		}
		if strings.Contains(lower, "ai") || strings.Contains(lower, "assistant") || strings.Contains(lower, "chat") {
			featureSet["AI Assistant / Chat"] = true
		}
		if strings.Contains(lower, "notification") || strings.Contains(lower, "alert") {
			featureSet["Alerting / Notifications"] = true
		}

		return nil
	})

	var features []string
	for f := range featureSet {
		features = append(features, f)
	}
	if len(features) > 0 {
		di.Features = fmt.Sprintf("## Features\n\nDetected from project structure:\n%s\n", strings.Join(features, "\n"))
	}

	if len(apiRoutes) > 0 {
		di.API = fmt.Sprintf("## API Endpoints\n\nDetected routes:\n%s\n", strings.Join(apiRoutes, "\n"))
	}

	if len(dbModels) > 0 {
		di.Database = fmt.Sprintf("## Database / Models\n\nDetected model/schema files:\n%s\n", strings.Join(dbModels, "\n"))
	}

	todoFiles := findPatternInFiles(projectPath, "TODO", "FIXME", "HACK")
	if len(todoFiles) > 0 {
		di.Todos = fmt.Sprintf("## TODOs from Codebase\n\n%s\n", strings.Join(todoFiles, "\n"))
	}

	return di
}

func findPatternInFiles(root string, patterns ...string) []string {
	var results []string
	filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() {
			if fi != nil {
				n := fi.Name()
				if n == ".git" || n == "node_modules" || n == "vendor" || n == ".dbasement" || n == "dist" {
					return filepath.SkipDir
				}
			}
			return nil
		}
		ext := filepath.Ext(path)
		if ext == ".go" || ext == ".js" || ext == ".ts" || ext == ".jsx" || ext == ".tsx" || ext == ".py" || ext == ".md" {
			data, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			content := string(data)
			lines := strings.Split(content, "\n")
			for _, pattern := range patterns {
				for i, line := range lines {
					if strings.Contains(line, pattern) {
						rel, _ := filepath.Rel(root, path)
						results = append(results, fmt.Sprintf("- %s:%d: %s", rel, i+1, strings.TrimSpace(line)))
					}
				}
			}
		}
		return nil
	})
	return results
}

func InitializeFromGit(info *ProjectInfo, g *git.Client) (string, string, error) {
	summary := GenerateSummary(info)
	architecture := GenerateArchitecture(info)

	if g != nil && g.IsRepo() {
		branch, _ := g.CurrentBranch()
		summary += fmt.Sprintf("\n- Git branch: %s", branch)

		commits, _ := g.RecentCommits(3)
		if len(commits) > 0 {
			summary += "\n\n## Recent Commits\n"
			for _, c := range commits {
				summary += fmt.Sprintf("- %s: %s\n", c.Hash[:7], c.Message)
			}
		}
	}

	return summary, architecture, nil
}

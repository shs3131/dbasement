package memory

import (
	"path/filepath"
	"testing"

	"github.com/shs3131/dbasement/internal/storage"
)

func setupTestManager(t *testing.T) *Manager {
	t.Helper()
	dir := t.TempDir()
	db, err := storage.New(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("storage.New() error = %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return New(db, dir)
}

func TestIsInitialized(t *testing.T) {
	m := setupTestManager(t)
	if m.IsInitialized() {
		t.Error("New manager should not be initialized")
	}

	m.MarkInitialized()
	if !m.IsInitialized() {
		t.Error("Manager should be initialized after MarkInitialized()")
	}
}

func TestProjectSummary(t *testing.T) {
	m := setupTestManager(t)

	summary, _ := m.GetProjectSummary()
	if summary != "" {
		t.Error("New project should have empty summary")
	}

	m.SetProjectSummary("A test project")
	summary, _ = m.GetProjectSummary()
	if summary != "A test project" {
		t.Errorf("GetProjectSummary() = %q, want %q", summary, "A test project")
	}
}

func TestArchitecture(t *testing.T) {
	m := setupTestManager(t)

	m.SetArchitecture("Frontend: React\nBackend: Go")
	arch, _ := m.GetArchitecture()
	if arch != "Frontend: React\nBackend: Go" {
		t.Errorf("GetArchitecture() = %q", arch)
	}

	m.SetArchitectureComponent("frontend", "React with TypeScript")
	comp, _ := m.GetArchitectureComponent("frontend")
	if comp != "React with TypeScript" {
		t.Errorf("GetArchitectureComponent() = %q", comp)
	}
}

func TestFeatures(t *testing.T) {
	m := setupTestManager(t)

	m.SetFeatures("- User authentication\n- File upload")
	features, _ := m.GetFeatures()
	if features != "- User authentication\n- File upload" {
		t.Errorf("GetFeatures() = %q", features)
	}
}

func TestAPI(t *testing.T) {
	m := setupTestManager(t)

	m.SetAPI("POST /api/login")
	api, _ := m.GetAPI()
	if api != "POST /api/login" {
		t.Errorf("GetAPI() = %q", api)
	}

	m.SetAPIEndpoint("/users", "GET /api/users")
	endpoint, _ := m.GetAPIEndpoint("/users")
	if endpoint != "GET /api/users" {
		t.Errorf("GetAPIEndpoint() = %q", endpoint)
	}
}

func TestDatabaseSchema(t *testing.T) {
	m := setupTestManager(t)

	m.SetDatabaseSchema("Table: users")
	schema, _ := m.GetDatabaseSchema()
	if schema != "Table: users" {
		t.Errorf("GetDatabaseSchema() = %q", schema)
	}
}

func TestDependencies(t *testing.T) {
	m := setupTestManager(t)

	m.SetDependencies("- Go 1.21: core language")
	deps, _ := m.GetDependencies()
	if deps != "- Go 1.21: core language" {
		t.Errorf("GetDependencies() = %q", deps)
	}
}

func TestDesignDecisions(t *testing.T) {
	m := setupTestManager(t)

	m.AddDesignDecision("Use SQLite", "Zero dependencies")
	decisions, _ := m.GetDesignDecisions()

	if len(decisions) != 1 {
		t.Fatalf("GetDesignDecisions() returned %d, want 1", len(decisions))
	}
	if decisions[0].Decision != "Use SQLite" {
		t.Errorf("Decision = %q, want %q", decisions[0].Decision, "Use SQLite")
	}
}

func TestGlossary(t *testing.T) {
	m := setupTestManager(t)

	m.SetGlossary("Dbasement: Project memory engine")
	glossary, _ := m.GetGlossary()
	if glossary != "Dbasement: Project memory engine" {
		t.Errorf("GetGlossary() = %q", glossary)
	}
}

func TestChangelog(t *testing.T) {
	m := setupTestManager(t)

	m.AddChangelog("Initial setup")
	m.AddChangelog("Added auth")

	entries, _ := m.GetChangelog(5)
	if len(entries) != 2 {
		t.Fatalf("GetChangelog() returned %d, want 2", len(entries))
	}
	if entries[0].Entry != "Added auth" {
		t.Errorf("Latest entry = %q, want %q", entries[0].Entry, "Added auth")
	}
}

func TestTodo(t *testing.T) {
	m := setupTestManager(t)

	m.AddTodo("Write tests", "ai")
	m.AddTodo("Deploy", "user")

	todos, _ := m.GetTodos(false)
	if len(todos) != 2 {
		t.Errorf("GetTodos() = %d, want 2", len(todos))
	}

	m.MarkTodoDone(todos[0].ID)
	open, _ := m.GetTodos(false)
	if len(open) != 1 {
		t.Errorf("Open todos after done = %d, want 1", len(open))
	}
}

func TestKnownIssues(t *testing.T) {
	m := setupTestManager(t)

	m.AddKnownIssue("Bug in login", 85)
	issues, _ := m.GetKnownIssues(false)
	if len(issues) != 1 {
		t.Fatalf("GetKnownIssues() = %d, want 1", len(issues))
	}
	if issues[0].Confidence != 85 {
		t.Errorf("Confidence = %d, want 85", issues[0].Confidence)
	}
}

func TestSearchMemory(t *testing.T) {
	m := setupTestManager(t)

	m.SetProjectSummary("A web application with REST API")
	m.SetArchitecture("React frontend, Go backend")
	m.SetFeatures("User authentication via JWT")

	results, _ := m.SearchMemory("REST")
	if len(results) == 0 {
		t.Error("SearchMemory('REST') should find results")
	}
}

func TestUpdateSection(t *testing.T) {
	m := setupTestManager(t)

	if err := m.UpdateSection("project_summary", "new summary"); err != nil {
		t.Fatalf("UpdateSection() error = %v", err)
	}

	summary, _ := m.GetProjectSummary()
	if summary != "new summary" {
		t.Errorf("GetProjectSummary() = %q", summary)
	}

	if err := m.UpdateSection("unknown", "value"); err == nil {
		t.Error("UpdateSection(unknown) should error")
	}
}

func TestGetRecentChanges(t *testing.T) {
	m := setupTestManager(t)

	changes, _ := m.GetRecentChanges()
	if changes != "No recent changes recorded." {
		t.Errorf("Empty changes = %q", changes)
	}

	m.AddChangelog("Added feature X")
	changes, _ = m.GetRecentChanges()
	if changes == "No recent changes recorded." {
		t.Errorf("Changes should not be empty after adding")
	}
}

func TestLastAnalysis(t *testing.T) {
	m := setupTestManager(t)

	hash, _ := m.GetLastAnalysis()
	if hash != "" {
		t.Error("Last analysis should be empty initially")
	}

	m.SetLastAnalysis("abc123")
	hash, _ = m.GetLastAnalysis()
	if hash != "abc123" {
		t.Errorf("Last analysis = %q, want %q", hash, "abc123")
	}
}

func TestGetAllSections(t *testing.T) {
	m := setupTestManager(t)

	m.SetProjectSummary("Test project")
	m.SetArchitecture("Simple monolith")

	sections, _ := m.GetAllSections()

	if _, ok := sections["project_summary"]; !ok {
		t.Error("GetAllSections() should include project_summary")
	}
	if _, ok := sections["architecture"]; !ok {
		t.Error("GetAllSections() should include architecture")
	}

	if sections["project_summary"] != "Test project" {
		t.Errorf("Summary = %q", sections["project_summary"])
	}
}

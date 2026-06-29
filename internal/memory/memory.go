package memory

import (
	"fmt"
	"strings"
	"sync"

	"github.com/shs3131/dbasement/internal/storage"
)

type Manager struct {
	db  *storage.DB
	mu  sync.RWMutex
	dir string
}

func New(db *storage.DB, dir string) *Manager {
	return &Manager{db: db, dir: dir}
}

func (m *Manager) IsInitialized() bool {
	v, _ := m.db.GetMeta("initialized")
	return v == "true"
}

func (m *Manager) MarkInitialized() error {
	return m.db.SetMeta("initialized", "true")
}

func (m *Manager) SetProjectSummary(summary string) error {
	return m.db.SetSection("project_summary", "", summary)
}

func (m *Manager) GetProjectSummary() (string, error) {
	return m.db.GetSection("project_summary", "")
}

func (m *Manager) SetArchitecture(data string) error {
	return m.db.SetSection("architecture", "_main", data)
}

func (m *Manager) GetArchitecture() (string, error) {
	return m.db.GetSection("architecture", "_main")
}

func (m *Manager) SetArchitectureComponent(component, details string) error {
	return m.db.SetSection("architecture", component, details)
}

func (m *Manager) GetArchitectureComponent(component string) (string, error) {
	return m.db.GetSection("architecture", component)
}

func (m *Manager) SetFeatures(features string) error {
	return m.db.SetSection("features", "_all", features)
}

func (m *Manager) GetFeatures() (string, error) {
	return m.db.GetSection("features", "_all")
}

func (m *Manager) SetAPI(data string) error {
	return m.db.SetSection("api", "_all", data)
}

func (m *Manager) GetAPI() (string, error) {
	return m.db.GetSection("api", "_all")
}

func (m *Manager) SetAPIEndpoint(endpoint, details string) error {
	return m.db.SetSection("api", endpoint, details)
}

func (m *Manager) GetAPIEndpoint(endpoint string) (string, error) {
	return m.db.GetSection("api", endpoint)
}

func (m *Manager) SetDatabaseSchema(schema string) error {
	return m.db.SetSection("database", "_all", schema)
}

func (m *Manager) GetDatabaseSchema() (string, error) {
	return m.db.GetSection("database", "_all")
}

func (m *Manager) SetDependencies(deps string) error {
	return m.db.SetSection("dependencies", "_all", deps)
}

func (m *Manager) GetDependencies() (string, error) {
	return m.db.GetSection("dependencies", "_all")
}

func (m *Manager) AddDesignDecision(decision, reason string) error {
	return m.db.AddDesignDecision(decision, reason)
}

func (m *Manager) GetDesignDecisions() ([]storage.DesignDecision, error) {
	return m.db.GetDesignDecisions()
}

func (m *Manager) SetGlossary(data string) error {
	return m.db.SetSection("glossary", "_all", data)
}

func (m *Manager) GetGlossary() (string, error) {
	return m.db.GetSection("glossary", "_all")
}

func (m *Manager) AddChangelog(entry string) error {
	return m.db.AddChangelog(entry)
}

func (m *Manager) GetChangelog(limit int) ([]storage.ChangelogEntry, error) {
	return m.db.GetChangelog(limit)
}

func (m *Manager) AddTodo(item, source string) error {
	return m.db.AddTodo(item, source)
}

func (m *Manager) GetTodos(includeDone bool) ([]storage.TodoItem, error) {
	return m.db.GetTodos(includeDone)
}

func (m *Manager) MarkTodoDone(id int) error {
	return m.db.MarkTodoDone(id)
}

func (m *Manager) AddKnownIssue(issue string, confidence int) error {
	return m.db.AddKnownIssue(issue, confidence)
}

func (m *Manager) GetKnownIssues(includeResolved bool) ([]storage.KnownIssue, error) {
	return m.db.GetKnownIssues(includeResolved)
}

func (m *Manager) ResolveKnownIssue(id int) error {
	return m.db.ResolveKnownIssue(id)
}

func (m *Manager) SearchMemory(query string) (map[string][]string, error) {
	return m.db.SearchMemory(query)
}

func (m *Manager) UpdateSection(section string, content string) error {
	switch section {
	case "project_summary":
		return m.SetProjectSummary(content)
	case "architecture":
		return m.SetArchitecture(content)
	case "features":
		return m.SetFeatures(content)
	case "api":
		return m.SetAPI(content)
	case "database":
		return m.SetDatabaseSchema(content)
	case "dependencies":
		return m.SetDependencies(content)
	case "glossary":
		return m.SetGlossary(content)
	default:
		return fmt.Errorf("unknown section: %s", section)
	}
}

func (m *Manager) GetSection(section string) (string, error) {
	switch section {
	case "project_summary":
		return m.GetProjectSummary()
	case "architecture":
		return m.GetArchitecture()
	case "features":
		return m.GetFeatures()
	case "api":
		return m.GetAPI()
	case "database":
		return m.GetDatabaseSchema()
	case "dependencies":
		return m.GetDependencies()
	case "glossary":
		return m.GetGlossary()
	default:
		return "", fmt.Errorf("unknown section: %s", section)
	}
}

func (m *Manager) GetAllSections() (map[string]string, error) {
	sections := []string{
		"project_summary",
		"architecture",
		"features",
		"api",
		"database",
		"dependencies",
		"glossary",
	}

	result := make(map[string]string)
	for _, s := range sections {
		content, err := m.GetSection(s)
		if err != nil {
			continue
		}
		if content != "" {
			result[s] = content
		}
	}

	return result, nil
}

func (m *Manager) GetRecentChanges() (string, error) {
	entries, err := m.db.GetChangelog(10)
	if err != nil {
		return "", err
	}

	if len(entries) == 0 {
		return "No recent changes recorded.", nil
	}

	var b strings.Builder
	for _, e := range entries {
		b.WriteString(fmt.Sprintf("[%s] %s\n", e.Timestamp.Format("2006-01-02 15:04"), e.Entry))
	}
	return b.String(), nil
}

func (m *Manager) GetLastAnalysis() (string, error) {
	return m.db.GetMeta("last_analysis")
}

func (m *Manager) SetLastAnalysis(hash string) error {
	return m.db.SetMeta("last_analysis", hash)
}

func (m *Manager) GetLastDiffHash() (string, error) {
	return m.db.GetMeta("last_diff_hash")
}

func (m *Manager) SetLastDiffHash(hash string) error {
	return m.db.SetMeta("last_diff_hash", hash)
}

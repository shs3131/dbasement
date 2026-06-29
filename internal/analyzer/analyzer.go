package analyzer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/shs3131/dbasement/internal/git"
	"github.com/shs3131/dbasement/internal/watcher"
)

type Relevance int

const (
	RelevanceLow Relevance = iota
	RelevanceMedium
	RelevanceHigh
)

type AnalysisResult struct {
	Relevant      bool
	Relevance     Relevance
	Confidence    int
	Summary       string
	Section       string
	Content       string
	ShouldUpdate  bool
	AddChangelog  bool
	ChangelogText string
	NewIssues     []string
	NewTodos      []string
	NewDecisions  []DecisionUpdate
}

type DecisionUpdate struct {
	Decision string
	Reason   string
}

type Analyzer struct {
	projectPath string
}

func New(projectPath string) *Analyzer {
	return &Analyzer{projectPath: projectPath}
}

var smallChangePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?m)^[\s]*$`),
	regexp.MustCompile(`(?m)^[\s]*\/\/.*$`),
	regexp.MustCompile(`(?m)^[\s]*#.*$`),
	regexp.MustCompile(`(?m)^[\s]*\/\*.*?\*\/`),
	regexp.MustCompile(`(?m)^[\s]*\*.*$`),
}

func (a *Analyzer) AnalyzeGitChanges(g *git.Client) (*AnalysisResult, error) {
	diff, err := g.DiffWithHEAD()
	if err != nil {
		diff, err = g.Diff()
		if err != nil {
			return nil, fmt.Errorf("get diff: %w", err)
		}
	}

	if diff == "" {
		staged, err := g.DiffStaged()
		if err == nil && staged != "" {
			diff = staged
		}
	}

	if diff == "" {
		return &AnalysisResult{Relevant: false}, nil
	}

	if a.isSmallChange(diff) {
		return &AnalysisResult{Relevant: false}, nil
	}

	commits, _ := g.RecentCommits(1)
	commitMsg := ""
	if len(commits) > 0 {
		commitMsg = commits[0].Message
	}

	files, _ := g.ChangedFilesWithStatus()

	relevance := a.judgeRelevance(diff, commitMsg, files)

	return &AnalysisResult{
		Relevant:     relevance >= RelevanceMedium,
		Relevance:    relevance,
		Confidence:   a.calculateConfidence(diff, commitMsg),
		Summary:      fmt.Sprintf("Changes detected: %s", truncate(commitMsg, 100)),
		ShouldUpdate: relevance >= RelevanceMedium,
		AddChangelog: relevance >= RelevanceHigh,
	}, nil
}

func (a *Analyzer) AnalyzeFileChanges(changes []watcher.FileChange) (*AnalysisResult, error) {
	if len(changes) == 0 {
		return &AnalysisResult{Relevant: false}, nil
	}

	meaningful := false
	var details []string
	for _, c := range changes {
		switch c.Type {
		case watcher.ChangeCreated:
			if a.isMeaningfulFile(c.Path) {
				meaningful = true
				details = append(details, fmt.Sprintf("Created: %s", c.Path))
			}
		case watcher.ChangeDeleted:
			if a.isMeaningfulFile(c.Path) {
				meaningful = true
				details = append(details, fmt.Sprintf("Deleted: %s", c.Path))
			}
		case watcher.ChangeModified:
			if a.isMeaningfulFile(c.Path) {
				meaningful = true
				details = append(details, fmt.Sprintf("Modified: %s", c.Path))
			}
		}
	}

	if !meaningful {
		return &AnalysisResult{Relevant: false}, nil
	}

	return &AnalysisResult{
		Relevant:     true,
		Relevance:    RelevanceMedium,
		Confidence:   70,
		Summary:      strings.Join(details, "; "),
		ShouldUpdate: true,
		AddChangelog: false,
	}, nil
}

var meaningfulFilePatterns = []*regexp.Regexp{
	regexp.MustCompile(`\.go$`),
	regexp.MustCompile(`\.py$`),
	regexp.MustCompile(`\.js$`),
	regexp.MustCompile(`\.ts$`),
	regexp.MustCompile(`\.tsx$`),
	regexp.MustCompile(`\.jsx$`),
	regexp.MustCompile(`\.rs$`),
	regexp.MustCompile(`\.java$`),
	regexp.MustCompile(`\.kt$`),
	regexp.MustCompile(`\.swift$`),
	regexp.MustCompile(`\.rb$`),
	regexp.MustCompile(`\.php$`),
	regexp.MustCompile(`\.cs$`),
	regexp.MustCompile(`\.sql$`),
	regexp.MustCompile(`\.yaml$`),
	regexp.MustCompile(`\.yml$`),
	regexp.MustCompile(`\.json$`),
	regexp.MustCompile(`\.toml$`),
	regexp.MustCompile(`\.xml$`),
	regexp.MustCompile(`\.proto$`),
	regexp.MustCompile(`\.graphql$`),
	regexp.MustCompile(`Dockerfile`),
	regexp.MustCompile(`docker-compose`),
	regexp.MustCompile(`Makefile`),
	regexp.MustCompile(`\.mod$`),
	regexp.MustCompile(`\.sum$`),
	regexp.MustCompile(`go\.mod`),
	regexp.MustCompile(`go\.sum`),
	regexp.MustCompile(`package\.json`),
	regexp.MustCompile(`requirements\.txt`),
	regexp.MustCompile(`Cargo\.toml`),
	regexp.MustCompile(`pom\.xml`),
	regexp.MustCompile(`build\.gradle`),
}

func (a *Analyzer) isMeaningfulFile(path string) bool {
	for _, pattern := range meaningfulFilePatterns {
		if pattern.MatchString(path) {
			return true
		}
	}
	return false
}

func (a *Analyzer) isSmallChange(diff string) bool {
	lines := strings.Split(diff, "\n")
	if len(lines) < 5 {
		return true
	}

	added := 0
	removed := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			added++
		}
		if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			removed++
		}
	}

	totalChanges := added + removed
	if totalChanges < 3 {
		return true
	}

	codeLines := 0
	commentLines := 0
	whitespaceLines := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			whitespaceLines++
			continue
		}
		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "/*") {
			commentLines++
			continue
		}
		codeLines++
	}

	if codeLines == 0 {
		return true
	}

	commentRatio := float64(commentLines) / float64(max(1, codeLines+commentLines))
	if commentRatio > 0.8 && totalChanges < 10 {
		return true
	}

	return false
}

var meaningfulPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)feature|new|add|create|implement`),
	regexp.MustCompile(`(?i)refactor|rewrite|restructure|redesign`),
	regexp.MustCompile(`(?i)api|endpoint|route|middleware`),
	regexp.MustCompile(`(?i)database|migration|schema|table|index`),
	regexp.MustCompile(`(?i)auth|login|password|token|session|oauth`),
	regexp.MustCompile(`(?i)config|configuration|setting`),
	regexp.MustCompile(`(?i)dependency|package|module|library`),
	regexp.MustCompile(`(?i)deploy|ci|cd|pipeline|action`),
	regexp.MustCompile(`(?i)breaking|incompatible|deprecat`),
	regexp.MustCompile(`(?i)security|vuln|cve|exploit`),
	regexp.MustCompile(`(?i)perf|optimize|bottleneck|slow`),
	regexp.MustCompile(`(?i)remove|delete|deprecate`),
	regexp.MustCompile(`(?i)rename|move|restructur`),
}

func (a *Analyzer) judgeRelevance(diff, commitMsg string, files map[string]string) Relevance {
	score := 0

	for _, pattern := range meaningfulPatterns {
		if pattern.MatchString(commitMsg) {
			score += 2
		}
		if pattern.MatchString(diff) {
			score += 1
		}
	}

	for file := range files {
		if strings.Contains(file, "migration") || strings.Contains(file, "schema") {
			score += 3
		}
		if strings.Contains(file, "api") || strings.Contains(file, "route") || strings.Contains(file, "endpoint") {
			score += 2
		}
		if strings.Contains(file, "config") || strings.Contains(file, "setting") {
			score += 2
		}
		if strings.Contains(file, "docker") || strings.Contains(file, "Dockerfile") {
			score += 2
		}
		if strings.Contains(file, "go.mod") || strings.Contains(file, "package.json") || strings.Contains(file, "Cargo.toml") {
			score += 3
		}
		if strings.Contains(file, "test") || strings.Contains(file, "spec") {
			score -= 1
		}
	}

	switch {
	case score >= 6:
		return RelevanceHigh
	case score >= 3:
		return RelevanceMedium
	default:
		return RelevanceLow
	}
}

func (a *Analyzer) calculateConfidence(diff, commitMsg string) int {
	score := 70

	diffLines := strings.Split(diff, "\n")
	if len(diffLines) > 10 {
		score += 10
	}
	if len(diffLines) > 50 {
		score += 5
	}

	if commitMsg != "" && len(commitMsg) > 10 {
		score += 10
	}

	if score > 98 {
		score = 98
	}
	if score < 50 {
		score = 50
	}

	return score
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

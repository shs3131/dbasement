package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/shs3131/dbasement/internal/analyzer"
	"github.com/shs3131/dbasement/internal/git"
	"github.com/shs3131/dbasement/internal/memory"
)

type Server struct {
	mem     *memory.Manager
	git     *git.Client
	analyze *analyzer.Analyzer
	reader  *bufio.Reader
	writer  io.Writer
	mu      sync.Mutex
	seq     int64
	closed  bool
	initialized bool
}

func New(mem *memory.Manager, g *git.Client, an *analyzer.Analyzer) *Server {
	return &Server{
		mem:     mem,
		git:     g,
		analyze: an,
		reader:  bufio.NewReader(os.Stdin),
		writer:  os.Stdout,
	}
}

func (s *Server) SetIO(r io.Reader, w io.Writer) {
	s.reader = bufio.NewReader(r)
	s.writer = w
}

func (s *Server) Run() error {
	for {
		msg, err := s.readMessage()
		if err != nil {
			if err == io.EOF || s.closed {
				return nil
			}
			return fmt.Errorf("read: %w", err)
		}

		var base struct {
			ID     json.RawMessage `json:"id"`
			Method string          `json:"method"`
		}

		if err := json.Unmarshal(msg, &base); err != nil {
			continue
		}

		s.handleMessage(msg, base)
	}
}

func (s *Server) readMessage() (json.RawMessage, error) {
	for {
		line, err := s.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimSuffix(line, "\r")
		line = strings.TrimSuffix(line, "\n")

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "Content-Length:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				continue
			}
			n := 0
			fmt.Sscanf(strings.TrimSpace(parts[1]), "%d", &n)
			if n <= 0 {
				continue
			}

			for {
				blank, err := s.reader.ReadString('\n')
				if err != nil {
					return nil, err
				}
				if blank == "\n" || blank == "\r\n" {
					break
				}
			}

			body := make([]byte, n)
			_, err := io.ReadFull(s.reader, body)
			if err != nil {
				return nil, err
			}

			var raw json.RawMessage
			if err := json.Unmarshal(body, &raw); err != nil {
				continue
			}
			return raw, nil
		}

		var raw json.RawMessage
		if err := json.Unmarshal([]byte(line), &raw); err != nil {
			continue
		}
		return raw, nil
	}
}

func (s *Server) handleMessage(msg json.RawMessage, base struct {
	ID     json.RawMessage `json:"id"`
	Method string          `json:"method"`
}) {
	switch base.Method {
	case "initialize":
		s.handleInitialize(base.ID, msg)
	case "notifications/initialized":
		s.handleInitialized(base.ID)
	case "tools/list":
		s.handleToolList(base.ID)
	case "tools/call":
		s.handleToolCall(base.ID, msg)
	default:
		s.sendError(base.ID, -32601, fmt.Sprintf("Method not found: %s", base.Method))
	}
}

func (s *Server) handleInitialize(id json.RawMessage, msg json.RawMessage) {
	var params struct {
		ProtocolVersion string          `json:"protocolVersion"`
		Capabilities    json.RawMessage `json:"capabilities"`
		ClientInfo      json.RawMessage `json:"clientInfo"`
	}
	json.Unmarshal(msg, &params)

	s.sendResponse(id, map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{},
		},
		"serverInfo": map[string]string{
			"name":    "dbasement",
			"version": "1.0.0",
		},
	})
}

func (s *Server) handleInitialized(id json.RawMessage) {
	s.initialized = true
}

func (s *Server) handleToolList(id json.RawMessage) {
	tools := []map[string]interface{}{
		{
			"name":        "initialize_project",
			"description": "Initialize project memory database. Scans repository and creates initial project snapshot.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project_path": map[string]interface{}{
						"type":        "string",
						"description": "Path to the project root",
					},
					"summary": map[string]interface{}{
						"type":        "string",
						"description": "Brief project summary (200-400 words)",
					},
					"architecture": map[string]interface{}{
						"type":        "string",
						"description": "Project architecture description",
					},
				},
				"required": []string{"project_path", "summary"},
			},
		},
		{
			"name":        "get_project_summary",
			"description": "Get a concise project summary (~200-400 words). Use this to understand what the project is about.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "get_architecture",
			"description": "Get the project architecture breakdown: frontend, backend, services, modules, communication patterns.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "get_features",
			"description": "Get the list of project features.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "get_api",
			"description": "Get API documentation: endpoints, authentication, request/response formats.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "get_database",
			"description": "Get database schema: tables, collections, relations, indexes.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "get_dependencies",
			"description": "Get project dependencies and why each major dependency exists.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "get_recent_changes",
			"description": "Get recent meaningful project changes from the changelog.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "get_known_issues",
			"description": "Get known project issues and their resolution status.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "get_todo",
			"description": "Get TODO/FIXME items collected from codebase and AI observations.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"include_done": map[string]interface{}{
						"type":        "boolean",
						"description": "Include completed TODO items",
					},
				},
			},
		},
		{
			"name":        "get_design_decisions",
			"description": "Get chronological design decisions with reasoning.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "get_glossary",
			"description": "Get project-specific terminology and definitions.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "search_memory",
			"description": "Search all project memory for specific information.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "Search query",
					},
				},
				"required": []string{"query"},
			},
		},
		{
			"name":        "update_memory",
			"description": "Update a specific section of project memory. Provide confidence score (0-100). Updates with confidence >= 85 are auto-applied; below 70 are ignored.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"section": map[string]interface{}{
						"type":        "string",
						"description": "Memory section to update: project_summary, architecture, features, api, database, dependencies, glossary",
					},
					"content": map[string]interface{}{
						"type":        "string",
						"description": "New content for the section",
					},
					"confidence": map[string]interface{}{
						"type":        "integer",
						"description": "Confidence score 0-100",
					},
					"changelog": map[string]interface{}{
						"type":        "string",
						"description": "Optional changelog entry describing the update",
					},
				},
				"required": []string{"section", "content", "confidence"},
			},
		},
		{
			"name":        "add_design_decision",
			"description": "Record a design decision with reasoning for future reference.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"decision": map[string]interface{}{
						"type":        "string",
						"description": "The design decision made",
					},
					"reason": map[string]interface{}{
						"type":        "string",
						"description": "Why this decision was made",
					},
				},
				"required": []string{"decision", "reason"},
			},
		},
		{
			"name":        "add_todo",
			"description": "Add a TODO item to project memory.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"item": map[string]interface{}{
						"type":        "string",
						"description": "TODO item description",
					},
					"source": map[string]interface{}{
						"type":        "string",
						"description": "Source: 'ai', 'code', 'user'",
					},
				},
				"required": []string{"item"},
			},
		},
		{
			"name":        "add_known_issue",
			"description": "Record a known issue or bug in project memory.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"issue": map[string]interface{}{
						"type":        "string",
						"description": "Issue description",
					},
					"confidence": map[string]interface{}{
						"type":        "integer",
						"description": "Confidence score 0-100",
					},
				},
				"required": []string{"issue", "confidence"},
			},
		},
		{
			"name":        "refresh_project",
			"description": "Check for uncommitted changes and analyze if they are meaningful. Does NOT rescan the whole repository.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			"name":        "resolve_known_issue",
			"description": "Mark a known issue as resolved.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "integer",
						"description": "Issue ID to resolve",
					},
				},
				"required": []string{"id"},
			},
		},
		{
			"name":        "mark_todo_done",
			"description": "Mark a TODO item as completed.",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "integer",
						"description": "TODO item ID",
					},
				},
				"required": []string{"id"},
			},
		},
	}

	s.sendResponse(id, map[string]interface{}{
		"tools": tools,
	})
}

func (s *Server) handleToolCall(id json.RawMessage, msg json.RawMessage) {
	var req struct {
		Params struct {
			Name      string          `json:"name"`
			Arguments json.RawMessage `json:"arguments"`
		} `json:"params"`
	}

	if err := json.Unmarshal(msg, &req); err != nil {
		s.sendError(id, -32602, "Invalid tool call parameters")
		return
	}

	switch req.Params.Name {
	case "initialize_project":
		s.callInitializeProject(id, req.Params.Arguments)
	case "get_project_summary":
		s.callGetProjectSummary(id)
	case "get_architecture":
		s.callGetArchitecture(id)
	case "get_features":
		s.callGetFeatures(id)
	case "get_api":
		s.callGetAPI(id)
	case "get_database":
		s.callGetDatabase(id)
	case "get_dependencies":
		s.callGetDependencies(id)
	case "get_recent_changes":
		s.callGetRecentChanges(id)
	case "get_known_issues":
		s.callGetKnownIssues(id)
	case "get_todo":
		s.callGetTodos(id, req.Params.Arguments)
	case "get_design_decisions":
		s.callGetDesignDecisions(id)
	case "get_glossary":
		s.callGetGlossary(id)
	case "search_memory":
		s.callSearchMemory(id, req.Params.Arguments)
	case "update_memory":
		s.callUpdateMemory(id, req.Params.Arguments)
	case "add_design_decision":
		s.callAddDesignDecision(id, req.Params.Arguments)
	case "add_todo":
		s.callAddTodo(id, req.Params.Arguments)
	case "add_known_issue":
		s.callAddKnownIssue(id, req.Params.Arguments)
	case "refresh_project":
		s.callRefreshProject(id)
	case "resolve_known_issue":
		s.callResolveKnownIssue(id, req.Params.Arguments)
	case "mark_todo_done":
		s.callMarkTodoDone(id, req.Params.Arguments)
	default:
		s.sendError(id, -32601, fmt.Sprintf("Tool not found: %s", req.Params.Name))
	}
}

func (s *Server) callInitializeProject(id json.RawMessage, args json.RawMessage) {
	var params struct {
		ProjectPath  string `json:"project_path"`
		Summary      string `json:"summary"`
		Architecture string `json:"architecture"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		s.sendError(id, -32602, "Invalid arguments")
		return
	}

	if s.mem.IsInitialized() {
		s.sendResult(id, "Project already initialized. Use update_memory to update sections.")
		return
	}

	if params.Summary != "" {
		s.mem.SetProjectSummary(params.Summary)
	}
	if params.Architecture != "" {
		s.mem.SetArchitecture(params.Architecture)
	}

	if params.Summary != "" {
		s.mem.AddChangelog("Project initialized with memory tracking")
	}

	s.mem.MarkInitialized()

	s.sendResult(id, map[string]interface{}{
		"status":  "initialized",
		"message": "Project memory initialized. You can now use other Dbasement tools.",
	})
}

func (s *Server) callGetProjectSummary(id json.RawMessage) {
	summary, _ := s.mem.GetProjectSummary()
	if summary == "" {
		s.sendResult(id, "No project summary available. Use initialize_project to set one.")
		return
	}
	s.sendResult(id, summary)
}

func (s *Server) callGetArchitecture(id json.RawMessage) {
	arch, _ := s.mem.GetArchitecture()
	if arch == "" {
		s.sendResult(id, "No architecture documentation available.")
		return
	}
	s.sendResult(id, arch)
}

func (s *Server) callGetFeatures(id json.RawMessage) {
	features, _ := s.mem.GetFeatures()
	if features == "" {
		s.sendResult(id, "No features documented yet.")
		return
	}
	s.sendResult(id, features)
}

func (s *Server) callGetAPI(id json.RawMessage) {
	api, _ := s.mem.GetAPI()
	if api == "" {
		s.sendResult(id, "No API documentation available.")
		return
	}
	s.sendResult(id, api)
}

func (s *Server) callGetDatabase(id json.RawMessage) {
	db, _ := s.mem.GetDatabaseSchema()
	if db == "" {
		s.sendResult(id, "No database documentation available.")
		return
	}
	s.sendResult(id, db)
}

func (s *Server) callGetDependencies(id json.RawMessage) {
	deps, _ := s.mem.GetDependencies()
	if deps == "" {
		s.sendResult(id, "No dependency documentation available.")
		return
	}
	s.sendResult(id, deps)
}

func (s *Server) callGetRecentChanges(id json.RawMessage) {
	changes, _ := s.mem.GetRecentChanges()
	s.sendResult(id, changes)
}

func (s *Server) callGetKnownIssues(id json.RawMessage) {
	issues, err := s.mem.GetKnownIssues(false)
	if err != nil {
		s.sendError(id, -32603, err.Error())
		return
	}
	if len(issues) == 0 {
		s.sendResult(id, "No known issues.")
		return
	}
	var b strings.Builder
	for _, issue := range issues {
		b.WriteString(fmt.Sprintf("- [ID:%d] (confidence: %d%%) %s\n", issue.ID, issue.Confidence, issue.Issue))
	}
	s.sendResult(id, b.String())
}

func (s *Server) callGetTodos(id json.RawMessage, args json.RawMessage) {
	var params struct {
		IncludeDone bool `json:"include_done"`
	}
	json.Unmarshal(args, &params)

	todos, err := s.mem.GetTodos(params.IncludeDone)
	if err != nil {
		s.sendError(id, -32603, err.Error())
		return
	}
	if len(todos) == 0 {
		s.sendResult(id, "No TODO items.")
		return
	}
	var b strings.Builder
	for _, t := range todos {
		status := " "
		if t.Done {
			status = "x"
		}
		b.WriteString(fmt.Sprintf("- [%s] [ID:%d] (%s) %s\n", status, t.ID, t.Source, t.Item))
	}
	s.sendResult(id, b.String())
}

func (s *Server) callGetDesignDecisions(id json.RawMessage) {
	decisions, err := s.mem.GetDesignDecisions()
	if err != nil {
		s.sendError(id, -32603, err.Error())
		return
	}
	if len(decisions) == 0 {
		s.sendResult(id, "No design decisions recorded.")
		return
	}
	var b strings.Builder
	for _, d := range decisions {
		b.WriteString(fmt.Sprintf("[%s] %s\n", d.Timestamp.Format("2006-01-02 15:04"), d.Decision))
		if d.Reason != "" {
			b.WriteString(fmt.Sprintf("  Reason: %s\n", d.Reason))
		}
	}
	s.sendResult(id, b.String())
}

func (s *Server) callGetGlossary(id json.RawMessage) {
	glossary, _ := s.mem.GetGlossary()
	if glossary == "" {
		s.sendResult(id, "No glossary defined.")
		return
	}
	s.sendResult(id, glossary)
}

func (s *Server) callSearchMemory(id json.RawMessage, args json.RawMessage) {
	var params struct {
		Query string `json:"query"`
	}
	if err := json.Unmarshal(args, &params); err != nil || params.Query == "" {
		s.sendError(id, -32602, "Query is required")
		return
	}

	results, err := s.mem.SearchMemory(params.Query)
	if err != nil {
		s.sendError(id, -32603, err.Error())
		return
	}

	if len(results) == 0 {
		s.sendResult(id, "No results found.")
		return
	}

	var b strings.Builder
	for section, contents := range results {
		b.WriteString(fmt.Sprintf("=== %s ===\n", section))
		for _, content := range contents {
			b.WriteString(truncate(content, 500))
			b.WriteString("\n")
		}
	}
	s.sendResult(id, b.String())
}

func (s *Server) callUpdateMemory(id json.RawMessage, args json.RawMessage) {
	var params struct {
		Section    string `json:"section"`
		Content    string `json:"content"`
		Confidence int    `json:"confidence"`
		Changelog  string `json:"changelog"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		s.sendError(id, -32602, "Invalid arguments")
		return
	}

	if params.Confidence < 70 {
		s.sendResult(id, fmt.Sprintf("Ignored: confidence %d%% is below 70%% threshold", params.Confidence))
		return
	}

	marked := ""
	if params.Confidence < 85 {
		marked = " [AI-inferred]"
	}

	content := params.Content
	if marked != "" {
		content = content + "\n\n--- " + marked
	}

	if err := s.mem.UpdateSection(params.Section, content); err != nil {
		s.sendError(id, -32603, err.Error())
		return
	}

	if params.Changelog != "" {
		s.mem.AddChangelog(params.Changelog)
	}

	status := "applied"
	if marked != "" {
		status = "applied (AI-inferred)"
	}

	s.sendResult(id, fmt.Sprintf("Memory section '%s' updated (%s) with %d%% confidence", params.Section, status, params.Confidence))
}

func (s *Server) callAddDesignDecision(id json.RawMessage, args json.RawMessage) {
	var params struct {
		Decision string `json:"decision"`
		Reason   string `json:"reason"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		s.sendError(id, -32602, "Invalid arguments")
		return
	}

	if err := s.mem.AddDesignDecision(params.Decision, params.Reason); err != nil {
		s.sendError(id, -32603, err.Error())
		return
	}

	s.mem.AddChangelog(fmt.Sprintf("Design decision: %s", params.Decision))
	s.sendResult(id, "Design decision recorded.")
}

func (s *Server) callAddTodo(id json.RawMessage, args json.RawMessage) {
	var params struct {
		Item   string `json:"item"`
		Source string `json:"source"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		s.sendError(id, -32602, "Invalid arguments")
		return
	}

	if params.Source == "" {
		params.Source = "ai"
	}

	if err := s.mem.AddTodo(params.Item, params.Source); err != nil {
		s.sendError(id, -32603, err.Error())
		return
	}

	s.sendResult(id, "TODO item added.")
}

func (s *Server) callAddKnownIssue(id json.RawMessage, args json.RawMessage) {
	var params struct {
		Issue      string `json:"issue"`
		Confidence int    `json:"confidence"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		s.sendError(id, -32602, "Invalid arguments")
		return
	}

	if err := s.mem.AddKnownIssue(params.Issue, params.Confidence); err != nil {
		s.sendError(id, -32603, err.Error())
		return
	}

	s.sendResult(id, "Known issue recorded.")
}

func (s *Server) callRefreshProject(id json.RawMessage) {
	if !s.git.IsRepo() {
		s.sendResult(id, "Not a git repository. File watcher is monitoring for changes.")
		return
	}

	result, err := s.analyze.AnalyzeGitChanges(s.git)
	if err != nil {
		s.sendError(id, -32603, fmt.Sprintf("Analysis error: %v", err))
		return
	}

	if !result.Relevant {
		s.sendResult(id, "No meaningful changes detected.")
		return
	}

	response := fmt.Sprintf("Changes detected (confidence: %d%%): %s", result.Confidence, result.Summary)

	if result.AddChangelog {
		s.mem.AddChangelog(result.Summary)
		response += "\nChangelog updated."
	}

	hash, _ := s.git.LatestCommitHash()
	if hash != "" {
		s.mem.SetLastAnalysis(hash)
	}

	s.sendResult(id, response)
}

func (s *Server) callResolveKnownIssue(id json.RawMessage, args json.RawMessage) {
	var params struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		s.sendError(id, -32602, "Invalid arguments")
		return
	}

	if err := s.mem.ResolveKnownIssue(params.ID); err != nil {
		s.sendError(id, -32603, err.Error())
		return
	}

	s.sendResult(id, "Issue resolved.")
}

func (s *Server) callMarkTodoDone(id json.RawMessage, args json.RawMessage) {
	var params struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		s.sendError(id, -32602, "Invalid arguments")
		return
	}

	if err := s.mem.MarkTodoDone(params.ID); err != nil {
		s.sendError(id, -32603, err.Error())
		return
	}

	s.sendResult(id, "TODO marked as done.")
}

type rpcResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *rpcError   `json:"error,omitempty"`
}

type rpcError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (s *Server) sendResponse(id interface{}, result interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	resp := rpcResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}

	s.writeJSON(resp)
}

func formatResult(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case map[string]interface{}:
		if msg, ok := val["message"]; ok {
			return fmt.Sprintf("%v", msg)
		}
		data, _ := json.MarshalIndent(val, "", "  ")
		return string(data)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (s *Server) sendResult(id interface{}, result interface{}) {
	s.sendResponse(id, map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": formatResult(result),
			},
		},
	})
}

func (s *Server) sendError(id interface{}, code int, message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	resp := rpcResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &rpcError{
			Code:    code,
			Message: message,
		},
	}

	s.writeJSON(resp)
}

func (s *Server) writeJSON(v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return
	}
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))
	fmt.Fprint(s.writer, header)
	fmt.Fprint(s.writer, string(data))
}

func (s *Server) Close() {
	s.closed = true
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

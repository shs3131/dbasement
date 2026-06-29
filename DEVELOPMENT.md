# Development Guide

## Architecture Overview

Dbasement is structured as a clean architecture Go project with independent
modules that communicate through well-defined interfaces.

```
cmd/dbasement/main.go        # Entry point: wires components together
internal/
  storage/storage.go         # SQLite persistence (pure Go, no CGO)
  memory/memory.go           # High-level memory section API
  git/git.go                 # Git operations via CLI
  watcher/watcher.go         # File system change detection
  analyzer/analyzer.go       # Change relevance analysis
  mcp/server.go              # MCP protocol server (JSON-RPC 2.0)
  summarizer/summarizer.go   # Auto-generated project info
```

### Data Flow

```
AI Agent  <--JSON-RPC 2.0 (stdio)-->  mcp.Server
                                         |
                                    memory.Manager
                                         |
                                    storage.DB (SQLite)
```

### Change Detection Flow

```
Git diff / File watcher
        |
    analyzer.Analyzer
        |
    (asks AI if relevant)
        |
    memory.Manager.UpdateSection()
        |
    storage.DB
```

## Project Layout

```
dbasement/
  cmd/dbasement/          # Main binary entry point
  internal/
    storage/              # Database operations (CRUD, search)
    memory/               # Memory section management
    git/                  # Git repository operations
    watcher/              # File system change detection
    analyzer/             # Change relevance scoring
    mcp/                  # MCP protocol implementation
    summarizer/           # Project auto-analysis
  pkg/                    # Public packages (future)
  .github/                # CI, issue templates
```

## Module Responsibilities

### storage

- Pure Go SQLite via `modernc.org/sqlite` (no CGO)
- Thread-safe with `sync.RWMutex` (per-method locking)
- WAL mode for concurrent reads
- Tables: `meta`, `sections`, `changelog`, `design_decisions`, `todo`,
  `known_issues`
- Single connection (`MaxOpenConns = 1`) to avoid SQLite locking issues

### memory

- Wraps `storage.DB` with section-specific methods
- Manages initialization state
- Routes section names to correct storage keys
- Provides `UpdateSection()` and `GetSection()` by section name

### git

- Shells out to `git` CLI (no Go git library dependency)
- Returns parsed diff, commit info, changed files
- All methods are read-only (never stages, commits, or pushes)
- Gracefully handles non-git directories

### watcher

- Poll-based file system monitoring (default 5s interval)
- SHA-256 hash comparison (first 8 hex chars)
- Ignores `.git`, `node_modules`, `vendor`, `.dbasement`, binaries
- Reports created, modified, deleted files

### analyzer

- Detects small changes (formatting, comments, whitespace)
- Scores relevance based on commit message and file patterns
- Calculates confidence (50-98 range)
- Flags meaningful patterns: auth, database, API, config, dependencies

### mcp

- JSON-RPC 2.0 over stdio
- Handles: `initialize`, `notifications/initialized`, `tools/list`,
  `tools/call`
- 20 tools defined with full JSON Schema input validation
- Confidence system: >=85 auto-apply, 70-84 AI-inferred, <70 ignored

### summarizer

- Walks project directory to gather metadata
- Detects language, Docker, CI, database, frontend, backend, tests
- Generates initial summary and architecture descriptions

## Development Workflow

### Setup

```bash
git clone https://github.com/shs3131/dbasement.git
cd dbasement
go mod tidy
```

### Building

```bash
# Development build
go build -o dbasement ./cmd/dbasement/

# Cross-compile for release
GOOS=linux GOARCH=amd64 go build -o dbasement-linux ./cmd/dbasement/
GOOS=darwin GOARCH=amd64 go build -o dbasement-darwin ./cmd/dbasement/
GOOS=windows GOARCH=amd64 go build -o dbasement.exe ./cmd/dbasement/
```

### Testing

```bash
# All tests
go test -count=1 -race ./...

# Verbose tests for a package
go test -v ./internal/storage/...

# With coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Manual Testing with MCP

```bash
# Start Dbasement for any test project
./dbasement --project /tmp/test-project

# In another terminal, send MCP messages
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | ./dbasement --project /tmp/test-project
```

## Code Conventions

- **No comments** in production code (per project convention)
- Idiomatic Go: `gofmt`, `go vet`, `golangci-lint`
- Thread safety: `sync.RWMutex` for all storage operations
- Error handling: return errors, never panic
- Testing: table-driven tests, temp directories for DB tests
- Imports: standard library first, third-party second, internal third

## Adding a New Tool

1. Define the tool in `internal/mcp/server.go` (`handleToolList`)
2. Add a handler function in `internal/mcp/server.go`
3. Register in `handleToolCall` switch statement
4. Add memory methods in `internal/memory/memory.go` if needed
5. Add storage methods in `internal/storage/storage.go` if needed
6. Write tests
7. Update `README.md` tool table

## Performance Targets

| Operation | Target |
|-----------|--------|
| Memory retrieval | <20ms |
| Memory update | <2s |
| Initialization | <30s |
| Idle RAM | ~10-15 MB |
| Idle CPU | ~0% |

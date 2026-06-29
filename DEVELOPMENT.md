# Development Guide

**Source builds are for contributors only. Users should download pre-built binaries from [GitHub Releases](https://github.com/shs3131/dbasement/releases).**

## Prerequisites

- Go 1.26+
- Git

## Quick Start

```bash
git clone https://github.com/shs3131/dbasement.git
cd dbasement
go build -o dbasement ./cmd/dbasement/
```

## Architecture Overview

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
  scripts/                # Install and release scripts
  .github/                # CI, issue templates
```

## Module Responsibilities

### storage
- Pure Go SQLite via `modernc.org/sqlite` (no CGO)
- Thread-safe with `sync.RWMutex`
- WAL mode for concurrent reads
- Tables: `meta`, `sections`, `changelog`, `design_decisions`, `todo`, `known_issues`

### memory
- Wraps `storage.DB` with section-specific methods
- Manages initialization state

### git
- Shells out to `git` CLI (no Go git library)
- All methods are read-only

### watcher
- Poll-based file system monitoring (5s interval)
- SHA-256 hash comparison

### analyzer
- Detects small changes (formatting, comments, whitespace)
- Scores relevance and confidence

### mcp
- JSON-RPC 2.0 over stdio
- 20 tools with JSON Schema input validation
- Confidence system: >=85 auto-apply, 70-84 AI-inferred, <70 ignored

### summarizer
- Walks project directory to gather metadata
- Detects language, Docker, CI, database, etc.

## Building

```bash
# Development build
go build -o dbasement ./cmd/dbasement/

# All platforms (for release)
GOOS=linux GOARCH=amd64 go build -o dbasement-linux-amd64 ./cmd/dbasement/
GOOS=darwin GOARCH=amd64 go build -o dbasement-darwin-amd64 ./cmd/dbasement/
GOOS=windows GOARCH=amd64 go build -o dbasement-windows-amd64.exe ./cmd/dbasement/
```

## Testing

```bash
# All tests
go test -count=1 -race ./...

# Specific package
go test -v ./internal/storage/...

# With coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Manual Testing with MCP

```bash
./dbasement --project /tmp/test-project
```

In another terminal:
```bash
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | ./dbasement --project /tmp/test-project
```

## Code Conventions
- No comments in production code
- Idiomatic Go: `gofmt`, `go vet`, `golangci-lint`
- Thread safety: `sync.RWMutex` for all storage operations
- Error handling: return errors, never panic
- Testing: table-driven tests, temp directories for DB tests

## Adding a New Tool

1. Define tool in `internal/mcp/server.go` (`handleToolList`)
2. Add handler function in `internal/mcp/server.go`
3. Register in `handleToolCall` switch statement
4. Add memory methods in `internal/memory/memory.go` if needed
5. Add storage methods in `internal/storage/storage.go` if needed
6. Write tests
7. Update README.md tool table if adding a new public tool

## Performance Targets

| Operation | Target |
|-----------|--------|
| Memory retrieval | <20ms |
| Memory update | <2s |
| Initialization | <30s |
| Idle RAM | ~10-15 MB |
| Idle CPU | ~0% |

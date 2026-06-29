# Contributing to Dbasement

Thank you for your interest in contributing! Dbasement is designed to be a
minimal, focused tool. Every addition should justify itself.

## Core Principles

Before contributing, understand Dbasement's philosophy:

1. **Zero infrastructure** — No external databases, no vector stores, no
   network services.
2. **Minimal dependencies** — Prefer the standard library. Every dependency
   must justify its weight.
3. **Token efficiency** — Dbasement reduces context, never increases it.
4. **Performance** — Sub-millisecond queries, single-digit millisecond updates.
5. **Portability** — Works on every OS, every terminal, every MCP client.

## How to Contribute

### Reporting Bugs

1. Check existing issues to avoid duplicates.
2. Use the bug report template.
3. Include your OS, Go version, and the exact MCP messages exchanged.
4. Include steps to reproduce.

### Suggesting Features

1. Open a feature request issue first to discuss.
2. Explain how the feature aligns with Dbasement's core principles.
3. Small, focused features are more likely to be accepted.

### Pull Requests

1. Fork the repository.
2. Create a feature branch: `git checkout -b feature/my-feature`.
3. Follow the code style (see [DEVELOPMENT.md](DEVELOPMENT.md)).
4. Add or update tests as needed.
5. Ensure all tests pass: `go test ./...`.
6. Run `go vet ./...` and fix any issues.
7. Update documentation if your change adds or modifies MCP tools.
8. Submit a PR against the `main` branch.

## Development Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/dbasement.git
cd dbasement

# Build
go build ./cmd/dbasement/

# Test
go test ./...

# Run
./dbasement --project /path/to/test-project
```

## Code Style

- Idiomatic Go (follow `gofmt`, `go vet`, `golangci-lint`)
- No unnecessary abstractions
- No external dependencies without strong justification
- Thread-safe by default (use `sync.RWMutex`)
- Prefer explicit error handling over panics
- Test coverage for all public APIs

## Adding a New MCP Tool

1. Add the tool definition in `handleToolList()` in `internal/mcp/server.go`.
2. Add the handler method in the same file.
3. Add the corresponding memory operations in `internal/memory/memory.go`.
4. Add storage operations in `internal/storage/storage.go` if needed.
5. Add tests for all layers.
6. Update the tool table in `README.md`.

## Commit Messages

Follow conventional commits:

```
feat: add get_dependencies tool
fix: handle empty diff in refresh_project
docs: update README with new tool list
test: add storage concurrency tests
refactor: extract formatResult helper
```

## Testing

```bash
# Run all tests
go test -count=1 -race ./...

# Run specific package tests
go test -v ./internal/storage/...

# Run benchmarks
go test -bench=. ./...
```

## Questions?

Open a [Discussion](https://github.com/shs3131/dbasement/discussions) for
questions, ideas, or help getting started.

## License

By contributing, you agree that your contributions will be licensed under the
MIT License.

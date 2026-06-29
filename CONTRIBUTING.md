# Contributing to Dbasement

Thank you for your interest! Dbasement is designed to be minimal, focused, and portable.

## Core Principles

1. **Zero infrastructure** — No external databases, vector stores, or network services
2. **Minimal dependencies** — Prefer standard library; every dependency must justify itself
3. **Token efficiency** — Dbasement reduces context, never increases it
4. **Performance** — Sub-millisecond queries, single-digit millisecond updates
5. **Portability** — Works on every OS, every terminal, every MCP client

## Important: Distribution Model

Dbasement's primary distribution channel is **GitHub Releases with pre-built binaries**.

- Users should never need to build from source
- AI agents should download release binaries automatically
- Source builds are only for contributors adding features or fixing bugs

When adding documentation, always lead with release-based installation. Source build instructions are supplementary.

## How to Contribute

### Reporting Bugs

1. Check existing issues first
2. Use the bug report template
3. Include OS, version, and the exact MCP messages exchanged

### Suggesting Features

1. Open a feature request issue first
2. Explain alignment with core principles
3. Small, focused features preferred

### Pull Requests

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Follow code style (see [DEVELOPMENT.md](DEVELOPMENT.md))
4. Add or update tests
5. Run `go test ./...` and `go vet ./...`
6. Update documentation if changing MCP tools
7. Submit PR against `main` branch

## Development Setup

```bash
git clone https://github.com/YOUR_USERNAME/dbasement.git
cd dbasement
go build ./cmd/dbasement/
go test ./...
```

## Code Style
- Idiomatic Go (`gofmt`, `go vet`, `golangci-lint`)
- No unnecessary abstractions
- No external dependencies without strong justification
- Thread-safe by default (`sync.RWMutex`)
- Prefer explicit error handling over panics
- Test coverage for all public APIs

## Adding a New MCP Tool

1. Add tool definition in `handleToolList()` in `internal/mcp/server.go`
2. Add handler method in the same file
3. Add memory operations in `internal/memory/memory.go`
4. Add storage operations in `internal/storage/storage.go` if needed
5. Add tests for all layers
6. Update tool table in README.md

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
go test -count=1 -race ./...
go test -v ./internal/storage/...
go test -bench=. ./...
```

## Questions?

Open a [Discussion](https://github.com/shs3131/dbasement/discussions).

## License

By contributing, you agree your contributions will be licensed under the MIT License.

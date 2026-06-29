# Roadmap

## v1.1 — Stability & Polish

- [ ] Add integration tests simulating full MCP sessions
- [ ] Add benchmarks for memory operations
- [ ] Improve error messages with actionable guidance
- [ ] Add `--version` flag
- [ ] Add `--help` with full usage documentation
- [ ] Support `DBASEMENT_PROJECT` environment variable

## v1.2 — Deeper Git Integration

- [ ] Track memory per git branch (branch-scoped knowledge)
- [ ] Auto-detect merge conflicts and flag them in memory
- [ ] Link changelog entries to specific commit hashes
- [ ] Support git hooks for automatic memory refresh on commit
- [ ] Detect branch creation/deletion and adjust memory scope

## v1.3 — Smarter Analysis

- [ ] Add diff stat analysis (file type distribution in changes)
- [ ] Improve small-change detection with language-aware parsing
- [ ] Add file-level confidence scoring based on change entropy
- [ ] Detect dependency changes (go.mod, package.json, Cargo.toml diffs)
- [ ] Recognize common refactoring patterns

## v1.4 — Memory Export & Sharing

- [ ] Export memory as Markdown for documentation generation
- [ ] Import memory from a shared `.dbasement/memory.db`
- [ ] Memory diff between branches or commits
- [ ] Human-readable memory snapshots (`dbasement snapshot`)

## v1.5 — Advanced Features

- [ ] Watch mode with configurable polling interval
- [ ] Daemon mode (run in background, communicate via socket)
- [ ] Memory compaction and garbage collection
- [ ] Encrypted memory for sensitive projects
- [ ] Plugin system for custom analyzers

## v2.0 — Community & Ecosystem

- [ ] VS Code extension sidebar with memory browser
- [ ] Pre-built binaries for all platforms via GitHub Releases
- [ ] Homebrew formula, Scoop manifest, APT repository
- [ ] Official website with documentation
- [ ] Community-contributed analyzer plugins
- [ ] Integration test suite against all major MCP clients

## Non-Goals

These items are explicitly out of scope:

- Vector database or embedding search
- Cloud synchronization
- Telemetry or analytics
- Inference or AI model integration
- GUI application
- Real-time collaborative memory

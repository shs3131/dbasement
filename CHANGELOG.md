# Changelog

All notable changes to Dbasement will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-06-29

### Added

- Initial release of Dbasement
- MCP protocol server over stdio (JSON-RPC 2.0)
- 20 MCP tools for project memory management
- SQLite-based persistent storage (pure Go, no CGO)
- Project memory with 7 structured sections:
  - Project Summary, Architecture, Features, API, Database, Dependencies, Glossary
- Changlog with chronological project updates
- Design decisions with reasoning
- Known issues tracking with confidence scoring
- TODO items from codebase and AI observations
- Full-text search across all memory sections
- Git integration for change detection
- File system watcher (poll-based, configurable interval)
- Change relevance analysis (confidence system: >=85 auto-apply, <70 ignore)
- Small change filtering (formatting, comments, whitespace)
- Automatic initialization when AI first discovers Dbasement
- Multi-agent support (Claude Code, Codex CLI, Cursor, Gemini CLI, Cline, Roo Code, Aider)
- Cross-platform (Windows, macOS, Linux)
- Single native binary, zero runtime dependencies
- 38 unit tests

### Performance

- Initialization: <30 seconds
- Memory retrieval: <20ms
- Memory update: <2s
- RAM usage: ~10-15 MB idle
- Binary size: ~11 MB

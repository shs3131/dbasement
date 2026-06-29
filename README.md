<div align="center">
  <img src="https://raw.githubusercontent.com/shs3131/dbasement/master/assets/logo.svg" alt="Dbasement Logo" width="120" height="120">
  <h1>Dbasement</h1>
  <p><strong>Persistent project memory for AI coding agents.</strong></p>
  <p>Never explain your project twice. Your AI remembers everything.</p>
</div>

<br>

<div align="center">

[![CI](https://github.com/shs3131/dbasement/actions/workflows/ci.yml/badge.svg)](https://github.com/shs3131/dbasement/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/shs3131/dbasement?logo=github)](https://github.com/shs3131/dbasement/releases)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

</div>

---

> **Are you an AI agent?**
>
> Read `AI_QUICKSTART.md` first. If more detail is needed, continue with `AGENTS.md`.

---

## What is Dbasement?

Every AI coding session starts from zero. The AI has no memory of what the project does, how it's structured, what changed yesterday, or why decisions were made. Dbasement eliminates this by maintaining a compact, structured, persistent memory of the project — accessible instantly by any MCP-compatible AI agent.

**Dbasement is not an AI. It never performs inference.** When analysis is needed, it asks the connected AI. This keeps it lightweight, fast, and focused.

## Features

- **Persistent memory** — Project knowledge survives between AI sessions
- **Structured knowledge** — Summary, architecture, API, database, dependencies, design decisions, glossary, known issues, TODOs, changelog
- **Incremental updates** — Never regenerate the entire memory; update only what changed
- **Git-aware** — Detects changes via `git diff`; no full rescans
- **File watching** — Optional poll-based watcher for non-git changes
- **Confidence system** — 0-100% scoring; low-confidence updates flagged or ignored
- **Small change filtering** — Ignores formatting, comments, whitespace, typo fixes
- **Multi-agent** — Works with Claude Code, Cursor, Codex CLI, Gemini CLI, Cline, Roo Code, Aider, and any MCP client
- **Zero infrastructure** — SQLite database, no vector DB, no Elasticsearch, no Redis, no cloud
- **Single binary** — ~11 MB, no runtime dependencies
- **Cross-platform** — Windows, macOS, Linux

---

## Quick Start

### 1. Install

Download the binary from the [latest release](https://github.com/shs3131/dbasement/releases/latest) and extract it into your project root:

```bash
# Linux / macOS
curl -sL https://github.com/shs3131/dbasement/releases/latest/download/dbasement-linux-amd64.tar.gz | tar xz

# Windows (PowerShell)
curl.exe -LO https://github.com/shs3131/dbasement/releases/latest/download/dbasement-windows-amd64.zip
Expand-Archive -Path dbasement-windows-amd64.zip -DestinationPath . -Force
```

Or use the auto-download scripts (binary downloads itself on first run):

```bash
# Linux / macOS
bash scripts/start.sh .

# Windows
pwsh -ExecutionPolicy Bypass -File scripts/start.ps1 --project .
```

### 2. Configure your AI client

See [MCP_CONFIGS.md](MCP_CONFIGS.md) for configuration instructions for every supported AI client.

### 3. Start a new AI session

The AI detects Dbasement and initializes project memory automatically. From that point on, every session picks up where the last one left off.

---

## How It Works

### Before Dbasement

```
Session 1: AI scans 500 files to understand the project
Session 2: AI scans 500 files again (no memory)
Session 3: AI scans 500 files again
Each session: tokens wasted, time wasted, context full
```

### With Dbasement

```
Session 1: AI scans 500 files → stores structured knowledge in 2KB
Session 2: AI retrieves 2KB of knowledge (20ms)
Session 3: AI retrieves 2KB of knowledge (20ms)
After a change: AI updates only the affected section (<2s)
```

### Memory Structure

```
.dbasement/
  memory.db              # SQLite database
    ├── project_summary  # 200-400 word explanation
    ├── architecture     # Frontend, backend, services, modules
    ├── features         # Feature list
    ├── api              # Endpoints, auth, requests
    ├── database         # Tables, relations, indexes
    ├── dependencies     # Why each dependency exists
    ├── design_decisions # Chronological decisions with reasoning
    ├── known_issues     # Tracked with confidence scores
    ├── todo             # From codebase and AI observations
    ├── changelog        # Meaningful project updates
    └── glossary         # Project terminology
```

---

## Available MCP Tools

| Tool | Description |
|------|-------------|
| `get_project_summary` | Get concise project description |
| `get_architecture` | Get architecture breakdown |
| `get_features` | Get feature list |
| `get_api` | Get API documentation |
| `get_database` | Get database schema |
| `get_dependencies` | Get dependency documentation |
| `get_recent_changes` | Get recent changelog |
| `get_known_issues` | Get unresolved issues |
| `get_todo` | Get TODO items |
| `get_design_decisions` | Get decision history |
| `get_glossary` | Get project terminology |
| `search_memory` | Full-text search across all memory |
| `initialize_project` | Initialize memory with project summary and architecture |
| `update_memory` | Update a memory section |
| `add_design_decision` | Record a design decision |
| `add_todo` | Add a TODO item |
| `add_known_issue` | Add a known issue |
| `refresh_project` | Check for meaningful changes |
| `resolve_known_issue` | Mark issue resolved |
| `mark_todo_done` | Mark TODO complete |

---

## Release Assets

Pre-built binaries are available for every [release](https://github.com/shs3131/dbasement/releases):

| Platform | Architecture | Archive |
|----------|-------------|---------|
| Linux | x86_64 | `dbasement-linux-amd64.tar.gz` |
| Linux | ARM64 | `dbasement-linux-arm64.tar.gz` |
| macOS | Intel | `dbasement-darwin-amd64.tar.gz` |
| macOS | Apple Silicon | `dbasement-darwin-arm64.tar.gz` |
| Windows | x86_64 | `dbasement-windows-amd64.zip` |

Each archive contains a single binary named `dbasement` (or `dbasement.exe` on Windows).

---

## FAQ

### Does Dbasement send my code anywhere?

No. Dbasement runs entirely locally. It never makes network requests, sends telemetry, or communicates with anything other than your AI agent via stdio.

### Does Dbasement modify my code?

No. Dbasement is read-only with respect to your project files. It only writes to `.dbasement/memory.db`.

### Can I use Dbasement without Git?

Yes. The file watcher detects changes via SHA-256 hash comparison. Git integration is preferred but optional.

### How big does the database get?

The `.dbasement/memory.db` file typically stays under 1 MB for most projects.

### Can multiple AI agents use Dbasement simultaneously?

Yes. SQLite handles read concurrency natively with WAL mode. Write operations are serialized.

### How do I remove Dbasement?

```bash
rm -rf .dbasement/
```

---

## Documentation

| File | Purpose |
|------|---------|
| [AI_QUICKSTART.md](AI_QUICKSTART.md) | First file AI agents should read |
| [AGENTS.md](AGENTS.md) | Comprehensive AI agent guide |
| [MCP_CONFIGS.md](MCP_CONFIGS.md) | MCP configuration for every AI client |
| [INSTALL.md](INSTALL.md) | Platform-specific installation instructions |
| [DEVELOPMENT.md](DEVELOPMENT.md) | Building from source (contributors only) |
| [CONTRIBUTING.md](CONTRIBUTING.md) | Contributing guidelines |

---

## License

MIT License. See [LICENSE](LICENSE) for details.

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
[![Go Report Card](https://goreportcard.com/badge/github.com/shs3131/dbasement)](https://goreportcard.com/report/github.com/shs3131/dbasement)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

</div>

---

## Setup (30 seconds, no Go required)

**1. Download the binary** from the [latest release](https://github.com/shs3131/dbasement/releases/latest):

| Platform | Download |
|----------|----------|
| Windows | `dbasement-windows-amd64.zip` |
| Linux | `dbasement-linux-amd64.tar.gz` |
| macOS (Intel) | `dbasement-darwin-amd64.tar.gz` |
| macOS (Apple Silicon) | `dbasement-darwin-arm64.tar.gz` |

**2. Extract it into your project root:**

```powershell
# Windows
Expand-Archive -Path dbasement-windows-amd64.zip -DestinationPath .
```

```bash
# Linux / macOS
tar xzf dbasement-linux-amd64.tar.gz && chmod +x dbasement
```

**3. Done.** Your AI will auto-discover Dbasement and use it as a memory server.

> If your AI client uses `.vscode/mcp.json`, it's already included in the project — just open the folder and Dbasement is ready.

---

## AI Client Configuration

### VS Code (auto-discovered)

Open the project folder. `.vscode/mcp.json` is already configured. The server appears automatically in VS Code's MCP list.

### Claude Code

```bash
claude --mcp "C:\path\to\dbasement.exe --project C:\path\to\project"
```

Or in `~/.claude/settings.json`:

```json
{
  "mcpServers": {
    "dbasement": {
      "command": "C:\\full\\path\\to\\dbasement.exe",
      "args": ["--project", "C:\\full\\path\\to\\project"]
    }
  }
}
```

### Cursor

Settings → MCP Servers → Add Server:

- **Name**: `dbasement`
- **Type**: `command`
- **Command**: `C:\path\to\dbasement.exe --project C:\path\to\project`

### Cline / Roo Code

`~/.config/cline/mcp.json` or `~/.config/roo/mcp.json`:

```json
{
  "mcpServers": {
    "dbasement": {
      "command": "/full/path/to/dbasement",
      "args": ["--project", "/full/path/to/project"]
    }
  }
}
```

### Codex CLI

`.codex/mcp.json` in the project root:

```json
{
  "mcpServers": {
    "dbasement": {
      "command": "./dbasement",
      "args": ["--project", "."]
    }
  }
}
```

### Gemini CLI

`~/.config/gemini/mcp.json`:

```json
{
  "mcpServers": {
    "dbasement": {
      "command": "/full/path/to/dbasement",
      "args": ["--project", "/full/path/to/project"]
    }
  }
}
```

### Aider

```bash
aider --mcp "/full/path/to/dbasement --project /full/path/to/project"
```

Or in `.aider.conf.yml`:

```yaml
mcp: /full/path/to/dbasement --project /full/path/to/project
```

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

## Features

- **Zero dependencies** — Single ~11 MB binary, no Go, no Node, no Python
- **Structured knowledge** — 10 sections: summary, architecture, API, DB, deps, decisions, TODOs, issues, changelog, glossary
- **Incremental updates** — Never regenerate the entire memory; update only what changed
- **Git-aware** — Detects changes via `git diff`; no full rescans
- **Confidence system** — ≥85 auto-apply, 70-84 AI-inferred, <70 ignore
- **Multi-agent** — Works with Claude Code, Cursor, Codex CLI, Gemini CLI, Cline, Roo Code, Aider, and any MCP client
- **Cross-platform** — Windows, macOS, Linux
- **Private** — Runs entirely locally, no network, no telemetry, no cloud

---

## Available MCP Tools (20 tools)

### Retrieval

| Tool | What it returns |
|------|----------------|
| `get_project_summary` | 200-400 word project description |
| `get_architecture` | Frontend/backend/service breakdown |
| `get_features` | List of features |
| `get_api` | Endpoints, auth, request/response |
| `get_database` | Tables, collections, relations |
| `get_dependencies` | Why each dependency exists |
| `get_recent_changes` | Recent changelog entries |
| `get_known_issues` | Open issues with confidence |
| `get_todo` | TODO/FIXME items |
| `get_design_decisions` | Decision history with reasons |
| `get_glossary` | Project terminology |

### Mutation

| Tool | What it does |
|------|-------------|
| `initialize_project` | Set up memory for a new project |
| `update_memory` | Update a section (with confidence score) |
| `add_design_decision` | Record a decision with reasoning |
| `add_todo` | Add a TODO item |
| `add_known_issue` | Add a known issue |
| `resolve_known_issue` | Mark issue as resolved |
| `mark_todo_done` | Mark TODO as complete |

### Analysis

| Tool | What it does |
|------|-------------|
| `search_memory` | Full-text search across all sections |
| `refresh_project` | Check for meaningful git changes |

---

## Performance

| Operation | Target |
|-----------|--------|
| Memory retrieval | <20ms |
| Memory update | <2s |
| Idle RAM | ~10-15 MB |
| Binary size | ~11 MB |

---

## Building from Source

Requires [Go 1.26+](https://go.dev/dl).

```bash
git clone https://github.com/shs3131/dbasement.git
cd dbasement
go build -o dbasement ./cmd/dbasement/
```

---

## FAQ

### Does Dbasement send my code anywhere?

No. Runs entirely locally. Never makes network requests.

### Does Dbasement modify my code?

No. Read-only. Only writes to `.dbasement/memory.db`.

### How do I remove it?

```bash
rm -rf .dbasement/
```

### Can I use it without Git?

Yes. File watcher uses SHA-256 hash comparison instead.

### Can multiple AI agents use it simultaneously?

Yes. SQLite WAL mode handles read concurrency.

---

## Project Status

Active development. Core functionality is stable. See [ROADMAP.md](ROADMAP.md).

## License

MIT

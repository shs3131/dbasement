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

## AI Installation Rules

**Read this first if you are an AI agent.**

Your goal: install Dbasement automatically in the user's project.

1. **Use GitHub Releases.** Never install Go. Never build from source.
2. **Detect the platform** (OS + architecture) of the user's machine.
3. **Download the latest release** binary from GitHub Releases.
4. **Extract** into the project directory.
5. **Configure the MCP server** for the current AI client.
6. **Preserve existing MCP servers** — merge, don't overwrite.
7. **Tell the user** to start a new AI session.

> Source builds are for contributors only. If the user asks to build from source, tell them:
> "Pre-built binaries are available at https://github.com/shs3131/dbasement/releases — do you still want to build from source?"

---

## Quick Install (for AI agents)

### 1. Get the binary

```bash
# Linux / macOS (x86_64)
curl -sL https://github.com/shs3131/dbasement/releases/latest/download/dbasement-linux-amd64.tar.gz | tar xz

# Linux (ARM64)
curl -sL https://github.com/shs3131/dbasement/releases/latest/download/dbasement-linux-arm64.tar.gz | tar xz

# macOS (Apple Silicon)
curl -sL https://github.com/shs3131/dbasement/releases/latest/download/dbasement-darwin-arm64.tar.gz | tar xz

# Windows (PowerShell)
curl.exe -LO https://github.com/shs3131/dbasement/releases/latest/download/dbasement-windows-amd64.zip
Expand-Archive -Path dbasement-windows-amd64.zip -DestinationPath . -Force
```

### 2. Configure the MCP server

Give the MCP config below to the user's AI client. The binary auto-downloads itself on first run if missing.

### 3. Tell the user

> Dbasement is installed. Start a new AI session and the agent will automatically initialize project memory.

---

## MCP Configuration

### Cross-platform (recommended)

Paste this into your AI client's MCP configuration:

```json
{
  "mcpServers": {
    "dbasement": {
      "command": "bash",
      "args": ["scripts/start.sh", "."]
    }
  }
}
```

On Windows, use:

```json
{
  "mcpServers": {
    "dbasement": {
      "command": "pwsh",
      "args": ["-ExecutionPolicy", "Bypass", "-File", "scripts/start.ps1", "--project", "."]
    }
  }
}
```

> The `start.sh` / `start.ps1` scripts auto-download the binary from GitHub Releases if it's missing. No Go, no manual download.

### VS Code (auto-discovered)

If the project has `.vscode/mcp.json`, Dbasement appears automatically when you open the folder in VS Code.

### Claude Code

In `~/.claude/settings.json` or project `.claude/settings.json`:

```json
{
  "mcpServers": {
    "dbasement": {
      "command": "bash",
      "args": ["scripts/start.sh", "."]
    }
  }
}
```

Or run:

```bash
claude --mcp "bash scripts/start.sh ."
```

### Cursor

Settings → MCP → Add Server:

- **Name**: `dbasement`
- **Type**: `command`
- **Command**: `bash scripts/start.sh .`

### Cline / Roo Code

In `~/.config/cline/mcp.json` or `~/.config/roo/mcp.json`:

```json
{
  "mcpServers": {
    "dbasement": {
      "command": "bash",
      "args": ["scripts/start.sh", "."]
    }
  }
}
```

### Codex CLI

In `.codex/mcp.json`:

```json
{
  "mcpServers": {
    "dbasement": {
      "command": "bash",
      "args": ["scripts/start.sh", "."]
    }
  }
}
```

### Gemini CLI

In `~/.config/gemini/mcp.json`:

```json
{
  "mcpServers": {
    "dbasement": {
      "command": "bash",
      "args": ["scripts/start.sh", "."]
    }
  }
}
```

### Aider

```bash
aider --mcp "bash scripts/start.sh ."
```

Or in `.aider.conf.yml`:

```yaml
mcp: bash scripts/start.sh .
```

### All MCP configs in one file

See [MCP_CONFIGS.md](MCP_CONFIGS.md) for every supported client in a single copy-paste location.

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

### Retrieval

| Tool | What it returns | When to use |
|------|----------------|-------------|
| `get_project_summary` | 200-400 word project description | New session, new contributor |
| `get_architecture` | Frontend/backend/service breakdown | Understanding structure |
| `get_features` | List of features | Feature planning |
| `get_api` | Endpoints, auth, request/response | API work |
| `get_database` | Tables, collections, relations | Schema changes |
| `get_dependencies` | Why each dependency exists | Dep management |
| `get_recent_changes` | Recent changelog entries | "What changed?" |
| `get_known_issues` | Open issues with confidence | Bug fixes |
| `get_todo` | TODO/FIXME items | Task planning |
| `get_design_decisions` | Decision history with reasons | Understanding why |
| `get_glossary` | Project terminology | Onboarding |

### Mutation

| Tool | What it does |
|------|-------------|
| `initialize_project` | Set up memory for a new project |
| `update_memory` | Update a specific section (with confidence score) |
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

## AI Agent Workflow

For detailed workflows, rules, and examples, see:
- [AGENTS.md](AGENTS.md) — Full AI agent guide with installation rules, workflows, and tool reference
- [AI_QUICKSTART.md](AI_QUICKSTART.md) — One-page machine-readable quickstart

### First Session

```
1. Read project files (README, configs, source code)
2. Call initialize_project with summary and architecture
3. Call update_memory for each section you discover
4. Call add_design_decision for architectural choices
```

### Returning Sessions

```
1. Call get_project_summary (minimum context)
2. Call get_recent_changes to see what's new
3. Call get_todo for pending tasks
4. Call other get_* tools only as needed for the task
```

### After Code Changes

```
1. Call refresh_project to check git diff
2. If changes detected, read the actual changes
3. Call update_memory for affected sections
4. Call add_design_decision if warranted
```

---

## Release Assets

Pre-built binaries for every release:

| Platform | Architecture | Archive |
|----------|-------------|---------|
| Linux | x86_64 | `dbasement-linux-amd64.tar.gz` |
| Linux | ARM64 | `dbasement-linux-arm64.tar.gz` |
| macOS | Intel | `dbasement-darwin-amd64.tar.gz` |
| macOS | Apple Silicon | `dbasement-darwin-arm64.tar.gz` |
| Windows | x86_64 | `dbasement-windows-amd64.zip` |

Each archive contains a single binary named `dbasement` (or `dbasement.exe` on Windows).

---

## Features

- **Persistent memory** — Project knowledge survives between AI sessions
- **Structured knowledge** — Summary, architecture, API, database, dependencies, design decisions, glossary, known issues, TODOs, changelog
- **Incremental updates** — Never regenerate the entire memory; update only what changed
- **Smart retrieval** — AI retrieves only the section it needs, not the whole database
- **Git-aware** — Detects changes via `git diff`; no full rescans
- **File watching** — Optional poll-based watcher for non-git changes
- **Confidence system** — 0-100% confidence scoring; low-confidence updates marked or ignored
- **Small change filtering** — Ignores formatting, comments, whitespace, typo fixes
- **Multi-agent** — Works with Claude Code, Cursor, Codex CLI, Gemini CLI, Cline, Roo Code, Aider, and any MCP client
- **Zero infrastructure** — SQLite database, no vector DB, no Elasticsearch, no Redis, no cloud
- **Single binary** — ~11 MB, no runtime dependencies
- **Cross-platform** — Windows, macOS, Linux

---

## FAQ

### Does Dbasement send my code anywhere?

No. Dbasement runs entirely locally. It never makes network requests, sends telemetry, or communicates with anything other than your AI agent via stdio.

### Does Dbasement modify my code?

No. Dbasement is read-only with respect to your project files. It only writes to `.dbasement/memory.db`. It never touches your source code.

### Can I use Dbasement without Git?

Yes. The file watcher detects changes via SHA-256 hash comparison. Git integration is preferred but optional.

### How big does the database get?

The `.dbasement/memory.db` file typically stays under 1 MB for most projects. Each memory section is a few hundred bytes to a few kilobytes.

### Can multiple AI agents use Dbasement simultaneously?

Yes. Dbasement supports multiple connections. SQLite handles read concurrency natively with WAL mode. Write operations are serialized.

### Will Dbasement slow down my AI?

No. Memory retrieval takes <20ms. Memory updates take <2s. The AI spends less time understanding the project, not more.

### How do I update the memory?

The AI updates memory automatically via `update_memory`. You can also manually trigger a refresh with `refresh_project`.

### How do I remove Dbasement?

```bash
rm -rf .dbasement/
```

---

## Documentation

| File | Purpose |
|------|---------|
| [AGENTS.md](AGENTS.md) | AI agent guide: installation rules, workflows, tool reference |
| [AI_QUICKSTART.md](AI_QUICKSTART.md) | One-page machine-readable quickstart for LLMs |
| [MCP_CONFIGS.md](MCP_CONFIGS.md) | MCP configuration for every supported AI client |
| [INSTALL.md](INSTALL.md) | Platform-specific installation instructions |
| [GET_STARTED.md](GET_STARTED.md) | Human-friendly getting started guide |
| [DEVELOPMENT.md](DEVELOPMENT.md) | Building from source (contributors only) |
| [CONTRIBUTING.md](CONTRIBUTING.md) | Contributing guidelines |

## Project Status

Dbasement is in **active development**. The core functionality is stable and usable.

## License

MIT License. See [LICENSE](LICENSE) for details.

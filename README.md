<div align="center">
  <img src="https://raw.githubusercontent.com/shs3131/dbasement/main/assets/logo.svg" alt="Dbasement Logo" width="120" height="120">
  <h1>Dbasement</h1>
  <p><strong>Persistent project memory for AI coding agents.</strong></p>
  <p>Never explain your project twice. Your AI remembers everything.</p>
</div>

<br>

<div align="center">

[![CI](https://github.com/shs3131/dbasement/actions/workflows/ci.yml/badge.svg)](https://github.com/shs3131/dbasement/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/shs3131/dbasement?logo=github)](https://github.com/shs3131/dbasement/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/shs3131/dbasement)](https://goreportcard.com/report/github.com/shs3131/dbasement)
[![Go Version](https://img.shields.io/badge/go-1.26+-00ADD8?logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen)](CONTRIBUTING.md)

</div>

---

## Mission

Every AI coding session starts from zero. The AI has no memory of what the
project does, how it's structured, what changed yesterday, or why decisions
were made. Dbasement eliminates this by maintaining a compact, structured,
persistent memory of the project — accessible instantly by any MCP-compatible
AI agent.

**Dbasement is not another AI. It never performs inference.** When analysis is
needed, it asks the connected AI. This keeps it lightweight, fast, and
focused.

## Features

- **Persistent memory** — Project knowledge survives between AI sessions
- **Structured knowledge** — Summary, architecture, API, database,
  dependencies, design decisions, glossary, known issues, TODOs, changelog
- **Incremental updates** — Never regenerate the entire memory; update only
  what changed
- **Smart retrieval** — AI retrieves only the section it needs, not the whole
  database
- **Git-aware** — Detects changes via `git diff`; no full rescans
- **File watching** — Optional poll-based watcher for non-git changes
- **Confidence system** — 0-100% confidence scoring; low-confidence updates
  marked or ignored
- **Small change filtering** — Ignores formatting, comments, whitespace,
  typo fixes
- **Multi-agent** — Works with Claude Code, Cursor, Codex CLI, Gemini CLI,
  Cline, Roo Code, Aider, and any MCP client
- **Zero infrastructure** — SQLite database, no vector DB, no Elasticsearch,
  no Redis, no cloud
- **Single binary** — ~11 MB, no runtime dependencies
- **Cross-platform** — Windows, macOS, Linux

## Installation

You don't need Go installed to use Dbasement. Pre-built binaries are available
for all major platforms.

### Download a Release

1. Go to the [Releases page](https://github.com/shs3131/dbasement/releases).
2. Download the archive for your operating system:
   - **Windows**: `dbasement-windows-amd64.zip`
   - **Linux (x86_64)**: `dbasement-linux-amd64.tar.gz`
   - **Linux (ARM64)**: `dbasement-linux-arm64.tar.gz`
   - **macOS (Intel)**: `dbasement-darwin-amd64.tar.gz`
   - **macOS (Apple Silicon)**: `dbasement-darwin-arm64.tar.gz`
3. Extract the archive and run the binary.

### Quick Start

**Windows (PowerShell):**

```powershell
# Extract
Expand-Archive -Path dbasement-windows-amd64.zip -DestinationPath .

# Run
.\dbasement.exe --project C:\path\to\your\project
```

**Linux / macOS:**

```bash
# Extract
tar xzf dbasement-linux-amd64.tar.gz

# Make executable and run
chmod +x dbasement
./dbasement --project /path/to/your/project
```

### Building from Source

Requires [Go 1.26+](https://go.dev/dl).

```bash
git clone https://github.com/shs3131/dbasement.git
cd dbasement
go build -o dbasement ./cmd/dbasement/
```

The binary is placed in the current directory. Move it anywhere in your `PATH`
for system-wide access.

### AI Client Configuration

Configure your AI client to use Dbasement as an MCP server:

```json
{
  "mcpServers": {
    "dbasement": {
      "command": "/absolute/path/to/dbasement",
      "args": ["--project", "/path/to/your/project"]
    }
  }
}
```

See [GET_STARTED.md](GET_STARTED.md) for per-client configuration examples and
detailed setup instructions.

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

## Performance

| Operation | Target |
|-----------|--------|
| Initialization | <30 seconds |
| Memory retrieval | <20 milliseconds |
| Memory update | <2 seconds |
| Idle RAM | ~10-15 MB |
| Idle CPU | ~0% |

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

## AI Client Configuration

### Claude Code

`~/.claude/settings.json`:

```json
{
  "mcpServers": {
    "dbasement": {
      "command": "dbasement",
      "args": ["--project", "/path/to/your/project"]
    }
  }
}
```

### Codex CLI

`.codex/mcp.json` in your project root:

```json
{
  "mcpServers": {
    "dbasement": {
      "command": "dbasement",
      "args": ["--project", "."]
    }
  }
}
```

### Cursor

Settings → MCP Servers → Add Server:

- **Name**: `dbasement`
- **Type**: `command`
- **Command**: `dbasement --project /path/to/your/project`

### Gemini CLI

`~/.config/gemini/mcp.json`:

```json
{
  "mcpServers": {
    "dbasement": {
      "command": "dbasement",
      "args": ["--project", "/path/to/your/project"]
    }
  }
}
```

### Cline / Roo Code

`~/.config/cline/mcp.json` or `~/.config/roo/mcp.json`:

```json
{
  "mcpServers": {
    "dbasement": {
      "command": "dbasement",
      "args": ["--project", "/path/to/your/project"]
    }
  }
}
```

### Aider

Aider supports MCP via the `--mcp` flag:

```bash
aider --mcp "dbasement --project /path/to/your/project"
```

Or via config file:

```yaml
# .aider.conf.yml
mcp: dbasement --project /path/to/your/project
```

### Generic MCP Client

Any client supporting the MCP stdio transport:

```json
{
  "mcpServers": {
    "dbasement": {
      "command": "dbasement",
      "args": ["--project", "/path/to/your/project"]
    }
  }
}
```

## Complete End-to-End Walkthrough

### Step 1: Install Dbasement

**Option A — Download a release binary (recommended, no Go needed):**

Download the latest binary from the
[Releases page](https://github.com/shs3131/dbasement/releases) and extract it.

**Option B — Build from source:**

```bash
git clone https://github.com/shs3131/dbasement.git
cd dbasement
go build -o dbasement ./cmd/dbasement/
```

### Step 2: Configure AI Client

Add Dbasement to your AI client's MCP config (see above).

### Step 3: First AI Session

The AI reads a project file and discovers Dbasement. It automatically
initializes the memory database and performs the first project analysis.

```
User: Read GET_STARTED.md

AI: (sees Dbasement is configured, initializes memory)
    Tool: initialize_project
      project_path: .
      summary: "This project is a web application..."
      architecture: "React frontend, Go backend, PostgreSQL database"

    Tool: update_memory
      section: features
      content: "User authentication, file upload..."
      confidence: 90

    Tool: update_memory
      section: api
      content: "POST /api/login, GET /api/files..."
      confidence: 90

    Tool: add_design_decision
      decision: "Using PostgreSQL for data persistence"
      reason: "Required for complex query support"

    Project memory initialized. I understand the project.
```

### Step 4: Memory is Stored

Dbasement writes to `.dbasement/memory.db`. This is a single SQLite file.

### Step 5: Second AI Session (Brand New Context)

The AI connects to Dbasement and retrieves knowledge instantly.

```
AI: Hello! I see Dbasement is set up. Let me retrieve the project context.

    Tool: get_project_summary
    -> "This project is a web application..."

    Tool: get_architecture
    -> "React frontend, Go backend, PostgreSQL database"

    Tool: get_recent_changes
    -> "No recent changes recorded."

    I'm ready. This is a web application with React frontend,
    Go backend, and PostgreSQL. What would you like to work on?
```

No file scanning. No token waste. The AI understands the project in <20ms.

### Step 6: User Modifies Code

```bash
# User adds a new API endpoint
echo 'func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
  // new endpoint
}' >> api/handlers.go
```

### Step 7: AI Detects and Updates Memory

```
User: What changed?

AI: Let me check.

    Tool: refresh_project
    -> "Changes detected (confidence: 85%): Added new API endpoint"

    Tool: update_memory
      section: api
      content: "...GET /api/users: Get user by ID..."
      confidence: 85
      changelog: "Added GET /api/users endpoint"
```

### Step 8: Third AI Session

```
AI: I remember this project.

    Tool: get_recent_changes
    -> "[2026-06-29 14:30] Added GET /api/users endpoint"

    The project recently added a user endpoint. The stack is
    React + Go + PostgreSQL. What would you like to do?
```

The AI is aware of the update without rescanning anything.

### How to Remove Dbasement

```bash
rm -rf .dbasement/
```

This completely removes all stored memory. The project is unaffected.

## Architecture

```
┌──────────────────────────────────────────────────┐
│                   AI Agent                        │
│   (Claude Code, Cursor, Codex CLI, Gemini, etc.) │
└──────────────┬───────────────────────────────────┘
               │  JSON-RPC 2.0 over stdio
               ▼
┌──────────────────────────────────────────────────┐
│              Dbasement MCP Server                 │
│                                                   │
│  ┌─────────┐  ┌──────────┐  ┌──────────────────┐ │
│  │  git    │  │ watcher  │  │    analyzer       │ │
│  │ client  │  │ (poll)   │  │ (relevance/conf)  │ │
│  └────┬────┘  └────┬─────┘  └────────┬─────────┘ │
│       │            │                  │           │
│       ▼            ▼                  ▼           │
│  ┌─────────────────────────────────────────────┐  │
│  │              memory.Manager                  │  │
│  └───────────────────┬─────────────────────────┘  │
│                      │                            │
│                      ▼                            │
│  ┌─────────────────────────────────────────────┐  │
│  │           storage.DB (SQLite)                │  │
│  │           .dbasement/memory.db               │  │
│  └─────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────┘
```

## FAQ

### Does Dbasement send my code anywhere?

No. Dbasement runs entirely locally. It never makes network requests, sends
telemetry, or communicates with anything other than your AI agent via stdio.

### Does Dbasement modify my code?

No. Dbasement is read-only with respect to your project files. It only writes
to `.dbasement/memory.db`. It never touches your source code.

### Can I use Dbasement without Git?

Yes. The file watcher detects changes via SHA-256 hash comparison. Git
integration is preferred but optional.

### How big does the database get?

The `.dbasement/memory.db` file typically stays under 1 MB for most projects.
Each memory section is a few hundred bytes to a few kilobytes.

### Can multiple AI agents use Dbasement simultaneously?

Yes. Dbasement supports multiple connections. SQLite handles read concurrency
natively with WAL mode. Write operations are serialized.

### Will Dbasement slow down my AI?

No. Memory retrieval takes <20ms. Memory updates take <2s. The AI spends less
time understanding the project, not more.

### How do I update the memory?

The AI updates memory automatically via `update_memory`. You can also manually
trigger a refresh with `refresh_project`.

### Is the memory shared between developers?

By default, each developer has their own `.dbasement/` directory. You can share
memory by committing the database (not recommended) or using the export feature
(future).

## Project Status

Dbasement is in **active development**. The core functionality is stable and
usable. See [ROADMAP.md](ROADMAP.md) for planned features.

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for
guidelines.

## License

MIT License. See [LICENSE](LICENSE) for details.

# Getting Started with Dbasement

Dbasement gives your AI coding agent persistent project memory. Once set up,
every AI session remembers what the project is about, how it works, what
changed, and why decisions were made.

## Quick Start (no Go required)

**1. Download** the binary from the [latest release](https://github.com/shs3131/dbasement/releases/latest).

**2. Extract** into your project root:

```bash
# Linux / macOS
tar xzf dbasement-*.tar.gz && chmod +x dbasement

# Windows (PowerShell)
Expand-Archive -Path dbasement-windows-amd64.zip -DestinationPath .
```

**3. Tell your AI** to use it ([README.md](README.md#ai-client-configuration) has configs for all clients).

**4. Done.** First AI session auto-initializes project memory.

## Quick Start

### 1. Install

```bash
go install github.com/shs3131/dbasement/cmd/dbasement@latest
```

Or download a pre-built binary from the [Releases](https://github.com/shs3131/dbasement/releases)
page. See [INSTALL.md](INSTALL.md) for OS-specific instructions.

### 2. Verify

```bash
dbasement --help
```

You should see:

```
Usage of dbasement:
  -project string
        Path to the project root (default: current directory)
```

### 3. Configure Your AI Client

Add Dbasement as an MCP server in your AI client's configuration.

#### Claude Code

Add to `~/.claude/settings.json` or your project's `.claude/settings.json`:

```json
{
  "mcpServers": {
    "dbasement": {
      "command": "dbasement",
      "args": ["--project", "/path/to/your/project"],
      "env": {}
    }
  }
}
```

#### Cursor

In Cursor settings, add a new MCP server:

- **Name**: `dbasement`
- **Type**: `command`
- **Command**: `dbasement --project /path/to/your/project`

#### Cline / Roo Code

Add to your MCP settings file (usually `~/.config/cline/mcp.json` or
`~/.config/roo/mcp.json`):

```json
{
  "mcpServers": {
    "dbasement": {
      "command": "dbasement",
      "args": ["--project", "/path/to/your/project"],
      "env": {}
    }
  }
}
```

#### Codex CLI

Codex CLI automatically discovers MCP servers. Place a config file at
`.codex/mcp.json` in your project:

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

#### Generic MCP Client

Any MCP-compatible client can connect over stdio:

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

### 4. Initialize Project Memory

In your first AI session, run the `initialize_project` tool:

```
Tool: initialize_project
Arguments:
  project_path: .
  summary: A brief 200-400 word description of your project.
  architecture: Frontend/backend breakdown.
```

Or, more commonly, **the AI will detect Dbasement automatically** when it
reads the project files and initialize itself.

### 5. Start Using Memory

Once initialized, you can retrieve knowledge instantly:

```
Tool: get_project_summary
-> Returns: "Your project's 200-400 word summary"
```

## How It Works

### First Session

1. AI discovers Dbasement in the project (e.g., by reading a file)
2. AI calls `initialize_project` with project summary and architecture
3. Dbasement stores this in `.dbasement/memory.db`
4. AI calls `update_memory` for each knowledge section

### Subsequent Sessions

1. AI connects to Dbasement (already initialized)
2. AI retrieves only what it needs:
   - `get_project_summary` to understand the project
   - `get_architecture` to understand structure
   - `get_api` to understand endpoints
   - etc.
3. AI never needs to rescan the repository
4. When AI makes changes, it calls `update_memory` to keep knowledge current

### After Git Changes

1. AI or user runs `refresh_project`
2. Dbasement checks `git diff HEAD`
3. If changes are meaningful, AI is notified and can update memory

## Available Tools

| Tool | Description |
|------|-------------|
| `initialize_project` | Initialize memory with project summary and architecture |
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
| `update_memory` | Update a memory section |
| `add_design_decision` | Record a design decision |
| `add_todo` | Add a TODO item |
| `add_known_issue` | Add a known issue |
| `refresh_project` | Check for meaningful changes |
| `resolve_known_issue` | Mark issue resolved |
| `mark_todo_done` | Mark TODO complete |

## Confidence System

Every memory update has a confidence score:

| Score | Behavior |
|-------|----------|
| >= 95 | Applied automatically |
| 85-94 | Applied automatically |
| 70-84 | Applied, marked AI-inferred |
| < 70 | Ignored |

## Removing Dbasement

To completely remove Dbasement from a project:

```bash
rm -rf .dbasement/
```

This deletes the SQLite database containing all stored memory. The project
itself is unaffected.

## Example Session

```
User: Read GET_STARTED.md

AI: (discovers Dbasement, initializes automatically)

AI: I see Dbasement is set up for this project.
    Let me check what we know about it.

Tool: get_project_summary
Tool: get_architecture
Tool: get_recent_changes
Tool: get_todo

AI: This project is a web application with React frontend
    and Go backend. Recent changes include adding JWT auth.
    There's an open TODO to write API tests.
    How can I help?
```

The AI understands the project instantly, without scanning hundreds of files.

# Getting Started with Dbasement

## Install (no Go required)

**Download the binary from GitHub Releases and extract into your project root.**

```bash
# Linux / macOS
curl -sL https://github.com/shs3131/dbasement/releases/latest/download/dbasement-linux-amd64.tar.gz | tar xz

# macOS (Apple Silicon)
curl -sL https://github.com/shs3131/dbasement/releases/latest/download/dbasement-darwin-arm64.tar.gz | tar xz

# Windows (PowerShell)
curl.exe -LO https://github.com/shs3131/dbasement/releases/latest/download/dbasement-windows-amd64.zip
Expand-Archive -Path dbasement-windows-amd64.zip -DestinationPath . -Force
```

Or use the automatic install scripts:

```bash
# Linux / macOS
bash scripts/start.sh .

# Windows (PowerShell)
pwsh -ExecutionPolicy Bypass -File scripts/start.ps1 --project .
```

The scripts auto-download the binary from GitHub Releases if it's missing.

## Configure AI Client

Add Dbasement as an MCP server. See [INSTALL.md](INSTALL.md) or [MCP_CONFIGS.md](MCP_CONFIGS.md) for your client.

## First AI Session

The AI detects Dbasement and initializes project memory automatically.

```
1. AI reads project files
2. AI calls initialize_project with summary and architecture
3. AI calls update_memory for sections it discovers
4. Memory persists in .dbasement/memory.db
```

## How It Works

### Session 1 (initialization)

AI reads the project and stores structured knowledge in `.dbasement/memory.db`:
- Project summary, architecture, features, API, database
- Dependencies, design decisions, glossary
- Known issues, TODOs, changelog

### Session 2+ (retrieval)

AI retrieves knowledge in 20ms instead of scanning files:
- `get_project_summary` to re-establish context
- Other `get_*` tools for specific sections as needed

### After code changes

AI runs `refresh_project` → detects meaningful changes → updates memory.

## Removing Dbasement

```bash
rm -rf .dbasement/
```

## Tools

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

## Confidence System

| Score | Behavior |
|-------|----------|
| >= 95 | Applied automatically |
| 85-94 | Applied automatically |
| 70-84 | Applied, marked AI-inferred |
| < 70 | Ignored |

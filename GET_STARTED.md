# Getting Started

## 1. Download the binary

See [INSTALL.md](INSTALL.md) — no Go needed.

## 2. Place it in your project root

```bash
./dbasement --help
```

## 3. Add MCP config for your AI client

See [README.md](README.md#ai-client-configuration) for all clients.

Quick examples:

**VS Code** — `.vscode/mcp.json` is already in the project. Just open the folder.

**Claude Code**:
```bash
claude --mpc "./dbasement --project ."
```

**Cline / Roo Code** — add to `mcp.json`:
```json
{
  "mcpServers": {
    "dbasement": {
      "command": "C:\\path\\to\\dbasement.exe",
      "args": ["--project", "C:\\path\\to\\project"]
    }
  }
}
```

## 4. First AI session

The AI detects Dbasement and runs `initialize_project` automatically. You get persistent project memory from that point on.

## Available Tools

| Tool | What it does |
|------|-------------|
| `get_project_summary` | 200-400 word project description |
| `get_architecture` | Architecture breakdown |
| `get_features` | Feature list |
| `get_api` | API documentation |
| `get_database` | Database schema |
| `get_dependencies` | Dependency documentation |
| `get_recent_changes` | Recent changelog |
| `get_known_issues` | Unresolved issues |
| `get_todo` | TODO items |
| `get_design_decisions` | Decision history |
| `get_glossary` | Project terminology |
| `search_memory` | Full-text search |
| `update_memory` | Update a memory section |
| `initialize_project` | Initialize project memory |
| `add_design_decision` | Record a decision |
| `add_todo` | Add a TODO |
| `add_known_issue` | Add an issue |
| `refresh_project` | Check for git changes |
| `resolve_known_issue` | Resolve an issue |
| `mark_todo_done` | Complete a TODO |

## Remove

```bash
rm -rf .dbasement/
```

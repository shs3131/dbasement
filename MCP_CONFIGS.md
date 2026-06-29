# MCP Configuration for All AI Clients

This file contains MCP server configuration for every supported AI client.
Copy-paste the relevant section into your client.

---

## Cross-Platform (Recommended)

This uses the auto-download scripts. The binary is downloaded from GitHub Releases
on first run — no Go, no manual download.

### Unix (Linux / macOS)

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

### Windows

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

---

## VS Code

Uses `.vscode/mcp.json` — auto-discovered when you open the project folder.
No manual configuration needed.

---

## Claude Code

### Global config (`~/.claude/settings.json`)

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

Or Windows:

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

### CLI flag

```bash
claude --mcp "bash scripts/start.sh ."
```

### With binary in PATH

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

---

## Cursor

Settings → MCP → Add Server:

| Field | Value |
|-------|-------|
| Name | `dbasement` |
| Type | `command` |
| Command | `bash scripts/start.sh .` |

Or with binary in PATH: `dbasement --project .`

---

## Cline

`~/.config/cline/mcp.json`:

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

---

## Roo Code

`~/.config/roo/mcp.json`:

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

---

## Codex CLI

`.codex/mcp.json` (project-level):

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

---

## Gemini CLI

`~/.config/gemini/mcp.json`:

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

---

## Aider

### CLI flag

```bash
aider --mcp "bash scripts/start.sh ."
```

### Config file (`.aider.conf.yml`)

```yaml
mcp: bash scripts/start.sh .
```

---

## Generic MCP Client

Any MCP-compatible client can connect over stdio:

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

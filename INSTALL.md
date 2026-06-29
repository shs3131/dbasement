# Installation Guide

**Prerequisite:** Git must be installed (for change detection). No other dependencies.

---

## Recommended: Pre-built Binary

Download from [GitHub Releases](https://github.com/shs3131/dbasement/releases).

### macOS (Apple Silicon)

```bash
curl -sL https://github.com/shs3131/dbasement/releases/latest/download/dbasement-darwin-arm64.tar.gz | tar xz
sudo mv dbasement /usr/local/bin/
```

### macOS (Intel)

```bash
curl -sL https://github.com/shs3131/dbasement/releases/latest/download/dbasement-darwin-amd64.tar.gz | tar xz
sudo mv dbasement /usr/local/bin/
```

### Linux (x86_64)

```bash
curl -sL https://github.com/shs3131/dbasement/releases/latest/download/dbasement-linux-amd64.tar.gz | tar xz
sudo mv dbasement /usr/local/bin/
```

### Linux (ARM64)

```bash
curl -sL https://github.com/shs3131/dbasement/releases/latest/download/dbasement-linux-arm64.tar.gz | tar xz
sudo mv dbasement /usr/local/bin/
```

### Windows (PowerShell)

```powershell
curl.exe -LO https://github.com/shs3131/dbasement/releases/latest/download/dbasement-windows-amd64.zip
Expand-Archive -Path dbasement-windows-amd64.zip -DestinationPath . -Force
Move-Item .\dbasement.exe C:\Windows\System32\
```

### Verify

```bash
dbasement --help
```

---

## Alternative: Install with `go install` (requires Go 1.26+)

```bash
go install github.com/shs3131/dbasement/cmd/dbasement@latest
```

---

## MCP Configuration

After installing the binary, configure your AI client:

### VS Code

Project has `.vscode/mcp.json` — auto-discovered on open.

### Claude Code

```json
// ~/.claude/settings.json or .claude/settings.json
{
  "mcpServers": {
    "dbasement": {
      "command": "dbasement",
      "args": ["--project", "/path/to/your/project"]
    }
  }
}
```

### Cursor

Settings → MCP → Add Server:
- Name: `dbasement`
- Type: `command`
- Command: `dbasement --project /path/to/your/project`

### Cline / Roo Code

```json
// ~/.config/cline/mcp.json or ~/.config/roo/mcp.json
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

```json
// .codex/mcp.json
{
  "mcpServers": {
    "dbasement": {
      "command": "dbasement",
      "args": ["--project", "."]
    }
  }
}
```

### Gemini CLI

```json
// ~/.config/gemini/mcp.json
{
  "mcpServers": {
    "dbasement": {
      "command": "dbasement",
      "args": ["--project", "/path/to/your/project"]
    }
  }
}
```

### Generic

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

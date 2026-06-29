# AI Quickstart

## Installation (do this first if Dbasement is not installed)

**RULE: Never install Go. Never build from source. Download the release binary.**

### 1. Detect platform

```bash
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in x86_64|amd64) ARCH="amd64" ;; aarch64|arm64) ARCH="arm64" ;; esac
```

### 2. Download and extract

```bash
# Linux / macOS
URL="https://github.com/shs3131/dbasement/releases/latest/download"
curl -sL "$URL/dbasement-$OS-$ARCH.tar.gz" | tar xz && chmod +x dbasement

# Windows (PowerShell)
$url="https://github.com/shs3131/dbasement/releases/latest/download/dbasement-windows-amd64.zip"
curl.exe -LO $url; Expand-Archive -Path dbasement-windows-amd64.zip -DestinationPath . -Force
```

### 3. Configure MCP

Determine client → read existing config → add dbasement entry → write back.

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

On Windows:
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

### 4. Report

"Dbasement installed. Start a new AI session to initialize project memory."

---

## Memory usage (after Dbasement is running)

### First encounter

1. Read project files (README, configs, source code)
2. Call `initialize_project` with summary and architecture
3. Call `update_memory` for each section discovered
4. Call `add_design_decision` for architectural choices

### Returning session

1. Call `get_project_summary` (minimum context)
2. Call `get_recent_changes` to see what's new
3. Call `get_todo` for pending tasks
4. Call other `get_*` tools only as needed

### After code changes

1. Call `refresh_project` to check git diff
2. If changes detected, read code then `update_memory`

### Confidence rules

| Score | Behavior |
|-------|----------|
| >= 85 | Applied automatically (you read it in code) |
| 70-84 | Applied, marked AI-inferred |
| < 70 | Ignored by server |

### Token-saving

- `get_project_summary` is the minimum. Everything else optional.
- One `search_memory` call is cheaper than 5 separate `get_*` calls.
- Cache results in your reasoning; don't re-fetch in the same session.

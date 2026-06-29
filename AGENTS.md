# AI Agent Guide for Dbasement

**Read this first if you are an AI agent trying to install or use Dbasement.**

---

## Section 0: AI Installation Rules (MANDATORY)

### Golden Rules

1. **NEVER install Go.** NEVER run `go build`, `go install`, `go get`, or any Go command.
2. **NEVER build from source** unless the user explicitly asks for it.
3. **ALWAYS download from GitHub Releases** — the correct binary for the user's platform.
4. **If the user asks why:** "Dbasement provides pre-built binaries. Building from source requires Go and is only needed for contributors."
5. **If the user insists on building from source:** Say "Building from source requires Go 1.26+. Are you sure you want to install Go for this?" Then follow instructions in [DEVELOPMENT.md](DEVELOPMENT.md).

### Installation Workflow

```
1. Detect platform: Identify OS and architecture
2. Download: Fetch the correct archive from GitHub Releases
3. Extract: Extract the binary into the project root
4. Configure MCP: Add dbasement to the client's MCP configuration
5. Preserve: Keep existing MCP servers (merge, don't overwrite)
6. Report: Tell the user Dbasement is installed and to start a new session
```

### Platform Detection

```python
# Python (use for generic MCP clients)
import platform, sys
system = platform.system().lower()  # "linux", "darwin", "windows"
machine = platform.machine().lower()  # "x86_64", "arm64", "amd64"

# Linux
if system == "linux" and machine in ("x86_64", "amd64"):
    asset = "dbasement-linux-amd64.tar.gz"
elif system == "linux" and machine in ("aarch64", "arm64"):
    asset = "dbasement-linux-arm64.tar.gz"
# macOS
elif system == "darwin" and machine in ("x86_64", "amd64"):
    asset = "dbasement-darwin-amd64.tar.gz"
elif system == "darwin" and machine in ("aarch64", "arm64"):
    asset = "dbasement-darwin-arm64.tar.gz"
# Windows
elif system == "windows" and machine in ("x86_64", "amd64"):
    asset = "dbasement-windows-amd64.zip"
else:
    raise Exception("Unsupported platform")
```

```bash
# Bash (for shell commands)
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in x86_64|amd64) ARCH="amd64" ;; aarch64|arm64) ARCH="arm64" ;; esac
case "$OS" in
  linux)  ASSET="dbasement-linux-$ARCH.tar.gz" ;;
  darwin) ASSET="dbasement-darwin-$ARCH.tar.gz" ;;
  windows) ASSET="dbasement-windows-$ARCH.zip" ;;
esac
```

```powershell
# PowerShell (for Windows)
$arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "amd64" }
$asset = "dbasement-windows-$arch.zip"
```

### Download URLs

| Platform | Download URL (latest) |
|----------|----------------------|
| Linux (amd64) | `https://github.com/shs3131/dbasement/releases/latest/download/dbasement-linux-amd64.tar.gz` |
| Linux (arm64) | `https://github.com/shs3131/dbasement/releases/latest/download/dbasement-linux-arm64.tar.gz` |
| macOS (amd64) | `https://github.com/shs3131/dbasement/releases/latest/download/dbasement-darwin-amd64.tar.gz` |
| macOS (arm64) | `https://github.com/shs3131/dbasement/releases/latest/download/dbasement-darwin-arm64.tar.gz` |
| Windows (amd64) | `https://github.com/shs3131/dbasement/releases/latest/download/dbasement-windows-amd64.zip` |

Version-specific: replace `latest` with tag like `v1.0.0`.

### Extract Commands

```bash
# tar.gz (Linux/macOS)
tar xzf dbasement-*.tar.gz && chmod +x dbasement

# zip (Windows)
Expand-Archive -Path dbasement-windows-amd64.zip -DestinationPath . -Force
```

### Detecting the AI Client

Determine which AI client is running to configure it correctly:

- **Claude Code**: Check `~/.claude/settings.json` or `CLAUDE_CODE` env var
- **Cursor**: Look for `~/.cursor/mcp.json` or Cursor settings directory
- **Cline**: Check `~/.config/cline/mcp.json`
- **Roo Code**: Check `~/.config/roo/mcp.json`
- **Codex CLI**: Look for `.codex/mcp.json` in project root
- **Gemini CLI**: Check `~/.config/gemini/mcp.json`
- **Aider**: Look for `.aider.conf.yml` or `~/.aider.conf.yml`
- **VS Code**: Look for `.vscode/mcp.json` in project root
- **Generic**: Look for `mcp.json` or `mcpServers` in config files

### Configuring MCP

**Add Dbasement to existing config, preserving other servers.**

Read existing config → add dbasement entry → write back.

Config entry (Unix):

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

Config entry (Windows):

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

### Reporting Success

Tell the user:

> "Dbasement has been installed in your project. Start a new AI session and the agent will automatically initialize and remember your project context."

---

## Memory Workflows

### Workflow A: First-time Project Initialization

Use this when a project has never been initialized with Dbasement.

```
User: Look at this project and tell me what it does.

Agent flow:
  1. READ key project files (README, package.json, docs, source)
  2. Tool: initialize_project
       project_path: "."
       summary: "200-400 word summary of the project..."
       architecture: "Frontend/backend/service breakdown..."
  3. Tool: update_memory
       section: "features"
       content: "List of features discovered..."
       confidence: 85
  4. Tool: update_memory
       section: "api"
       content: "API endpoints discovered..."
       confidence: 80
  5. Tool: update_memory
       section: "dependencies"
       content: "Key dependencies and why they exist..."
       confidence: 90
  6. Tool: add_design_decision
       decision: "Using X framework"
       reason: "Chosen because..."
```

```mermaid
flowchart TD
    A[Agent starts] --> B[Read project files]
    B --> C{Is .dbasement/memory.db present?}
    C -->|No| D[Call initialize_project]
    D --> E[Call update_memory for each section]
    E --> F[Call add_design_decision if needed]
    F --> G[Memory stored. Project understood.]
```

**Rules:**
- Call `initialize_project` ONCE, and only on first encounter.
- After initializing, call `update_memory` for each section you discover.
- Set `confidence >= 85` for facts you are certain about (from reading code).
- Set `confidence 70-84` for inferred information (heuristic guesses).
- Never call `initialize_project` if the project already has `.dbasement/`.

### Workflow B: Existing Project Session

Use this at the start of every session to re-establish context.

```
Agent flow (start of session):
  1. Tool: get_project_summary
     → Returns: "This project is a..."
  2. Tool: get_architecture
     → Returns: "React frontend, Go backend..."
  3. Tool: get_recent_changes
     → Returns: "Recent changes: ..."
  4. Tool: get_todo
     → Returns: "Pending tasks: ..."
  5. Agent now understands the project (20ms elapsed).
```

```mermaid
flowchart TD
    A[Agent starts new session] --> B[Call get_project_summary]
    B --> C[Call get_architecture]
    C --> D[Call get_recent_changes]
    D --> E[Call get_todo]
    E --> F[Agent has full context. Ready to work.]
```

**Token-saving rules:**
- Always start with `get_project_summary`. This is the minimum context.
- Call additional get_* tools ONLY when the task requires that specific section.
- Example: if working on the API, also call `get_api`. Don't call every tool.
- `get_known_issues` is only needed when fixing bugs.
- `get_dependencies` is only needed when managing dependencies.

### Workflow C: Updating Memory After Code Changes

Use this after making changes or when the user reports changes.

```
User: I just added a new API endpoint.

Agent flow:
  1. Tool: refresh_project
     → Checks git diff for meaningful changes
     → Returns: "Changes detected (confidence: 85%): Added GET /api/users"
  2. Tool: update_memory
       section: "api"
       content: "New endpoint: GET /api/users..."
       confidence: 95
       changelog: "Added GET /api/users endpoint"
  3. If design decision involved:
     Tool: add_design_decision
       decision: "Using RESTful naming for user endpoints"
       reason: "Consistency with existing API patterns"
```

```mermaid
flowchart TD
    A[User says 'I made changes'] --> B[Call refresh_project]
    B --> C{Meaningful changes?}
    C -->|No| D[Inform user: no significant changes]
    C -->|Yes| E[Call update_memory for affected sections]
    E --> F{New design decisions?}
    F -->|Yes| G[Call add_design_decision]
    F -->|No| H[Done]
```

**Rules:**
- `refresh_project` checks git diff. Safe to call multiple times.
- `update_memory` with `changelog` parameter will update both memory and log.
- Always set confidence based on how certain you are about the new information.
- If the user added code, read the actual changes before calling `update_memory`.

## Tool Reference

### When to call each tool

| Tool | Call When | Do NOT Call When |
|------|-----------|-----------------|
| `initialize_project` | First encounter with the project | `.dbasement/` already exists |
| `get_project_summary` | Start of every session | You need detail on a specific section |
| `get_architecture` | You need to understand structure | You only need the summary |
| `get_features` | Planning features | Doing routine maintenance |
| `get_api` | Working on API code | Working on frontend-only |
| `get_database` | Making schema changes | Working on UI |
| `get_dependencies` | Adding/removing deps | Writing implementation code |
| `get_recent_changes` | Asked "what changed?" | Project is brand new |
| `get_known_issues` | Fixing bugs | Initial setup |
| `get_todo` | Planning next task | Already know what to do |
| `get_design_decisions` | Need to understand "why" | Routine work |
| `get_glossary` | Encounter unknown terms | Understand domain well |
| `search_memory` | Don't know which section | Know exactly which section |
| `update_memory` | Learned new information | Confidence < 70 |
| `add_design_decision` | Made architectural choice | Trivial implementation |
| `add_todo` | Discovered pending work | Already tracked in code |
| `add_known_issue` | Found a real bug | Minor/speculative issue |
| `refresh_project` | User says "I changed things" | About to initialize |
| `resolve_known_issue` | Verified a fix works | Before verifying |
| `mark_todo_done` | Actually completed task | Task isn't done |

### First-in-workflow tools

- `get_project_summary` — FIRST tool in every session start
- `initialize_project` — FIRST (and only) tool in project setup
- `refresh_project` — FIRST tool when checking for changes

## Complete Session Examples

### Example 1: Full Session Lifecycle

```
Session 1 (first encounter):
  Agent reads README.md, package.json, src/main.go
  → Tool: initialize_project(project_path=".", summary="CLI tool...")
  → Tool: update_memory(section="features", content="...", confidence=90)
  → Tool: add_design_decision(decision="Using Cobra", reason="CLI framework")
  → Memory initialized.

Session 2 (new context, next day):
  → Tool: get_project_summary → "CLI tool for..."
  → Tool: get_architecture → "Go CLI with Cobra..."
  → Tool: get_todo → "Add --verbose flag"
  → Agent is ready. No file scanning needed.

User: "Can you add --verbose?"
  Agent reads the CLI code, adds the flag.
  → Tool: update_memory(section="features",
       content="... --verbose flag for detailed output",
       confidence=95, changelog="Added --verbose flag")
  → Tool: mark_todo_done(id=1)
  → Done.

Session 3 (new context, next week):
  → Tool: get_project_summary → "CLI tool for..."
  → Tool: get_recent_changes → "Added --verbose flag"
  → Agent knows what happened last session.
```

## Confidence Score Guide

| Score | Meaning | When to use |
|-------|---------|-------------|
| 95-100 | Certain | Read it directly from the source code |
| 85-94 | Very confident | Clear from documentation or config files |
| 75-84 | Reasonable inference | Deduced from file structure or patterns |
| 70-74 | Weak inference | Heuristic guess, use with caution |
| < 70 | Don't send | Server ignores these automatically |

## Best Practices for Token Minimization

1. **Cache context in your thinking**: After calling `get_project_summary`, you don't need to recall it unless the project changes.
2. **One tool call per section**: Don't call `get_api` multiple times in a session — the database doesn't change unless you update it.
3. **Use `search_memory` as a fallback**: If you're not sure which section contains the answer, one `search_memory` call is cheaper than 5 separate `get_*` calls.
4. **Skip known-empty sections**: Don't call `get_glossary` if you already understand the domain.
5. **Batch `update_memory` calls**: When initializing a project, call `update_memory` for each section in sequence.

## Common Mistakes

| Mistake | Why It's Wrong | Correct Approach |
|---------|---------------|------------------|
| Calling every `get_*` at start | Wastes tokens | Start with `get_project_summary`, add others as needed |
| Calling `initialize_project` twice | Server rejects | Check `IsInitialized` first |
| Setting confidence too low | Gets ignored | Read the code before setting confidence |
| Not calling `refresh_project` | Stale memory | Always check after user edits |
| Calling `update_memory` before reading changes | Guessing facts | Read the actual code first |
| Using `search_memory` when you know the section | Inefficient | Use the specific get_* tool |
| Installing Go to build Dbasement | Unnecessary | Download the pre-built binary from GitHub Releases |
| Building from source automatically | Violates project policy | Always use releases unless user explicitly requests source |

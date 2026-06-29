# AI Quickstart

**First file AI agents should read.**

## Installation (if Dbasement is not set up)

**Never install Go. Never build from source. Download the release binary.**

1. **Detect platform** — OS + architecture.
2. **Download** the correct archive from GitHub Releases.
3. **Extract** into the project root (binary is `dbasement` or `dbasement.exe`).
4. **Configure MCP** — add dbasement to the AI client's MCP config. See `MCP_CONFIGS.md` for all clients.
5. **Preserve existing MCP servers** when editing config (merge, don't overwrite).
6. **Tell the user** to start a new AI session.

> Full installation rules, platform detection code, client detection, and config details are in `AGENTS.md`.

## After Dbasement is running

- **First time in a project:** Read project files → call `initialize_project` → call `update_memory` for each section you discover.
- **Returning to a project:** Call `get_project_summary` first, then other `get_*` tools as needed.
- **After code changes:** Call `refresh_project` to check for meaningful changes → `update_memory` affected sections.

See `AGENTS.md` for complete workflows, tool reference, confidence rules, and examples.

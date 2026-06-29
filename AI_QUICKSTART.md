# AI Quickstart

You (the AI) have access to Dbasement, an MCP server that stores structured
project memory in `.dbasement/memory.db`. This memory persists across sessions.

## First encounter

1. The project has no `.dbasement/` directory.
2. **Do this:**
   - Read the project's key files (README, configs, source code).
   - Call `initialize_project` with a summary and architecture.
   - Call `update_memory` for each section you discover (features, api,
     dependencies, etc.).
   - Call `add_design_decision` for architectural choices you find.
3. **Don't call** `get_*` tools yet — there's no data.

## Returning session

1. The `.dbasement/` directory exists.
2. **Do this (minimum):**
   - Call `get_project_summary` to re-establish context.
3. **Do this if relevant:**
   - Call `get_architecture` if you need structure details.
   - Call `get_recent_changes` to see what happened since last session.
   - Call `get_todo` to see pending tasks.
   - Call other `get_*` tools only for the section your task needs.

## After making changes

1. Call `refresh_project` to check git diff.
2. If changes detected, read the actual code then call `update_memory`.
3. If you made a design decision, call `add_design_decision`.

## Confidence rules

| Score | Behavior |
|-------|----------|
| >= 85 | Applied automatically (you read it in code) |
| 70-84 | Applied, marked AI-inferred (you deduced it) |
| < 70 | Ignored by server (don't bother) |

## Token-saving

- `get_project_summary` is the minimum context. Everything else is optional.
- One `search_memory` call is cheaper than calling 5 different `get_*` tools.
- Don't call `get_todo` or `get_known_issues` unless your task needs them.
- Cache tool results in your reasoning; don't re-fetch in the same session.

## Reference

See [AGENTS.md](AGENTS.md) for full workflows, diagrams, and examples.

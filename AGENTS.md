# AGENTS

Go CLI + Charm TUI: given the current shell line and cwd, suggest the **next token** in a fuzzy list (`Label`, `Detail`, `Value`). Shell keybinding is the main integration; any terminal.

## v1

- Path provider for path-taking commands (`mv`, `cp`, `cd`, `rm`, `cat`, …)
- Package scripts provider for `npm` / `pnpm` / `yarn` / `bun` run-style prefixes; parse nearest `package.json` in Go (no Node)
- Generic TUI + CLI wiring + zsh/bash keybinding

## Do not

- Build a terminal emulator or Warp-style overlay
- Put provider / filesystem / manifest logic in the TUI
- Put rendering in provider packages
- Require Node (or other runtimes) to read project files
- Aim to complete every CLI in v1

## Layout

```text
cmd/tabui/              # main — wire only
internal/path/          # path provider
internal/package_scripts/
internal/ui/            # Bubble Tea / Lip Gloss only
internal/fuzzy/         # example shared helper — only if 2+ callers
```

Colocate feature-local types and helpers under `internal/<feature>/`. Promote to a named `internal/<name>` package only on a second caller — no dump `utils`.

## Providers

Under `internal/<feature>/`:

1. **Match** — does this provider apply to the cmdline context?
2. **Suggest** — return `[]Suggestion`; never render

| Field | Role |
|-------|------|
| `Label` | Primary (name, filename, subcommand) |
| `Detail` | Secondary (path, script body, description) |
| `Value` | Token inserted into the line |

TUI stays generic: only render and return the chosen `Value`.

## Stack

Go · Bubble Tea / Lip Gloss · fuzzy filter · `encoding/json` for manifests.

## Working style

Small focused changes. Test match/suggest without the TUI. Skill: `shell-integration` when editing shell bindings.

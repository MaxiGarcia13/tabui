# AGENTS

Go CLI + Charm TUI: given the current shell line and cwd, suggest the **next token** in a fuzzy list (`Label`, `Detail`, `Value`). Shell keybinding is the main integration; any terminal.

## Layers

1. **Package scripts** — `npm` / `pnpm` / `yarn` / `bun` **run** / `run-script` slots; nearest `package.json` in Go (no Node)
2. **Path** — path-taking commands (`mv`, `cp`, `cd`, `rm`, `cat`, …) plus curated flags
3. **CLI help** — any other `LookPath` command: `cmd [sub…] --help`, sanitize/parse subcommands + flags, optional SQLite cache; async load in the TUI on cache miss

Specialized layers win first; help is the fallback. Builtins / tools without useful `--help` stay on curated path flags.

## Do not

- Put provider / filesystem / manifest / help-exec logic in the TUI
- Put rendering in provider packages
- Require Node (or other runtimes) to read project files
- Aim to complete every CLI dialect of `--help` on day one

## Layout

```text
cmd/tabui/              # main — wire only
internal/path/          # path provider
internal/package_scripts/
internal/clihelp/       # generic --help fallback + cache
internal/shellinit/     # embedded zsh / bash widgets (`tabui init`)
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

TUI stays generic: only render and return the chosen `Value` (sync `Run` or `RunLoading`).

## Stack

Go · Bubble Tea / Lip Gloss · fuzzy filter · `encoding/json` for manifests · optional SQLite (`modernc.org/sqlite`) for help cache.

## Working style

Small focused changes. Test match/suggest without the TUI. Skill: `shell-integration` when editing shell bindings.

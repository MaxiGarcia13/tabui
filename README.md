# tabui

A Go CLI + TUI that suggests the **next thing to type** while you compose a shell command — paths, scripts, subcommands, or flags, depending on context.

Works in any terminal (Cursor, VS Code, iTerm, …) as a fuzzy picker with **label** + **detail**. Not a custom terminal overlay.

## How it feels

1. You are building a command (shell keybinding, or invoke the tool with the current line).
2. tabui reads **context**: current command line + cwd (+ optional project files).
3. A fuzzy TUI opens with suggestions for **the next argument**.
4. Each row shows **label** (filename, script name, subcommand) and **detail** (path, script body, description) dimmed on the right.
5. You pick one → that token is inserted (or the command is completed / run, depending on mode).

| You are typing | tabui suggests |
|----------------|----------------|
| `mv ` / `cp ` / `rm ` / `cd ` | files and folders (paths) |
| `pnpm ` / `pnpm run ` / `npm run ` | scripts from nearest `package.json`, plus tool subcommands when relevant |
| other known commands | whatever that command needs next (flags, subcommands, paths, …) |

Examples:

- `mv ` → files/folders in the current (or relevant) directory
- `pnpm ` / `npm run ` → package scripts (`dev`, `build`, …) with command body as detail

## Architecture

Provider model — easy to extend without rewriting the TUI:

- **Path provider** — files/folders for commands that take paths (`mv`, `cp`, `cd`, `rm`, `cat`, …)
- **Package scripts provider** — `npm run` / `pnpm` / `pnpm run` / `yarn` / `bun run`; nearest `package.json` parsed in Go (no Node at runtime)
- **Subcommand provider** (later) — known flags/subcommands for selected tools

The TUI only renders `label`, `detail`, and the value to insert.

Layout (standard Go):

```text
cmd/tabui/           # binary entry
internal/<feature>/  # providers, matching, domain types
internal/ui/         # Charm TUI
internal/<name>/     # shared helpers only when reused (named packages)
```

## Integration

Primary target: **shell-friendly**, terminal-agnostic.

- Go binary accepts the current command line / cursor word (args or stdin).
- Bind a key in zsh/bash (e.g. Ctrl+Space) that:
  1. grabs the current line
  2. calls tabui
  3. inserts the selected suggestion back into the line

Also supports standalone modes where useful (e.g. script picker only).

Does **not** require a specific terminal app (not Warp-only, not VS Code-extension-only).

### Shell keybinding (sketch)

```zsh
# ~/.zshrc — example; finalize once the binary CLI flags are stable
tabui-complete() {
  local selected
  selected="$(tabui --line="$BUFFER" --cursor="$CURSOR")" || return
  [[ -n "$selected" ]] || return
  LBUFFER="${LBUFFER}${selected}"
}
zle -N tabui-complete
bindkey '^@' tabui-complete  # Ctrl+Space (may vary by terminal)
```

## Tech

- Go
- Charm TUI (Bubble Tea / Lip Gloss)
- Fuzzy filter
- JSON parsed in Go for `package.json`

## Development

Requires Go 1.22+ and [golangci-lint](https://golangci-lint.run/) for linting.

```bash
# run unit tests
go test ./...

# lint (see .golangci.yml)
golangci-lint run ./...

# build the binary
go build -o tabui ./cmd/tabui

# install onto PATH
go install ./cmd/tabui
```

Try the package-scripts picker (interactive TTY required):

```bash
go run ./cmd/tabui --line="pnpm run " --cwd=/path/to/node/project
```

On success, the chosen script name is printed on stdout (for shell insert). Esc / Ctrl+C cancels with no output.

Package-scripts provider tests live next to the code under `internal/package_scripts/` (fixtures in `testdata/`).

## Out of scope (for now)

- Rebuilding a full terminal emulator or Warp’s floating UI chrome
- Replacing the entire shell
- Completing every CLI on day one — ship path + package-script providers first

## Success

A Go tool you can bind in your shell so that, depending on what you are typing (`mv`, `pnpm`, `npm run`, …), you get a fast fuzzy TUI of the **right next options** — files/folders, scripts, or commands — with a clear label + detail, in any terminal including Cursor’s.

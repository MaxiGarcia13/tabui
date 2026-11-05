# tabui

A Go CLI + TUI that suggests the **next thing to type** while you compose a shell command — paths, scripts, subcommands, or flags, depending on context.

Press **Tab** in your shell: tabui opens a fuzzy picker with **label** + **detail**, you pick a row, and that token is inserted into the line. Works in any terminal. If nothing matches, the shell falls back to normal completion.

| You are typing | tabui suggests |
|----------------|----------------|
| `pnpm ` / `pnpm run ` / `npm run ` | scripts from nearest `package.json`, plus tool commands when relevant |
| `mv ` / `cp ` / `rm ` / `cd ` | files and folders (paths) |

## Status

Target: suggest the next token from context (package scripts, paths, flags, …) via a Tab fuzzy picker in any terminal.

**Status:** ✅ ready · 🚧 in progress · ❌ not started

| Piece | Notes | Status |
|-------|-------|--------|
| Package scripts provider | `npm` / `pnpm` / `yarn` / `bun`; nearest `package.json` in Go | ✅ |
| Generic TUI | Fuzzy list with label + detail | ✅ |
| CLI wiring | `--line` / `--cwd` / `--cursor` | ✅ |
| zsh + bash Tab keybinding | `tabui init`; falls back when no match | ✅ |
| Path provider | Files/folders for `mv` / `cp` / `cd` / `rm` / `cat` | ✅ |
| Subcommand / flags provider | Known flags and subcommands for selected tools | ❌ |
| More CLI providers | Beyond path + package scripts — not every tool on day one | ❌ |

## Install

Requires Go 1.24+.

1. Install the binary onto your `PATH`:

```bash
go install github.com/MaxiGarcia13/tabui/cmd/tabui@latest
```

Or from a clone:

```bash
go install ./cmd/tabui
```

2. Enable the Tab widget — add one line to your shell rc:

**zsh** — `~/.zshrc`:

```zsh
eval "$(tabui init zsh)"
```

**bash** — `~/.bashrc`:

```bash
eval "$(tabui init bash)"
```

3. Open a new shell (or re-source your rc file).

Try: `pnpm ` + **Tab** → tool commands and scripts; `mv ` + **Tab** → files and folders.

### Disable

| Scope | zsh | bash |
|-------|-----|------|
| This session | `bindkey '^I' expand-or-complete` | `bind '"\t": complete'` |
| Permanent | Remove the `eval "$(tabui init …)"` line from your rc file | same |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, tests, and layout.

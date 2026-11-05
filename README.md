# tabui

A Go CLI + TUI that suggests the **next thing to type** while you compose a shell command — paths, scripts, subcommands, or flags, depending on context.

Press **Tab** in your shell: tabui opens a fuzzy picker with **label** + **detail**, you pick a row, and that token is inserted into the line. Works in any terminal. If nothing matches, the shell falls back to normal completion.

| You are typing | tabui suggests |
|----------------|----------------|
| `pnpm ` / `pnpm run ` / `npm run ` | scripts from nearest `package.json`, plus tool commands when relevant |
| `mv ` / `cp ` / `rm ` / `cd ` | files and folders (paths) — planned |

## Status

### Done

| Piece | Notes |
|-------|-------|
| Package scripts provider | `npm` / `pnpm` / `yarn` / `bun`; nearest `package.json` in Go |
| Generic TUI | Fuzzy list with label + detail |
| CLI wiring | `--line` / `--cwd` / `--cursor` |
| zsh + bash Tab keybinding | Falls back to normal completion when no match |

### Left to do

| Piece | Notes |
|-------|-------|
| Path provider | Files/folders for `mv` / `cp` / `cd` / `rm` / `cat` / … |
| Subcommand / flags provider | Known flags and subcommands for selected tools |
| More CLI providers | Beyond path + package scripts — not every tool on day one |

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

2. Source the shell widget (adjust the path if you cloned the repo):

**zsh** — add to `~/.zshrc`:

```zsh
source /path/to/tabui/shell/tabui.zsh
```

**bash** — add to `~/.bashrc`:

```bash
source /path/to/tabui/shell/tabui.bash
```

3. Open a new shell (or re-source your rc file).

Try in a Node project: `pnpm ` + **Tab** → tool commands and scripts; `pnpm run ` + **Tab** → scripts only.

### Disable

| Scope | zsh | bash |
|-------|-----|------|
| This session | `bindkey '^I' expand-or-complete` | `bind '"\t": complete'` |
| Permanent | Remove the `source …/shell/tabui.*` line from your rc file | same |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, tests, and layout.

---
name: shell-integration
description: >-
  Documents and edits zsh/bash keybindings for tabui. Use when wiring Tab
  (or similar), grabbing BUFFER/READLINE_LINE, invoking the tabui binary, inserting
  the selected token, or fixing shell completion widgets.
---

# Shell integration

Goal: keybinding grabs the current line → runs `tabui` → inserts the chosen
`Value` (or falls back / no-ops when empty).

Canonical scripts live in the repo:

- `shell/tabui.zsh`
- `shell/tabui.bash`

Keep README enable/disable docs in sync with those files and with CLI flags.

## Contract with the binary

- Flags: `--line`, `--cursor`, `--cwd` (see `cmd/tabui`).
- On success: print the **new left-of-cursor line** on stdout
  (e.g. `pnpm run build`), not only the bare token.
- Non-zero exit or empty stdout → leave the line unchanged (zsh: fall back to
  `expand-or-complete`).
- TUI uses alt screen + `/dev/tty` so the list clears and stdout stays clean.

## zsh (preferred)

Source `shell/tabui.zsh` from `~/.zshrc`. Widget:

1. Calls `tabui --line="$BUFFER" --cursor="$CURSOR" --cwd="$PWD"`
2. Sets `LBUFFER` to tabui’s stdout (full left side, e.g. `pnpm run build`)
3. Bound to `^I` (Tab)
4. On empty/error: `zle expand-or-complete` so normal Tab completion still works

Disable this session: `bindkey '^I' expand-or-complete`

## bash

Source `shell/tabui.bash`. Same flags via `READLINE_LINE` / `READLINE_POINT`.
Bound with `bind -x '"\t": tabui-complete'` (replaces default Tab completion).

Disable this session: `bind '"\t": complete'`

## Do not

- Hard-require Warp, VS Code, or a specific terminal app
- Print prompts or logs on stdout from the binary during insert mode
- Change the line when the user cancels the TUI

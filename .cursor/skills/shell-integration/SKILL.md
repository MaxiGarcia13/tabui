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

Canonical scripts (embedded into the binary via `tabui init`):

- `internal/shellinit/tabui.zsh`
- `internal/shellinit/tabui.bash`

Users enable with `eval "$(tabui init zsh)"` or `eval "$(tabui init bash)"`.
Keep README enable/disable docs in sync with those files and with CLI flags.

## Contract with the binary

- Flags: `--line`, `--cursor`, `--cwd` (see `cmd/tabui`).
- On success: print the **new left-of-cursor line** on stdout
  (e.g. `pnpm run build`), not only the bare token.
- Exit `0` + empty stdout → no match / empty load (zsh: `expand-or-complete`).
- Exit `2` → user canceled the TUI (Ctrl+C / Esc); leave the line alone,
  do **not** fall back to completion.
- Other non-zero → leave the line unchanged (zsh: `expand-or-complete`).
- TUI uses alt screen + `/dev/tty` so the list clears and stdout stays clean.

## zsh (preferred)

`tabui init zsh` prints the widget script. Widget:

1. Calls `tabui --line="$BUFFER" --cursor="$CURSOR" --cwd="$PWD"`
2. Sets `LBUFFER` to tabui’s stdout (full left side, e.g. `pnpm run build`)
3. Bound to `^I` (Tab)
4. On exit `2` (cancel): return without changing the line
5. On empty stdout or other error: `zle expand-or-complete` so normal Tab
   completion still works

Disable this session: `bindkey '^I' expand-or-complete`

## bash

`tabui init bash` prints the widget script. Same flags via `READLINE_LINE` /
`READLINE_POINT`. Bound with `bind -x '"\t": tabui-complete'` (replaces default
Tab completion). Cancel (exit `2`) and empty/error leave the line unchanged.

Disable this session: `bind '"\t": complete'`

## Do not

- Hard-require Warp, VS Code, or a specific terminal app
- Print prompts or logs on stdout from the binary during insert mode
- Change the line when the user cancels the TUI
- Fall back to `expand-or-complete` after a TUI cancel

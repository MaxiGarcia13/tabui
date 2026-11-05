---
name: shell-integration
description: >-
  Documents and edits zsh/bash keybindings for tabui. Use when wiring Ctrl+Space
  (or similar), grabbing BUFFER/READLINE_LINE, invoking the tabui binary, inserting
  the selected token, or fixing shell completion widgets.
---

# Shell integration

Goal: keybinding grabs the current line → runs `tabui` → inserts the chosen
`Value` (or no-ops on cancel).

Canonical scripts live in the repo:

- `shell/tabui.zsh`
- `shell/tabui.bash`

Keep README enable/disable docs in sync with those files and with CLI flags.

## Contract with the binary

- Flags: `--line`, `--cursor`, `--cwd` (see `cmd/tabui`).
- On success: print **only** the insert token on stdout (no UI chrome).
- Non-zero exit or empty stdout → leave the line unchanged.
- TUI uses `/dev/tty` (fallback stderr) so stdout stays clean for the shell.

## zsh (preferred)

Source `shell/tabui.zsh` from `~/.zshrc`. Widget:

1. Calls `tabui --line="$BUFFER" --cursor="$CURSOR" --cwd="$PWD"`
2. Replaces the partial token at the cursor (or appends after a prefix word /
   trailing space)
3. Bound to `^@` (Ctrl+Space on many terminals)

Disable this session: `bindkey -r '^@'`

## bash

Source `shell/tabui.bash`. Same flags via `READLINE_LINE` / `READLINE_POINT`.
Bound with `bind -x '"\C-@": tabui-complete'`.

Disable this session: `bind -r "\C-@"`

## Do not

- Hard-require Warp, VS Code, or a specific terminal app
- Print prompts or logs on stdout from the binary during insert mode
- Change the line when the user cancels the TUI

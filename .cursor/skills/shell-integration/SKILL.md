---
name: shell-integration
description: >-
  Documents and edits zsh/bash keybindings for tabui. Use when wiring Ctrl+Space
  (or similar), grabbing BUFFER/READLINE_LINE, invoking the tabui binary, inserting
  the selected token, or fixing shell completion widgets.
---

# Shell integration

Goal: keybinding grabs the current line → runs `tabui` → inserts the chosen `Value` (or no-ops on cancel).

## Contract with the binary

- Pass current line + cursor (flags or stdin; match whatever `cmd/tabui` exposes).
- On success: print **only** the insert token on stdout (no UI chrome).
- Non-zero exit or empty stdout → leave the line unchanged.
- TUI uses stderr / the terminal directly so stdout stays clean for the shell.

## zsh (preferred)

```zsh
tabui-complete() {
  local selected
  selected="$(tabui --line="$BUFFER" --cursor="$CURSOR")" || return
  [[ -n "$selected" ]] || return
  LBUFFER="${LBUFFER}${selected}"
}
zle -N tabui-complete
bindkey '^@' tabui-complete
```

Notes:

- `^@` is often Ctrl+Space; some terminals need a different sequence — document both if needed.
- Prefer replacing/appending at the cursor word once CLI supports partial-token replace; until then append to `LBUFFER` is the sketch default.
- Keep the widget in repo docs (README) in sync with real flags.

## bash

```bash
tabui-complete() {
  local selected
  selected="$(tabui --line="$READLINE_LINE" --cursor="$READLINE_POINT")" || return
  [[ -n "$selected" ]] || return
  READLINE_LINE="${READLINE_LINE:0:READLINE_POINT}${selected}${READLINE_LINE:READLINE_POINT}"
  READLINE_POINT=$((READLINE_POINT + ${#selected}))
}
bind -x '"\C-@": tabui-complete'
```

## Do not

- Hard-require Warp, VS Code, or a specific terminal app
- Print prompts or logs on stdout from the binary during insert mode
- Change the line when the user cancels the TUI

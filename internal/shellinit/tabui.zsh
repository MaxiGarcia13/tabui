# tabui zsh integration
# Usage: eval "$(tabui init zsh)"
# Requires: `tabui` on PATH
#
# On success, tabui prints the full new LBUFFER (e.g. "pnpm run build").

tabui-complete() {
  local new_left exit_code
  new_left="$(
    tabui --line="$BUFFER" --cursor="$CURSOR" --cwd="$PWD"
  )"
  exit_code=$?
  # Exit 2: user canceled the TUI — leave the line alone.
  if (( exit_code == 2 )); then
    return
  fi
  if (( exit_code != 0 )) || [[ -z "$new_left" ]]; then
    zle expand-or-complete
    return
  fi

  LBUFFER="$new_left"
  zle reset-prompt
}

zle -N tabui-complete
bindkey '^I' tabui-complete

# Disable in this session: bindkey '^I' expand-or-complete

# tabui zsh integration
# Usage: eval "$(tabui init zsh)"
# Requires: `tabui` on PATH
#
# On success, tabui prints the full new LBUFFER (e.g. "pnpm run build").

tabui-complete() {
  local new_left
  new_left="$(
    tabui --line="$BUFFER" --cursor="$CURSOR" --cwd="$PWD"
  )" || {
    zle expand-or-complete
    return
  }
  if [[ -z "$new_left" ]]; then
    zle expand-or-complete
    return
  fi

  LBUFFER="$new_left"
  zle reset-prompt
}

zle -N tabui-complete
bindkey '^I' tabui-complete

# Disable in this session: bindkey '^I' expand-or-complete

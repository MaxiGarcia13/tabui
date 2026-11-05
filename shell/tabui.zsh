# tabui zsh integration
# Usage: source /path/to/tabui/shell/tabui.zsh
# Requires: `tabui` on PATH (go install ./cmd/tabui)

tabui-complete() {
  local selected
  selected="$(
    tabui --line="$BUFFER" --cursor="$CURSOR" --cwd="$PWD"
  )" || return
  [[ -n "$selected" ]] || return

  local last_word="${LBUFFER##* }"
  if [[ -z "$LBUFFER" || "$LBUFFER" == *[[:space:]] ]]; then
    LBUFFER="${LBUFFER}${selected}"
  elif [[ "$last_word" == (run|run-script|npm|pnpm|yarn|bun) ]]; then
    LBUFFER="${LBUFFER} ${selected}"
  else
    LBUFFER="${LBUFFER%%[![:space:]]#}${selected}"
  fi
}

zle -N tabui-complete
bindkey '^@' tabui-complete

# Disable in this session: bindkey -r '^@'

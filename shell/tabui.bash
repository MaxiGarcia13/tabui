# tabui bash integration
# Usage: source /path/to/tabui/shell/tabui.bash
# Requires: `tabui` on PATH (go install ./cmd/tabui)
# Requires: bash with readline (bind -x)
#
# Note: binding Tab replaces default bash completion for that key.

tabui-complete() {
  local selected
  selected="$(
    tabui --line="$READLINE_LINE" --cursor="$READLINE_POINT" --cwd="$PWD"
  )" || return
  [[ -n "$selected" ]] || return

  local left="${READLINE_LINE:0:READLINE_POINT}"
  local right="${READLINE_LINE:READLINE_POINT}"
  local last_word="${left##* }"
  local new_left

  if [[ -z "$left" || "$left" == *[[:space:]] ]]; then
    new_left="${left}${selected}"
  elif [[ "$last_word" == run || "$last_word" == run-script \
    || "$last_word" == npm || "$last_word" == pnpm \
    || "$last_word" == yarn || "$last_word" == bun ]]; then
    new_left="${left} ${selected}"
  else
    new_left="${left%"${last_word}"}${selected}"
  fi

  READLINE_LINE="${new_left}${right}"
  READLINE_POINT=${#new_left}
}

bind -x '"\t": tabui-complete'

# Disable in this session: bind '"\t": complete'

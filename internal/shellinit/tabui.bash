# tabui bash integration
# Usage: eval "$(tabui init bash)"
# Requires: `tabui` on PATH
# Requires: bash with readline (bind -x)
#
# On success, tabui prints the full new left-of-cursor text
# (e.g. "pnpm run build").
#
# Note: binding Tab replaces default bash completion for that key.

tabui-complete() {
  local new_left exit_code
  local right="${READLINE_LINE:READLINE_POINT}"
  new_left="$(
    tabui --line="$READLINE_LINE" --cursor="$READLINE_POINT" --cwd="$PWD"
  )"
  exit_code=$?
  # Exit 2: user canceled the TUI — leave the line alone.
  (( exit_code == 2 )) && return
  (( exit_code != 0 )) && return
  [[ -n "$new_left" ]] || return

  READLINE_LINE="${new_left}${right}"
  READLINE_POINT=${#new_left}
}

bind -x '"\t": tabui-complete'

# Disable in this session: bind '"\t": complete'

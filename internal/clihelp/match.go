package clihelp

import (
  "os/exec"

  "github.com/MaxiGarcia13/tabui/internal/suggest"
)

var lookPath = exec.LookPath

func Match(ctx suggest.Context) bool {
  return parseTarget(lineUpToCursor(ctx)).ok
}

func lineUpToCursor(ctx suggest.Context) string {
  line := ctx.Line
  if ctx.Cursor < 0 {
    return line
  }
  if ctx.Cursor < len(line) {
    return line[:ctx.Cursor]
  }
  return line
}

func commandOnPath(command string) bool {
  _, err := lookPath(command)
  return err == nil
}

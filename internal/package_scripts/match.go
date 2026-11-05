package package_scripts

import (
  "strings"

  "github.com/MaxiGarcia13/tabui/internal/suggest"
)

// Longer prefixes first so "pnpm run" wins over bare "pnpm".
var scriptPrefixes = []string{
  "npm run-script ",
  "pnpm run-script ",
  "npm run ",
  "pnpm run ",
  "yarn run ",
  "bun run ",
  "pnpm ",
  "yarn ",
}

func Match(ctx suggest.Context) bool {
  _, ok := scriptPartial(lineUpToCursor(ctx))
  return ok
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

func scriptPartial(line string) (partial string, ok bool) {
  line = strings.TrimLeft(line, " \t")
  for _, prefix := range scriptPrefixes {
    if p, matched := queryAfter(line, prefix); matched {
      return p, true
    }
  }
  return "", false
}

func queryAfter(line, prefix string) (partial string, ok bool) {
  trimmed := strings.TrimRight(prefix, " ")
  if line == trimmed {
    return "", true
  }
  if !strings.HasPrefix(line, prefix) {
    return "", false
  }
  rest := line[len(prefix):]
  if strings.ContainsAny(rest, " \t") {
    return "", false
  }
  return rest, true
}

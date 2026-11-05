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
  _, matched := scriptPartial(lineUpToCursor(ctx))
  return matched
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

func scriptPartial(line string) (partial string, matched bool) {
  line = strings.TrimLeft(line, " \t")
  for _, prefix := range scriptPrefixes {
    if token, ok := queryAfter(line, prefix); ok {
      return token, true
    }
  }
  return "", false
}

func queryAfter(line, prefix string) (partial string, matched bool) {
  prefixWithoutSpace := strings.TrimRight(prefix, " ")
  if line == prefixWithoutSpace {
    return "", true
  }
  if !strings.HasPrefix(line, prefix) {
    return "", false
  }
  scriptToken := line[len(prefix):]
  if strings.ContainsAny(scriptToken, " \t") {
    return "", false
  }
  return scriptToken, true
}

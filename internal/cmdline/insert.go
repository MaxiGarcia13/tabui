package cmdline

import "strings"

var keepPrefixWords = map[string]struct{}{
  "run":        {},
  "run-script": {},
  "npm":        {},
  "pnpm":       {},
  "yarn":       {},
  "bun":        {},
}

// InsertToken returns left-of-cursor text after inserting token.
func InsertToken(left, token string) string {
  if left == "" ||
    strings.HasSuffix(left, " ") ||
    strings.HasSuffix(left, "\t") {
    return left + token
  }

  lastWord := left
  if index := strings.LastIndexAny(left, " \t"); index >= 0 {
    lastWord = left[index+1:]
  }

  if _, keep := keepPrefixWords[lastWord]; keep {
    return left + " " + token
  }

  if index := strings.LastIndexAny(left, " \t"); index >= 0 {
    return left[:index+1] + token
  }
  return token
}

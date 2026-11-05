package package_scripts

import (
  "strings"

  "github.com/MaxiGarcia13/tabui/internal/suggest"
)

type completionKind int

const (
  completionNone completionKind = iota
  completionToolCommands
  completionScripts
)

type completion struct {
  kind    completionKind
  tool    string
  partial string
}

func Match(ctx suggest.Context) bool {
  return parseCompletion(lineUpToCursor(ctx)).kind != completionNone
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

func parseCompletion(line string) completion {
  line = strings.TrimLeft(line, " \t")
  if line == "" {
    return completion{}
  }

  trailingSpace := strings.HasSuffix(line, " ") ||
    strings.HasSuffix(line, "\t")
  words := strings.Fields(line)
  if len(words) == 0 {
    return completion{}
  }

  tool := words[0]
  if _, ok := toolCommands[tool]; !ok {
    return completion{}
  }

  switch {
  case len(words) == 1:
    // "npm" or "npm " → first-level commands/flags.
    return completion{
      kind:    completionToolCommands,
      tool:    tool,
      partial: "",
    }

  case len(words) == 2 && !trailingSpace:
    second := words[1]
    if isRunSubcommand(second) {
      // "npm run" → scripts with empty partial.
      return completion{
        kind:    completionScripts,
        tool:    tool,
        partial: "",
      }
    }
    // "npm in" → filter tool commands.
    return completion{
      kind:    completionToolCommands,
      tool:    tool,
      partial: second,
    }

  case len(words) == 2 && trailingSpace:
    second := words[1]
    if isRunSubcommand(second) {
      // "npm run " → scripts.
      return completion{
        kind:    completionScripts,
        tool:    tool,
        partial: "",
      }
    }
    // "npm install " → deeper (not implemented yet).
    return completion{}

  case len(words) >= 3 && isRunSubcommand(words[1]):
    if len(words) == 3 && !trailingSpace {
      return completion{
        kind:    completionScripts,
        tool:    tool,
        partial: words[2],
      }
    }
    if len(words) == 3 && trailingSpace {
      // "npm run build " → past script name.
      return completion{}
    }
    // More tokens after script → no match.
    return completion{}
  }

  return completion{}
}

func isRunSubcommand(word string) bool {
  return word == "run" || word == "run-script"
}

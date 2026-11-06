package package_scripts

import (
  "strings"

  "github.com/MaxiGarcia13/tabui/internal/suggest"
)

type completion struct {
  ok       bool
  tool     string
  partial  string
  needsRun bool
}

func Match(ctx suggest.Context) bool {
  current := parseCompletion(lineUpToCursor(ctx))
  return current.ok && !current.needsRun
}

// ToolLevel covers bare package-tool slots merged with clihelp in main.
func ToolLevel(ctx suggest.Context) bool {
  current := parseCompletion(lineUpToCursor(ctx))
  return current.ok && current.needsRun
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
  if !isPackageTool(tool) {
    return completion{}
  }

  // Tool-level: `<tool>` / `<tool> ` / `<tool> partial`
  // Selecting a script inserts `run <script>`.
  if len(words) == 1 {
    return completion{
      ok:       true,
      tool:     tool,
      partial:  "",
      needsRun: true,
    }
  }

  if isRunSubcommand(words[1]) {
    switch {
    case len(words) == 2:
      return completion{
        ok:      true,
        tool:    tool,
        partial: "",
      }

    case len(words) == 3 && !trailingSpace:
      return completion{
        ok:      true,
        tool:    tool,
        partial: words[2],
      }

    default:
      return completion{}
    }
  }

  if len(words) == 2 && !trailingSpace {
    return completion{
      ok:       true,
      tool:     tool,
      partial:  words[1],
      needsRun: true,
    }
  }

  return completion{}
}

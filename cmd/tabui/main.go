package main

import (
  "errors"
  "flag"
  "fmt"
  "os"

  "github.com/MaxiGarcia13/tabui/internal/clihelp"
  "github.com/MaxiGarcia13/tabui/internal/cmdline"
  "github.com/MaxiGarcia13/tabui/internal/package_scripts"
  "github.com/MaxiGarcia13/tabui/internal/path"
  "github.com/MaxiGarcia13/tabui/internal/shellinit"
  "github.com/MaxiGarcia13/tabui/internal/suggest"
  "github.com/MaxiGarcia13/tabui/internal/ui"
)

const version = "0.5.1"

func main() {
  if len(os.Args) >= 2 && os.Args[1] == "init" {
    shell := ""
    if len(os.Args) >= 3 {
      shell = os.Args[2]
    }
    if err := shellinit.Write(os.Stdout, shell); err != nil {
      fmt.Fprintf(os.Stderr, "tabui: %v\n", err)
      os.Exit(1)
    }
    return
  }

  showVersion := flag.Bool("version", false, "print version and exit")
  line := flag.String("line", "", "current command line")
  cursor := flag.Int(
    "cursor",
    -1,
    "cursor byte offset (default: end of line)",
  )
  cwd := flag.String("cwd", "", "working directory (default: cwd)")
  flag.Parse()

  if *showVersion {
    fmt.Println(version)
    return
  }

  workingDir := *cwd
  if workingDir == "" {
    var err error
    workingDir, err = os.Getwd()
    if err != nil {
      fmt.Fprintf(os.Stderr, "tabui: %v\n", err)
      os.Exit(1)
    }
  }

  cursorOffset := *cursor
  if cursorOffset < 0 {
    cursorOffset = len(*line)
  }

  context := suggest.Context{
    Line:   *line,
    Cursor: cursorOffset,
    Cwd:    workingDir,
  }

  linePrefix := *line
  if cursorOffset < len(linePrefix) {
    linePrefix = linePrefix[:cursorOffset]
  }

  suggestions := suggest.FirstMatch(
    context,
    suggest.Funcs{
      MatchFn:   package_scripts.Match,
      SuggestFn: package_scripts.Suggest,
    },
    suggest.Funcs{
      MatchFn:   path.Match,
      SuggestFn: path.Suggest,
    },
  )

  var (
    selected string
    err      error
  )
  switch {
  case len(suggestions) > 0:
    selected, err = ui.Run(suggestions, linePrefix)
  default:
    items, matched, load := clihelp.Lookup(context)
    if !matched {
      return
    }
    if load == nil {
      selected, err = ui.Run(items, linePrefix)
    } else {
      selected, err = ui.RunLoading(linePrefix, load)
    }
  }
  if errors.Is(err, ui.ErrCanceled) {
    os.Exit(2)
  }
  if err != nil {
    fmt.Fprintf(os.Stderr, "tabui: %v\n", err)
    os.Exit(1)
  }
  if selected == "" {
    return
  }

  // Print the new left-of-cursor line so the shell shows e.g. "pnpm run build".
  fmt.Print(cmdline.InsertToken(linePrefix, selected))
}

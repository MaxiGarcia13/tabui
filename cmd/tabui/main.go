package main

import (
  "flag"
  "fmt"
  "os"

  "github.com/MaxiGarcia13/tabui/internal/package_scripts"
  "github.com/MaxiGarcia13/tabui/internal/suggest"
  "github.com/MaxiGarcia13/tabui/internal/ui"
)

func main() {
  line := flag.String("line", "", "current command line")
  cursor := flag.Int(
    "cursor",
    -1,
    "cursor byte offset (default: end of line)",
  )
  cwd := flag.String("cwd", "", "working directory (default: cwd)")
  flag.Parse()

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

  if !package_scripts.Match(context) {
    return
  }

  suggestions := package_scripts.Suggest(context)
  if len(suggestions) == 0 {
    return
  }

  selected, err := ui.Run(suggestions)
  if err != nil {
    fmt.Fprintf(os.Stderr, "tabui: %v\n", err)
    os.Exit(1)
  }
  if selected == "" {
    return
  }

  fmt.Print(selected)
}

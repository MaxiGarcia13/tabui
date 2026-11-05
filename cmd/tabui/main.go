package main

import (
  "flag"
  "fmt"
  "os"

  "github.com/MaxiGarcia13/tabui/internal/cmdline"
  "github.com/MaxiGarcia13/tabui/internal/package_scripts"
  "github.com/MaxiGarcia13/tabui/internal/suggest"
  "github.com/MaxiGarcia13/tabui/internal/ui"
)

const version = "0.1.1"

func main() {
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

  if !package_scripts.Match(context) {
    return
  }

  suggestions := package_scripts.Suggest(context)
  if len(suggestions) == 0 {
    return
  }

  linePrefix := *line
  if cursorOffset < len(linePrefix) {
    linePrefix = linePrefix[:cursorOffset]
  }

  selected, err := ui.Run(suggestions, linePrefix)
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

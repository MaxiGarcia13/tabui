package clihelp

import (
  "errors"
  "testing"

  "github.com/MaxiGarcia13/tabui/internal/suggest"
)

func TestMatch(t *testing.T) {
  previous := lookPath
  lookPath = func(file string) (string, error) {
    switch file {
    case "docker", "git", "npm":
      return "/usr/bin/" + file, nil
    default:
      return "", errors.New("not found")
    }
  }
  t.Cleanup(func() { lookPath = previous })

  tests := []struct {
    name string
    line string
    want bool
  }{
    {name: "bare incomplete", line: "dock", want: false},
    {name: "exact command no space", line: "docker", want: false},
    {name: "command space", line: "docker ", want: true},
    {name: "partial subcommand", line: "docker ru", want: true},
    {name: "unknown binary", line: "not-a-real-cli ", want: false},
    {name: "git space", line: "git ", want: true},
  }

  for _, testCase := range tests {
    t.Run(testCase.name, func(t *testing.T) {
      got := Match(suggest.Context{
        Line:   testCase.line,
        Cursor: len(testCase.line),
      })
      if got != testCase.want {
        t.Fatalf("Match(%q) = %v, want %v", testCase.line, got, testCase.want)
      }
    })
  }
}

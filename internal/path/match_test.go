package path

import (
  "testing"

  "github.com/MaxiGarcia13/tabui/internal/suggest"
)

func TestMatch(t *testing.T) {
  tests := []struct {
    name string
    line string
    want bool
  }{
    {name: "mv alone", line: "mv", want: false},
    {name: "mv space", line: "mv ", want: true},
    {name: "mv partial", line: "mv src", want: true},
    {name: "mv path with slash", line: "mv src/", want: true},
    {name: "cp with flag then path", line: "cp -r src", want: true},
    {name: "cp with flag space", line: "cp -r ", want: true},
    {name: "cp flag partial", line: "cp -", want: true},
    {name: "cp flag name", line: "cp -r", want: true},
    {name: "cd flag partial", line: "cd --", want: true},
    {name: "cat partial", line: "cat ./fi", want: true},
    {name: "rm second arg", line: "rm a b", want: true},
    {name: "cd space", line: "cd ", want: true},
    {name: "pnpm space", line: "pnpm ", want: false},
    {name: "leading space", line: "  mv ", want: true},
    {name: "unknown command", line: "echo ", want: false},
  }

  for _, testCase := range tests {
    t.Run(testCase.name, func(t *testing.T) {
      matched := Match(suggest.Context{
        Line:   testCase.line,
        Cursor: len(testCase.line),
      })
      if matched != testCase.want {
        t.Fatalf(
          "Match(%q) = %v, want %v",
          testCase.line,
          matched,
          testCase.want,
        )
      }
    })
  }
}

func TestParseCompletion(t *testing.T) {
  tests := []struct {
    name         string
    line         string
    wantOK       bool
    wantCommand  string
    wantPartial  string
    wantDirsOnly bool
  }{
    {
      name:         "empty path slot",
      line:         "mv ",
      wantOK:       true,
      wantCommand:  "mv",
      wantPartial:  "",
      wantDirsOnly: false,
    },
    {
      name:         "partial path",
      line:         "cat ./fi",
      wantOK:       true,
      wantCommand:  "cat",
      wantPartial:  "./fi",
      wantDirsOnly: false,
    },
    {
      name:         "cd dirs only",
      line:         "cd src/",
      wantOK:       true,
      wantCommand:  "cd",
      wantPartial:  "src/",
      wantDirsOnly: true,
    },
    {
      name:         "flag token",
      line:         "rm -",
      wantOK:       true,
      wantCommand:  "rm",
      wantPartial:  "-",
      wantDirsOnly: false,
    },
    {
      name:         "cd long flag",
      line:         "cd --",
      wantOK:       true,
      wantCommand:  "cd",
      wantPartial:  "--",
      wantDirsOnly: true,
    },
  }

  for _, testCase := range tests {
    t.Run(testCase.name, func(t *testing.T) {
      current := parseCompletion(testCase.line)
      if current.ok != testCase.wantOK {
        t.Fatalf("ok = %v, want %v", current.ok, testCase.wantOK)
      }
      if current.command != testCase.wantCommand {
        t.Fatalf(
          "command = %q, want %q",
          current.command,
          testCase.wantCommand,
        )
      }
      if current.partial != testCase.wantPartial {
        t.Fatalf(
          "partial = %q, want %q",
          current.partial,
          testCase.wantPartial,
        )
      }
      if current.dirsOnly != testCase.wantDirsOnly {
        t.Fatalf(
          "dirsOnly = %v, want %v",
          current.dirsOnly,
          testCase.wantDirsOnly,
        )
      }
    })
  }
}

func TestMatchUsesCursor(t *testing.T) {
  line := "mv src/foo && echo done"
  ctx := suggest.Context{Line: line, Cursor: len("mv src/f")}
  if !Match(ctx) {
    t.Fatal("expected match at cursor inside path token")
  }
  current := parseCompletion(lineUpToCursor(ctx))
  if !current.ok || current.partial != "src/f" {
    t.Fatalf("got ok=%v partial=%q", current.ok, current.partial)
  }
}

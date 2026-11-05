package package_scripts

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
    {name: "npm alone", line: "npm", want: true},
    {name: "npm space", line: "npm ", want: true},
    {name: "npm partial command", line: "npm in", want: true},
    {name: "npm install word", line: "npm install", want: true},
    {name: "npm install deeper", line: "npm install ", want: false},
    {name: "npm run space", line: "npm run ", want: true},
    {name: "npm run partial", line: "npm run de", want: true},
    {name: "npm run-script", line: "npm run-script ", want: true},
    {name: "npm run exact", line: "npm run", want: true},
    {name: "pnpm run", line: "pnpm run ", want: true},
    {name: "pnpm bare", line: "pnpm ", want: true},
    {name: "pnpm bare partial", line: "pnpm ad", want: true},
    {name: "pnpm exact", line: "pnpm", want: true},
    {name: "yarn run", line: "yarn run ", want: true},
    {name: "yarn bare", line: "yarn ", want: true},
    {name: "bun run", line: "bun run ", want: true},
    {name: "bun run partial", line: "bun run te", want: true},
    {name: "bun bare", line: "bun ", want: true},
    {name: "mv", line: "mv ", want: false},
    {name: "past script args", line: "pnpm run build --watch", want: false},
    {name: "leading space", line: "  pnpm run ", want: true},
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

func TestParseCompletionLevels(t *testing.T) {
  tests := []struct {
    name        string
    line        string
    wantKind    completionKind
    wantPartial string
  }{
    {
      name:        "tool level",
      line:        "npm ",
      wantKind:    completionToolCommands,
      wantPartial: "",
    },
    {
      name:        "tool partial",
      line:        "pnpm in",
      wantKind:    completionToolCommands,
      wantPartial: "in",
    },
    {
      name:        "scripts after run",
      line:        "npm run ",
      wantKind:    completionScripts,
      wantPartial: "",
    },
    {
      name:        "script partial",
      line:        "npm run bu",
      wantKind:    completionScripts,
      wantPartial: "bu",
    },
  }

  for _, testCase := range tests {
    t.Run(testCase.name, func(t *testing.T) {
      current := parseCompletion(testCase.line)
      if current.kind != testCase.wantKind {
        t.Fatalf(
          "kind = %v, want %v",
          current.kind,
          testCase.wantKind,
        )
      }
      if current.partial != testCase.wantPartial {
        t.Fatalf(
          "partial = %q, want %q",
          current.partial,
          testCase.wantPartial,
        )
      }
    })
  }
}

func TestMatchUsesCursor(t *testing.T) {
  line := "pnpm run build && echo done"
  ctx := suggest.Context{Line: line, Cursor: len("pnpm run bu")}
  if !Match(ctx) {
    t.Fatal("expected match at cursor inside script token")
  }
  current := parseCompletion(lineUpToCursor(ctx))
  if current.kind != completionScripts || current.partial != "bu" {
    t.Fatalf("got kind=%v partial=%q", current.kind, current.partial)
  }
}

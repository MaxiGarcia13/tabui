package package_scripts

import (
  "testing"

  "github.com/MaxiGarcia13/tabui/internal/suggest"
)

func TestMatch(t *testing.T) {
  tests := []struct {
    name  string
    line  string
    want  bool
  }{
    {name: "npm run space", line: "npm run ", want: true},
    {name: "npm run partial", line: "npm run de", want: true},
    {name: "npm run-script", line: "npm run-script ", want: true},
    {name: "npm run exact", line: "npm run", want: true},
    {name: "pnpm run", line: "pnpm run ", want: true},
    {name: "pnpm bare", line: "pnpm ", want: true},
    {name: "pnpm bare partial", line: "pnpm de", want: true},
    {name: "pnpm exact", line: "pnpm", want: true},
    {name: "yarn run", line: "yarn run ", want: true},
    {name: "yarn bare", line: "yarn ", want: true},
    {name: "bun run", line: "bun run ", want: true},
    {name: "bun run partial", line: "bun run te", want: true},
    {name: "mv", line: "mv ", want: false},
    {name: "npm install", line: "npm install", want: false},
    {name: "npm alone", line: "npm", want: false},
    {name: "bun bare", line: "bun ", want: false},
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

func TestMatchUsesCursor(t *testing.T) {
  line := "pnpm run build && echo done"
  ctx := suggest.Context{Line: line, Cursor: len("pnpm run bu")}
  if !Match(ctx) {
    t.Fatal("expected match at cursor inside script token")
  }
  partial, matched := scriptPartial(lineUpToCursor(ctx))
  if !matched || partial != "bu" {
    t.Fatalf("partial = %q, matched = %v", partial, matched)
  }
}

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
		{name: "npm alone", line: "npm", want: false},
		{name: "npm space", line: "npm ", want: false},
		{name: "npm partial command", line: "npm in", want: false},
		{name: "npm install word", line: "npm install", want: false},
		{name: "npm install deeper", line: "npm install ", want: false},
		{name: "npm run space", line: "npm run ", want: true},
		{name: "npm run partial", line: "npm run de", want: true},
		{name: "npm run-script", line: "npm run-script ", want: true},
		{name: "npm run exact", line: "npm run", want: true},
		{name: "pnpm run", line: "pnpm run ", want: true},
		{name: "pnpm bare", line: "pnpm ", want: false},
		{name: "pnpm bare partial", line: "pnpm ad", want: false},
		{name: "pnpm exact", line: "pnpm", want: false},
		{name: "yarn run", line: "yarn run ", want: true},
		{name: "yarn bare", line: "yarn ", want: false},
		{name: "bun run", line: "bun run ", want: true},
		{name: "bun run partial", line: "bun run te", want: true},
		{name: "bun bare", line: "bun ", want: false},
		{name: "mv", line: "mv ", want: false},
		{name: "past script args", line: "pnpm run build --watch", want: false},
		{name: "past script space", line: "pnpm run build ", want: false},
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

func TestParseCompletionScripts(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		wantOK      bool
		wantPartial string
	}{
		{
			name:        "tool level falls through",
			line:        "npm ",
			wantOK:      false,
			wantPartial: "",
		},
		{
			name:        "tool partial falls through",
			line:        "pnpm in",
			wantOK:      false,
			wantPartial: "",
		},
		{
			name:        "scripts after run",
			line:        "npm run ",
			wantOK:      true,
			wantPartial: "",
		},
		{
			name:        "script partial",
			line:        "npm run bu",
			wantOK:      true,
			wantPartial: "bu",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			current := parseCompletion(testCase.line)
			if current.ok != testCase.wantOK {
				t.Fatalf("ok = %v, want %v", current.ok, testCase.wantOK)
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
	if !current.ok || current.partial != "bu" {
		t.Fatalf("got ok=%v partial=%q", current.ok, current.partial)
	}
}

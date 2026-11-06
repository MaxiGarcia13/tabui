package package_scripts

import (
	"strings"

	"github.com/MaxiGarcia13/tabui/internal/suggest"
)

type completion struct {
	ok      bool
	tool    string
	partial string
}

func Match(ctx suggest.Context) bool {
	return parseCompletion(lineUpToCursor(ctx)).ok
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

	// Only script slots: `<tool> run` / `run-script` [partial].
	// Bare tool-level tokens fall through to clihelp.
	if len(words) == 1 {
		return completion{}
	}

	if !isRunSubcommand(words[1]) {
		return completion{}
	}

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

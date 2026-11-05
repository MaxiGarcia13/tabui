package path

import (
	"strings"

	"github.com/MaxiGarcia13/tabui/internal/suggest"
)

type completion struct {
	ok       bool
	command  string
	partial  string
	dirsOnly bool
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

	command := words[0]
	if !isPathCommand(command) {
		return completion{}
	}

	// Still typing the command name alone — not a path slot yet.
	if len(words) == 1 && !trailingSpace {
		return completion{}
	}

	partial := ""
	if !trailingSpace {
		partial = words[len(words)-1]
	}

	if strings.HasPrefix(partial, "-") {
		return completion{}
	}

	return completion{
		ok:       true,
		command:  command,
		partial:  partial,
		dirsOnly: isDirOnlyCommand(command),
	}
}

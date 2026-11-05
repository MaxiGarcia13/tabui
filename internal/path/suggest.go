package path

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/MaxiGarcia13/tabui/internal/suggest"
)

func Suggest(ctx suggest.Context) []suggest.Suggestion {
	current := parseCompletion(lineUpToCursor(ctx))
	if !current.ok {
		return nil
	}
	return suggestPaths(ctx.Cwd, current.partial, current.dirsOnly)
}

func suggestPaths(
	cwd string,
	partial string,
	dirsOnly bool,
) []suggest.Suggestion {
	dirPart, namePrefix := splitPartial(partial)
	listDir, err := resolveListDir(cwd, dirPart)
	if err != nil {
		return nil
	}

	entries, err := os.ReadDir(listDir)
	if err != nil {
		return nil
	}

	absCwd, err := filepath.Abs(cwd)
	if err != nil {
		absCwd = cwd
	}

	showHidden := strings.HasPrefix(namePrefix, ".")
	matched := make([]suggest.Suggestion, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if !showHidden && strings.HasPrefix(name, ".") {
			continue
		}
		if !matchesPrefix(name, namePrefix) {
			continue
		}

		isDir := entry.IsDir()
		if !isDir && entry.Type()&os.ModeSymlink != 0 {
			isDir = isSymlinkDir(listDir, name)
		}
		if dirsOnly && !isDir {
			continue
		}

		label := name
		if isDir {
			label += "/"
		}
		value := dirPart + label
		matched = append(matched, suggest.Suggestion{
			Label:  label,
			Detail: detailFor(absCwd, listDir, name, isDir, value),
			Value:  value,
		})
	}

	sort.Slice(matched, func(left, right int) bool {
		leftDir := strings.HasSuffix(matched[left].Label, "/")
		rightDir := strings.HasSuffix(matched[right].Label, "/")
		if leftDir != rightDir {
			return leftDir
		}
		return matched[left].Label < matched[right].Label
	})
	return matched
}

func matchesPrefix(name, prefix string) bool {
	if prefix == "" {
		return true
	}
	return strings.HasPrefix(
		strings.ToLower(name),
		strings.ToLower(prefix),
	)
}

func isSymlinkDir(parent, name string) bool {
	info, err := os.Stat(filepath.Join(parent, name))
	return err == nil && info.IsDir()
}

func detailFor(
	absCwd string,
	listDir string,
	name string,
	isDir bool,
	fallback string,
) string {
	full := filepath.Join(listDir, name)
	rel, err := filepath.Rel(absCwd, full)
	if err != nil || strings.HasPrefix(rel, "..") {
		return fallback
	}
	detail := filepath.ToSlash(rel)
	if isDir {
		detail += "/"
	}
	return detail
}

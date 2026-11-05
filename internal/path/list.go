package path

import (
	"os"
	"path/filepath"
	"strings"
)

var lookUpHome = os.UserHomeDir

func splitPartial(partial string) (dirPart, namePrefix string) {
	if partial == "" {
		return "", ""
	}
	index := strings.LastIndexAny(partial, `/\`)
	if index < 0 {
		return "", partial
	}
	return partial[:index+1], partial[index+1:]
}

func expandHome(path string) (string, error) {
	if path == "~" {
		return lookUpHome()
	}
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, `~\`) {
		home, err := lookUpHome()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, path[2:]), nil
	}
	return path, nil
}

func resolveListDir(cwd, dirPart string) (string, error) {
	base := cwd
	if dirPart != "" {
		expanded, err := expandHome(dirPart)
		if err != nil {
			return "", err
		}
		if filepath.IsAbs(expanded) {
			base = expanded
		} else {
			base = filepath.Join(cwd, expanded)
		}
	}
	return filepath.Abs(base)
}

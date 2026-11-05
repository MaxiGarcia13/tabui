package package_scripts

import (
  "encoding/json"
  "os"
  "path/filepath"
)

type packageJSON struct {
  Scripts map[string]string `json:"scripts"`
}

func findNearestPackageJSON(cwd string) (string, error) {
  dir, err := filepath.Abs(cwd)
  if err != nil {
    return "", err
  }
  for {
    candidate := filepath.Join(dir, "package.json")
    info, err := os.Stat(candidate)
    if err == nil && !info.IsDir() {
      return candidate, nil
    }
    parent := filepath.Dir(dir)
    if parent == dir {
      return "", os.ErrNotExist
    }
    dir = parent
  }
}

func readScripts(path string) (map[string]string, error) {
  data, err := os.ReadFile(path)
  if err != nil {
    return nil, err
  }
  var pkg packageJSON
  if err := json.Unmarshal(data, &pkg); err != nil {
    return nil, err
  }
  if pkg.Scripts == nil {
    return map[string]string{}, nil
  }
  return pkg.Scripts, nil
}

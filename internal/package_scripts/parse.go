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
  directory, err := filepath.Abs(cwd)
  if err != nil {
    return "", err
  }
  for {
    candidate := filepath.Join(directory, "package.json")
    fileInfo, statErr := os.Stat(candidate)
    if statErr == nil && !fileInfo.IsDir() {
      return candidate, nil
    }
    parent := filepath.Dir(directory)
    if parent == directory {
      return "", os.ErrNotExist
    }
    directory = parent
  }
}

func readScripts(path string) (map[string]string, error) {
  rawJSON, err := os.ReadFile(path)
  if err != nil {
    return nil, err
  }
  var manifest packageJSON
  if err := json.Unmarshal(rawJSON, &manifest); err != nil {
    return nil, err
  }
  if manifest.Scripts == nil {
    return map[string]string{}, nil
  }
  return manifest.Scripts, nil
}

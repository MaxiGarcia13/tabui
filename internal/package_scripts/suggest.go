package package_scripts

import (
  "sort"
  "strings"

  "github.com/MaxiGarcia13/tabui/internal/suggest"
)

func Suggest(ctx suggest.Context) []suggest.Suggestion {
  current := parseCompletion(lineUpToCursor(ctx))
  if !current.ok {
    return nil
  }
  return suggestScripts(ctx.Cwd, current.partial, current.needsRun)
}

func suggestScripts(
  workingDir string,
  partial string,
  needsRun bool,
) []suggest.Suggestion {
  packageJSONPath, err := findNearestPackageJSON(workingDir)
  if err != nil {
    return nil
  }

  scripts, err := readScripts(packageJSONPath)
  if err != nil {
    return nil
  }

  scriptNames := make([]string, 0, len(scripts))
  for scriptName := range scripts {
    if matchesPartial(scriptName, partial) {
      scriptNames = append(scriptNames, scriptName)
    }
  }
  sort.Strings(scriptNames)

  suggestions := make([]suggest.Suggestion, 0, len(scriptNames))
  for _, scriptName := range scriptNames {
    value := scriptName
    if needsRun {
      value = "run " + scriptName
    }
    suggestions = append(suggestions, suggest.Suggestion{
      Label:  scriptName,
      Detail: scripts[scriptName],
      Value:  value,
    })
  }
  return suggestions
}

func matchesPartial(name, partial string) bool {
  if partial == "" {
    return true
  }
  return strings.HasPrefix(
    strings.ToLower(name),
    strings.ToLower(partial),
  )
}

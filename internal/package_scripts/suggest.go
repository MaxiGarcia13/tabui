package package_scripts

import (
  "sort"
  "strings"

  "github.com/MaxiGarcia13/tabui/internal/suggest"
)

func Suggest(ctx suggest.Context) []suggest.Suggestion {
  partial, matched := scriptPartial(lineUpToCursor(ctx))
  if !matched {
    return nil
  }

  packageJSONPath, err := findNearestPackageJSON(ctx.Cwd)
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
    suggestions = append(suggestions, suggest.Suggestion{
      Label:  scriptName,
      Detail: scripts[scriptName],
      Value:  scriptName,
    })
  }
  return suggestions
}

func matchesPartial(scriptName, partial string) bool {
  if partial == "" {
    return true
  }
  return strings.HasPrefix(
    strings.ToLower(scriptName),
    strings.ToLower(partial),
  )
}

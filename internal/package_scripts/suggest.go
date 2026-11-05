package package_scripts

import (
  "sort"
  "strings"

  "github.com/MaxiGarcia13/tabui/internal/suggest"
)

func Suggest(ctx suggest.Context) []suggest.Suggestion {
  current := parseCompletion(lineUpToCursor(ctx))
  switch current.kind {
  case completionToolCommands:
    return suggestToolCommands(current.tool, current.partial)
  case completionScripts:
    return suggestScripts(ctx.Cwd, current.partial)
  default:
    return nil
  }
}

func suggestToolCommands(
  tool string,
  partial string,
) []suggest.Suggestion {
  options, ok := toolCommands[tool]
  if !ok {
    return nil
  }

  matched := make([]commandOption, 0, len(options))
  for _, option := range options {
    if matchesPartial(option.value, partial) {
      matched = append(matched, option)
    }
  }
  sort.Slice(matched, func(left, right int) bool {
    return matched[left].value < matched[right].value
  })

  suggestions := make([]suggest.Suggestion, 0, len(matched))
  for _, option := range matched {
    suggestions = append(suggestions, suggest.Suggestion{
      Label:  option.value,
      Detail: option.detail,
      Value:  option.value,
    })
  }
  return suggestions
}

func suggestScripts(
  workingDir string,
  partial string,
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
    suggestions = append(suggestions, suggest.Suggestion{
      Label:  scriptName,
      Detail: scripts[scriptName],
      Value:  scriptName,
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

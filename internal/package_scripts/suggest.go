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
    return suggestToolLevel(
      ctx.Cwd,
      current.tool,
      current.partial,
    )
  case completionScripts:
    return suggestScripts(ctx.Cwd, current.partial, false)
  default:
    return nil
  }
}

func suggestToolLevel(
  workingDir string,
  tool string,
  partial string,
) []suggest.Suggestion {
  suggestions := suggestToolCommands(tool, partial)
  commandNames := commandNameSet(tool)
  scripts := suggestScripts(workingDir, partial, true)
  for _, script := range scripts {
    if _, exists := commandNames[script.Label]; exists {
      continue
    }
    suggestions = append(suggestions, script)
  }
  return suggestions
}

func commandNameSet(tool string) map[string]struct{} {
  options, ok := toolCommands[tool]
  if !ok {
    return nil
  }
  names := make(map[string]struct{}, len(options))
  for _, option := range options {
    names[option.value] = struct{}{}
  }
  return names
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

// suggestScripts lists package scripts. When withRunPrefix is true (tool
// level), Value is "run <name>" so the line becomes e.g. "pnpm run build".
func suggestScripts(
  workingDir string,
  partial string,
  withRunPrefix bool,
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
    detail := scripts[scriptName]
    if withRunPrefix {
      value = "run " + scriptName
      detail = "script: " + detail
    }
    suggestions = append(suggestions, suggest.Suggestion{
      Label:  scriptName,
      Detail: detail,
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

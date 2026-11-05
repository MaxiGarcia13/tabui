package package_scripts

import (
  "sort"
  "strings"

  "github.com/MaxiGarcia13/tabui/internal/suggest"
)

func Suggest(ctx suggest.Context) []suggest.Suggestion {
  partial, ok := scriptPartial(lineUpToCursor(ctx))
  if !ok {
    return nil
  }

  path, err := findNearestPackageJSON(ctx.Cwd)
  if err != nil {
    return nil
  }

  scripts, err := readScripts(path)
  if err != nil {
    return nil
  }

  names := make([]string, 0, len(scripts))
  for name := range scripts {
    if matchesPartial(name, partial) {
      names = append(names, name)
    }
  }
  sort.Strings(names)

  out := make([]suggest.Suggestion, 0, len(names))
  for _, name := range names {
    out = append(out, suggest.Suggestion{
      Label:  name,
      Detail: scripts[name],
      Value:  name,
    })
  }
  return out
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

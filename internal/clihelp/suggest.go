package clihelp

import (
  "os"
  "strings"
  "time"

  "github.com/MaxiGarcia13/tabui/internal/suggest"
)

// Lookup returns suggestions for a generic CLI help target.
// matched is false when this provider does not apply.
// On a cache hit, load is nil and items are ready for a sync UI.
// On a cache miss, items is nil and load fetches, parses, caches, and filters.
func Lookup(
  ctx suggest.Context,
) (
  items []suggest.Suggestion,
  matched bool,
  load func() []suggest.Suggestion,
) {
  current := parseTarget(lineUpToCursor(ctx))
  if !current.ok {
    return nil, false, nil
  }

  resolved, err := lookPath(current.command)
  if err != nil {
    return nil, false, nil
  }

  key := helpCacheKey(resolved, current.helpArgs)
  store := defaultStore()
  if cached, ok := store.Get(key); ok {
    return filterSuggestions(cached, current.partial), true, nil
  }

  partial := current.partial
  helpArgs := append([]string{}, current.helpArgs...)
  return nil, true, func() []suggest.Suggestion {
    text, err := fetchHelp(resolved, helpArgs)
    if err != nil || text == "" {
      return nil
    }
    parsed := Parse(text)
    if len(parsed) == 0 {
      return nil
    }
    store.Put(key, parsed)
    return filterSuggestions(parsed, partial)
  }
}

func helpCacheKey(resolved string, helpArgs []string) string {
  modTime := time.Time{}
  if info, err := os.Stat(resolved); err == nil {
    modTime = info.ModTime()
  }
  return cacheKey(resolved, helpArgs, modTime)
}

func filterSuggestions(
  items []suggest.Suggestion,
  partial string,
) []suggest.Suggestion {
  if partial == "" {
    return items
  }
  matched := make([]suggest.Suggestion, 0, len(items))
  for _, item := range items {
    if matchesPartial(item.Value, partial) {
      matched = append(matched, item)
    }
  }
  return matched
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

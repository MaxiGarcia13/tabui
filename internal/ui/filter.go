package ui

import (
	"strings"

	"github.com/MaxiGarcia13/tabui/internal/suggest"
)

func filterSuggestions(
  items []suggest.Suggestion,
  query string,
) []suggest.Suggestion {
  query = strings.TrimSpace(query)
  if query == "" {
    allItems := make([]suggest.Suggestion, len(items))
    copy(allItems, items)
    return allItems
  }
  queryLower := strings.ToLower(query)
  matched := make([]suggest.Suggestion, 0, len(items))
  for _, item := range items {
    labelLower := strings.ToLower(item.Label)
    if strings.Contains(labelLower, queryLower) {
      matched = append(matched, item)
    }
  }
  return matched
}

package ui

import (
  "testing"

  "github.com/MaxiGarcia13/tabui/internal/suggest"
)

func TestFilterSuggestions(t *testing.T) {
  items := []suggest.Suggestion{
    {Label: "build", Detail: "tsc", Value: "build"},
    {Label: "dev", Detail: "vite", Value: "dev"},
    {Label: "test", Detail: "vitest", Value: "test"},
  }

  unfiltered := filterSuggestions(items, "")
  if len(unfiltered) != 3 {
    t.Fatalf("empty query len = %d", len(unfiltered))
  }

  matched := filterSuggestions(items, "te")
  if len(matched) != 1 || matched[0].Label != "test" {
    t.Fatalf("got %+v", matched)
  }

  noMatches := filterSuggestions(items, "zzz")
  if len(noMatches) != 0 {
    t.Fatalf("got %+v", noMatches)
  }
}

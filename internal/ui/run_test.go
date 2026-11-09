package ui

import (
  "testing"

  "github.com/MaxiGarcia13/tabui/internal/suggest"
)

func TestRunEmptyReturnsEmpty(t *testing.T) {
  selected, err := Run(nil, "pnpm run ")
  if err != nil {
    t.Fatalf("err = %v", err)
  }
  if selected != "" {
    t.Fatalf("selected = %q, want empty", selected)
  }
}

func TestRunSingleReturnsValueWithoutPicker(t *testing.T) {
  selected, err := Run(
    []suggest.Suggestion{{Label: "test", Value: "test"}},
    "pnpm run ",
  )
  if err != nil {
    t.Fatalf("err = %v", err)
  }
  if selected != "test" {
    t.Fatalf("selected = %q, want test", selected)
  }
}

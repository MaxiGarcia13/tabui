package ui

import (
  "testing"

  "github.com/MaxiGarcia13/tabui/internal/suggest"
)

func TestEnsureCursorVisibleScrollsWindow(t *testing.T) {
  items := make([]suggest.Suggestion, 20)
  for index := range items {
    items[index] = suggest.Suggestion{
      Label: "item",
      Value: "item",
    }
  }
  state := newPickerState(items, "pnpm ")
  state.height = 8 // capacity = 8 - 2 - 1 = 5

  state.cursor = 10
  state.ensureCursorVisible()
  if state.offset != 6 {
    t.Fatalf("offset = %d, want 6", state.offset)
  }

  state.cursor = 2
  state.ensureCursorVisible()
  if state.offset != 2 {
    t.Fatalf("offset = %d, want 2", state.offset)
  }
}

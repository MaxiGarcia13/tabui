package ui

import (
  tea "github.com/charmbracelet/bubbletea"

  "github.com/MaxiGarcia13/tabui/internal/suggest"
)

type pickerState struct {
  allItems []suggest.Suggestion
  filtered []suggest.Suggestion
  cursor   int
  query    string
  chosen   string
  canceled bool
}

func newPickerState(items []suggest.Suggestion) pickerState {
  filtered := filterSuggestions(items, "")
  return pickerState{
    allItems: items,
    filtered: filtered,
  }
}

func (state pickerState) Init() tea.Cmd {
  return nil
}

func (state *pickerState) refilter() {
  state.filtered = filterSuggestions(state.allItems, state.query)
  if state.cursor >= len(state.filtered) {
    state.cursor = len(state.filtered) - 1
  }
  if state.cursor < 0 {
    state.cursor = 0
  }
}

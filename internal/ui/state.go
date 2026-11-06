package ui

import (
  tea "github.com/charmbracelet/bubbletea"

  "github.com/MaxiGarcia13/tabui/internal/suggest"
)

const (
  defaultPickerHeight = 16
  pickerHeaderLines   = 2
  pickerFooterLines   = 1
)

type loadedMsg struct {
  items []suggest.Suggestion
}

type pickerState struct {
  allItems   []suggest.Suggestion
  filtered   []suggest.Suggestion
  cursor     int
  offset     int
  height     int
  query      string
  chosen     string
  canceled   bool
  linePrefix string
  loading    bool
  loadFn     func() []suggest.Suggestion
}

func newPickerState(
  items []suggest.Suggestion,
  linePrefix string,
) pickerState {
  filtered := filterSuggestions(items, "")
  return pickerState{
    allItems:   items,
    filtered:   filtered,
    linePrefix: linePrefix,
    height:     defaultPickerHeight,
  }
}

func newLoadingPickerState(
  linePrefix string,
  load func() []suggest.Suggestion,
) pickerState {
  return pickerState{
    linePrefix: linePrefix,
    height:     defaultPickerHeight,
    loading:    true,
    loadFn:     load,
  }
}

func (state pickerState) Init() tea.Cmd {
  if !state.loading || state.loadFn == nil {
    return nil
  }
  loadFn := state.loadFn
  return func() tea.Msg {
    return loadedMsg{items: loadFn()}
  }
}

func (state *pickerState) applyLoaded(items []suggest.Suggestion) {
  state.loading = false
  state.loadFn = nil
  if len(items) == 0 {
    state.canceled = true
    return
  }
  state.allItems = items
  state.refilter()
}

func (state *pickerState) refilter() {
  state.filtered = filterSuggestions(state.allItems, state.query)
  if state.cursor >= len(state.filtered) {
    state.cursor = len(state.filtered) - 1
  }
  if state.cursor < 0 {
    state.cursor = 0
  }
  state.offset = 0
  state.ensureCursorVisible()
}

func (state *pickerState) visibleCapacity() int {
  capacity := state.height - pickerHeaderLines - pickerFooterLines
  if capacity < 1 {
    return 1
  }
  return capacity
}

func (state *pickerState) ensureCursorVisible() {
  capacity := state.visibleCapacity()
  if state.cursor < state.offset {
    state.offset = state.cursor
  }
  if state.cursor >= state.offset+capacity {
    state.offset = state.cursor - capacity + 1
  }
  if state.offset < 0 {
    state.offset = 0
  }
  maxOffset := len(state.filtered) - capacity
  if maxOffset < 0 {
    maxOffset = 0
  }
  if state.offset > maxOffset {
    state.offset = maxOffset
  }
}

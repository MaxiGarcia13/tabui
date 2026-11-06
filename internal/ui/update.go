package ui

import tea "github.com/charmbracelet/bubbletea"

func (state pickerState) Update(
  message tea.Msg,
) (tea.Model, tea.Cmd) {
  switch message := message.(type) {
  case loadedMsg:
    state.applyLoaded(message.items)
    if state.canceled {
      return state, tea.Quit
    }
    return state, nil
  case tea.WindowSizeMsg:
    if message.Height > 0 {
      state.height = message.Height
    }
    state.ensureCursorVisible()
  case tea.KeyMsg:
    if state.loading {
      switch message.Type {
      case tea.KeyCtrlC, tea.KeyEsc:
        state.canceled = true
        return state, tea.Quit
      }
      return state, nil
    }
    switch message.Type {
    case tea.KeyCtrlC, tea.KeyEsc:
      state.canceled = true
      return state, tea.Quit
    case tea.KeyEnter:
      if len(state.filtered) == 0 {
        state.canceled = true
        return state, tea.Quit
      }
      state.chosen = state.filtered[state.cursor].Value
      return state, tea.Quit
    case tea.KeyUp:
      if state.cursor > 0 {
        state.cursor--
        state.ensureCursorVisible()
      }
    case tea.KeyDown:
      if state.cursor < len(state.filtered)-1 {
        state.cursor++
        state.ensureCursorVisible()
      }
    case tea.KeyPgUp:
      state.cursor -= state.visibleCapacity()
      if state.cursor < 0 {
        state.cursor = 0
      }
      state.ensureCursorVisible()
    case tea.KeyPgDown:
      state.cursor += state.visibleCapacity()
      if state.cursor > len(state.filtered)-1 {
        state.cursor = len(state.filtered) - 1
      }
      if state.cursor < 0 {
        state.cursor = 0
      }
      state.ensureCursorVisible()
    case tea.KeyHome:
      state.cursor = 0
      state.ensureCursorVisible()
    case tea.KeyEnd:
      if len(state.filtered) > 0 {
        state.cursor = len(state.filtered) - 1
      }
      state.ensureCursorVisible()
    case tea.KeyBackspace:
      if len(state.query) > 0 {
        state.query = state.query[:len(state.query)-1]
        state.refilter()
      }
    case tea.KeySpace:
      state.query += " "
      state.refilter()
    case tea.KeyRunes:
      state.query += string(message.Runes)
      state.refilter()
    }
  }
  return state, nil
}

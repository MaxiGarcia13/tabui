package ui

import (
  "os"

  tea "github.com/charmbracelet/bubbletea"

  "github.com/MaxiGarcia13/tabui/internal/suggest"
)

// Run shows a fuzzy list and returns the chosen Value.
// Empty string means the user canceled.
func Run(items []suggest.Suggestion) (string, error) {
  if len(items) == 0 {
    return "", nil
  }

  state := newPickerState(items)
  programOptions := []tea.ProgramOption{
    tea.WithOutput(os.Stderr),
  }
  tty, ttyErr := os.OpenFile("/dev/tty", os.O_RDWR, 0)
  if ttyErr == nil {
    defer tty.Close()
    programOptions = []tea.ProgramOption{
      tea.WithInput(tty),
      tea.WithOutput(tty),
    }
  }

  finalModel, err := tea.NewProgram(state, programOptions...).Run()
  if err != nil {
    return "", err
  }
  result, ok := finalModel.(pickerState)
  if !ok || result.canceled || result.chosen == "" {
    return "", nil
  }
  return result.chosen, nil
}

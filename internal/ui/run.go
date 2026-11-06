package ui

import (
  "errors"
  "os"

  tea "github.com/charmbracelet/bubbletea"

  "github.com/MaxiGarcia13/tabui/internal/suggest"
)

// ErrCanceled means the user dismissed the picker (Ctrl+C / Esc).
// The shell should leave the line unchanged and not fall back to completion.
var ErrCanceled = errors.New("canceled")

// Run shows a fuzzy list and returns the chosen Value.
// Empty string means no selection (shell may fall back).
// ErrCanceled means the user dismissed the picker.
// linePrefix is the command line left of the cursor (for preview).
func Run(
  items []suggest.Suggestion,
  linePrefix string,
) (string, error) {
  if len(items) == 0 {
    return "", nil
  }
  return runPicker(newPickerState(items, linePrefix))
}

// RunLoading opens the picker with a loading state, runs load in a Bubble Tea
// Cmd, then fills the list. Empty load results return "" (shell fallback).
func RunLoading(
  linePrefix string,
  load func() []suggest.Suggestion,
) (string, error) {
  if load == nil {
    return "", nil
  }
  return runPicker(newLoadingPickerState(linePrefix, load))
}

func runPicker(state pickerState) (string, error) {
  programOptions := []tea.ProgramOption{
    tea.WithAltScreen(),
    tea.WithOutput(os.Stderr),
  }
  tty, ttyErr := os.OpenFile("/dev/tty", os.O_RDWR, 0)
  if ttyErr == nil {
    defer tty.Close()
    programOptions = []tea.ProgramOption{
      tea.WithAltScreen(),
      tea.WithInput(tty),
      tea.WithOutput(tty),
    }
  }

  finalModel, err := tea.NewProgram(state, programOptions...).Run()
  if err != nil {
    return "", err
  }
  result, ok := finalModel.(pickerState)
  if !ok {
    return "", nil
  }
  if result.canceled {
    return "", ErrCanceled
  }
  return result.chosen, nil
}

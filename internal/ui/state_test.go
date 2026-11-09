package ui

import (
  "strings"
  "testing"

  tea "github.com/charmbracelet/bubbletea"

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

func TestLoadingInitReturnsCmd(t *testing.T) {
  state := newLoadingPickerState("docker ", func() []suggest.Suggestion {
    return []suggest.Suggestion{{Label: "run", Value: "run"}}
  })
  if !state.loading {
    t.Fatal("expected loading")
  }
  if state.Init() == nil {
    t.Fatal("expected load cmd")
  }
}

func TestLoadedMsgSingleAutoChooses(t *testing.T) {
  state := newLoadingPickerState("docker ", nil)
  next, cmd := state.Update(loadedMsg{
    items: []suggest.Suggestion{
      {Label: "run", Value: "run", Detail: "Create a container"},
    },
  })
  result := next.(pickerState)
  if cmd == nil {
    t.Fatal("expected quit cmd")
  }
  if result.loading {
    t.Fatal("expected loading false")
  }
  if result.chosen != "run" {
    t.Fatalf("chosen = %q, want run", result.chosen)
  }
  if result.canceled {
    t.Fatal("auto-choose should not mark cancel")
  }
}

func TestLoadedMsgMultipleFillsList(t *testing.T) {
  state := newLoadingPickerState("docker ", nil)
  next, cmd := state.Update(loadedMsg{
    items: []suggest.Suggestion{
      {Label: "run", Value: "run", Detail: "Create a container"},
      {Label: "build", Value: "build", Detail: "Build an image"},
    },
  })
  result := next.(pickerState)
  if cmd != nil {
    t.Fatal("expected no quit cmd")
  }
  if result.loading {
    t.Fatal("expected loading false")
  }
  if len(result.filtered) != 2 {
    t.Fatalf("filtered = %+v", result.filtered)
  }
  view := result.View()
  if strings.Contains(view, "loading…") {
    t.Fatalf("view still loading: %q", view)
  }
  if !strings.Contains(view, "run") {
    t.Fatalf("view missing run: %q", view)
  }
}

func TestLoadedMsgEmptyQuitsWithoutCancel(t *testing.T) {
  state := newLoadingPickerState("docker ", nil)
  next, cmd := state.Update(loadedMsg{})
  result := next.(pickerState)
  if result.canceled {
    t.Fatal("empty load should not mark user cancel")
  }
  if cmd == nil {
    t.Fatal("expected quit cmd")
  }
  if result.chosen != "" {
    t.Fatalf("chosen = %q, want empty", result.chosen)
  }
}

func TestLoadingEscCancels(t *testing.T) {
  state := newLoadingPickerState("docker ", nil)
  next, cmd := state.Update(tea.KeyMsg{Type: tea.KeyEsc})
  result := next.(pickerState)
  if !result.canceled {
    t.Fatal("expected canceled")
  }
  if cmd == nil {
    t.Fatal("expected quit cmd")
  }
}

func TestLoadingIgnoresTyping(t *testing.T) {
  state := newLoadingPickerState("docker ", nil)
  next, _ := state.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")})
  result := next.(pickerState)
  if result.query != "" {
    t.Fatalf("query = %q, want empty while loading", result.query)
  }
  if !result.loading {
    t.Fatal("expected still loading")
  }
}

func TestLoadingViewShowsStatus(t *testing.T) {
  state := newLoadingPickerState("docker ", nil)
  view := state.View()
  if !strings.Contains(view, "loading…") {
    t.Fatalf("view = %q", view)
  }
}

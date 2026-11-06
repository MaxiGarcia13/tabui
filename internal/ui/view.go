package ui

import (
	"fmt"
	"strings"

	"github.com/MaxiGarcia13/tabui/internal/cmdline"
)

func (state pickerState) View() string {
  var view strings.Builder

  preview := state.linePrefix
  if !state.loading && len(state.filtered) > 0 {
    preview = cmdline.InsertToken(
      state.linePrefix,
      state.filtered[state.cursor].Value,
    )
  }
  view.WriteString(queryStyle.Render("→ " + preview))
  view.WriteString("\n")
  view.WriteString(queryStyle.Render("filter: " + state.query))
  view.WriteString("\n")

  if state.loading {
    view.WriteString(detailStyle.Render("loading…"))
    view.WriteString("\n")
    return view.String()
  }

  if len(state.filtered) == 0 {
    view.WriteString(detailStyle.Render("no matches"))
    view.WriteString("\n")
    return view.String()
  }

  capacity := state.visibleCapacity()
  start := state.offset
  end := start + capacity
  if end > len(state.filtered) {
    end = len(state.filtered)
  }

  for index := start; index < end; index++ {
    item := state.filtered[index]
    prefix := padGlyph
    label := normalStyle.Render(item.Label)
    if index == state.cursor {
      prefix = cursorGlyph
      label = selectedStyle.Render(item.Label)
    }
    if item.Detail == "" {
      fmt.Fprintf(&view, "%s%s\n", prefix, label)
      continue
    }
    detail := detailStyle.Render(item.Detail)
    fmt.Fprintf(&view, "%s%s  %s\n", prefix, label, detail)
  }

  status := fmt.Sprintf(
    "%d/%d",
    state.cursor+1,
    len(state.filtered),
  )
  if start > 0 || end < len(state.filtered) {
    status += "  ↑↓ scroll"
  }
  view.WriteString(detailStyle.Render(status))
  view.WriteString("\n")
  return view.String()
}

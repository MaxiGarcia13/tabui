package ui

import (
  "fmt"
  "strings"
)

func (state pickerState) View() string {
  var view strings.Builder
  view.WriteString(queryStyle.Render("filter: " + state.query))
  view.WriteString("\n")
  if len(state.filtered) == 0 {
    view.WriteString(detailStyle.Render("no matches"))
    view.WriteString("\n")
    return view.String()
  }
  for index, item := range state.filtered {
    prefix := padGlyph
    label := normalStyle.Render(item.Label)
    if index == state.cursor {
      prefix = cursorGlyph
      label = selectedStyle.Render(item.Label)
    }
    detail := detailStyle.Render(item.Detail)
    fmt.Fprintf(&view, "%s%s  %s\n", prefix, label, detail)
  }
  return view.String()
}

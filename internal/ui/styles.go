package ui

import "github.com/charmbracelet/lipgloss"

const (
  selectedColor = "212"
  mutedColor    = "243"
  queryColor    = "244"
)

var (
  selectedStyle = lipgloss.NewStyle().Bold(true).Foreground(
    lipgloss.Color(selectedColor),
  )
  normalStyle = lipgloss.NewStyle()
  detailStyle = lipgloss.NewStyle().Foreground(
    lipgloss.Color(mutedColor),
  )
  queryStyle = lipgloss.NewStyle().Foreground(
    lipgloss.Color(queryColor),
  )
  cursorGlyph = "> "
  padGlyph    = "  "
)

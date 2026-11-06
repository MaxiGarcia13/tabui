package clihelp

import (
  "regexp"
  "strings"

  "github.com/MaxiGarcia13/tabui/internal/suggest"
)

var (
  ansiEscapeRE   = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
  sectionRE      = regexp.MustCompile(`^[A-Za-z][^:]*:$`)
  commandRowRE   = regexp.MustCompile(`^\s{2,}([a-zA-Z][\w-]*)\s{2,}(\S.*)$`)
  commandNameRE  = regexp.MustCompile(`^[a-zA-Z][\w-]*$`)
  flagRowRE = regexp.MustCompile(
    `^\s*((?:-[a-zA-Z0-9],\s*)?` +
      `--?[a-zA-Z0-9][\w-]*` +
      `(?:,\s*--[a-zA-Z0-9][\w-]*)?` +
      `(?:\s+[A-Za-z][\w\[\]-]*)?)` +
      `\s{2,}(.*)$`,
  )
  flagTokenRE = regexp.MustCompile(`--?[a-zA-Z0-9][\w-]*`)
)

func Parse(helpText string) []suggest.Suggestion {
  text := sanitizeHelp(helpText)
  if text == "" {
    return nil
  }

  seen := make(map[string]struct{})
  suggestions := make([]suggest.Suggestion, 0)

  for _, line := range strings.Split(text, "\n") {
    trimmed := strings.TrimSpace(line)
    if trimmed == "" || sectionRE.MatchString(trimmed) {
      continue
    }
    if strings.HasPrefix(trimmed, "Usage:") ||
      strings.HasPrefix(trimmed, "usage:") {
      continue
    }

    if matches := flagRowRE.FindStringSubmatch(line); matches != nil {
      detail := strings.TrimSpace(matches[2])
      for _, token := range flagTokenRE.FindAllString(matches[1], -1) {
        addSuggestion(&suggestions, seen, token, detail)
      }
      continue
    }

    if matches := commandRowRE.FindStringSubmatch(line); matches != nil {
      name := matches[1]
      detail := strings.TrimSpace(matches[2])
      addSuggestion(&suggestions, seen, name, detail)
      continue
    }

    for _, name := range commaSeparatedCommands(line) {
      addSuggestion(&suggestions, seen, name, "")
    }
  }

  return suggestions
}

func commaSeparatedCommands(line string) []string {
  trimmed := strings.TrimSpace(line)
  trimmed = strings.TrimSuffix(trimmed, ",")
  if !strings.Contains(trimmed, ",") {
    return nil
  }
  parts := strings.Split(trimmed, ",")
  names := make([]string, 0, len(parts))
  for _, part := range parts {
    name := strings.TrimSpace(part)
    if name == "" {
      continue
    }
    if !commandNameRE.MatchString(name) {
      return nil
    }
    names = append(names, name)
  }
  if len(names) < 2 {
    return nil
  }
  return names
}

func sanitizeHelp(helpText string) string {
  text := strings.ReplaceAll(helpText, "\r\n", "\n")
  text = strings.ReplaceAll(text, "\r", "\n")
  text = ansiEscapeRE.ReplaceAllString(text, "")
  return text
}

func addSuggestion(
  suggestions *[]suggest.Suggestion,
  seen map[string]struct{},
  value string,
  detail string,
) {
  if value == "" {
    return
  }
  if _, exists := seen[value]; exists {
    return
  }
  seen[value] = struct{}{}
  *suggestions = append(*suggestions, suggest.Suggestion{
    Label:  value,
    Detail: detail,
    Value:  value,
  })
}

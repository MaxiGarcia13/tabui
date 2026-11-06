package clihelp

import "strings"

type target struct {
  ok       bool
  command  string
  helpArgs []string
  partial  string
}

func parseTarget(line string) target {
  line = strings.TrimLeft(line, " \t")
  if line == "" {
    return target{}
  }

  trailingSpace := strings.HasSuffix(line, " ") ||
    strings.HasSuffix(line, "\t")
  words := strings.Fields(line)
  if len(words) == 0 {
    return target{}
  }

  command := words[0]
  if !commandOnPath(command) {
    return target{}
  }

  // Still typing the command name — leave to the shell.
  if len(words) == 1 && !trailingSpace {
    return target{}
  }

  partial := ""
  complete := words
  if !trailingSpace {
    partial = words[len(words)-1]
    complete = words[:len(words)-1]
  }

  helpArgs := make([]string, 0, len(complete))
  for _, word := range complete[1:] {
    if strings.HasPrefix(word, "-") {
      continue
    }
    helpArgs = append(helpArgs, word)
  }

  return target{
    ok:       true,
    command:  command,
    helpArgs: helpArgs,
    partial:  partial,
  }
}

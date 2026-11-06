package clihelp

import (
  "bytes"
  "context"
  "os/exec"
  "time"
)

const helpTimeout = 2 * time.Second

var fetchHelp = fetchHelpOutput

func fetchHelpOutput(command string, helpArgs []string) (string, error) {
  ctx, cancel := context.WithTimeout(context.Background(), helpTimeout)
  defer cancel()

  args := append(append([]string{}, helpArgs...), "--help")
  cmd := exec.CommandContext(ctx, command, args...)
  var output bytes.Buffer
  cmd.Stdout = &output
  cmd.Stderr = &output
  err := cmd.Run()
  text := output.String()
  if text == "" && err != nil {
    return "", err
  }
  return text, nil
}

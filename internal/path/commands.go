package path

var pathCommands = map[string]struct{}{
  "mv":  {},
  "cp":  {},
  "cd":  {},
  "rm":  {},
  "cat": {},
}

var dirOnlyCommands = map[string]struct{}{
  "cd": {},
}

func isPathCommand(command string) bool {
  _, ok := pathCommands[command]
  return ok
}

func isDirOnlyCommand(command string) bool {
  _, ok := dirOnlyCommands[command]
  return ok
}

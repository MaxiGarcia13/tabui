package package_scripts

var packageTools = map[string]struct{}{
  "npm":  {},
  "pnpm": {},
  "yarn": {},
  "bun":  {},
}

func isPackageTool(command string) bool {
  _, ok := packageTools[command]
  return ok
}

func isRunSubcommand(word string) bool {
  return word == "run" || word == "run-script"
}

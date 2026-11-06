package clihelp

import (
  "os"
  "path/filepath"
  "testing"

  "github.com/MaxiGarcia13/tabui/internal/suggest"
)

func TestParseDockerHelp(t *testing.T) {
  text := readFixture(t, "docker_help.txt")
  items := Parse(text)
  if !containsValue(items, "run") {
    t.Fatalf("missing run in %+v", items)
  }
  if !containsValue(items, "compose") {
    t.Fatalf("missing compose in %+v", items)
  }
  if !containsValue(items, "--debug") {
    t.Fatalf("missing --debug in %+v", items)
  }
  if !containsValue(items, "-D") {
    t.Fatalf("missing -D in %+v", items)
  }
  run := findValue(items, "run")
  if run == nil || run.Detail == "" {
    t.Fatalf("run detail missing: %+v", run)
  }
}

func TestParseGitHelp(t *testing.T) {
  text := readFixture(t, "git_help.txt")
  items := Parse(text)
  for _, name := range []string{"clone", "init", "commit", "branch"} {
    if !containsValue(items, name) {
      t.Fatalf("missing %s in %+v", name, items)
    }
  }
}

func TestParseNpmHelpCommaList(t *testing.T) {
  text := readFixture(t, "npm_help.txt")
  items := Parse(text)
  for _, name := range []string{"install", "run-script", "publish", "exec"} {
    if !containsValue(items, name) {
      t.Fatalf("missing %s in %+v", name, items)
    }
  }
}

func TestParseStripsANSI(t *testing.T) {
  text := "\x1b[1mCommands:\x1b[0m\n  run         Create a container\n"
  items := Parse(text)
  if len(items) != 1 || items[0].Value != "run" {
    t.Fatalf("got %+v", items)
  }
}

func readFixture(t *testing.T, name string) string {
  t.Helper()
  bytes, err := os.ReadFile(filepath.Join("testdata", name))
  if err != nil {
    t.Fatal(err)
  }
  return string(bytes)
}

func containsValue(items []suggest.Suggestion, value string) bool {
  return findValue(items, value) != nil
}

func findValue(items []suggest.Suggestion, value string) *suggest.Suggestion {
  for index := range items {
    if items[index].Value == value {
      return &items[index]
    }
  }
  return nil
}

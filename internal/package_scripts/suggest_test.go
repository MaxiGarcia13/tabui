package package_scripts

import (
  "path/filepath"
  "testing"

  "github.com/MaxiGarcia13/tabui/internal/suggest"
)

func TestSuggestToolCommands(t *testing.T) {
  suggestions := Suggest(suggest.Context{
    Line:   "npm ",
    Cursor: len("npm "),
    Cwd:    t.TempDir(),
  })
  if len(suggestions) == 0 {
    t.Fatal("expected tool command suggestions")
  }
  if !containsLabel(suggestions, "run") {
    t.Fatalf("missing run in %+v", suggestions)
  }
  if !containsLabel(suggestions, "--help") {
    t.Fatalf("missing --help in %+v", suggestions)
  }
}

func TestSuggestToolCommandsFilter(t *testing.T) {
  suggestions := Suggest(suggest.Context{
    Line:   "pnpm in",
    Cursor: len("pnpm in"),
    Cwd:    t.TempDir(),
  })
  if len(suggestions) == 0 {
    t.Fatal("expected filtered commands")
  }
  for _, item := range suggestions {
    if item.Label != "install" && item.Label != "init" {
      t.Fatalf("unexpected label %q in %+v", item.Label, suggestions)
    }
  }
}

func TestSuggestScriptsFromNearestPackageJSON(t *testing.T) {
  workingDir := filepath.Join("testdata", "project")
  suggestions := Suggest(suggest.Context{
    Line:   "pnpm run ",
    Cursor: len("pnpm run "),
    Cwd:    workingDir,
  })
  if len(suggestions) != 4 {
    t.Fatalf("len = %d, want 4: %+v", len(suggestions), suggestions)
  }
  first := suggestions[0]
  if first.Label != "build" || first.Detail != "tsc -b" {
    t.Fatalf("first = %+v", first)
  }
  if first.Value != "build" {
    t.Fatalf("value = %q", first.Value)
  }
}

func TestSuggestScriptsWalksUpFromNested(t *testing.T) {
  workingDir := filepath.Join("testdata", "project", "nested")
  suggestions := Suggest(suggest.Context{
    Line:   "npm run ",
    Cursor: len("npm run "),
    Cwd:    workingDir,
  })
  if len(suggestions) != 4 {
    t.Fatalf("len = %d, want 4", len(suggestions))
  }
}

func TestSuggestScriptsFiltersPartial(t *testing.T) {
  workingDir := filepath.Join("testdata", "project")
  suggestions := Suggest(suggest.Context{
    Line:   "pnpm run te",
    Cursor: len("pnpm run te"),
    Cwd:    workingDir,
  })
  if len(suggestions) != 1 || suggestions[0].Label != "test" {
    t.Fatalf("got %+v", suggestions)
  }
}

func TestSuggestScriptsMissingPackageJSON(t *testing.T) {
  workingDir := t.TempDir()
  suggestions := Suggest(suggest.Context{
    Line:   "pnpm run ",
    Cursor: len("pnpm run "),
    Cwd:    workingDir,
  })
  if suggestions != nil {
    t.Fatalf("got %+v, want nil", suggestions)
  }
}

func TestSuggestBarePnpmIncludesScriptsWithRun(t *testing.T) {
  workingDir := filepath.Join("testdata", "project")
  suggestions := Suggest(suggest.Context{
    Line:   "pnpm ",
    Cursor: len("pnpm "),
    Cwd:    workingDir,
  })
  if !containsLabel(suggestions, "run") {
    t.Fatalf("bare pnpm should list run: %+v", suggestions)
  }
  build := findByLabel(suggestions, "build")
  if build == nil {
    t.Fatalf("bare pnpm should list scripts: %+v", suggestions)
  }
  if build.Value != "run build" {
    t.Fatalf("build value = %q, want %q", build.Value, "run build")
  }
}

func findByLabel(
  suggestions []suggest.Suggestion,
  label string,
) *suggest.Suggestion {
  for index := range suggestions {
    if suggestions[index].Label == label {
      return &suggestions[index]
    }
  }
  return nil
}

func TestSuggestNoMatch(t *testing.T) {
  workingDir := filepath.Join("testdata", "project")
  suggestions := Suggest(suggest.Context{
    Line:   "mv ",
    Cursor: len("mv "),
    Cwd:    workingDir,
  })
  if suggestions != nil {
    t.Fatalf("got %+v, want nil", suggestions)
  }
}

func containsLabel(
  suggestions []suggest.Suggestion,
  label string,
) bool {
  for _, item := range suggestions {
    if item.Label == label {
      return true
    }
  }
  return false
}

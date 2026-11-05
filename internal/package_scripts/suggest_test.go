package package_scripts

import (
  "path/filepath"
  "testing"

  "github.com/MaxiGarcia13/tabui/internal/suggest"
)

func TestSuggestFromNearestPackageJSON(t *testing.T) {
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

func TestSuggestWalksUpFromNested(t *testing.T) {
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

func TestSuggestFiltersPartial(t *testing.T) {
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

func TestSuggestMissingPackageJSON(t *testing.T) {
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

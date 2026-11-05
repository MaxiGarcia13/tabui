package package_scripts

import (
  "path/filepath"
  "testing"

  "github.com/MaxiGarcia13/tabui/internal/suggest"
)

func TestSuggestFromNearestPackageJSON(t *testing.T) {
  cwd := filepath.Join("testdata", "project")
  got := Suggest(suggest.Context{
    Line:   "pnpm run ",
    Cursor: len("pnpm run "),
    Cwd:    cwd,
  })
  if len(got) != 4 {
    t.Fatalf("len = %d, want 4: %+v", len(got), got)
  }
  if got[0].Label != "build" || got[0].Detail != "tsc -b" {
    t.Fatalf("first = %+v", got[0])
  }
  if got[0].Value != "build" {
    t.Fatalf("value = %q", got[0].Value)
  }
}

func TestSuggestWalksUpFromNested(t *testing.T) {
  cwd := filepath.Join("testdata", "project", "nested")
  got := Suggest(suggest.Context{
    Line:   "npm run ",
    Cursor: len("npm run "),
    Cwd:    cwd,
  })
  if len(got) != 4 {
    t.Fatalf("len = %d, want 4", len(got))
  }
}

func TestSuggestFiltersPartial(t *testing.T) {
  cwd := filepath.Join("testdata", "project")
  got := Suggest(suggest.Context{
    Line:   "pnpm run te",
    Cursor: len("pnpm run te"),
    Cwd:    cwd,
  })
  if len(got) != 1 || got[0].Label != "test" {
    t.Fatalf("got %+v", got)
  }
}

func TestSuggestMissingPackageJSON(t *testing.T) {
  cwd := t.TempDir()
  got := Suggest(suggest.Context{
    Line:   "pnpm run ",
    Cursor: len("pnpm run "),
    Cwd:    cwd,
  })
  if got != nil {
    t.Fatalf("got %+v, want nil", got)
  }
}

func TestSuggestNoMatch(t *testing.T) {
  cwd := filepath.Join("testdata", "project")
  got := Suggest(suggest.Context{
    Line:   "mv ",
    Cursor: len("mv "),
    Cwd:    cwd,
  })
  if got != nil {
    t.Fatalf("got %+v, want nil", got)
  }
}

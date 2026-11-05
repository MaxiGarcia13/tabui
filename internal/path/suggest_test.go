package path

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MaxiGarcia13/tabui/internal/suggest"
)

func fixtureRoot(t *testing.T) string {
	t.Helper()
	workingDir, err := filepath.Abs(filepath.Join("testdata", "root"))
	if err != nil {
		t.Fatal(err)
	}
	return workingDir
}

func TestSuggestListsCwd(t *testing.T) {
	workingDir := fixtureRoot(t)
	suggestions := Suggest(suggest.Context{
		Line:   "mv ",
		Cursor: len("mv "),
		Cwd:    workingDir,
	})
	if !containsLabel(suggestions, "dir/") {
		t.Fatalf("missing dir/ in %+v", suggestions)
	}
	if !containsLabel(suggestions, "file.txt") {
		t.Fatalf("missing file.txt in %+v", suggestions)
	}
	if containsLabel(suggestions, ".hidden") {
		t.Fatalf("unexpected .hidden in %+v", suggestions)
	}
	if suggestions[0].Label != "dir/" {
		t.Fatalf("dirs should sort first: %+v", suggestions)
	}
}

func TestSuggestFiltersPartial(t *testing.T) {
	workingDir := fixtureRoot(t)
	suggestions := Suggest(suggest.Context{
		Line:   "mv di",
		Cursor: len("mv di"),
		Cwd:    workingDir,
	})
	if len(suggestions) != 1 || suggestions[0].Value != "dir/" {
		t.Fatalf("got %+v", suggestions)
	}
}

func TestSuggestInsideDirectory(t *testing.T) {
	workingDir := fixtureRoot(t)
	suggestions := Suggest(suggest.Context{
		Line:   "cat dir/",
		Cursor: len("cat dir/"),
		Cwd:    workingDir,
	})
	if len(suggestions) != 1 {
		t.Fatalf("len = %d, want 1: %+v", len(suggestions), suggestions)
	}
	item := suggestions[0]
	if item.Label != "nested.txt" || item.Value != "dir/nested.txt" {
		t.Fatalf("got %+v", item)
	}
	if item.Detail != "dir/nested.txt" {
		t.Fatalf("detail = %q", item.Detail)
	}
}

func TestSuggestOptions(t *testing.T) {
	suggestions := Suggest(suggest.Context{
		Line:   "cd --",
		Cursor: len("cd --"),
		Cwd:    t.TempDir(),
	})
	if len(suggestions) != 1 {
		t.Fatalf("len = %d, want 1: %+v", len(suggestions), suggestions)
	}
	item := suggestions[0]
	if item.Label != "--" || item.Value != "--" {
		t.Fatalf("got %+v", item)
	}
	if item.Detail == "" {
		t.Fatal("expected option detail")
	}
}

func TestSuggestOptionsPartial(t *testing.T) {
	suggestions := Suggest(suggest.Context{
		Line:   "cp -r",
		Cursor: len("cp -r"),
		Cwd:    t.TempDir(),
	})
	if !containsLabel(suggestions, "-r") {
		t.Fatalf("missing -r in %+v", suggestions)
	}
	if !containsLabel(suggestions, "-R") {
		t.Fatalf("missing -R in %+v", suggestions)
	}
	for _, item := range suggestions {
		if item.Detail == "" {
			t.Fatalf("expected detail for %+v", item)
		}
	}
}

func TestSuggestCdDirsOnly(t *testing.T) {
	workingDir := fixtureRoot(t)
	suggestions := Suggest(suggest.Context{
		Line:   "cd ",
		Cursor: len("cd "),
		Cwd:    workingDir,
	})
	if !containsLabel(suggestions, "dir/") {
		t.Fatalf("missing dir/ in %+v", suggestions)
	}
	if containsLabel(suggestions, "file.txt") {
		t.Fatalf("unexpected file.txt in %+v", suggestions)
	}
	for _, item := range suggestions {
		if item.Label == "dir/" && item.Detail != "" {
			t.Fatalf("cwd dir detail should be empty, got %+v", item)
		}
	}
}

func TestSuggestShowsHiddenWhenPartialHasDot(t *testing.T) {
	workingDir := fixtureRoot(t)
	suggestions := Suggest(suggest.Context{
		Line:   "mv .",
		Cursor: len("mv ."),
		Cwd:    workingDir,
	})
	if !containsLabel(suggestions, ".hidden") {
		t.Fatalf("missing .hidden in %+v", suggestions)
	}
}

func TestSuggestAbsolutePath(t *testing.T) {
	workingDir := fixtureRoot(t)
	dirPart := workingDir + "/"
	line := "cat " + dirPart
	suggestions := Suggest(suggest.Context{
		Line:   line,
		Cursor: len(line),
		Cwd:    t.TempDir(),
	})
	if !containsValue(suggestions, dirPart+"dir/") {
		t.Fatalf("missing %q in %+v", dirPart+"dir/", suggestions)
	}
}

func TestSuggestHome(t *testing.T) {
	home := t.TempDir()
	if err := os.WriteFile(
		filepath.Join(home, "notes.txt"),
		[]byte("x"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	previous := lookUpHome
	lookUpHome = func() (string, error) { return home, nil }
	t.Cleanup(func() { lookUpHome = previous })

	suggestions := Suggest(suggest.Context{
		Line:   "cat ~/",
		Cursor: len("cat ~/"),
		Cwd:    t.TempDir(),
	})
	if len(suggestions) != 1 {
		t.Fatalf("len = %d, want 1: %+v", len(suggestions), suggestions)
	}
	if suggestions[0].Value != "~/notes.txt" {
		t.Fatalf("value = %q", suggestions[0].Value)
	}
}

func TestSplitPartial(t *testing.T) {
	tests := []struct {
		partial  string
		wantDir  string
		wantName string
	}{
		{partial: "", wantDir: "", wantName: ""},
		{partial: "fi", wantDir: "", wantName: "fi"},
		{partial: "dir/", wantDir: "dir/", wantName: ""},
		{partial: "dir/n", wantDir: "dir/", wantName: "n"},
		{partial: "./fi", wantDir: "./", wantName: "fi"},
		{partial: "~/Do", wantDir: "~/", wantName: "Do"},
	}
	for _, testCase := range tests {
		dirPart, namePrefix := splitPartial(testCase.partial)
		if dirPart != testCase.wantDir || namePrefix != testCase.wantName {
			t.Fatalf(
				"splitPartial(%q) = %q %q, want %q %q",
				testCase.partial,
				dirPart,
				namePrefix,
				testCase.wantDir,
				testCase.wantName,
			)
		}
	}
}

func TestExpandHome(t *testing.T) {
	home := t.TempDir()
	previous := lookUpHome
	lookUpHome = func() (string, error) { return home, nil }
	t.Cleanup(func() { lookUpHome = previous })

	got, err := expandHome("~/Documents")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(home, "Documents")
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func containsLabel(items []suggest.Suggestion, label string) bool {
	for _, item := range items {
		if item.Label == label {
			return true
		}
	}
	return false
}

func containsValue(items []suggest.Suggestion, value string) bool {
	for _, item := range items {
		if item.Value == value {
			return true
		}
	}
	return false
}

package clihelp

import (
  "errors"
  "os"
  "path/filepath"
  "testing"

  "github.com/MaxiGarcia13/tabui/internal/suggest"
)

func TestLookupCacheMissThenHit(t *testing.T) {
  dir := t.TempDir()
  bin := filepath.Join(dir, "fakectl")
  if err := os.WriteFile(bin, []byte("#!/bin/sh\n"), 0o755); err != nil {
    t.Fatal(err)
  }

  previousLook := lookPath
  previousFetch := fetchHelp
  previousDir := cacheDirFn
  lookPath = func(file string) (string, error) {
    if file == "fakectl" {
      return bin, nil
    }
    return "", errors.New("not found")
  }
  fetchHelp = func(command string, helpArgs []string) (string, error) {
    return readFixture(t, "docker_help.txt"), nil
  }
  cacheDirFn = func() (string, error) { return dir, nil }
  useStoreForTest(openStore())
  t.Cleanup(func() {
    lookPath = previousLook
    fetchHelp = previousFetch
    cacheDirFn = previousDir
    useStoreForTest(noopStore{})
  })

  items, matched, load := Lookup(suggest.Context{
    Line:   "fakectl ",
    Cursor: len("fakectl "),
  })
  if !matched {
    t.Fatal("expected match")
  }
  if items != nil || load == nil {
    t.Fatalf(
      "want cache miss with load, items=%v load=%v",
      items,
      load != nil,
    )
  }

  loaded := load()
  if !containsValue(loaded, "run") {
    t.Fatalf("loaded missing run: %+v", loaded)
  }

  items, matched, load = Lookup(suggest.Context{
    Line:   "fakectl ru",
    Cursor: len("fakectl ru"),
  })
  if !matched || load != nil {
    t.Fatalf("expected cache hit, matched=%v load=%v", matched, load != nil)
  }
  if len(items) == 0 || items[0].Value != "run" {
    t.Fatalf("filtered hit = %+v", items)
  }
}

func TestLookupNoMatch(t *testing.T) {
  previous := lookPath
  lookPath = func(string) (string, error) { return "", errors.New("missing") }
  t.Cleanup(func() { lookPath = previous })

  _, matched, _ := Lookup(suggest.Context{Line: "nope ", Cursor: 5})
  if matched {
    t.Fatal("expected no match")
  }
}

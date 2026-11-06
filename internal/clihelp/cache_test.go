package clihelp

import (
  "os"
  "path/filepath"
  "testing"

  "github.com/MaxiGarcia13/tabui/internal/suggest"
)

func TestSQLiteStoreRoundTrip(t *testing.T) {
  dir := t.TempDir()
  previousDir := cacheDirFn
  cacheDirFn = func() (string, error) { return dir, nil }
  t.Cleanup(func() {
    cacheDirFn = previousDir
    useStoreForTest(noopStore{})
  })

  store := openStore()
  useStoreForTest(store)
  if _, ok := store.(*sqliteStore); !ok {
    t.Fatalf("expected sqlite store, got %T", store)
  }

  items := []suggest.Suggestion{
    {Label: "run", Detail: "run a container", Value: "run"},
    {Label: "ps", Detail: "list containers", Value: "ps"},
  }
  store.Put("key-1", items)

  got, ok := store.Get("key-1")
  if !ok {
    t.Fatal("expected cache hit")
  }
  if len(got) != 2 {
    t.Fatalf("len = %d, want 2: %+v", len(got), got)
  }
  if _, err := os.Stat(filepath.Join(dir, "help.db")); err != nil {
    t.Fatalf("db file: %v", err)
  }
}

func TestNoopStoreSkipsPersist(t *testing.T) {
  useStoreForTest(noopStore{})
  t.Cleanup(func() { useStoreForTest(noopStore{}) })

  store := defaultStore()
  store.Put("key", []suggest.Suggestion{
    {Label: "run", Value: "run"},
  })
  if _, ok := store.Get("key"); ok {
    t.Fatal("noop store should not persist")
  }
}

func TestOpenStoreFallsBackWhenDirUnusable(t *testing.T) {
  previousDir := cacheDirFn
  cacheDirFn = func() (string, error) {
    return filepath.Join(t.TempDir(), "missing", "nested"), nil
  }
  // Parent is a file so MkdirAll fails.
  root := t.TempDir()
  blocker := filepath.Join(root, "blocker")
  if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
    t.Fatal(err)
  }
  cacheDirFn = func() (string, error) {
    return filepath.Join(blocker, "tabui"), nil
  }
  t.Cleanup(func() { cacheDirFn = previousDir })

  store := openStore()
  if _, ok := store.(noopStore); !ok {
    t.Fatalf("expected noop store, got %T", store)
  }
}

package clihelp

import (
  "crypto/sha256"
  "database/sql"
  "encoding/hex"
  "fmt"
  "os"
  "path/filepath"
  "sync"
  "time"

  "github.com/MaxiGarcia13/tabui/internal/suggest"

  _ "modernc.org/sqlite" // pure-Go SQLite driver for database/sql
)

type cacheStore interface {
  Get(key string) ([]suggest.Suggestion, bool)
  Put(key string, items []suggest.Suggestion)
}

type noopStore struct{}

func (noopStore) Get(string) ([]suggest.Suggestion, bool) { return nil, false }
func (noopStore) Put(string, []suggest.Suggestion)        {}

type sqliteStore struct {
  db *sql.DB
}

var (
  storeOnce sync.Once
  globalStore cacheStore = noopStore{}
  cacheDirFn = defaultCacheDir
)

func defaultCacheDir() (string, error) {
  root, err := os.UserCacheDir()
  if err != nil {
    return "", err
  }
  return filepath.Join(root, "tabui"), nil
}

func defaultStore() cacheStore {
  storeOnce.Do(func() {
    globalStore = openStore()
  })
  return globalStore
}

func openStore() cacheStore {
  dir, err := cacheDirFn()
  if err != nil {
    return noopStore{}
  }
  if err := os.MkdirAll(dir, 0o755); err != nil {
    return noopStore{}
  }

  db, err := sql.Open("sqlite", filepath.Join(dir, "help.db"))
  if err != nil {
    return noopStore{}
  }
  if err := db.Ping(); err != nil {
    _ = db.Close()
    return noopStore{}
  }

  _, err = db.Exec(`
CREATE TABLE IF NOT EXISTS help_entries (
  command_key TEXT NOT NULL,
  label TEXT NOT NULL,
  detail TEXT NOT NULL,
  value TEXT NOT NULL,
  fetched_at INTEGER NOT NULL,
  PRIMARY KEY (command_key, value)
);
`)
  if err != nil {
    _ = db.Close()
    return noopStore{}
  }

  return &sqliteStore{db: db}
}

func (store *sqliteStore) Get(key string) ([]suggest.Suggestion, bool) {
  rows, err := store.db.Query(
    `SELECT label, detail, value
     FROM help_entries
     WHERE command_key = ?
     ORDER BY value`,
    key,
  )
  if err != nil {
    return nil, false
  }
  defer rows.Close()

  items := make([]suggest.Suggestion, 0)
  for rows.Next() {
    var item suggest.Suggestion
    if err := rows.Scan(&item.Label, &item.Detail, &item.Value); err != nil {
      return nil, false
    }
    items = append(items, item)
  }
  if err := rows.Err(); err != nil {
    return nil, false
  }
  if len(items) == 0 {
    return nil, false
  }
  return items, true
}

func (store *sqliteStore) Put(key string, items []suggest.Suggestion) {
  if len(items) == 0 {
    return
  }
  tx, err := store.db.Begin()
  if err != nil {
    return
  }
  defer func() { _ = tx.Rollback() }()

  if _, err := tx.Exec(
    `DELETE FROM help_entries WHERE command_key = ?`,
    key,
  ); err != nil {
    return
  }

  statement, err := tx.Prepare(`
INSERT INTO help_entries (command_key, label, detail, value, fetched_at)
VALUES (?, ?, ?, ?, ?)
`)
  if err != nil {
    return
  }
  defer statement.Close()

  now := time.Now().Unix()
  for _, item := range items {
    if _, err := statement.Exec(
      key,
      item.Label,
      item.Detail,
      item.Value,
      now,
    ); err != nil {
      return
    }
  }
  _ = tx.Commit()
}

func cacheKey(commandPath string, helpArgs []string, modTime time.Time) string {
  hash := sha256.New()
  _, _ = fmt.Fprintf(hash, "%s\n", commandPath)
  _, _ = fmt.Fprintf(hash, "%d\n", modTime.UnixNano())
  for _, arg := range helpArgs {
    _, _ = fmt.Fprintf(hash, "%s\n", arg)
  }
  return hex.EncodeToString(hash.Sum(nil))
}

func useStoreForTest(store cacheStore) {
  storeOnce = sync.Once{}
  if store == nil {
    store = noopStore{}
  }
  globalStore = store
  storeOnce.Do(func() {})
}

package store_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/dharundds/jobtrack/internal/store"
)

func TestStore_OpenAndMigrate(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	st, err := store.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	if err := store.RunMigrations(st.DB, filepath.Join("..", "..", "migrations")); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	if err := st.Ping(ctx); err != nil {
		t.Fatal(err)
	}

	var name string
	err = st.DB.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='accounts'`).Scan(&name)
	if err != nil {
		t.Fatal(err)
	}
}

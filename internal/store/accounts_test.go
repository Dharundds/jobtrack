package store_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/dharundds/jobtrack/internal/crypto"
	"github.com/dharundds/jobtrack/internal/store"
)

func TestAccountUpsertAndDelete(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	st, err := store.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	migrations := filepath.Join("..", "..", "migrations")
	if err := store.RunMigrations(st.DB, migrations); err != nil {
		t.Fatal(err)
	}

	key := make([]byte, 32)
	enc, err := crypto.EncryptToken(key, "refresh-abc")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	acc, err := st.UpsertAccount(ctx, "user@gmail.com", enc)
	if err != nil {
		t.Fatal(err)
	}
	if acc.Email != "user@gmail.com" {
		t.Fatalf("email: %q", acc.Email)
	}

	got, err := st.GetAccount(ctx)
	if err != nil {
		t.Fatal(err)
	}

	plain, err := crypto.DecryptToken(key, got.EncRefreshToken)
	if err != nil || plain != "refresh-abc" {
		t.Fatalf("decrypt: %v plain=%q", err, plain)
	}

	if err := st.DeleteAllAccounts(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := st.GetAccount(ctx); err != store.ErrNoAccount {
		t.Fatalf("expected ErrNoAccount, got %v", err)
	}
}

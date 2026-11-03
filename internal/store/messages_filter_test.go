package store_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/dharundds/jobtrack/internal/store"
)

func TestFilterNewGmailIDs(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	if err := store.RunMigrations(st.DB, filepath.Join("..", "..", "migrations")); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	_, err = st.InsertMessage(ctx, store.MessageInsert{
		GmailID:    "id-a",
		ThreadID:   "t1",
		Direction:  "in",
		ReceivedAt: time.Now(),
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = st.InsertMessage(ctx, store.MessageInsert{
		GmailID:    "id-b",
		ThreadID:   "t1",
		Direction:  "in",
		ReceivedAt: time.Now(),
	})
	if err != nil {
		t.Fatal(err)
	}

	input := []string{"id-a", "id-b", "id-c", "id-d"}
	newIDs, skipped, err := st.FilterNewGmailIDs(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	if skipped != 2 {
		t.Fatalf("skipped = %d, want 2", skipped)
	}
	if len(newIDs) != 2 || newIDs[0] != "id-c" || newIDs[1] != "id-d" {
		t.Fatalf("newIDs = %v", newIDs)
	}

	newIDs2, skipped2, err := st.FilterNewGmailIDs(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	if skipped2 != 2 || len(newIDs2) != 2 {
		t.Fatalf("second pass skipped=%d new=%v", skipped2, newIDs2)
	}

	_, err = st.InsertMessage(ctx, store.MessageInsert{GmailID: "id-c", ThreadID: "t2", Direction: "in", ReceivedAt: time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	_, skipped3, err := st.FilterNewGmailIDs(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	if skipped3 != 3 {
		t.Fatalf("third pass skipped=%d, want 3", skipped3)
	}
}

package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type MessageInsert struct {
	GmailID     string
	ThreadID    string
	Direction   string
	FromAddr    string
	ToAddr      string
	Subject     string
	Snippet     string
	BodyExcerpt string
	ReceivedAt  time.Time
}

func (s *Store) HasMessage(ctx context.Context, gmailID string) (bool, error) {
	var id int64
	err := s.DB.QueryRowContext(ctx, `SELECT id FROM messages WHERE gmail_id = ?`, gmailID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Store) InsertMessage(ctx context.Context, m MessageInsert) (int64, error) {
	res, err := s.DB.ExecContext(ctx, `
		INSERT INTO messages (
			gmail_id, thread_id, direction, from_addr, to_addr, subject,
			snippet, body_excerpt, received_at, event_type, classified_by
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 'other', 'sync')
	`, m.GmailID, m.ThreadID, m.Direction, m.FromAddr, m.ToAddr, m.Subject,
		m.Snippet, m.BodyExcerpt, m.ReceivedAt.UTC().Format(time.RFC3339))
	if err != nil {
		return 0, fmt.Errorf("insert message: %w", err)
	}
	return res.LastInsertId()
}

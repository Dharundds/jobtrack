package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dharundds/jobtrack/internal/config"
)

func (s *Store) GetSetting(ctx context.Context, key string) (string, error) {
	var value string
	err := s.DB.QueryRowContext(ctx, `SELECT value FROM settings WHERE key = ?`, key).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return value, nil
}

func (s *Store) GmailQuery(ctx context.Context) (string, error) {
	q, err := s.GetSetting(ctx, "gmail_query")
	if err != nil || q == "" {
		return config.DefaultGmailQuery, nil
	}
	return q, nil
}

func (s *Store) UpdateAccountSync(ctx context.Context, historyID string, syncedAt string) error {
	_, err := s.DB.ExecContext(ctx, `
		UPDATE accounts SET history_id = ?, last_synced_at = ?
		WHERE id = (SELECT id FROM accounts ORDER BY id ASC LIMIT 1)
	`, historyID, syncedAt)
	if err != nil {
		return fmt.Errorf("update account sync: %w", err)
	}
	return nil
}

package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var ErrNoAccount = errors.New("no account configured")

type Account struct {
	ID              int64
	Email           string
	EncRefreshToken string
	HistoryID       sql.NullString
	LastSyncedAt    sql.NullTime
}

func (s *Store) UpsertAccount(ctx context.Context, email, encRefreshToken string) (*Account, error) {
	res, err := s.DB.ExecContext(ctx, `
		INSERT INTO accounts (email, enc_refresh_token)
		VALUES (?, ?)
		ON CONFLICT(email) DO UPDATE SET enc_refresh_token = excluded.enc_refresh_token
	`, email, encRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("upsert account: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil || id == 0 {
		var accID int64
		err = s.DB.QueryRowContext(ctx, `SELECT id FROM accounts WHERE email = ?`, email).Scan(&accID)
		if err != nil {
			return nil, fmt.Errorf("lookup account id: %w", err)
		}
		id = accID
	}
	return s.GetAccountByID(ctx, id)
}

func (s *Store) GetAccount(ctx context.Context) (*Account, error) {
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, email, enc_refresh_token, history_id, last_synced_at
		FROM accounts
		ORDER BY id ASC
		LIMIT 1
	`)
	return scanAccount(row)
}

func (s *Store) GetAccountByID(ctx context.Context, id int64) (*Account, error) {
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, email, enc_refresh_token, history_id, last_synced_at
		FROM accounts WHERE id = ?
	`, id)
	return scanAccount(row)
}

func (s *Store) DeleteAllAccounts(ctx context.Context) error {
	_, err := s.DB.ExecContext(ctx, `DELETE FROM accounts`)
	return err
}

// ClearSyncHistory removes the Gmail history cursor so the next sync re-lists mail via the search query.
func (s *Store) ClearSyncHistory(ctx context.Context) error {
	_, err := s.DB.ExecContext(ctx, `
		UPDATE accounts SET history_id = NULL
		WHERE id = (SELECT id FROM accounts ORDER BY id ASC LIMIT 1)
	`)
	if err != nil {
		return fmt.Errorf("clear sync history: %w", err)
	}
	return nil
}

func scanAccount(row *sql.Row) (*Account, error) {
	var a Account
	var historyID, lastSynced sql.NullString
	err := row.Scan(&a.ID, &a.Email, &a.EncRefreshToken, &historyID, &lastSynced)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoAccount
	}
	if err != nil {
		return nil, err
	}
	a.HistoryID = historyID
	if lastSynced.Valid {
		t, parseErr := time.Parse(time.RFC3339, lastSynced.String)
		if parseErr == nil {
			a.LastSyncedAt = sql.NullTime{Time: t, Valid: true}
		}
	}
	return &a, nil
}

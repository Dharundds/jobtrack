package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (s *Store) GetMessage(ctx context.Context, messageID int64) (*MessageRow, error) {
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, gmail_id, thread_id, application_id, direction, from_addr, to_addr,
			subject, snippet, body_excerpt, received_at, event_type, classified_by
		FROM messages WHERE id = ? AND ignored = 0
	`, messageID)
	var m MessageRow
	var appID sql.NullInt64
	var received string
	err := row.Scan(
		&m.ID, &m.GmailID, &m.ThreadID, &appID, &m.Direction, &m.FromAddr, &m.ToAddr,
		&m.Subject, &m.Snippet, &m.BodyExcerpt, &received, &m.EventType, &m.ClassifiedBy,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrMessageNotFound
	}
	if err != nil {
		return nil, err
	}
	m.ApplicationID = appID
	t, _ := time.Parse(time.RFC3339, received)
	m.ReceivedAt = t
	return &m, nil
}

func (s *Store) CountMessagesForApplication(ctx context.Context, applicationID int64) (int, error) {
	var n int
	err := s.DB.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM messages WHERE application_id = ? AND ignored = 0
	`, applicationID).Scan(&n)
	return n, err
}

// IgnoreMessage removes a message from tracking (won't be re-classified).
func (s *Store) IgnoreMessage(ctx context.Context, applicationID, messageID int64) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var appID sql.NullInt64
	if err := tx.QueryRowContext(ctx, `SELECT application_id FROM messages WHERE id = ? AND ignored = 0`, messageID).Scan(&appID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrMessageNotFound
		}
		return err
	}
	if !appID.Valid || appID.Int64 != applicationID {
		return ErrMessageNotFound
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM events WHERE message_id = ?`, messageID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE messages SET application_id = NULL, ignored = 1, classified_by = 'manual_ignore'
		WHERE id = ?
	`, messageID); err != nil {
		return err
	}
	if err := s.deleteApplicationIfEmptyTx(ctx, tx, applicationID); err != nil {
		return err
	}
	return tx.Commit()
}

// DeleteApplication removes a listing and ignores its messages so sync won't re-group them.
func (s *Store) DeleteApplication(ctx context.Context, applicationID int64) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var exists int
	if err := tx.QueryRowContext(ctx, `SELECT 1 FROM applications WHERE id = ?`, applicationID).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrApplicationNotFound
		}
		return err
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE messages SET application_id = NULL, ignored = 1, classified_by = 'manual_ignore'
		WHERE application_id = ?
	`, applicationID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM applications WHERE id = ?`, applicationID); err != nil {
		return err
	}
	return tx.Commit()
}

// SplitMessageToApplication moves one message (and its events) onto a new application row.
func (s *Store) SplitMessageToApplication(
	ctx context.Context,
	sourceAppID, messageID int64,
	company, role, status string,
	appliedAt *time.Time,
) (int64, error) {
	if company == "" {
		return 0, fmt.Errorf("company is required")
	}
	if status == "" {
		status = "applied"
	}

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	var appID sql.NullInt64
	var received string
	var eventType string
	if err := tx.QueryRowContext(ctx, `
		SELECT application_id, received_at, event_type FROM messages
		WHERE id = ? AND ignored = 0
	`, messageID).Scan(&appID, &received, &eventType); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrMessageNotFound
		}
		return 0, err
	}
	if !appID.Valid || appID.Int64 != sourceAppID {
		return 0, ErrMessageNotFound
	}

	now := time.Now().UTC().Format(time.RFC3339)
	var applied sql.NullString
	if appliedAt != nil {
		applied = sql.NullString{String: appliedAt.UTC().Format(time.RFC3339), Valid: true}
	} else if eventType == "applied" {
		if t, parseErr := time.Parse(time.RFC3339, received); parseErr == nil {
			applied = sql.NullString{String: t.UTC().Format(time.RFC3339), Valid: true}
		}
	}

	res, err := tx.ExecContext(ctx, `
		INSERT INTO applications (company, role, status, source, applied_at, created_at, updated_at)
		VALUES (?, ?, ?, 'email', ?, ?, ?)
	`, company, role, status, applied, now, now)
	if err != nil {
		return 0, err
	}
	newID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if _, err := tx.ExecContext(ctx, `UPDATE messages SET application_id = ? WHERE id = ?`, newID, messageID); err != nil {
		return 0, err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE events SET application_id = ? WHERE message_id = ?`, newID, messageID); err != nil {
		return 0, err
	}

	if err := s.deleteApplicationIfEmptyTx(ctx, tx, sourceAppID); err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return newID, nil
}

func (s *Store) deleteApplicationIfEmptyTx(ctx context.Context, tx *sql.Tx, applicationID int64) error {
	var n int
	if err := tx.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM messages WHERE application_id = ? AND ignored = 0
	`, applicationID).Scan(&n); err != nil {
		return err
	}
	if n > 0 {
		now := time.Now().UTC().Format(time.RFC3339)
		_, err := tx.ExecContext(ctx, `UPDATE applications SET updated_at = ? WHERE id = ?`, now, applicationID)
		return err
	}
	_, err := tx.ExecContext(ctx, `DELETE FROM applications WHERE id = ?`, applicationID)
	return err
}

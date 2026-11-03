package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type RuleRow struct {
	ID        int64
	Field     string
	Pattern   string
	EventType string
	Priority  int
	Enabled   bool
}

func (s *Store) ListEnabledRules(ctx context.Context) ([]RuleRow, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, field, pattern, event_type, priority, enabled
		FROM rules WHERE enabled = 1
		ORDER BY priority DESC, id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []RuleRow
	for rows.Next() {
		var r RuleRow
		var enabled int
		if err := rows.Scan(&r.ID, &r.Field, &r.Pattern, &r.EventType, &r.Priority, &enabled); err != nil {
			return nil, err
		}
		r.Enabled = enabled == 1
		out = append(out, r)
	}
	return out, rows.Err()
}

type MessageRow struct {
	ID             int64
	GmailID        string
	ThreadID       string
	ApplicationID  sql.NullInt64
	Direction      string
	FromAddr       string
	ToAddr         string
	Subject        string
	Snippet        string
	BodyExcerpt    string
	ReceivedAt     time.Time
	EventType      string
	ClassifiedBy   string
}

func (s *Store) ListUnlinkedMessages(ctx context.Context) ([]MessageRow, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, gmail_id, thread_id, application_id, direction, from_addr, to_addr,
			subject, snippet, body_excerpt, received_at, event_type, classified_by
		FROM messages
		WHERE application_id IS NULL AND ignored = 0
		ORDER BY received_at ASC, id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []MessageRow
	for rows.Next() {
		var m MessageRow
		var appID sql.NullInt64
		var received string
		if err := rows.Scan(
			&m.ID, &m.GmailID, &m.ThreadID, &appID, &m.Direction, &m.FromAddr, &m.ToAddr,
			&m.Subject, &m.Snippet, &m.BodyExcerpt, &received, &m.EventType, &m.ClassifiedBy,
		); err != nil {
			return nil, err
		}
		m.ApplicationID = appID
		t, err := time.Parse(time.RFC3339, received)
		if err != nil {
			t = time.Now().UTC()
		}
		m.ReceivedAt = t
		out = append(out, m)
	}
	return out, rows.Err()
}

func (s *Store) ApplicationIDByThread(ctx context.Context, threadID string) (int64, bool, error) {
	if threadID == "" {
		return 0, false, nil
	}
	var id int64
	err := s.DB.QueryRowContext(ctx, `
		SELECT application_id FROM messages
		WHERE thread_id = ? AND application_id IS NOT NULL
		LIMIT 1
	`, threadID).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return id, true, nil
}

type ApplicationRow struct {
	ID                int64
	Company           string
	Role              string
	Status            string
	StatusOverridden  bool
	AppliedAt         sql.NullString
}

func (s *Store) ListApplications(ctx context.Context) ([]ApplicationRow, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, company, role, status, status_overridden, applied_at
		FROM applications ORDER BY id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ApplicationRow
	for rows.Next() {
		var a ApplicationRow
		var overridden int
		if err := rows.Scan(&a.ID, &a.Company, &a.Role, &a.Status, &overridden, &a.AppliedAt); err != nil {
			return nil, err
		}
		a.StatusOverridden = overridden == 1
		out = append(out, a)
	}
	return out, rows.Err()
}

func (s *Store) CreateApplication(ctx context.Context, company, role, status string, appliedAt *time.Time) (int64, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	var applied sql.NullString
	if appliedAt != nil {
		applied = sql.NullString{String: appliedAt.UTC().Format(time.RFC3339), Valid: true}
	}
	res, err := s.DB.ExecContext(ctx, `
		INSERT INTO applications (company, role, status, source, applied_at, created_at, updated_at)
		VALUES (?, ?, ?, 'email', ?, ?, ?)
	`, company, role, status, applied, now, now)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) UpdateMessageClassification(
	ctx context.Context,
	messageID int64,
	applicationID int64,
	eventType, classifiedBy string,
) error {
	_, err := s.DB.ExecContext(ctx, `
		UPDATE messages SET application_id = ?, event_type = ?, classified_by = ?
		WHERE id = ?
	`, applicationID, eventType, classifiedBy, messageID)
	return err
}

func (s *Store) EventExistsForMessage(ctx context.Context, messageID int64) (bool, error) {
	var n int
	err := s.DB.QueryRowContext(ctx, `SELECT 1 FROM events WHERE message_id = ? LIMIT 1`, messageID).Scan(&n)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *Store) InsertEvent(ctx context.Context, applicationID, messageID int64, eventType string, occurredAt time.Time, detailsJSON string) error {
	if detailsJSON == "" {
		detailsJSON = "{}"
	}
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO events (application_id, type, occurred_at, message_id, details_json)
		VALUES (?, ?, ?, ?, ?)
	`, applicationID, eventType, occurredAt.UTC().Format(time.RFC3339), messageID, detailsJSON)
	return err
}

func (s *Store) ListApplicationEvents(ctx context.Context, applicationID int64) ([]struct {
	Type       string
	OccurredAt time.Time
}, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT type, occurred_at FROM events
		WHERE application_id = ? ORDER BY occurred_at ASC, id ASC
	`, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []struct {
		Type       string
		OccurredAt time.Time
	}
	for rows.Next() {
		var e struct {
			Type       string
			OccurredAt time.Time
		}
		var at string
		if err := rows.Scan(&e.Type, &at); err != nil {
			return nil, err
		}
		t, _ := time.Parse(time.RFC3339, at)
		e.OccurredAt = t
		out = append(out, e)
	}
	return out, rows.Err()
}

func (s *Store) UpdateApplicationDerived(ctx context.Context, applicationID int64, company, role, status string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := s.DB.ExecContext(ctx, `
		UPDATE applications SET company = CASE WHEN company = '' THEN ? ELSE company END,
			role = CASE WHEN role = '' THEN ? ELSE role END,
			status = ?, updated_at = ?
		WHERE id = ?
	`, company, role, status, now, applicationID)
	return err
}

func (s *Store) UpdateApplicationStatusOnly(ctx context.Context, applicationID int64, status string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := s.DB.ExecContext(ctx, `
		UPDATE applications SET status = ?, updated_at = ? WHERE id = ? AND status_overridden = 0
	`, status, now, applicationID)
	return err
}

func (s *Store) PatchApplicationCompanyRole(ctx context.Context, applicationID int64, company, role string) error {
	_, err := s.DB.ExecContext(ctx, `
		UPDATE applications SET company = ?, role = ? WHERE id = ?
	`, company, role, applicationID)
	return err
}

func DetailsJSON(details map[string]any) (string, error) {
	b, err := json.Marshal(details)
	if err != nil {
		return "{}", err
	}
	return string(b), nil
}

func (s *Store) ListMessagesForApplication(ctx context.Context, applicationID int64) ([]MessageRow, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, gmail_id, thread_id, application_id, direction, from_addr, to_addr,
			subject, snippet, body_excerpt, received_at, event_type, classified_by
		FROM messages WHERE application_id = ? ORDER BY received_at ASC
	`, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []MessageRow
	for rows.Next() {
		var m MessageRow
		var appID sql.NullInt64
		var received string
		if err := rows.Scan(
			&m.ID, &m.GmailID, &m.ThreadID, &appID, &m.Direction, &m.FromAddr, &m.ToAddr,
			&m.Subject, &m.Snippet, &m.BodyExcerpt, &received, &m.EventType, &m.ClassifiedBy,
		); err != nil {
			return nil, err
		}
		m.ApplicationID = appID
		t, _ := time.Parse(time.RFC3339, received)
		m.ReceivedAt = t
		out = append(out, m)
	}
	return out, rows.Err()
}

func (s *Store) ListEventDetailsForApplication(ctx context.Context, applicationID int64) ([]string, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT details_json FROM events WHERE application_id = ?
	`, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var d string
		if err := rows.Scan(&d); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func ParseEventDetailsJSON(raw string) (interviewAt *time.Time, meetingURL string) {
	var m map[string]any
	if json.Unmarshal([]byte(raw), &m) != nil {
		return nil, ""
	}
	if u, ok := m["meeting_url"].(string); ok {
		meetingURL = u
	}
	if s, ok := m["interview_at"].(string); ok {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			interviewAt = &t
		}
	}
	return interviewAt, meetingURL
}

func (s *Store) MergeApplicationMeta(ctx context.Context, id int64, company, role string, appliedAt *time.Time) error {
	a, err := s.GetApplication(ctx, id)
	if err != nil {
		return err
	}
	if company != "" && a.Company == "" {
		a.Company = company
	}
	if role != "" && a.Role == "" {
		a.Role = role
	}
	_, err = s.DB.ExecContext(ctx, `
		UPDATE applications SET company = ?, role = ?,
			applied_at = COALESCE(applied_at, ?)
		WHERE id = ?
	`, a.Company, a.Role, nullableTime(appliedAt), id)
	return err
}

func nullableTime(t *time.Time) sql.NullString {
	if t == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: t.UTC().Format(time.RFC3339), Valid: true}
}

func (s *Store) GetApplication(ctx context.Context, id int64) (*ApplicationRow, error) {
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, company, role, status, status_overridden, applied_at FROM applications WHERE id = ?
	`, id)
	var a ApplicationRow
	var overridden int
	if err := row.Scan(&a.ID, &a.Company, &a.Role, &a.Status, &overridden, &a.AppliedAt); err != nil {
		return nil, fmt.Errorf("get application: %w", err)
	}
	a.StatusOverridden = overridden == 1
	return &a, nil
}

package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ApplicationFull struct {
	ApplicationRow
	Source    string
	Notes     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Store) ListApplicationsFull(ctx context.Context, status, q string) ([]ApplicationFull, error) {
	query := `
		SELECT id, company, role, status, status_overridden, applied_at, source, notes, created_at, updated_at
		FROM applications WHERE 1=1
	`
	var args []any
	if status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}
	if q != "" {
		query += ` AND (company LIKE ? OR role LIKE ? OR notes LIKE ?)`
		like := "%" + q + "%"
		args = append(args, like, like, like)
	}
	query += ` ORDER BY updated_at DESC, id DESC`

	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ApplicationFull
	for rows.Next() {
		a, err := scanApplicationFull(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (s *Store) GetApplicationFull(ctx context.Context, id int64) (*ApplicationFull, error) {
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, company, role, status, status_overridden, applied_at, source, notes, created_at, updated_at
		FROM applications WHERE id = ?
	`, id)
	a, err := scanApplicationFullRow(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrApplicationNotFound
	}
	if err != nil {
		return nil, err
	}
	return a, nil
}

func scanApplicationFull(rows *sql.Rows) (ApplicationFull, error) {
	var a ApplicationFull
	var overridden int
	var created, updated string
	if err := rows.Scan(
		&a.ID, &a.Company, &a.Role, &a.Status, &overridden, &a.AppliedAt,
		&a.Source, &a.Notes, &created, &updated,
	); err != nil {
		return ApplicationFull{}, err
	}
	a.StatusOverridden = overridden == 1
	a.CreatedAt, _ = time.Parse(time.RFC3339, created)
	a.UpdatedAt, _ = time.Parse(time.RFC3339, updated)
	return a, nil
}

func scanApplicationFullRow(row *sql.Row) (*ApplicationFull, error) {
	var a ApplicationFull
	var overridden int
	var created, updated string
	if err := row.Scan(
		&a.ID, &a.Company, &a.Role, &a.Status, &overridden, &a.AppliedAt,
		&a.Source, &a.Notes, &created, &updated,
	); err != nil {
		return nil, err
	}
	a.StatusOverridden = overridden == 1
	a.CreatedAt, _ = time.Parse(time.RFC3339, created)
	a.UpdatedAt, _ = time.Parse(time.RFC3339, updated)
	return &a, nil
}

func (s *Store) PatchApplication(
	ctx context.Context,
	id int64,
	company, role, status, notes *string,
	statusManual bool,
) error {
	app, err := s.GetApplicationFull(ctx, id)
	if err != nil {
		return err
	}
	if company != nil {
		app.Company = *company
	}
	if role != nil {
		app.Role = *role
	}
	if notes != nil {
		app.Notes = *notes
	}
	overridden := 0
	if app.StatusOverridden {
		overridden = 1
	}
	if status != nil {
		app.Status = *status
		if statusManual {
			overridden = 1
		}
	}
	now := time.Now().UTC().Format(time.RFC3339)
	res, err := s.DB.ExecContext(ctx, `
		UPDATE applications SET company = ?, role = ?, status = ?, notes = ?,
			status_overridden = ?, updated_at = ?
		WHERE id = ?
	`, app.Company, app.Role, app.Status, app.Notes, overridden, now, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrApplicationNotFound
	}
	return nil
}

func (s *Store) MergeApplications(ctx context.Context, targetID int64, sourceIDs []int64) error {
	if targetID <= 0 || len(sourceIDs) == 0 {
		return fmt.Errorf("invalid merge request")
	}
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	for _, src := range sourceIDs {
		if src == targetID {
			continue
		}
		if _, err := tx.ExecContext(ctx, `
			UPDATE messages SET application_id = ? WHERE application_id = ?
		`, targetID, src); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `
			UPDATE events SET application_id = ? WHERE application_id = ?
		`, targetID, src); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `DELETE FROM applications WHERE id = ?`, src); err != nil {
			return err
		}
	}
	now := time.Now().UTC().Format(time.RFC3339)
	if _, err := tx.ExecContext(ctx, `UPDATE applications SET updated_at = ? WHERE id = ?`, now, targetID); err != nil {
		return err
	}
	return tx.Commit()
}

type TimelineEventRow struct {
	ID          int64
	Type        string
	OccurredAt  time.Time
	MessageID   sql.NullInt64
	DetailsJSON string
}

func (s *Store) ListTimelineEvents(ctx context.Context, applicationID int64) ([]TimelineEventRow, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, type, occurred_at, message_id, details_json
		FROM events WHERE application_id = ? ORDER BY occurred_at ASC, id ASC
	`, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []TimelineEventRow
	for rows.Next() {
		var e TimelineEventRow
		var at string
		if err := rows.Scan(&e.ID, &e.Type, &at, &e.MessageID, &e.DetailsJSON); err != nil {
			return nil, err
		}
		e.OccurredAt, _ = time.Parse(time.RFC3339, at)
		out = append(out, e)
	}
	return out, rows.Err()
}

func (s *Store) ListAllRules(ctx context.Context) ([]RuleRow, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, field, pattern, event_type, priority, enabled FROM rules
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

func (s *Store) GetRule(ctx context.Context, id int64) (*RuleRow, error) {
	var r RuleRow
	var enabled int
	err := s.DB.QueryRowContext(ctx, `
		SELECT id, field, pattern, event_type, priority, enabled FROM rules WHERE id = ?
	`, id).Scan(&r.ID, &r.Field, &r.Pattern, &r.EventType, &r.Priority, &enabled)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("rule not found")
	}
	if err != nil {
		return nil, err
	}
	r.Enabled = enabled == 1
	return &r, nil
}

func (s *Store) InsertRule(ctx context.Context, field, pattern, eventType string, priority int, enabled bool) (int64, error) {
	en := 0
	if enabled {
		en = 1
	}
	res, err := s.DB.ExecContext(ctx, `
		INSERT INTO rules (field, pattern, event_type, priority, enabled)
		VALUES (?, ?, ?, ?, ?)
	`, field, pattern, eventType, priority, en)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) UpdateRule(ctx context.Context, id int64, field, pattern, eventType *string, priority *int, enabled *bool) error {
	r, err := s.GetRule(ctx, id)
	if err != nil {
		return err
	}
	if field != nil {
		r.Field = *field
	}
	if pattern != nil {
		r.Pattern = *pattern
	}
	if eventType != nil {
		r.EventType = *eventType
	}
	if priority != nil {
		r.Priority = *priority
	}
	if enabled != nil {
		r.Enabled = *enabled
	}
	en := 0
	if r.Enabled {
		en = 1
	}
	res, err := s.DB.ExecContext(ctx, `
		UPDATE rules SET field = ?, pattern = ?, event_type = ?, priority = ?, enabled = ? WHERE id = ?
	`, r.Field, r.Pattern, r.EventType, r.Priority, en, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("rule not found")
	}
	return nil
}

func (s *Store) DeleteRule(ctx context.Context, id int64) error {
	res, err := s.DB.ExecContext(ctx, `DELETE FROM rules WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("rule not found")
	}
	return nil
}

func (s *Store) SetSetting(ctx context.Context, key, value string) error {
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO settings (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value
	`, key, value)
	return err
}

type StatsRow struct {
	ByStatus            map[string]int
	ResponseRate        float64
	AvgDaysToFirstReply float64
}

func (s *Store) ComputeStats(ctx context.Context) (StatsRow, error) {
	var out StatsRow
	out.ByStatus = map[string]int{
		"applied": 0, "replied": 0, "interview": 0, "offer": 0, "rejected": 0,
	}
	rows, err := s.DB.QueryContext(ctx, `SELECT status, COUNT(*) FROM applications GROUP BY status`)
	if err != nil || rows.Err() != nil {
		return out, err
	}
	defer rows.Close()
	for rows.Next() {
		var st string
		var n int
		if err := rows.Scan(&st, &n); err != nil {
			rows.Close()
			return out, err
		}
		out.ByStatus[st] = n
	}
	rows.Close()

	var withInbound, withReply int
	err = s.DB.QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT application_id) FROM messages WHERE direction = 'in' AND application_id IS NOT NULL
	`).Scan(&withInbound)
	if err != nil {
		return out, err
	}
	err = s.DB.QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT m1.application_id) FROM messages m1
		WHERE m1.direction = 'out' AND m1.application_id IS NOT NULL
		AND EXISTS (
			SELECT 1 FROM messages m2 WHERE m2.application_id = m1.application_id AND m2.direction = 'in'
		)
	`).Scan(&withReply)
	if err != nil {
		return out, err
	}
	if withInbound > 0 {
		out.ResponseRate = float64(withReply) / float64(withInbound)
	}

	// avg days from first inbound to first outbound per application
	rows2, err := s.DB.QueryContext(ctx, `
		SELECT application_id,
			MIN(CASE WHEN direction = 'in' THEN received_at END) AS first_in,
			MIN(CASE WHEN direction = 'out' THEN received_at END) AS first_out
		FROM messages WHERE application_id IS NOT NULL
		GROUP BY application_id
	`)
	if err != nil || rows2.Err() != nil {
		return out, err
	}
	defer rows2.Close()
	var totalDays float64
	var count int
	for rows2.Next() {
		var appID int64
		var firstIn, firstOut sql.NullString
		if err := rows2.Scan(&appID, &firstIn, &firstOut); err != nil {
			return out, err
		}
		if !firstIn.Valid || !firstOut.Valid {
			continue
		}
		tIn, e1 := time.Parse(time.RFC3339, firstIn.String)
		tOut, e2 := time.Parse(time.RFC3339, firstOut.String)
		if e1 != nil || e2 != nil || !tOut.After(tIn) {
			continue
		}
		totalDays += tOut.Sub(tIn).Hours() / 24
		count++
	}
	if count > 0 {
		out.AvgDaysToFirstReply = totalDays / float64(count)
	}
	return out, nil
}

func ParseSyncIntervalMinutes(raw string, defaultMin int) int {
	if raw == "" {
		return defaultMin
	}
	n, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || n < 1 {
		return defaultMin
	}
	return n
}

func (s *Store) FirstThreadIDForApplication(ctx context.Context, applicationID int64) (string, error) {
	var tid string
	err := s.DB.QueryRowContext(ctx, `
		SELECT thread_id FROM messages WHERE application_id = ? AND thread_id != '' LIMIT 1
	`, applicationID).Scan(&tid)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	return tid, err
}

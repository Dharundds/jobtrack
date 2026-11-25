package store

import (
	"context"
	"database/sql"
)

// CountApplicationsNeedingLLMReview counts pipeline applications with at least one message not classified by the LLM.
func (s *Store) CountApplicationsNeedingLLMReview(ctx context.Context) (int, error) {
	var n int
	err := s.DB.QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT a.id)
		FROM applications a
		JOIN messages m ON m.application_id = a.id AND m.ignored = 0
		WHERE COALESCE(LOWER(TRIM(m.classified_by)), '') != 'llm'
	`).Scan(&n)
	return n, err
}

// ListApplicationIDsNeedingLLMReview returns application ids newest activity first. limit <= 0 means no limit.
func (s *Store) ListApplicationIDsNeedingLLMReview(ctx context.Context, limit int) ([]int64, error) {
	q := `
		SELECT a.id
		FROM applications a
		WHERE EXISTS (
			SELECT 1 FROM messages m
			WHERE m.application_id = a.id AND m.ignored = 0
			  AND COALESCE(LOWER(TRIM(m.classified_by)), '') != 'llm'
		)
		ORDER BY (
			SELECT MAX(received_at) FROM messages
			WHERE application_id = a.id AND ignored = 0
		) DESC, a.id DESC`
	var rows *sql.Rows
	var err error
	if limit > 0 {
		rows, err = s.DB.QueryContext(ctx, q+` LIMIT ?`, limit)
	} else {
		rows, err = s.DB.QueryContext(ctx, q)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, rows.Err()
}
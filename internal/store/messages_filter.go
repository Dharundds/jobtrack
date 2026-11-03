package store

import (
	"context"
	"fmt"
	"strings"
)

const gmailIDChunkSize = 400

// FilterNewGmailIDs returns IDs not yet stored and the number skipped (already present).
func (s *Store) FilterNewGmailIDs(ctx context.Context, ids []string) (newIDs []string, skipped int, err error) {
	if len(ids) == 0 {
		return nil, 0, nil
	}

	existing := make(map[string]struct{}, len(ids))
	for start := 0; start < len(ids); start += gmailIDChunkSize {
		end := start + gmailIDChunkSize
		if end > len(ids) {
			end = len(ids)
		}
		chunk := ids[start:end]
		placeholders := strings.Repeat("?,", len(chunk))
		placeholders = placeholders[:len(placeholders)-1]
		args := make([]any, len(chunk))
		for i, id := range chunk {
			args[i] = id
		}
		query := fmt.Sprintf(`SELECT gmail_id FROM messages WHERE gmail_id IN (%s)`, placeholders)
		rows, qerr := s.DB.QueryContext(ctx, query, args...)
		if qerr != nil {
			return nil, 0, qerr
		}
		for rows.Next() {
			var id string
			if scanErr := rows.Scan(&id); scanErr != nil {
				rows.Close()
				return nil, 0, scanErr
			}
			existing[id] = struct{}{}
		}
		if err := rows.Close(); err != nil {
			return nil, 0, err
		}
		if err := rows.Err(); err != nil {
			return nil, 0, err
		}
	}

	for _, id := range ids {
		if _, ok := existing[id]; ok {
			skipped++
			continue
		}
		newIDs = append(newIDs, id)
	}
	return newIDs, skipped, nil
}

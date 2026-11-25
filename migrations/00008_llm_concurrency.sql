-- +goose Up
INSERT OR IGNORE INTO settings (key, value) VALUES ('llm_concurrency', '2');

-- +goose Down
DELETE FROM settings WHERE key = 'llm_concurrency';
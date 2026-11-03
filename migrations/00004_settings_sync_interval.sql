-- +goose Up
INSERT OR IGNORE INTO settings (key, value) VALUES ('sync_interval_minutes', '5');

-- +goose Down
DELETE FROM settings WHERE key = 'sync_interval_minutes';

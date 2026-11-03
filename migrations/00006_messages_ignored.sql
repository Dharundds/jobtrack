-- +goose Up
ALTER TABLE messages ADD COLUMN ignored INTEGER NOT NULL DEFAULT 0;
CREATE INDEX idx_messages_ignored ON messages(ignored);

-- +goose Down
DROP INDEX IF EXISTS idx_messages_ignored;
ALTER TABLE messages DROP COLUMN ignored;

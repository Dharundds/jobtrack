-- +goose Up
INSERT OR IGNORE INTO settings (key, value) VALUES ('llm_base_url', '');
INSERT OR IGNORE INTO settings (key, value) VALUES ('llm_model', '');
INSERT OR IGNORE INTO settings (key, value) VALUES ('llm_api_key_enc', '');

-- +goose Down
DELETE FROM settings WHERE key IN ('llm_base_url', 'llm_model', 'llm_api_key_enc');
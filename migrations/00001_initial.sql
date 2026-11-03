-- +goose Up
CREATE TABLE accounts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    enc_refresh_token TEXT NOT NULL,
    history_id TEXT,
    last_synced_at TEXT
);

CREATE TABLE applications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    company TEXT NOT NULL DEFAULT '',
    role TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'applied',
    source TEXT NOT NULL DEFAULT 'email',
    applied_at TEXT,
    status_overridden INTEGER NOT NULL DEFAULT 0,
    notes TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    gmail_id TEXT NOT NULL UNIQUE,
    thread_id TEXT NOT NULL,
    application_id INTEGER REFERENCES applications(id) ON DELETE SET NULL,
    direction TEXT NOT NULL CHECK (direction IN ('in', 'out')),
    from_addr TEXT NOT NULL DEFAULT '',
    to_addr TEXT NOT NULL DEFAULT '',
    subject TEXT NOT NULL DEFAULT '',
    snippet TEXT NOT NULL DEFAULT '',
    body_excerpt TEXT NOT NULL DEFAULT '',
    received_at TEXT NOT NULL,
    event_type TEXT NOT NULL DEFAULT 'other',
    classified_by TEXT NOT NULL DEFAULT 'rule'
);

CREATE INDEX idx_messages_gmail_id ON messages(gmail_id);
CREATE INDEX idx_messages_thread_id ON messages(thread_id);
CREATE INDEX idx_messages_application_id ON messages(application_id);
CREATE INDEX idx_messages_received_at ON messages(received_at);

CREATE TABLE events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    application_id INTEGER NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    occurred_at TEXT NOT NULL,
    message_id INTEGER REFERENCES messages(id) ON DELETE SET NULL,
    details_json TEXT NOT NULL DEFAULT '{}'
);

CREATE INDEX idx_events_application_id ON events(application_id);

CREATE TABLE rules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    field TEXT NOT NULL CHECK (field IN ('subject', 'body', 'from')),
    pattern TEXT NOT NULL,
    event_type TEXT NOT NULL,
    priority INTEGER NOT NULL DEFAULT 0,
    enabled INTEGER NOT NULL DEFAULT 1
);

-- +goose Down
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS applications;
DROP TABLE IF EXISTS rules;
DROP TABLE IF EXISTS accounts;

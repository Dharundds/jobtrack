-- +goose Up
CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY NOT NULL,
    value TEXT NOT NULL DEFAULT ''
);

INSERT OR IGNORE INTO settings (key, value) VALUES (
    'gmail_query',
    '(application OR applied OR interview OR offer OR unfortunately OR "thank you for applying" OR assessment OR recruiter OR from:greenhouse.io OR from:lever.co OR from:ashbyhq.com OR from:myworkday OR from:smartrecruiters OR from:icims OR from:linkedin.com OR from:naukri.com OR from:indeed.com OR from:wellfound.com)'
);

-- +goose Down
DROP TABLE IF EXISTS settings;

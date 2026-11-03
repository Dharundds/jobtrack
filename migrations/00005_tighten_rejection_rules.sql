-- +goose Up
-- Tighten broad rejection subject rule (requires job-rejection phrasing, not bare "unfortunately").
UPDATE rules
SET pattern = '(?i)(not moving forward|regret to inform|not selected|decided to pursue other|will not be proceeding)'
WHERE field = 'subject' AND event_type = 'rejected' AND pattern LIKE '%unfortunately%';

-- +goose Down
UPDATE rules
SET pattern = '(?i)unfortunately|not moving forward|regret to inform|not selected|other candidates'
WHERE field = 'subject' AND event_type = 'rejected';

-- +goose Up
INSERT INTO rules (field, pattern, event_type, priority, enabled) VALUES
('subject', '(?i)application for .+ received', 'applied', 105, 1),
('subject', '(?i)received,? thank you', 'applied', 105, 1),
('body', '(?i)thank you for applying to|confirm that we have received your application|we have received your application', 'applied', 95, 1),
('from', '(?i)kekamail\.com', 'applied', 55, 1);

-- +goose Down
DELETE FROM rules WHERE field = 'subject' AND pattern LIKE '%application for .+ received%';
DELETE FROM rules WHERE field = 'subject' AND pattern LIKE '%received,? thank you%';
DELETE FROM rules WHERE field = 'body' AND pattern LIKE '%thank you for applying to%';
DELETE FROM rules WHERE field = 'from' AND pattern LIKE '%kekamail%';
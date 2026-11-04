-- +goose Up
INSERT INTO rules (field, pattern, event_type, priority, enabled) VALUES
('subject', '(?i)thank you for applying|application received|we received your application', 'applied', 100, 1),
('body', '(?i)thank you for applying|application received|successfully applied', 'applied', 90, 1),
('from', '(?i)greenhouse|lever\.co|ashbyhq|workday|smartrecruiters|icims', 'applied', 50, 1),

('subject', '(?i)acknowledg|confirm.*application|received your resume', 'acknowledged', 100, 1),
('body', '(?i)we have received your (application|resume)|thanks for your interest', 'acknowledged', 80, 1),

('subject', '(?i)(not moving forward|regret to inform|not selected|decided to pursue other|will not be proceeding)', 'rejected', 100, 1),
('body', '(?i)unfortunately|not moving forward|regret to inform|will not be proceeding|decided to pursue other', 'rejected', 90, 1),

('subject', '(?i)interview|phone screen|video call|schedule.*(call|chat|meet)|screening process|participate in (?:our |the )?screening|invites you to participate|invited to participate', 'interview_invite', 100, 1),
('body', '(?i)interview|schedule.*(call|meeting)|calendar invite|availability for a (call|chat)|screening process|participate in (?:our |the )?screening|invites you to participate|joining us as (?:a|an)\s+\w|interest in joining us as', 'interview_invite', 85, 1),
('body', '(?i)zoom\.us/j/|meet\.google\.com/|teams\.microsoft\.com/', 'interview_invite', 70, 1),

('subject', '(?i)assessment|take-home|coding challenge|hackerrank|codility|your screening for|screening for .+ successfully submitted|screening has been successfully submitted', 'assessment', 100, 1),
('body', '(?i)complete (the|this) assessment|online assessment|take-home', 'assessment', 90, 1),

('subject', '(?i)offer letter|pleased to offer|job offer|offer of employment', 'offer', 100, 1),
('body', '(?i)pleased to offer|offer letter|extend an offer', 'offer', 90, 1),

('subject', '(?i)opportunity at|hiring for|role at|opening at', 'recruiter_outreach', 80, 1),
('from', '(?i)recruiter|talent|careers|hiring', 'recruiter_outreach', 40, 1),
('body', '(?i)are you open to|reach out about|exciting opportunity', 'recruiter_outreach', 70, 1);

-- +goose Down
DELETE FROM rules;

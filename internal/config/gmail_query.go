package config

// FullSyncDays is how far back the first (or expired) Gmail list query looks.
const FullSyncDays = 15

// DefaultGmailQuery is used when settings.gmail_query is empty.
// Job-related mail is matched by keywords (subject/body); avoid employer-specific from: domains.
const DefaultGmailQuery = `(application OR applied OR interview OR offer OR unfortunately OR "thank you for applying" OR assessment OR recruiter OR subject:screening OR subject:"screening process" OR from:greenhouse.io OR from:lever.co OR from:ashbyhq.com OR from:myworkday OR from:smartrecruiters OR from:icims OR from:linkedin.com OR from:naukri.com OR from:indeed.com OR from:wellfound.com)`

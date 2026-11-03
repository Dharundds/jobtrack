package classify

import "time"

const (
	EventApplied            = "applied"
	EventAcknowledged       = "acknowledged"
	EventRejected           = "rejected"
	EventInterviewInvite    = "interview_invite"
	EventAssessment         = "assessment"
	EventOffer              = "offer"
	EventRecruiterOutreach  = "recruiter_outreach"
	EventOther              = "other"
)

// EmailInput is the data passed to a Classifier.
type EmailInput struct {
	From        string
	To          string
	Subject     string
	Body        string
	Snippet     string
	Direction   string // in | out
	AccountEmail string
	LabelIDs    []string
	ReceivedAt  time.Time
}

// Result holds classification and extracted metadata.
type Result struct {
	EventType   string
	Company     string
	Role        string
	Details     EventDetails
	ClassifiedBy string
}

type EventDetails struct {
	InterviewAt *time.Time `json:"interview_at,omitempty"`
	MeetingURL  string     `json:"meeting_url,omitempty"`
}

// Classifier maps an email to an event type and extracted fields.
type Classifier interface {
	Classify(in EmailInput) Result
}

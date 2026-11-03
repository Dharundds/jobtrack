package classify

import (
	"strings"
	"time"
)

// StatusFromEvent maps a classified event to an application status.
func StatusFromEvent(eventType string) string {
	switch eventType {
	case EventApplied, EventRecruiterOutreach:
		return "applied"
	case EventAcknowledged:
		return "replied"
	case EventInterviewInvite, EventAssessment:
		return "interview"
	case EventOffer:
		return "offer"
	case EventRejected:
		return "rejected"
	case EventOther:
		return "applied" // only used when attaching to an existing application thread
	default:
		return "replied"
	}
}

// LatestStatusFromEvents picks status from the newest meaningful event.
func LatestStatusFromEvents(events []EventRecord) string {
	for i := len(events) - 1; i >= 0; i-- {
		e := events[i]
		if e.Type == EventOther {
			continue
		}
		return StatusFromEvent(e.Type)
	}
	return "applied"
}

type EventRecord struct {
	Type       string
	OccurredAt time.Time
}

// ApplicationActivity derived from messages on an application.
type ApplicationActivity struct {
	LastInboundAt         *time.Time
	LastOutboundAt        *time.Time
	AwaitingMyReply       bool
	DaysSinceLastActivity int
	NextInterviewAt       *time.Time
}

type MessageActivity struct {
	Direction  string
	ReceivedAt time.Time
	EventType  string
}

func DeriveActivity(msgs []MessageActivity, events []EventDetails) ApplicationActivity {
	var act ApplicationActivity
	var lastActivity *time.Time

	for _, m := range msgs {
		t := m.ReceivedAt
		if m.Direction == "in" {
			if act.LastInboundAt == nil || t.After(*act.LastInboundAt) {
				act.LastInboundAt = &t
			}
		}
		if m.Direction == "out" {
			if act.LastOutboundAt == nil || t.After(*act.LastOutboundAt) {
				act.LastOutboundAt = &t
			}
		}
		if lastActivity == nil || t.After(*lastActivity) {
			lastActivity = &t
		}
	}

	if act.LastInboundAt != nil {
		if act.LastOutboundAt == nil || act.LastInboundAt.After(*act.LastOutboundAt) {
			// Last message is inbound — awaiting reply if recent inbound after last outbound
			act.AwaitingMyReply = true
		}
	}

	// If last message overall is outbound, not awaiting
	if len(msgs) > 0 {
		last := msgs[0]
		for _, m := range msgs[1:] {
			if m.ReceivedAt.After(last.ReceivedAt) {
				last = m
			}
		}
		if last.Direction == "out" {
			act.AwaitingMyReply = false
		} else if last.Direction == "in" {
			act.AwaitingMyReply = act.LastOutboundAt == nil || last.ReceivedAt.After(*act.LastOutboundAt)
		}
	}

	if lastActivity != nil {
		days := int(time.Since(*lastActivity).Hours() / 24)
		if days < 0 {
			days = 0
		}
		act.DaysSinceLastActivity = days
	}

	for _, e := range events {
		if e.InterviewAt != nil {
			if act.NextInterviewAt == nil || e.InterviewAt.Before(*act.NextInterviewAt) {
				if e.InterviewAt.After(time.Now().UTC()) {
					act.NextInterviewAt = e.InterviewAt
				}
			}
		}
	}

	return act
}

// SimilarCompanyRole returns true when two company/role pairs likely refer to the same application.
func SimilarCompanyRole(c1, r1, c2, r2 string) bool {
	c1, r1 = normalizeKey(c1), normalizeKey(r1)
	c2, r2 = normalizeKey(c2), normalizeKey(r2)
	// Shared ATS sender (e.g. greenhouse-mail.io) must not merge different employers.
	if isATSBrandName(c1) || isATSBrandName(c2) {
		return false
	}
	if c1 != "" && c2 != "" && stringsMatch(c1, c2) {
		if r1 == "" || r2 == "" {
			return true
		}
		return stringsMatch(r1, r2)
	}
	return false
}

func normalizeKey(s string) string {
	return stringsTrimLower(s)
}

func stringsMatch(a, b string) bool {
	if a == b {
		return true
	}
	return strings.Contains(a, b) || strings.Contains(b, a)
}

func stringsTrimLower(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

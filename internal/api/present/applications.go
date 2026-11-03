package present

import (
	"encoding/json"
	"time"

	"github.com/dharundds/jobtrack/internal/classify"
	"github.com/dharundds/jobtrack/internal/store"
)

type ApplicationJSON struct {
	ID                    int64      `json:"id"`
	Company               string     `json:"company"`
	Role                  string     `json:"role"`
	Status                string     `json:"status"`
	Source                string     `json:"source"`
	AppliedAt             *time.Time `json:"applied_at,omitempty"`
	Notes                 string     `json:"notes"`
	StatusOverridden      bool       `json:"status_overridden"`
	LastInboundAt         *time.Time `json:"last_inbound_at,omitempty"`
	LastOutboundAt        *time.Time `json:"last_outbound_at,omitempty"`
	AwaitingMyReply       bool       `json:"awaiting_my_reply"`
	DaysSinceLastActivity int        `json:"days_since_last_activity"`
	NextInterviewAt       *time.Time `json:"next_interview_at,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	ThreadID              string     `json:"thread_id,omitempty"`
}

type MessageJSON struct {
	ID          int64     `json:"id"`
	Direction   string    `json:"direction"`
	FromAddr    string    `json:"from_addr"`
	ToAddr      string    `json:"to_addr"`
	Subject     string    `json:"subject"`
	Snippet     string    `json:"snippet"`
	BodyExcerpt string    `json:"body_excerpt"`
	ReceivedAt  time.Time `json:"received_at"`
	EventType   string    `json:"event_type"`
	ThreadID    string    `json:"thread_id,omitempty"`
	GmailID     string    `json:"gmail_id,omitempty"`
}

type EventJSON struct {
	ID         int64          `json:"id"`
	Type       string         `json:"type"`
	OccurredAt time.Time      `json:"occurred_at"`
	MessageID  *int64         `json:"message_id,omitempty"`
	Details    map[string]any `json:"details"`
}

func ApplicationFromStore(
	app store.ApplicationFull,
	msgs []store.MessageRow,
	detailJSONs []string,
	threadID string,
) ApplicationJSON {
	actMsgs := make([]classify.MessageActivity, len(msgs))
	for i, m := range msgs {
		actMsgs[i] = classify.MessageActivity{
			Direction:  m.Direction,
			ReceivedAt: m.ReceivedAt,
			EventType:  m.EventType,
		}
	}
	var details []classify.EventDetails
	for _, raw := range detailJSONs {
		at, url := store.ParseEventDetailsJSON(raw)
		details = append(details, classify.EventDetails{InterviewAt: at, MeetingURL: url})
	}
	act := classify.DeriveActivity(actMsgs, details)

	out := ApplicationJSON{
		ID:                    app.ID,
		Company:               app.Company,
		Role:                  app.Role,
		Status:                app.Status,
		Source:                app.Source,
		Notes:                 app.Notes,
		StatusOverridden:      app.StatusOverridden,
		LastInboundAt:         act.LastInboundAt,
		LastOutboundAt:        act.LastOutboundAt,
		AwaitingMyReply:       act.AwaitingMyReply,
		DaysSinceLastActivity: act.DaysSinceLastActivity,
		NextInterviewAt:       act.NextInterviewAt,
		CreatedAt:             app.CreatedAt,
		UpdatedAt:             app.UpdatedAt,
		ThreadID:              threadID,
	}
	if app.AppliedAt.Valid {
		t, err := time.Parse(time.RFC3339, app.AppliedAt.String)
		if err == nil {
			out.AppliedAt = &t
		}
	}
	return out
}

func MessageFromStore(m store.MessageRow, includeGmail bool) MessageJSON {
	j := MessageJSON{
		ID:          m.ID,
		Direction:   m.Direction,
		FromAddr:    m.FromAddr,
		ToAddr:      m.ToAddr,
		Subject:     m.Subject,
		Snippet:     m.Snippet,
		BodyExcerpt: m.BodyExcerpt,
		ReceivedAt:  m.ReceivedAt,
		EventType:   m.EventType,
		ThreadID:    m.ThreadID,
	}
	if includeGmail {
		j.GmailID = m.GmailID
	}
	return j
}

func EventFromStore(e store.TimelineEventRow) EventJSON {
	details := map[string]any{}
	_ = json.Unmarshal([]byte(e.DetailsJSON), &details)
	if details == nil {
		details = map[string]any{}
	}
	var msgID *int64
	if e.MessageID.Valid {
		id := e.MessageID.Int64
		msgID = &id
	}
	return EventJSON{
		ID:         e.ID,
		Type:       e.Type,
		OccurredAt: e.OccurredAt,
		MessageID:  msgID,
		Details:    details,
	}
}

func RuleFromStore(r store.RuleRow) map[string]any {
	return map[string]any{
		"id":         r.ID,
		"field":      r.Field,
		"pattern":    r.Pattern,
		"event_type": r.EventType,
		"priority":   r.Priority,
		"enabled":    r.Enabled,
	}
}

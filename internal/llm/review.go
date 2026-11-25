package llm

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/dharundds/jobtrack/internal/store"
)

// ApplicationReview is the LLM verdict on whether a listing is a real job application.
type ApplicationReview struct {
	IsJobApplication bool
	Reason           string
	Company          string
	Role             string
	PipelineStatus   string // applied | replied | interview | offer | rejected
}

// ReviewApplication decides if linked emails are one job application listing.
func (c *Client) ReviewApplication(
	ctx context.Context,
	accountEmail string,
	app store.ApplicationFull,
	msgs []store.MessageRow,
) (ApplicationReview, error) {
	if !c.Configured() {
		return ApplicationReview{}, fmt.Errorf("llm not configured")
	}
	if len(msgs) == 0 {
		return ApplicationReview{}, fmt.Errorf("no messages to review")
	}

	sorted := append([]store.MessageRow(nil), msgs...)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].ReceivedAt.Before(sorted[j].ReceivedAt)
	})

	var b strings.Builder
	fmt.Fprintf(&b, "Current listing in database: company=%q role=%q status=%q\n\nEmails:\n",
		app.Company, app.Role, app.Status)
	for i, m := range sorted {
		body := m.BodyExcerpt
		if body == "" {
			body = m.Snippet
		}
		body = trimRunes(body, 800)
		fmt.Fprintf(&b, "--- email %d ---\ndirection: %s\ndate: %s\nfrom: %s\nto: %s\nsubject: %s\nclassified_as: %s (by %s)\nbody:\n%s\n\n",
			i+1, m.Direction, m.ReceivedAt.UTC().Format(time.RFC3339),
			m.FromAddr, m.ToAddr, m.Subject, m.EventType, m.ClassifiedBy, body)
	}

	user := b.String() + `
You review whether this listing should exist in a job-application tracker for account ` + accountEmail + `.

Reply with JSON only (no markdown):
{
  "is_job_application": true or false,
  "reason": "short explanation",
  "company": "employer name if is_job_application else empty",
  "role": "job title if known else empty",
  "pipeline_status": "applied|replied|interview|offer|rejected"
}

Set is_job_application to false for newsletters, shopping, LeetCode/contest mail, social notifications, unrelated LinkedIn noise, Instahyre/Naukri/Indeed job alerts ("has just posted", "you look like a good fit", "Apply here" with no confirmation you applied), or generic marketing — those listings should be removed.
Set true only when emails show you actually applied or are in an employer's candidacy thread (confirmation, ATS updates, interview, rejection).
When true, company and role must be the real employer and job title from the emails — never leave role empty if any position is mentioned.`

	content, err := c.chatComplete(ctx, classifyTimeout,
		"You judge whether grouped emails are one real job application. Output a single JSON object only.",
		user,
	)
	if err != nil {
		return ApplicationReview{}, err
	}
	raw, err := parseReviewJSON(content)
	if err != nil {
		return ApplicationReview{}, err
	}
	return ApplicationReview{
		IsJobApplication: raw.IsJobApplication,
		Reason:           strings.TrimSpace(raw.Reason),
		Company:          strings.TrimSpace(raw.Company),
		Role:             strings.TrimSpace(raw.Role),
		PipelineStatus:   normalizePipelineStatus(raw.PipelineStatus),
	}, nil
}


func normalizePipelineStatus(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case "applied", "replied", "interview", "offer", "rejected":
		return s
	default:
		return "applied"
	}
}
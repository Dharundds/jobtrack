package classify

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dharundds/jobtrack/internal/store"
	"github.com/rs/zerolog"
)

// Processor classifies unlinked messages and groups them into applications.
type Processor struct {
	store *store.Store
	log   zerolog.Logger
}

func NewProcessor(st *store.Store, log zerolog.Logger) *Processor {
	return &Processor{
		store: st,
		log:   log.With().Str("component", "classify").Logger(),
	}
}

// Run classifies and groups all messages without application_id.
func (p *Processor) Run(ctx context.Context, accountEmail string, opts RunOptions) (linked int, err error) {
	msgs, err := p.store.ListUnlinkedMessages(ctx)
	if err != nil {
		return 0, err
	}
	return p.RunMessages(ctx, accountEmail, msgs, opts)
}

// RunMessages classifies and groups the given unlinked messages (caller picks order/limit).
func (p *Processor) RunMessages(ctx context.Context, accountEmail string, msgs []store.MessageRow, opts RunOptions) (linked int, err error) {
	rules, err := p.store.ListEnabledRules(ctx)
	if err != nil {
		return 0, err
	}
	clRules := make([]Rule, len(rules))
	for i, r := range rules {
		clRules[i] = Rule{
			ID:        r.ID,
			Field:     r.Field,
			Pattern:   r.Pattern,
			EventType: r.EventType,
			Priority:  r.Priority,
			Enabled:   r.Enabled,
		}
	}
	ruleClassifier := NewRuleClassifier(clRules)
	if opts.LLMAvailable && opts.LLM != nil {
		p.log.Info().Int("llm_concurrency", opts.LLMConcurrency).Int("messages", len(msgs)).Msg("classify using llm with rule fallback")
	} else {
		p.log.Info().Int("messages", len(msgs)).Msg("classify using rules")
	}

	if len(msgs) == 0 {
		return 0, nil
	}

	apps, err := p.store.ListApplications(ctx)
	if err != nil {
		return 0, err
	}
	candidates := make([]ApplicationCandidate, len(apps))
	for i, a := range apps {
		candidates[i] = ApplicationCandidate{ID: a.ID, Company: a.Company, Role: a.Role}
	}
	threadMap := make(map[string]int64)
	touchedApps := make(map[int64]struct{})

	results := ClassifyMessages(ctx, p.log, msgs, accountEmail, ruleClassifier, opts)
	llmReviewed := countLLMInMap(results)

	for i, m := range msgs {
		res, ok := results[m.ID]
		if !ok {
			continue
		}
		in := emailInputFromMessage(m, accountEmail)
		res = RefineResult(in, res)
		if res.Company == "" || res.Role == "" {
			if c := ExtractCompany(m.FromAddr, m.Subject, m.BodyExcerpt); res.Company == "" {
				res.Company = c
			}
			if r := ExtractRole(m.Subject, m.BodyExcerpt); res.Role == "" {
				res.Role = r
			}
		}

		appID, ok := PickApplication(m.ThreadID, res.Company, res.Role, threadMap, candidates)
		if !ok {
			if tid, found, _ := p.store.ApplicationIDByThread(ctx, m.ThreadID); found {
				appID = tid
				ok = true
			}
		}
		if !ok && !ShouldLinkToApplication(res) {
			continue
		}
		if !ok {
			status := StatusFromEvent(res.EventType)
			var appliedAt *time.Time
			if res.EventType == EventApplied {
				t := m.ReceivedAt
				appliedAt = &t
			}
			appID, err = p.store.CreateApplication(ctx, res.Company, res.Role, status, appliedAt)
			if err != nil {
				return linked, err
			}
			candidates = append(candidates, ApplicationCandidate{ID: appID, Company: res.Company, Role: res.Role})
		} else {
			_ = p.store.MergeApplicationMeta(ctx, appID, res.Company, res.Role, appliedAtForEvent(res.EventType, m.ReceivedAt))
		}

		if m.ThreadID != "" {
			threadMap[m.ThreadID] = appID
		}
		touchedApps[appID] = struct{}{}

		if err := p.store.UpdateMessageClassification(ctx, m.ID, appID, res.EventType, res.ClassifiedBy); err != nil {
			return linked, err
		}

		exists, err := p.store.EventExistsForMessage(ctx, m.ID)
		if err != nil {
			return linked, err
		}
		if !exists {
			details := map[string]any{}
			if res.Details.MeetingURL != "" {
				details["meeting_url"] = res.Details.MeetingURL
			}
			if res.Details.InterviewAt != nil {
				details["interview_at"] = res.Details.InterviewAt.UTC().Format(time.RFC3339)
			}
			dj, _ := json.Marshal(details)
			if err := p.store.InsertEvent(ctx, appID, m.ID, res.EventType, m.ReceivedAt, string(dj)); err != nil {
				return linked, err
			}
		}

		linked++

		if opts.OnProgress != nil {
			opts.OnProgress(Progress{
				Phase:              "group",
				ClassifyTotal:      len(msgs),
				ClassifyDone:       i + 1,
				LLMReviewed:        llmReviewed,
				ApplicationsLinked: len(touchedApps),
			})
		}
	}

	for id := range touchedApps {
		if err := p.recomputeApplication(ctx, id); err != nil {
			p.log.Warn().Err(err).Int64("application_id", id).Msg("recompute application")
		}
	}

	p.log.Info().Int("messages_linked", linked).Msg("classify and group done")
	return linked, nil
}

func appliedAtForEvent(eventType string, at time.Time) *time.Time {
	if eventType == EventApplied {
		t := at
		return &t
	}
	return nil
}

func (p *Processor) recomputeApplication(ctx context.Context, applicationID int64) error {
	app, err := p.store.GetApplication(ctx, applicationID)
	if err != nil {
		return err
	}
	events, err := p.store.ListApplicationEvents(ctx, applicationID)
	if err != nil {
		return err
	}
	recs := make([]EventRecord, len(events))
	for i, e := range events {
		recs[i] = EventRecord{Type: e.Type, OccurredAt: e.OccurredAt}
	}
	if !app.StatusOverridden {
		status := LatestStatusFromEvents(recs)
		if err := p.store.UpdateApplicationStatusOnly(ctx, applicationID, status); err != nil {
			return err
		}
	}
	return nil
}

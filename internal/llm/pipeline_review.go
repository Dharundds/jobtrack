package llm

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dharundds/jobtrack/internal/classify"
	"github.com/dharundds/jobtrack/internal/store"
	"github.com/rs/zerolog"
)

// ReviewApplicationWithLLM re-checks one pipeline application and re-classifies its linked messages with the LLM.
// Returns true when the application was removed as not a job application.
func ReviewApplicationWithLLM(
	ctx context.Context,
	st *store.Store,
	log zerolog.Logger,
	client *Client,
	accountEmail string,
	applicationID int64,
	llmConcurrency int,
) (deleted bool, err error) {
	app, err := st.GetApplicationFull(ctx, applicationID)
	if err != nil {
		return false, err
	}
	msgs, err := st.ListMessagesForApplication(ctx, applicationID)
	if err != nil {
		return false, err
	}

	review, err := client.ReviewApplication(ctx, accountEmail, *app, msgs)
	if err != nil {
		return false, err
	}

	rules, err := st.ListEnabledRules(ctx)
	if err != nil {
		return false, err
	}
	clRules := make([]classify.Rule, len(rules))
	for i, r := range rules {
		clRules[i] = classify.Rule{
			ID: r.ID, Field: r.Field, Pattern: r.Pattern, EventType: r.EventType, Priority: r.Priority, Enabled: r.Enabled,
		}
	}
	ruleClassifier := classify.NewRuleClassifier(clRules)
	results := classify.ClassifyMessages(ctx, log, msgs, accountEmail, ruleClassifier, classify.RunOptions{
		LLMAvailable:   true,
		LLM:            client,
		LLMConcurrency: llmConcurrency,
	})

	if classify.ShouldRemoveListingAfterLLM(review.IsJobApplication, review.Reason, results, msgs) {
		log.Info().Str("reason", review.Reason).Int64("application_id", applicationID).Msg("llm review: removing listing")
		if err := st.DeleteApplication(ctx, applicationID); err != nil {
			return false, err
		}
		return true, nil
	}

	company, role := classify.BestListingCompanyRole(
		review.Company, review.Role, results, msgs, app.Company, app.Role,
	)

	var status *string
	statusManual := false
	appRow, err := st.GetApplication(ctx, applicationID)
	if err == nil && !appRow.StatusOverridden {
		status = &review.PipelineStatus
	}

	if err := st.PatchApplication(ctx, applicationID, &company, &role, status, nil, statusManual); err != nil {
		return false, err
	}

	for _, m := range msgs {
		res, ok := results[m.ID]
		if !ok {
			continue
		}
		in := classify.EmailInput{
			From: m.FromAddr, To: m.ToAddr, Subject: m.Subject, Body: m.BodyExcerpt, Snippet: m.Snippet,
			Direction: m.Direction, AccountEmail: accountEmail, ReceivedAt: m.ReceivedAt,
		}
		res = classify.RefineResult(in, res)
		if res.Company == "" {
			res.Company = company
		}
		if res.Role == "" {
			res.Role = role
		}
		if err := st.UpdateMessageClassification(ctx, m.ID, applicationID, res.EventType, "llm"); err != nil {
			continue
		}
		details := map[string]any{}
		if res.Details.MeetingURL != "" {
			details["meeting_url"] = res.Details.MeetingURL
		}
		if res.Details.InterviewAt != nil {
			details["interview_at"] = res.Details.InterviewAt.UTC().Format(time.RFC3339)
		}
		dj, _ := json.Marshal(details)
		_ = st.ReplaceEventForMessage(ctx, applicationID, m.ID, res.EventType, m.ReceivedAt, string(dj))
	}

	return false, classify.RecomputeApplicationStatusFromEvents(ctx, st, applicationID)
}
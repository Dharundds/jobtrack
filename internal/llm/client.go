package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/dharundds/jobtrack/internal/classify"
	"github.com/dharundds/jobtrack/internal/store"
)

const (
	pingTimeout     = 8 * time.Second
	classifyTimeout = 90 * time.Second
	maxBodyRunes    = 2500
)

// Client calls an OpenAI-compatible chat API (LM Studio, etc.).
type Client struct {
	baseURL string
	model   string
	apiKey  string
	http    *http.Client
}

func NewClient(cfg store.LLMConfig) *Client {
	return &Client{
		baseURL: NormalizeBaseURL(cfg.BaseURL),
		model:   strings.TrimSpace(cfg.Model),
		apiKey:  cfg.APIKey,
		http: &http.Client{
			Timeout: classifyTimeout,
		},
	}
}

func (c *Client) Configured() bool {
	return c.baseURL != "" && c.model != ""
}

func (c *Client) Ping(ctx context.Context) bool {
	if !c.Configured() {
		return false
	}
	pctx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(pctx, http.MethodGet, c.baseURL+"/models", nil)
	if err != nil {
		return false
	}
	c.setAuth(req)
	res, err := c.http.Do(req)
	if err != nil {
		return false
	}
	defer res.Body.Close()
	_, _ = io.Copy(io.Discard, res.Body)
	return res.StatusCode >= 200 && res.StatusCode < 300
}

type classifyJSON struct {
	EventType   string `json:"event_type"`
	Company     string `json:"company"`
	Role        string `json:"role"`
	MeetingURL  string `json:"meeting_url"`
	InterviewAt string `json:"interview_at"`
}

// Classify one email via chat completions.
func (c *Client) Classify(ctx context.Context, in classify.EmailInput) (classify.Result, error) {
	if !c.Configured() {
		return classify.Result{}, fmt.Errorf("llm not configured")
	}
	pctx, cancel := context.WithTimeout(ctx, classifyTimeout)
	defer cancel()

	bodyText := in.Body
	if bodyText == "" {
		bodyText = in.Snippet
	}
	bodyText = trimRunes(bodyText, maxBodyRunes)

	user := fmt.Sprintf(`Classify this job-search email.
direction: %s
account_email: %s
from: %s
to: %s
subject: %s
body:
%s

Reply with JSON only (no markdown). Keys: event_type, company, role, meeting_url, interview_at (RFC3339 or "").
event_type must be one of: applied, acknowledged, rejected, interview_invite, assessment, offer, recruiter_outreach, other.
Use recruiter_outreach only for personal recruiter mail. Use other for Instahyre/Naukri job blasts, "good fit" recommendations, and mail where the user did not apply.`,
		in.Direction, in.AccountEmail, in.From, in.To, in.Subject, bodyText)

	content, err := c.chatComplete(pctx, classifyTimeout,
		"You extract structured job-application signals from email. Output a single JSON object only.",
		user,
	)
	if err != nil {
		return classify.Result{}, err
	}
	parsed, err := parseClassifyJSON(content)
	if err != nil {
		return classify.Result{}, err
	}

	eventType := normalizeEventType(parsed.EventType)
	details := classify.EventDetails{MeetingURL: strings.TrimSpace(parsed.MeetingURL)}
	if parsed.InterviewAt != "" {
		if t, e := time.Parse(time.RFC3339, parsed.InterviewAt); e == nil {
			details.InterviewAt = &t
		}
	}

	return classify.Result{
		EventType:    eventType,
		Company:      strings.TrimSpace(parsed.Company),
		Role:         strings.TrimSpace(parsed.Role),
		Details:      details,
		ClassifiedBy: "llm",
	}, nil
}

func (c *Client) chatComplete(ctx context.Context, timeout time.Duration, system, user string) (string, error) {
	if !c.Configured() {
		return "", fmt.Errorf("llm not configured")
	}
	pctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	makePayload := func(jsonMode bool) (map[string]any, error) {
		p := map[string]any{
			"model": c.model,
			"messages": []map[string]string{
				{"role": "system", "content": system},
				{"role": "user", "content": user},
			},
			"temperature":      0.1,
			"max_tokens":       2048,
			"reasoning_effort": "none",
		}
		if jsonMode {
			p["response_format"] = map[string]string{"type": "json_object"}
		}
		return p, nil
	}

	post := func(jsonMode bool) ([]byte, int, error) {
		payload, err := makePayload(jsonMode)
		if err != nil {
			return nil, 0, err
		}
		raw, err := json.Marshal(payload)
		if err != nil {
			return nil, 0, err
		}
		req, err := http.NewRequestWithContext(pctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(raw))
		if err != nil {
			return nil, 0, err
		}
		req.Header.Set("Content-Type", "application/json")
		c.setAuth(req)
		res, err := c.http.Do(req)
		if err != nil {
			return nil, 0, err
		}
		defer res.Body.Close()
		respBody, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
		if err != nil {
			return nil, res.StatusCode, err
		}
		return respBody, res.StatusCode, nil
	}

	respBody, status, err := post(true)
	if err != nil {
		return "", err
	}
	if status == http.StatusBadRequest {
		respBody, status, err = post(false)
		if err != nil {
			return "", err
		}
	}
	if status < 200 || status >= 300 {
		return "", fmt.Errorf("llm http %d: %s", status, trimRunes(string(respBody), 200))
	}

	var chatResp struct {
		Choices []struct {
			FinishReason string `json:"finish_reason"`
			Message      struct {
				Content          string `json:"content"`
				ReasoningContent string `json:"reasoning_content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("decode chat response: %w", err)
	}
	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("empty llm response")
	}
	choice := chatResp.Choices[0]
	content := strings.TrimSpace(choice.Message.Content)
	reasoning := strings.TrimSpace(choice.Message.ReasoningContent)
	if content == "" {
		content = reasoning
	} else if reasoning != "" && extractJSONObject(extractJSONPayload(content)) == "" {
		content = content + "\n" + reasoning
	}
	if content == "" {
		return "", fmt.Errorf("empty llm message content")
	}
	if choice.FinishReason == "length" {
		// Still try parse; may be partial JSON inside thinking.
		if extractJSONObject(extractJSONPayload(content)) == "" {
			return "", fmt.Errorf("llm response truncated (finish_reason=length); increase max_tokens or disable thinking")
		}
	}
	return content, nil
}

func (c *Client) setAuth(req *http.Request) {
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
}

func parseClassifyJSON(content string) (classifyJSON, error) {
	return parseJSONContent[classifyJSON](content, "llm json")
}

func stripCodeFence(s string) string {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "```") {
		return s
	}
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}

func normalizeEventType(raw string) string {
	raw = strings.ToLower(strings.TrimSpace(raw))
	switch raw {
	case classify.EventApplied, classify.EventAcknowledged, classify.EventRejected,
		classify.EventInterviewInvite, classify.EventAssessment, classify.EventOffer,
		classify.EventRecruiterOutreach, classify.EventOther:
		return raw
	default:
		return classify.EventOther
	}
}

func trimRunes(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max])
}

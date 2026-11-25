package llm

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// flexBool unmarshals booleans that models sometimes emit as strings or 0/1.
type flexBool bool

func (b *flexBool) UnmarshalJSON(data []byte) error {
	s := strings.TrimSpace(string(data))
	if len(s) == 0 || s == "null" {
		*b = false
		return nil
	}
	var asBool bool
	if err := json.Unmarshal(data, &asBool); err == nil {
		*b = flexBool(asBool)
		return nil
	}
	var asInt int
	if err := json.Unmarshal(data, &asInt); err == nil {
		*b = flexBool(asInt != 0)
		return nil
	}
	var asStr string
	if err := json.Unmarshal(data, &asStr); err == nil {
		asStr = strings.TrimSpace(strings.ToLower(asStr))
		switch asStr {
		case "true", "yes", "1":
			*b = true
		default:
			*b = false
		}
		return nil
	}
	if n, err := strconv.ParseBool(s); err == nil {
		*b = flexBool(n)
		return nil
	}
	*b = false
	return nil
}

func (b flexBool) Bool() bool { return bool(b) }

type reviewJSON struct {
	IsJobApplication bool
	Reason           string
	Company          string
	Role             string
	PipelineStatus   string
}

func parseReviewJSON(raw string) (reviewJSON, error) {
	payload := extractJSONPayload(raw)
	if payload == "" {
		snippet := trimRunes(strings.TrimSpace(raw), 120)
		if snippet == "" {
			return reviewJSON{}, fmt.Errorf("parse llm review json: model returned empty content")
		}
		return reviewJSON{}, fmt.Errorf("parse llm review json: no JSON object in response: %q", snippet)
	}

	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(payload), &m); err != nil {
		return reviewJSON{}, fmt.Errorf("parse llm review json: %w", err)
	}

	var out reviewJSON
	for _, key := range []string{"is_job_application", "isJobApplication", "job_application"} {
		if v, ok := m[key]; ok {
			var fb flexBool
			if err := json.Unmarshal(v, &fb); err == nil {
				out.IsJobApplication = fb.Bool()
				break
			}
		}
	}
	for _, key := range []string{"reason", "explanation"} {
		if v, ok := m[key]; ok {
			_ = json.Unmarshal(v, &out.Reason)
		}
	}
	for _, key := range []string{"company", "employer", "company_name"} {
		if v, ok := m[key]; ok {
			var s string
			if err := json.Unmarshal(v, &s); err == nil && strings.TrimSpace(s) != "" {
				out.Company = s
				break
			}
		}
	}
	for _, key := range []string{"role", "job_title", "title", "position"} {
		if v, ok := m[key]; ok {
			var s string
			if err := json.Unmarshal(v, &s); err == nil && strings.TrimSpace(s) != "" {
				out.Role = s
				break
			}
		}
	}
	for _, key := range []string{"pipeline_status", "pipelineStatus", "status"} {
		if v, ok := m[key]; ok {
			var s string
			if err := json.Unmarshal(v, &s); err == nil && strings.TrimSpace(s) != "" {
				out.PipelineStatus = s
				break
			}
		}
	}
	return out, nil
}
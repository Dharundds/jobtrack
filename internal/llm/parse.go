package llm

import (
	"encoding/json"
	"fmt"
	"strings"
)

// extractJSONPayload strips model scaffolding and returns the first JSON object substring.
func extractJSONPayload(raw string) string {
	s := strings.TrimSpace(raw)
	s = stripCodeFence(s)
	s = stripThinkingBlocks(s)
	if obj := extractJSONObject(s); obj != "" {
		return obj
	}
	return s
}

const (
	thinkOpen  = "<" + "think" + ">"
	thinkClose = "<" + "/" + "think" + ">"
)

// stripThinkingBlocks removes Qwen-style reasoning sections from model output.
func stripThinkingBlocks(s string) string {
	for {
		lower := strings.ToLower(s)
		start := strings.Index(lower, thinkOpen)
		if start < 0 {
			break
		}
		relEnd := strings.Index(lower[start:], thinkClose)
		if relEnd < 0 {
			s = strings.TrimSpace(s[start+len(thinkOpen):])
			break
		}
		end := start + relEnd + len(thinkClose)
		s = strings.TrimSpace(s[:start] + s[end:])
	}
	return strings.TrimSpace(s)
}

func extractJSONObject(s string) string {
	start := strings.Index(s, "{")
	if start < 0 {
		return ""
	}
	depth := 0
	inString := false
	escape := false
	for i := start; i < len(s); i++ {
		c := s[i]
		if inString {
			if escape {
				escape = false
				continue
			}
			if c == '\\' {
				escape = true
				continue
			}
			if c == '"' {
				inString = false
			}
			continue
		}
		switch c {
		case '"':
			inString = true
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return s[start : i+1]
			}
		}
	}
	return ""
}

func parseJSONContent[T any](raw string, label string) (T, error) {
	var zero T
	payload := extractJSONPayload(raw)
	if payload == "" {
		snippet := trimRunes(strings.TrimSpace(raw), 120)
		if snippet == "" {
			return zero, fmt.Errorf("parse %s: model returned empty content (try disabling thinking mode or raising max tokens)", label)
		}
		return zero, fmt.Errorf("parse %s: no JSON object in response: %q", label, snippet)
	}
	var out T
	if err := json.Unmarshal([]byte(payload), &out); err != nil {
		return zero, fmt.Errorf("parse %s: %w", label, err)
	}
	return out, nil
}
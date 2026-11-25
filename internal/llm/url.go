package llm

import "strings"

// NormalizeBaseURL ensures an OpenAI-compatible /v1 prefix (e.g. LM Studio).
func NormalizeBaseURL(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimRight(raw, "/")
	if raw == "" {
		return ""
	}
	if strings.HasSuffix(raw, "/v1") {
		return raw
	}
	return raw + "/v1"
}
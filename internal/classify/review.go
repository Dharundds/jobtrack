package classify

import (
	"strings"

	"github.com/dharundds/jobtrack/internal/store"
)

// NeedsLLMReview is true when any linked message was not classified by the LLM.
func NeedsLLMReview(msgs []store.MessageRow) bool {
	for _, m := range msgs {
		if strings.ToLower(strings.TrimSpace(m.ClassifiedBy)) != "llm" {
			return true
		}
	}
	return false
}
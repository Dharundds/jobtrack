package classify

import (
	"strings"

	"github.com/dharundds/jobtrack/internal/store"
)

// IsPlausibleListingRole rejects rule/LLM junk titles that should not be stored on an application row.
func IsPlausibleListingRole(role string) bool {
	role = strings.TrimSpace(role)
	if role == "" {
		return false
	}
	lower := strings.ToLower(role)
	bad := []string{
		"a remote", "your submission", "a llm", "your interest", "relative", "bottom",
		"the role", "for the role", "the position", "unknown",
	}
	for _, b := range bad {
		if lower == b || strings.HasPrefix(lower, b+" ") {
			return false
		}
	}
	if strings.HasPrefix(lower, "a ") && len(role) < 14 {
		return false
	}
	return len(role) >= 2
}

func betterRole(candidate, current string) bool {
	candidate = strings.TrimSpace(candidate)
	if !IsPlausibleListingRole(candidate) {
		return false
	}
	if !IsPlausibleListingRole(current) {
		return true
	}
	return len(candidate) > len(current)
}

func betterCompany(candidate, current string) bool {
	candidate = sanitizeCompany(candidate)
	if candidate == "" {
		return false
	}
	if strings.TrimSpace(current) == "" {
		return true
	}
	return len(candidate) > len(sanitizeCompany(current))
}

// BestListingCompanyRole merges application-review JSON with per-message LLM results.
func BestListingCompanyRole(
	reviewCompany, reviewRole string,
	results map[int64]Result,
	msgs []store.MessageRow,
	fallbackCompany, fallbackRole string,
) (company, role string) {
	company = strings.TrimSpace(reviewCompany)
	role = strings.TrimSpace(reviewRole)

	for _, m := range msgs {
		res, ok := results[m.ID]
		if !ok {
			continue
		}
		if betterCompany(res.Company, company) {
			company = sanitizeCompany(res.Company)
		}
		if betterRole(res.Role, role) {
			role = sanitizeRole(res.Role)
		}
	}

	if company == "" {
		company = sanitizeCompany(fallbackCompany)
	}
	if company == "" {
		for i := len(msgs) - 1; i >= 0; i-- {
			m := msgs[i]
			c := ExtractCompany(m.FromAddr, m.Subject, m.BodyExcerpt)
			if c != "" {
				company = sanitizeCompany(c)
				break
			}
		}
	}

	if !IsPlausibleListingRole(role) {
		role = ""
	}
	if role == "" && IsPlausibleListingRole(fallbackRole) {
		role = sanitizeRole(fallbackRole)
	}
	if role == "" {
		for i := len(msgs) - 1; i >= 0; i-- {
			m := msgs[i]
			r := ExtractRole(m.Subject, m.BodyExcerpt)
			if IsPlausibleListingRole(r) {
				role = r
				break
			}
		}
	}
	return company, role
}

// ShouldRemoveListingAfterLLM is true when the listing is not a real job application.
func ShouldRemoveListingAfterLLM(
	isJob bool,
	reason string,
	results map[int64]Result,
	msgs []store.MessageRow,
) bool {
	if !isJob {
		return true
	}
	r := strings.ToLower(strings.TrimSpace(reason))
	for _, p := range []string{
		"not a job", "not job", "not an application", "newsletter", "marketing",
		"unrelated", "should be removed", "should not", "not belong", "linkedin noise",
		"social notification", "shopping", "leetcode contest", "promotional",
		"instahyre", "job alert", "job recommendation", "did not apply", "never applied",
	} {
		if strings.Contains(r, p) {
			return true
		}
	}
	if len(msgs) == 0 {
		return true
	}
	if len(msgs) > 0 {
		allAlerts := true
		for _, m := range msgs {
			if !isJobBoardAlert(m.FromAddr, m.Subject, m.BodyExcerpt, m.Snippet) {
				allAlerts = false
				break
			}
		}
		if allAlerts {
			return true
		}
	}
	hasSignal := false
	for _, m := range msgs {
		if isJobBoardAlert(m.FromAddr, m.Subject, m.BodyExcerpt, m.Snippet) {
			continue
		}
		res, ok := results[m.ID]
		if !ok {
			continue
		}
		switch res.EventType {
		case EventApplied, EventAcknowledged, EventRejected, EventInterviewInvite, EventAssessment, EventOffer:
			hasSignal = true
		case EventRecruiterOutreach:
			if hasApplicationSubmittedSignal(m.FromAddr, m.Subject, m.BodyExcerpt, m.Snippet) {
				hasSignal = true
			}
		}
		if hasSignal {
			break
		}
		if IsPlausibleListingRole(res.Role) || sanitizeCompany(res.Company) != "" {
			hasSignal = true
			break
		}
		if isJobApplicationContext(m.FromAddr, m.Subject, m.BodyExcerpt, m.Snippet) {
			hasSignal = true
			break
		}
	}
	return !hasSignal
}
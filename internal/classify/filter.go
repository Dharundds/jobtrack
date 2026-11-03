package classify

import "strings"

// RefineResult adjusts rule matches that are too broad (marketing mail, LinkedIn social, etc.).
func RefineResult(in EmailInput, res Result) Result {
	if res.EventType == EventRejected && !isPlausibleJobRejection(in.From, in.Subject, in.Body, in.Snippet) {
		res.EventType = EventOther
	}
	if isLinkedInSocialNoise(in) {
		res.EventType = EventOther
	}
	if res.EventType == EventRecruiterOutreach && isLinkedInSocialNoise(in) {
		res.EventType = EventOther
	}
	return res
}

func isPlausibleJobRejection(from, subject, body, snippet string) bool {
	if isJobApplicationContext(from, subject, body, snippet) {
		return true
	}
	combined := strings.ToLower(subject + "\n" + body + "\n" + snippet)
	phrases := []string{
		"not moving forward",
		"will not be moving forward",
		"will not be proceeding",
		"regret to inform",
		"not selected",
		"decided to pursue other",
		"other candidates",
	}
	for _, p := range phrases {
		if strings.Contains(combined, p) {
			return true
		}
	}
	return false
}

// ShouldLinkToApplication decides whether a message should create or join an application row.
// EventOther is only linked when merging into an existing thread/application (handled in processor).
func ShouldLinkToApplication(res Result) bool {
	switch res.EventType {
	case EventOther:
		return false
	case EventApplied, EventAcknowledged, EventRejected, EventInterviewInvite, EventAssessment, EventOffer, EventRecruiterOutreach:
		return hasMinimumIdentity(res) || isLikelyEmployerMail(res.EventType, res.Company)
	default:
		return false
	}
}

func hasMinimumIdentity(res Result) bool {
	return strings.TrimSpace(res.Company) != "" || strings.TrimSpace(res.Role) != ""
}

func isLikelyEmployerMail(eventType, company string) bool {
	if eventType == EventApplied || eventType == EventAcknowledged || eventType == EventRejected {
		return company != ""
	}
	return false
}

func isJobApplicationContext(from, subject, body, snippet string) bool {
	if isLikelyATSFrom(from) {
		return true
	}
	combined := strings.ToLower(subject + "\n" + body + "\n" + snippet)
	hints := []string{
		"application", "applied", "apply for", "you applied", "your application",
		"resume", " résumé", "candidacy", "candidate", "position", " requisition",
		"opening", "job you", "submitted your", "interview process", "hiring process",
		"thank you for applying", "your interest in the", "for the role", "for the position",
		"regret to inform you that we will not", "not moving forward with your application",
	}
	for _, h := range hints {
		if strings.Contains(combined, h) {
			return true
		}
	}
	return false
}

func isLikelyATSFrom(from string) bool {
	from = strings.ToLower(from)
	markers := []string{
		"greenhouse", "lever.co", "ashbyhq", "myworkday.com", "workday",
		"smartrecruiters", "icims", "jobs-noreply@linkedin.com",
	}
	for _, m := range markers {
		if strings.Contains(from, m) {
			return true
		}
	}
	return false
}

func isLinkedInSocialNoise(in EmailInput) bool {
	from := strings.ToLower(in.From)
	if !strings.Contains(from, "linkedin.com") {
		return false
	}
	// Easy Apply / job application mail uses jobs-noreply; keep those unless clearly social.
	if strings.Contains(from, "jobs-noreply") {
		return false
	}
	combined := strings.ToLower(in.Subject + "\n" + in.Body + "\n" + in.Snippet)
	social := []string{
		"invitation to connect",
		"invited you to connect",
		"wants to connect",
		"accepted your invitation",
		"connection request",
		"sent you an invitation",
		"people you may know",
		"viewed your profile",
		"views of your profile",
		"new message",
		"mentioned you in",
		"followed you",
		"job alert",
		"recommended jobs",
		"jobs for you",
		"top job picks",
	}
	for _, p := range social {
		if strings.Contains(combined, p) {
			return true
		}
	}
	// Generic LinkedIn notification without job-application signal → not an application.
	if !isJobApplicationContext(in.From, in.Subject, in.Body, in.Snippet) {
		return true
	}
	return false
}

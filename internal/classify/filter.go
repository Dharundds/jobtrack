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
	if isJobBoardAlert(in.From, in.Subject, in.Body, in.Snippet) {
		res.EventType = EventOther
	}
	if res.EventType == EventApplied && !hasApplicationSubmittedSignal(in.From, in.Subject, in.Body, in.Snippet) {
		if isJobBoardAlert(in.From, in.Subject, in.Body, in.Snippet) ||
			strings.Contains(strings.ToLower(in.Subject+"\n"+in.Body), "look like a good fit") {
			res.EventType = EventOther
		}
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
		"not move forward",
		"decided not to move forward",
		"will not be moving forward",
		"will not be proceeding",
		"regret to inform",
		"not selected",
		"decided to pursue other",
		"other candidates",
	}
	if strings.Contains(strings.ToLower(from), "jobs-noreply@linkedin.com") {
		combined := strings.ToLower(subject + "\n" + body + "\n" + snippet)
		if strings.Contains(combined, "not move forward") || strings.Contains(combined, "not moving forward") {
			return true
		}
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
	if isJobBoardAlert(from, subject, body, snippet) {
		return false
	}
	if isLikelyATSFrom(from) {
		return hasApplicationSubmittedSignal(from, subject, body, snippet)
	}
	return hasApplicationSubmittedSignal(from, subject, body, snippet)
}

func hasApplicationSubmittedSignal(from, subject, body, snippet string) bool {
	combined := strings.ToLower(subject + "\n" + body + "\n" + snippet)
	submitted := []string{
		"thank you for applying",
		"thanks for applying",
		"you applied",
		"your application was sent",
		"your application has been",
		"we received your application",
		"we've received your application",
		"confirm that we have received your application",
		"application submitted",
		"successfully applied",
		"confirming your application",
		"application received",
		"submitted your application",
		"update on your application",
		"regret to inform you that we will not",
		"not moving forward with your application",
		"interview", "assessment", "offer letter",
	}
	if strings.Contains(combined, "application for") && strings.Contains(combined, "received") &&
		strings.Contains(combined, "thank you") {
		return true
	}
	for _, p := range submitted {
		if strings.Contains(combined, p) {
			return true
		}
	}
	fromLower := strings.ToLower(from)
	for _, ats := range []string{"greenhouse", "lever.co", "ashbyhq", "myworkday.com", "workday", "icims", "smartrecruiters", "kekamail.com"} {
		if strings.Contains(fromLower, ats) {
			return strings.Contains(combined, "application")
		}
	}
	return false
}

// isJobBoardAlert is a job recommendation blast (Instahyre, Naukri, etc.), not an application you submitted.
func isJobBoardAlert(from, subject, body, snippet string) bool {
	fromLower := strings.ToLower(from)
	combined := strings.ToLower(subject + "\n" + body + "\n" + snippet)

	if hasApplicationSubmittedSignal(from, subject, body, snippet) {
		return false
	}

	boardSenders := []string{
		"instahyre.com", "naukri.com", "foundit.in", "shine.com", "iimjobs.com",
		"indeed.com", "glassdoor.com", "wellfound.com", "cutshort.io", "hirist.com",
	}
	for _, b := range boardSenders {
		if strings.Contains(fromLower, b) {
			return true
		}
	}

	alertPhrases := []string{
		"has just posted",
		"look like a good fit",
		"looks like a good fit",
		"based on your profile",
		"based on your instahyre",
		"you might be interested",
		"recommended for you",
		"jobs matching your profile",
		"new job matching",
		"matching jobs for you",
	}
	for _, p := range alertPhrases {
		if strings.Contains(combined, p) {
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

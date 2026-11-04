package classify

import (
	"net/mail"
	"regexp"
	"strings"
	"unicode"
)

var (
	reRoleForThePosition = regexp.MustCompile(`(?i)for the\s+(.+?)\s+(?:position|role)\b`)
	reRoleForPosition    = regexp.MustCompile(`(?i)(?:^|[^\w])for\s+(.+?)\s+(?:position|role)\b`)
	reRoleApplication    = regexp.MustCompile(`(?i)application:\s*(.+)`)
	reRoleTitle          = regexp.MustCompile(`(?i)(?:role|position|job title):\s*(.+)`)
	reAtCompanySubject   = regexp.MustCompile(`(?i)\bat\s+([A-Z][A-Za-z0-9 &.\-'/]{0,58}[A-Za-z0-9])\b`)
	reFromCompany        = regexp.MustCompile(`(?i)\bfrom\s+([A-Z][A-Za-z0-9 &.\-'/]{0,58}[A-Za-z0-9])\s`)
	reApplyingTo         = regexp.MustCompile(`(?i)(?:thank you for applying to|thanks for applying to|you applied to|your application was sent to)\s+([A-Za-z0-9][A-Za-z0-9 &.\-'/]{0,58})`)
	reApplyingForCompany = regexp.MustCompile(`(?i)(?:application for|applied for|your application to)\s+([A-Za-z0-9][A-Za-z0-9 &.\-'/]{0,58})`)
	reReceivedAppAt      = regexp.MustCompile(`(?i)received your application(?: for)?(?: at| with)?\s+([A-Za-z0-9][A-Za-z0-9 &.\-'/]{0,58})`)
	reWorkdayDash        = regexp.MustCompile(`(?i)^(.+?)\s*[-–]\s*workday\b`)
	reWorkdayAt          = regexp.MustCompile(`(?i)^(?:my)?workday@(.+)$`)
	reWorkdayNotification = regexp.MustCompile(`(?i)^(.+?)\s+workday\s+notification\b`)
	reSubjectBrandPrefix  = regexp.MustCompile(`(?i)^([A-Za-z0-9][A-Za-z0-9 &.\-'/]{0,40}?)\s*[-–]\s*(?:invites you|your screening)`)
	reJoiningUsAs         = regexp.MustCompile(`(?i)(?:interest in )?joining us as (?:a|an)\s+([A-Za-z0-9][^.!?\n]{2,60})`)
	reScreeningForRole    = regexp.MustCompile(`(?i)your screening for\s+(.+?)\s+has been`)
)

var atsDomains = map[string]string{
	"greenhouse.io":       "Greenhouse",
	"lever.co":            "Lever",
	"ashbyhq.com":         "Ashby",
	"myworkday.com":       "Workday",
	"smartrecruiters.com": "SmartRecruiters",
	"icims.com":           "iCIMS",
	"linkedin.com":        "LinkedIn",
	"naukri.com":          "Naukri",
	"indeed.com":          "Indeed",
	"wellfound.com":       "Wellfound",
}

var companyBlocklist = []string{
	"the next steps",
	"the position",
	"your application",
	"our team",
	"this role",
	"this position",
	"the role",
	"the company",
	"if your",
	"isn't selected",
	"is not selected",
	"we have received",
	"we've received",
	"thank you",
}

// ExtractCompany best-effort from From header, subject, and body.
func ExtractCompany(from, subject, body string) string {
	subject = strings.TrimSpace(subject)
	if c := extractCompanyFromText(subject); c != "" {
		return c
	}
	if c := companyFromDisplayName(from); c != "" {
		return c
	}
	if c := extractCompanyFromText(body); c != "" {
		return c
	}
	if addr := emailAddress(from); addr != "" {
		if c := companyFromEmployerDomain(addr); c != "" {
			return c
		}
	}
	return ""
}

func extractCompanyFromText(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	for _, re := range []*regexp.Regexp{
		reSubjectBrandPrefix,
		reApplyingTo,
		reApplyingForCompany,
		reReceivedAppAt,
		reAtCompanySubject,
	} {
		if m := re.FindStringSubmatch(text); len(m) > 1 {
			if c := sanitizeCompany(m[1]); c != "" {
				return c
			}
		}
	}
	if m := reFromCompany.FindStringSubmatch(text + " "); len(m) > 1 {
		if c := sanitizeCompany(m[1]); c != "" {
			return c
		}
	}
	return ""
}

func companyFromDisplayName(from string) string {
	name := displayName(from)
	if name == "" {
		return ""
	}
	for _, re := range []*regexp.Regexp{reWorkdayAt, reWorkdayDash, reWorkdayNotification} {
		if m := re.FindStringSubmatch(name); len(m) > 1 {
			if c := sanitizeCompany(m[1]); c != "" {
				return c
			}
		}
	}
	if strings.Contains(strings.ToLower(name), "workday") {
		name = regexp.MustCompile(`(?i)\s*workday\s*notification\s*`).ReplaceAllString(name, "")
		name = regexp.MustCompile(`(?i)\s*[-–@]\s*workday\s*`).ReplaceAllString(name, "")
		name = strings.TrimSpace(name)
	}
	if strings.Contains(strings.ToLower(name), "noreply") ||
		strings.Contains(strings.ToLower(name), "no-reply") ||
		strings.Contains(strings.ToLower(name), "do not reply") {
		return ""
	}
	return sanitizeCompany(name)
}

func sanitizeCompany(raw string) string {
	c := trimCompany(raw)
	if c == "" || !isPlausibleCompanyName(c) {
		return ""
	}
	if isATSBrandName(c) {
		return ""
	}
	return c
}

func isPlausibleCompanyName(s string) bool {
	s = strings.TrimSpace(s)
	if len(s) < 2 || len(s) > 60 {
		return false
	}
	lower := strings.ToLower(s)
	for _, bad := range companyBlocklist {
		if lower == bad || strings.HasPrefix(lower, bad+" ") {
			return false
		}
	}
	if strings.ContainsAny(s, ".?!;") {
		return false
	}
	if strings.EqualFold(s, "the") || strings.HasPrefix(lower, "the ") {
		return false
	}
	words := strings.Fields(s)
	if len(words) > 7 {
		return false
	}
	// Reject sentence-like fragments (many small grammar words).
	stop := 0
	for _, w := range words {
		switch strings.ToLower(w) {
		case "if", "your", "isn't", "isnt", "selected", "application", "are", "and", "or", "but", "we", "you", "our", "this", "that", "with", "for", "from", "next", "steps":
			stop++
		}
	}
	if stop >= 2 && len(words) >= 3 {
		return false
	}
	// Require at least one letter and mostly printable name characters.
	letters := 0
	for _, r := range s {
		if unicode.IsLetter(r) {
			letters++
		}
	}
	return letters >= 2
}

func isATSBrandName(name string) bool {
	n := strings.TrimSpace(strings.ToLower(name))
	for _, brand := range atsDomains {
		if n == strings.ToLower(brand) {
			return true
		}
	}
	return n == "greenhouse" || n == "lever" || n == "ashby" || n == "workday" ||
		strings.Contains(n, "workday notification") || strings.HasPrefix(n, "myworkday")
}

func ExtractRole(subject, body string) string {
	for _, text := range []string{subject, body} {
		if r := extractRoleFromText(text); r != "" {
			return r
		}
	}
	return ""
}

func extractRoleFromText(text string) string {
	for _, re := range []*regexp.Regexp{
		reScreeningForRole,
		reJoiningUsAs,
		reRoleTitle,
		reRoleApplication,
		reRoleForThePosition,
		reRoleForPosition,
	} {
		if m := re.FindStringSubmatch(text); len(m) > 1 {
			if r := sanitizeRole(m[1]); r != "" {
				return r
			}
		}
	}
	return ""
}

func sanitizeRole(raw string) string {
	r := trimRole(raw)
	if len(r) < 2 {
		return ""
	}
	if strings.EqualFold(r, "the") || strings.EqualFold(r, "a") || strings.EqualFold(r, "an") {
		return ""
	}
	if strings.ContainsAny(r, ".?!;") && len(r) > 40 {
		return ""
	}
	return r
}

func displayName(from string) string {
	addr, err := mail.ParseAddress(from)
	if err != nil {
		if i := strings.Index(from, "<"); i > 0 {
			return strings.TrimSpace(from[:i])
		}
		return strings.TrimSpace(from)
	}
	return strings.TrimSpace(addr.Name)
}

func emailAddress(from string) string {
	addr, err := mail.ParseAddress(from)
	if err != nil {
		if i := strings.Index(from, "<"); i >= 0 {
			if j := strings.Index(from, ">"); j > i {
				return strings.ToLower(strings.TrimSpace(from[i+1 : j]))
			}
		}
		return strings.ToLower(strings.TrimSpace(from))
	}
	return strings.ToLower(addr.Address)
}

// companyFromEmployerDomain uses the domain label only when it is not a known ATS host.
func companyFromEmployerDomain(email string) string {
	at := strings.LastIndex(email, "@")
	if at < 0 {
		return ""
	}
	domain := strings.ToLower(email[at+1:])
	for suffix := range atsDomains {
		if domain == suffix || strings.HasSuffix(domain, "."+suffix) {
			return ""
		}
	}
	if strings.Contains(domain, "workday") || strings.Contains(domain, "greenhouse") {
		return ""
	}
	parts := strings.Split(domain, ".")
	if len(parts) >= 3 && (parts[0] == "jobalerts" || parts[0] == "mail" || parts[0] == "email") {
		base := parts[1]
		if base != "" && base != "co" && base != "com" {
			name := strings.ToUpper(base[:1]) + base[1:]
			if sanitized := sanitizeCompany(name); sanitized != "" {
				return sanitized
			}
		}
	}
	if len(parts) < 2 {
		return ""
	}
	base := parts[len(parts)-2]
	if base == "co" || base == "com" || base == "org" || base == "mail" || base == "email" {
		return ""
	}
	c := strings.ToUpper(base[:1]) + base[1:]
	return sanitizeCompany(c)
}

func trimCompany(s string) string {
	s = strings.TrimSpace(s)
	// Cut at sentence boundaries or parenthetical ATS junk.
	if i := strings.IndexAny(s, ".?!\n("); i > 0 {
		s = strings.TrimSpace(s[:i])
	}
	s = strings.Trim(s, `"'.,;:`)
	if len(s) > 80 {
		s = s[:80]
	}
	return strings.TrimSpace(s)
}

func trimRole(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.IndexAny(s, ".?!\n"); i > 0 {
		s = strings.TrimSpace(s[:i])
	}
	s = strings.Trim(s, `"'.,;:`)
	if len(s) > 120 {
		s = s[:120]
	}
	return strings.TrimSpace(s)
}

// DetectDirection returns in|out using SENT label or From matching account email.
func DetectDirection(from, accountEmail string, labelIDs []string) string {
	for _, l := range labelIDs {
		if l == "SENT" {
			return "out"
		}
	}
	acc := strings.ToLower(strings.TrimSpace(accountEmail))
	fromLower := strings.ToLower(from)
	if acc != "" && strings.Contains(fromLower, acc) {
		return "out"
	}
	return "in"
}

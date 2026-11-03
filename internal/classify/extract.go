package classify

import (
	"net/mail"
	"regexp"
	"strings"
)

var (
	reRoleForPosition = regexp.MustCompile(`(?i)(?:for the|for)\s+(.+?)\s+(?:position|role)\b`)
	reRoleApplication = regexp.MustCompile(`(?i)application:\s*(.+)`)
	reRoleTitle       = regexp.MustCompile(`(?i)(?:role|position|job title):\s*(.+)`)
	reAtCompany       = regexp.MustCompile(`(?i)(?:at|@|with)\s+([A-Za-z0-9][A-Za-z0-9 &.\-'/]{1,60})`)
	reFromCompany     = regexp.MustCompile(`(?i)from\s+([A-Za-z0-9][A-Za-z0-9 &.\-'/]{1,60})`)
	reApplyingTo      = regexp.MustCompile(`(?i)(?:thank you for applying to|thanks for applying to|you applied to|application (?:to|for)|your application was sent to)\s+([A-Za-z0-9][A-Za-z0-9 &.\-'/]{1,80})`)
)

var atsDomains = map[string]string{
	"greenhouse.io":     "Greenhouse",
	"lever.co":          "Lever",
	"ashbyhq.com":       "Ashby",
	"myworkday.com":     "Workday",
	"smartrecruiters.com": "SmartRecruiters",
	"icims.com":         "iCIMS",
	"linkedin.com":      "LinkedIn",
	"naukri.com":        "Naukri",
	"indeed.com":        "Indeed",
	"wellfound.com":     "Wellfound",
}

// ExtractCompany best-effort from From header, domain, and subject/body.
func ExtractCompany(from, subject, body string) string {
	combined := subject + "\n" + body
	if m := reApplyingTo.FindStringSubmatch(combined); len(m) > 1 {
		if c := trimCompany(m[1]); c != "" && !isATSBrandName(c) {
			return c
		}
	}
	if m := reAtCompany.FindStringSubmatch(combined); len(m) > 1 {
		if c := trimCompany(m[1]); c != "" && !isATSBrandName(c) {
			return c
		}
	}
	if m := reFromCompany.FindStringSubmatch(combined); len(m) > 1 {
		if c := trimCompany(m[1]); c != "" && !isATSBrandName(c) {
			return c
		}
	}
	if name := displayName(from); name != "" && !strings.Contains(strings.ToLower(name), "noreply") {
		clean := trimCompany(name)
		if clean != "" && !isATSBrandName(clean) {
			return clean
		}
	}
	if addr := emailAddress(from); addr != "" {
		if c := companyFromDomain(addr); c != "" && !isATSBrandName(c) {
			return c
		}
	}
	return ""
}

func isATSBrandName(name string) bool {
	n := strings.TrimSpace(strings.ToLower(name))
	for _, brand := range atsDomains {
		if n == strings.ToLower(brand) {
			return true
		}
	}
	return n == "greenhouse" || n == "lever" || n == "ashby" || n == "workday"
}

func ExtractRole(subject, body string) string {
	combined := subject + "\n" + body
	for _, re := range []*regexp.Regexp{reRoleForPosition, reRoleApplication, reRoleTitle} {
		if m := re.FindStringSubmatch(combined); len(m) > 1 {
			return trimRole(m[1])
		}
	}
	return ""
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

func companyFromDomain(email string) string {
	at := strings.LastIndex(email, "@")
	if at < 0 {
		return ""
	}
	domain := strings.ToLower(email[at+1:])
	for suffix, name := range atsDomains {
		if domain == suffix || strings.HasSuffix(domain, "."+suffix) {
			return name
		}
	}
	parts := strings.Split(domain, ".")
	if len(parts) >= 2 {
		base := parts[len(parts)-2]
		if base == "co" || base == "com" || base == "org" {
			return ""
		}
		if len(base) > 0 {
			return strings.ToUpper(base[:1]) + base[1:]
		}
	}
	return ""
}

func trimCompany(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, `"'.,;:`)
	if len(s) > 80 {
		s = s[:80]
	}
	return s
}

func trimRole(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, `"'.,;:`)
	if len(s) > 120 {
		s = s[:120]
	}
	return s
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

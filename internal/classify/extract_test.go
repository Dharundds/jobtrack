package classify

import (
	"strings"
	"testing"
)

func TestExtractCompany_greenhouseSharedSender(t *testing.T) {
	from := "no-reply@us.greenhouse-mail.io"
	rubrik := ExtractCompany(from, "Thank you for applying to Rubrik", "Thanks for applying to Rubrik")
	if rubrik != "Rubrik" {
		t.Fatalf("Rubrik company = %q", rubrik)
	}
	prodigal := ExtractCompany(from, "Thank you for applying to Prodigal", "Thanks for applying to Prodigal")
	if prodigal != "Prodigal" {
		t.Fatalf("Prodigal company = %q", prodigal)
	}
	if SimilarCompanyRole(rubrik, "", prodigal, "") {
		t.Fatal("Rubrik and Prodigal should not match")
	}
}

func TestSimilarCompanyRole_ATSBrand(t *testing.T) {
	if SimilarCompanyRole("Greenhouse", "", "Greenhouse", "") {
		t.Fatal("should not merge on ATS brand alone")
	}
}

func TestExtractCompany_rejectsGarbage(t *testing.T) {
	body := "If your application isn't selected we will notify you. With the next steps in our process."
	c := ExtractCompany("noreply@workday.com", "Application update", body)
	if c == "the next steps" || strings.Contains(c, "isn't selected") {
		t.Fatalf("unexpected company = %q", c)
	}
}

func TestExtractCompany_workdayDisplayNames(t *testing.T) {
	tests := []struct {
		from, want string
	}{
		{"Motorola Solutions - Workday <noreply@myworkday.com>", "Motorola Solutions"},
		{"MyWorkday@Proofpoint <notifications@proofpoint.com>", "Proofpoint"},
		{"FOX Workday Notification <noreply@myworkday.com>", "FOX"},
	}
	for _, tt := range tests {
		c := ExtractCompany(tt.from, "Candidate Task", "Please sign in to Workday.")
		if c != tt.want {
			t.Fatalf("from %q: company = %q, want %q", tt.from, c, tt.want)
		}
	}
}

func TestExtractCompany_abnormalAck(t *testing.T) {
	subject := "Thank you for applying to Abnormal"
	body := "Thanks for applying to Abnormal AI. We've received your application."
	c := ExtractCompany("jobs@abnormal.ai", subject, body)
	if c != "Abnormal" && c != "Abnormal AI" {
		t.Fatalf("company = %q", c)
	}
}

func TestExtractCompany_thales(t *testing.T) {
	from := "Thales Group <recruiting@jobalerts.thalesgroup.com>"
	subject := "Thales - Invites you to participate in our screening process"
	body := "Thanks for your interest in joining us as a Golang Developer. In order to get to know you better"
	c := ExtractCompany(from, subject, body)
	if c != "Thales" && c != "Thales Group" {
		t.Fatalf("company = %q", c)
	}
	r := ExtractRole(subject, body)
	if r != "Golang Developer" {
		t.Fatalf("role = %q", r)
	}
}

func TestExtractCompany_thalesSubmitted(t *testing.T) {
	from := "Thales Group <recruiting@jobalerts.thalesgroup.com>"
	subject := "Your Screening for Golang Developer has been successfully submitted"
	c := ExtractCompany(from, subject, "")
	if c != "Thales Group" && c != "Thales" {
		t.Fatalf("company = %q", c)
	}
	r := ExtractRole(subject, "")
	if r != "Golang Developer" {
		t.Fatalf("role = %q", r)
	}
}

func TestRuleClassifier_thalesScreeningInvite(t *testing.T) {
	rules := testRules()
	rules = append(rules, Rule{
		Field: "subject", Pattern: `(?i)screening process|invites you to participate`, EventType: EventInterviewInvite, Priority: 100, Enabled: true,
	})
	c := NewRuleClassifier(rules)
	in := EmailInput{
		From:    "Thales Group <recruiting@jobalerts.thalesgroup.com>",
		Subject: "Thales - Invites you to participate in our screening process",
		Body:    "joining us as a Golang Developer",
	}
	res := RefineResult(in, c.Classify(in))
	if res.EventType != EventInterviewInvite {
		t.Fatalf("event = %q", res.EventType)
	}
}

func TestExtractRole_rejectsThe(t *testing.T) {
	body := "Thank you for your interest for the role of Software Engineer."
	r := ExtractRole("", body)
	if r == "the" {
		t.Fatalf("role = %q", r)
	}
}

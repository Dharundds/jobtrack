package classify

import (
	"context"
	"testing"
)

func TestBeBettaApplicationReceived(t *testing.T) {
	subject := "Application for Software Engineer II (SDE 2) - Backend received, Thank you!"
	body := `Thank you for applying to BeBetta for the role of Software Engineer II (SDE 2) - Backend through our career portal.
This email is to confirm that we have received your application.`
	from := `BeBetta Hiring team <no-reply@kekamail.com>`
	in := EmailInput{From: from, Subject: subject, Body: body, Snippet: body}

	c := NewRuleClassifier(testRules())
	res := RefineResult(in, c.Classify(context.Background(), in))
	if res.EventType != EventApplied && res.EventType != EventAcknowledged {
		t.Fatalf("EventType = %q, want applied or acknowledged", res.EventType)
	}
	if !ShouldLinkToApplication(res) {
		t.Fatal("should link BeBetta confirmation")
	}
	if !hasApplicationSubmittedSignal(from, subject, body, body) {
		t.Fatal("expected submission signal")
	}
	role := ExtractRole(subject, body)
	if role == "" {
		t.Fatal("expected role extracted")
	}
}

func TestBeBettaApplicationReceived_snippetOnlyBody(t *testing.T) {
	subject := "Application for Software Engineer II (SDE 2) - Backend received, Thank you!"
	snippet := "Thank you for applying to BeBetta for the role of Software Engineer II (SDE 2) - Backend"
	in := EmailInput{
		From: `BeBetta Hiring team <no-reply@kekamail.com>`,
		Subject: subject,
		Body:    "@media only screen { .mj-column { width: 100% } }",
		Snippet: snippet,
	}
	c := NewRuleClassifier(testRules())
	res := RefineResult(in, c.Classify(context.Background(), in))
	if res.EventType == EventOther {
		t.Fatalf("EventType = %q with snippet fallback", res.EventType)
	}
	if !ShouldLinkToApplication(res) {
		t.Fatal("should link with submission signal fallback")
	}
}
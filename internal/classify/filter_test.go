package classify

import "testing"

func TestRefineResult_rejectionNoise(t *testing.T) {
	c := NewRuleClassifier(testRules())

	news := EmailInput{
		From:    "deals@store.com",
		Subject: "Unfortunately, this offer expired",
		Body:    "Shop our sale today.",
	}
	res := RefineResult(news, c.Classify(news))
	if res.EventType != EventOther {
		t.Fatalf("EventType = %q, want other", res.EventType)
	}
	if ShouldLinkToApplication(res) {
		t.Fatal("marketing mail should not link")
	}
}

func TestRefineResult_workdayRejection(t *testing.T) {
	c := NewRuleClassifier(testRules())
	in := EmailInput{
		From:    "workday@bigcorp.com",
		Subject: "Update on your application",
		Body:    "Unfortunately we will not be moving forward.",
	}
	res := RefineResult(in, c.Classify(in))
	if res.EventType != EventRejected {
		t.Fatalf("EventType = %q, want rejected", res.EventType)
	}
}

func TestRefineResult_linkedinConnect(t *testing.T) {
	in := EmailInput{
		From:    "LinkedIn <invitations-noreply@linkedin.com>",
		Subject: "Alex wants to connect",
		Body:    "Accept Alex's invitation to connect.",
	}
	res := RefineResult(in, Result{EventType: EventOther, Company: "LinkedIn"})
	if res.EventType != EventOther {
		t.Fatalf("EventType = %q, want other", res.EventType)
	}
	if ShouldLinkToApplication(res) {
		t.Fatal("connect invite should not link")
	}
}

func TestRefineResult_linkedinJobRejection(t *testing.T) {
	c := NewRuleClassifier(testRules())
	in := EmailInput{
		From:    "jobs-noreply@linkedin.com",
		Subject: "Unfortunately, we decided not to move forward",
		Body:    "Thank you for your interest.",
	}
	res := RefineResult(in, c.Classify(in))
	if res.EventType != EventRejected {
		t.Fatalf("EventType = %q, want rejected", res.EventType)
	}
}

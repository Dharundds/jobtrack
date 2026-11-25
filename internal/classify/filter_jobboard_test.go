package classify

import (
	"context"
	"strings"
	"testing"

	"github.com/dharundds/jobtrack/internal/store"
)

func TestRefineResult_instahyreRecommendation(t *testing.T) {
	body := strings.TrimSpace(`Hey Dharun, Wissen Technology has just posted a Full Stack Developer (Python + React) job.
You look like a good fit for the position based on your Instahyre profile. Apply here`)
	in := EmailInput{
		From:    `Wissen Technology Careers <noreply@instahyre.com>`,
		Subject: "Job | Full Stack Developer (Python + React) at Wissen Technology (Bangalore)",
		Body:    body,
		Snippet: body,
	}
	res := Result{EventType: EventRecruiterOutreach, Company: "Wissen Technology", Role: "Full Stack Developer"}
	res = RefineResult(in, res)
	if res.EventType != EventOther {
		t.Fatalf("EventType = %q, want other", res.EventType)
	}
	if ShouldLinkToApplication(res) {
		t.Fatal("instahyre alert should not link")
	}
	if ShouldRemoveListingAfterLLM(true, "job alert", map[int64]Result{1: res}, []store.MessageRow{{
		ID: 1, FromAddr: in.From, Subject: in.Subject, BodyExcerpt: body, Snippet: body,
	}}) {
		// isJob true from bad LLM still removed via allAlerts
	} else if !ShouldRemoveListingAfterLLM(false, "instahyre recommendation", map[int64]Result{1: res}, []store.MessageRow{{
		ID: 1, FromAddr: in.From, Subject: in.Subject, BodyExcerpt: body, Snippet: body,
	}}) {
		t.Fatal("should remove instahyre-only listing")
	}
}

func TestHasApplicationSubmittedSignal_workdayAck(t *testing.T) {
	in := EmailInput{
		From:    "noreply@greenhouse.io",
		Subject: "Thank you for applying to Acme",
		Body:    "We received your application.",
	}
	c := NewRuleClassifier(testRules())
	res := RefineResult(in, c.Classify(context.Background(), in))
	if res.EventType == EventOther {
		t.Fatalf("expected real application mail, got other (%+v)", res)
	}
}
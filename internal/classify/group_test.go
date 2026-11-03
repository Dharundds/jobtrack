package classify

import "testing"

func TestPickApplication_byThread(t *testing.T) {
	byThread := map[string]int64{"thread-1": 42}
	id, ok := PickApplication("thread-1", "Acme", "Engineer", byThread, nil)
	if !ok || id != 42 {
		t.Fatalf("got id=%d ok=%v", id, ok)
	}
}

func TestPickApplication_byCompanyRole(t *testing.T) {
	apps := []ApplicationCandidate{
		{ID: 7, Company: "Stripe", Role: "Backend Engineer"},
	}
	id, ok := PickApplication("thread-new", "stripe", "backend engineer", map[string]int64{}, apps)
	if !ok || id != 7 {
		t.Fatalf("got id=%d ok=%v", id, ok)
	}
}

func TestSimilarCompanyRole(t *testing.T) {
	if !SimilarCompanyRole("Google", "Backend Engineer", "google", "backend engineer") {
		t.Fatal("expected similar")
	}
	if SimilarCompanyRole("Google", "PM", "Amazon", "PM") {
		t.Fatal("expected different companies")
	}
}

func TestDeriveActivity_awaitingReply(t *testing.T) {
	in := mustTime("2026-01-10")
	out := mustTime("2026-01-08")
	msgs := []MessageActivity{
		{Direction: "in", ReceivedAt: in},
		{Direction: "out", ReceivedAt: out},
	}
	act := DeriveActivity(msgs, nil)
	if !act.AwaitingMyReply {
		t.Fatal("expected awaiting reply")
	}
}

func TestExtractEventDetails_ics(t *testing.T) {
	body := "BEGIN:VCALENDAR\nDTSTART:20261115T140000Z\nEND:VCALENDAR"
	d := ExtractEventDetails(body)
	if d.InterviewAt == nil {
		t.Fatal("expected interview_at")
	}
}

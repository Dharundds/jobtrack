package classify

import (
	"testing"
	"time"
)

func testRules() []Rule {
	return []Rule{
		{ID: 1, Field: "subject", Pattern: `(?i)thank you for applying|application received`, EventType: EventApplied, Priority: 100, Enabled: true},
		{ID: 2, Field: "body", Pattern: `(?i)we have received your application|thanks for your interest`, EventType: EventAcknowledged, Priority: 90, Enabled: true},
		{ID: 3, Field: "subject", Pattern: `(?i)unfortunately|not moving forward|regret to inform`, EventType: EventRejected, Priority: 100, Enabled: true},
		{ID: 10, Field: "body", Pattern: `(?i)unfortunately|regret to inform`, EventType: EventRejected, Priority: 95, Enabled: true},
		{ID: 4, Field: "subject", Pattern: `(?i)interview|phone screen|screening process|participate in (?:our |the )?screening|invites you to participate`, EventType: EventInterviewInvite, Priority: 100, Enabled: true},
		{ID: 5, Field: "body", Pattern: `(?i)zoom\.us/j/|meet\.google\.com/|teams\.microsoft\.com/`, EventType: EventInterviewInvite, Priority: 80, Enabled: true},
		{ID: 6, Field: "subject", Pattern: `(?i)assessment|hackerrank|take-home`, EventType: EventAssessment, Priority: 100, Enabled: true},
		{ID: 11, Field: "body", Pattern: `(?i)take-home assessment`, EventType: EventAssessment, Priority: 90, Enabled: true},
		{ID: 7, Field: "subject", Pattern: `(?i)offer letter|pleased to offer`, EventType: EventOffer, Priority: 100, Enabled: true},
		{ID: 8, Field: "body", Pattern: `(?i)exciting opportunity|are you open to`, EventType: EventRecruiterOutreach, Priority: 70, Enabled: true},
		{ID: 9, Field: "from", Pattern: `(?i)greenhouse-mail|lever\.co`, EventType: EventApplied, Priority: 50, Enabled: true},
	}
}

func TestRuleClassifier_samples(t *testing.T) {
	c := NewRuleClassifier(testRules())

	tests := []struct {
		name  string
		in    EmailInput
		want  string
	}{
		{
			name: "greenhouse applied",
			in: EmailInput{
				From:    "Acme Hiring <noreply@greenhouse-mail.io>",
				Subject: "Thank you for applying to Software Engineer at Acme",
				Body:    "We received your application.",
			},
			want: EventApplied,
		},
		{
			name: "lever acknowledgment",
			in: EmailInput{
				From:    "jobs@mail.lever.co",
				Subject: "Confirm your application - Backend Engineer",
				Body:    "Thanks for your interest in Contoso.",
			},
			want: EventAcknowledged,
		},
		{
			name: "workday rejection",
			in: EmailInput{
				From:    "workday@bigcorp.com",
				Subject: "Update on your application",
				Body:    "Unfortunately we will not be moving forward.",
			},
			want: EventRejected,
		},
		{
			name: "linkedin rejection subject",
			in: EmailInput{
				From:    "jobs-noreply@linkedin.com",
				Subject: "Unfortunately, we decided not to move forward",
				Body:    "Thank you for your interest.",
			},
			want: EventRejected,
		},
		{
			name: "naukri applied",
			in: EmailInput{
				From:    "alert@naukri.com",
				Subject: "Application received for Java Developer",
				Body:    "Your application has been sent to the recruiter.",
			},
			want: EventApplied,
		},
		{
			name: "interview invite subject",
			in: EmailInput{
				From:    "recruiter@startup.io",
				Subject: "Phone screen for Product Manager role",
				Body:    "Please share your availability.",
			},
			want: EventInterviewInvite,
		},
		{
			name: "thales screening subject",
			in: EmailInput{
				From:    "Thales Careers <noreply@thalesgroup.com>",
				Subject: "Thales - Invites you to participate in our screening process",
				Body:    "Please complete the next step in your application.",
			},
			want: EventInterviewInvite,
		},
		{
			name: "zoom link body",
			in: EmailInput{
				From:    "calendar@company.com",
				Subject: "Updated invitation",
				Body:    "Join Zoom Meeting https://zoom.us/j/1234567890",
			},
			want: EventInterviewInvite,
		},
		{
			name: "google meet",
			in: EmailInput{
				Subject: "Interview schedule",
				Body:    "Meet link: https://meet.google.com/abc-defg-hij",
			},
			want: EventInterviewInvite,
		},
		{
			name: "assessment hackerrank",
			in: EmailInput{
				Subject: "Complete your HackerRank assessment",
				Body:    "Please complete the online assessment within 3 days.",
			},
			want: EventAssessment,
		},
		{
			name: "offer letter",
			in: EmailInput{
				Subject: "Offer letter - Senior Engineer",
				Body:    "We are pleased to offer you employment at ExampleCo.",
			},
			want: EventOffer,
		},
		{
			name: "recruiter outreach",
			in: EmailInput{
				From:    "Jane Recruiter <jane@talent.com>",
				Subject: "Opportunity at FinTech Co",
				Body:    "Are you open to exploring a new role?",
			},
			want: EventRecruiterOutreach,
		},
		{
			name: "plain recruiter mail",
			in: EmailInput{
				From:    "hiring@agency.com",
				Subject: "Full stack role",
				Body:    "I'd like to reach out about an exciting opportunity.",
			},
			want: EventRecruiterOutreach,
		},
		{
			name: "ashby interview",
			in: EmailInput{
				From:    "notifications@ashbyhq.com",
				Subject: "Schedule your interview with DataCo",
				Body:    "Click here to schedule a video call.",
			},
			want: EventInterviewInvite,
		},
		{
			name: "smartrecruiters ack body",
			in: EmailInput{
				Body: "We have received your application for UX Designer.",
			},
			want: EventAcknowledged,
		},
		{
			name: "indeed applied subject",
			in: EmailInput{
				From:    "noreply@indeed.com",
				Subject: "Application received for Warehouse Associate",
				Body:    "Good luck!",
			},
			want: EventApplied,
		},
		{
			name: "icims rejection",
			in: EmailInput{
				From:    "donotreply@icims.com",
				Subject: "Regret to inform",
				Body:    "After careful consideration we regret to inform you.",
			},
			want: EventRejected,
		},
		{
			name: "teams meeting",
			in: EmailInput{
				Body: "Join: https://teams.microsoft.com/l/meetup-join/19%3ameeting",
			},
			want: EventInterviewInvite,
		},
		{
			name: "take home assessment",
			in: EmailInput{
				Subject: "Take-home exercise",
				Body:    "Please complete this take-home assessment.",
			},
			want: EventAssessment,
		},
		{
			name: "wellfound outreach",
			in: EmailInput{
				From:    "hello@wellfound.com",
				Subject: "Startup hiring for Founding Engineer",
				Body:    "Exciting opportunity at a seed stage startup.",
			},
			want: EventRecruiterOutreach,
		},
		{
			name: "generic newsletter other",
			in: EmailInput{
				From:    "news@random.com",
				Subject: "Weekly digest",
				Body:    "Here are top stories this week.",
			},
			want: EventOther,
		},
		{
			name: "outbound user reply other",
			in: EmailInput{
				Direction: "out",
				Subject:   "Re: Tuesday availability",
				Body:      "Thanks, Tuesday works for me.",
			},
			want: EventOther,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := c.Classify(tt.in)
			if got.EventType != tt.want {
				t.Fatalf("EventType = %q, want %q", got.EventType, tt.want)
			}
		})
	}
}

func TestExtractCompanyRole(t *testing.T) {
	c := ExtractCompany("Stripe Careers <jobs@stripe.com>", "Interview at Stripe", "")
	if c != "Stripe Careers" && c != "Stripe" {
		t.Fatalf("company = %q", c)
	}
	r := ExtractRole("Application: Staff Engineer", "")
	if r != "Staff Engineer" {
		t.Fatalf("role = %q", r)
	}
}

func TestLatestStatusFromEvents(t *testing.T) {
	events := []EventRecord{
		{Type: EventApplied, OccurredAt: mustTime("2026-01-01")},
		{Type: EventInterviewInvite, OccurredAt: mustTime("2026-01-05")},
		{Type: EventRejected, OccurredAt: mustTime("2026-01-10")},
	}
	if got := LatestStatusFromEvents(events); got != "rejected" {
		t.Fatalf("status = %q", got)
	}
}

func mustTime(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}

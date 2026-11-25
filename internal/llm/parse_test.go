package llm

import "testing"

func TestExtractJSONPayload_thinkingThenJSON(t *testing.T) {
	thinkOpen := "<" + "think" + ">"
	thinkClose := "<" + "/" + "think" + ">"
	raw := thinkOpen + "reasoning here" + thinkClose + "\n{\"event_type\":\"applied\",\"company\":\"Acme\",\"role\":\"Engineer\",\"meeting_url\":\"\",\"interview_at\":\"\"}"
	out, err := parseClassifyJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	if out.EventType != "applied" || out.Company != "Acme" {
		t.Fatalf("parsed %+v", out)
	}
}

func TestExtractJSONPayload_codeFence(t *testing.T) {
	raw := "```json\n{\"is_job_application\":true,\"reason\":\"ok\",\"company\":\"X\",\"role\":\"Y\",\"pipeline_status\":\"applied\"}\n```"
	out, err := parseReviewJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	if !out.IsJobApplication || out.Company != "X" {
		t.Fatalf("parsed %+v", out)
	}
}

func TestExtractJSONObject_trailingText(t *testing.T) {
	s := "here is json {\"event_type\":\"other\",\"company\":\"\",\"role\":\"\",\"meeting_url\":\"\",\"interview_at\":\"\"} thanks"
	obj := extractJSONObject(extractJSONPayload(s))
	if obj == "" {
		t.Fatal("expected object")
	}
}
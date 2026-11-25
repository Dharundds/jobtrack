package llm

import "testing"

func TestParseReviewJSON_isJobApplicationVariants(t *testing.T) {
	raw := `{"isJobApplication":false,"reason":"newsletter","company":"","role":"","pipeline_status":"applied"}`
	out, err := parseReviewJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	if out.IsJobApplication {
		t.Fatal("expected false")
	}
	if out.Reason != "newsletter" {
		t.Fatalf("reason %q", out.Reason)
	}
}

func TestParseReviewJSON_roleAliases(t *testing.T) {
	raw := `{"is_job_application":true,"job_title":"Software Engineer","company_name":"Acme"}`
	out, err := parseReviewJSON(raw)
	if err != nil {
		t.Fatal(err)
	}
	if out.Role != "Software Engineer" || out.Company != "Acme" {
		t.Fatalf("%+v", out)
	}
}
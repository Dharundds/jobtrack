package llm

import "testing"

func TestNormalizeBaseURL(t *testing.T) {
	if got := NormalizeBaseURL("http://localhost:1234"); got != "http://localhost:1234/v1" {
		t.Fatalf("got %q", got)
	}
	if got := NormalizeBaseURL("http://localhost:1234/v1/"); got != "http://localhost:1234/v1" {
		t.Fatalf("got %q", got)
	}
}
package classify

import (
	"regexp"
	"sort"
	"strings"
)

// Rule is a single classification rule (from DB or tests).
type Rule struct {
	ID        int64
	Field     string
	Pattern   string
	EventType string
	Priority  int
	Enabled   bool
	compiled  *regexp.Regexp
}

func compileRule(r Rule) Rule {
	if r.compiled == nil && r.Pattern != "" {
		re, err := regexp.Compile(r.Pattern)
		if err == nil {
			r.compiled = re
		}
	}
	return r
}

// RuleClassifier applies regex rules by priority (highest first).
type RuleClassifier struct {
	rules []Rule
}

func NewRuleClassifier(rules []Rule) *RuleClassifier {
	sorted := make([]Rule, 0, len(rules))
	for _, r := range rules {
		if !r.Enabled {
			continue
		}
		sorted = append(sorted, compileRule(r))
	}
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Priority == sorted[j].Priority {
			return sorted[i].ID < sorted[j].ID
		}
		return sorted[i].Priority > sorted[j].Priority
	})
	return &RuleClassifier{rules: sorted}
}

func (c *RuleClassifier) Classify(in EmailInput) Result {
	texts := map[string]string{
		"subject": in.Subject,
		"body":    strings.Join([]string{in.Body, in.Snippet}, "\n"),
		"from":    in.From,
	}
	for _, r := range c.rules {
		text, ok := texts[r.Field]
		if !ok || text == "" || r.compiled == nil {
			continue
		}
		if r.compiled.MatchString(text) {
			company := ExtractCompany(in.From, in.Subject, in.Body)
			role := ExtractRole(in.Subject, in.Body)
			details := ExtractEventDetails(in.Body + "\n" + in.Snippet)
			return Result{
				EventType:    r.EventType,
				Company:      company,
				Role:         role,
				Details:      details,
				ClassifiedBy: "rule",
			}
		}
	}
	return Result{
		EventType:    EventOther,
		Company:      ExtractCompany(in.From, in.Subject, in.Body),
		Role:         ExtractRole(in.Subject, in.Body),
		Details:      ExtractEventDetails(in.Body + "\n" + in.Snippet),
		ClassifiedBy: "rule",
	}
}

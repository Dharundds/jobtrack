package classify

import (
	"regexp"
	"strings"
	"time"
)

var (
	reICSDate     = regexp.MustCompile(`(?i)DTSTART(?:;[^:\r\n]*)?:([0-9]{8}T[0-9]{6}Z?)`)
	reZoom        = regexp.MustCompile(`https?://(?:[\w-]+\.)?zoom\.us/j/[A-Za-z0-9?=&_-]+`)
	reMeet        = regexp.MustCompile(`https?://meet\.google\.com/[a-z-]+`)
	reTeams       = regexp.MustCompile(`https?://teams\.microsoft\.com/l/meetup-join/[^\s>"']+`)
)

// ExtractEventDetails parses interview time and meeting links from plain text / ICS snippets.
func ExtractEventDetails(text string) EventDetails {
	var d EventDetails
	if m := reZoom.FindString(text); m != "" {
		d.MeetingURL = m
	} else if m := reMeet.FindString(text); m != "" {
		d.MeetingURL = m
	} else if m := reTeams.FindString(text); m != "" {
		d.MeetingURL = m
	}
	if m := reICSDate.FindStringSubmatch(text); len(m) > 1 {
		if t, err := parseICSDate(m[1]); err == nil {
			d.InterviewAt = &t
		}
	}
	return d
}

func parseICSDate(raw string) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	if strings.HasSuffix(raw, "Z") {
		return time.Parse("20060102T150405Z", raw)
	}
	return time.Parse("20060102T150405", raw)
}

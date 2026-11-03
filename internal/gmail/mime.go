package gmail

import (
	"encoding/base64"
	"io"
	"mime"
	"strings"
	"unicode/utf8"

	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
	"google.golang.org/api/gmail/v1"
)

func extractBody(part *gmail.MessagePart) (plain string, hasCalendar bool) {
	if part == nil {
		return "", false
	}
	if part.MimeType == "text/calendar" {
		hasCalendar = true
	}
	if part.Body != nil && part.Body.Data != "" {
		if part.MimeType == "text/plain" || part.MimeType == "text/html" {
			decoded, err := decodePart(part)
			if err == nil && decoded != "" {
				if part.MimeType == "text/html" {
					decoded = htmlToText(decoded)
				}
				return decoded, hasCalendar
			}
		}
	}
	var plainParts, htmlParts []string
	for _, child := range part.Parts {
		p, cal := extractBody(child)
		if cal {
			hasCalendar = true
		}
		if p == "" {
			continue
		}
		if strings.HasPrefix(child.MimeType, "text/plain") {
			plainParts = append(plainParts, p)
		} else if strings.HasPrefix(child.MimeType, "text/html") {
			htmlParts = append(htmlParts, p)
		}
	}
	if len(plainParts) > 0 {
		return strings.Join(plainParts, "\n\n"), hasCalendar
	}
	if len(htmlParts) > 0 {
		return htmlToText(strings.Join(htmlParts, "\n\n")), hasCalendar
	}
	return "", hasCalendar
}

func decodePart(part *gmail.MessagePart) (string, error) {
	raw, err := base64.RawURLEncoding.DecodeString(part.Body.Data)
	if err != nil {
		raw, err = base64.StdEncoding.DecodeString(part.Body.Data)
		if err != nil {
			return "", err
		}
	}
	cs := "utf-8"
	if part.Body != nil && part.Headers != nil {
		for _, h := range part.Headers {
			if strings.EqualFold(h.Name, "Content-Type") {
				_, params, _ := mime.ParseMediaType(h.Value)
				if enc, ok := params["charset"]; ok {
					cs = enc
				}
			}
		}
	}
	r, err := charset.NewReaderLabel(cs, strings.NewReader(string(raw)))
	if err != nil {
		return string(raw), nil
	}
	b, err := io.ReadAll(r)
	if err != nil {
		return string(raw), nil
	}
	if !utf8.Valid(b) {
		return string(raw), nil
	}
	return string(b), nil
}

func htmlToText(src string) string {
	doc, err := html.Parse(strings.NewReader(src))
	if err != nil {
		return stripTagsSimple(src)
	}
	var b strings.Builder
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.TextNode {
			b.WriteString(n.Data)
			b.WriteString(" ")
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
		if n.Type == html.ElementNode && (n.Data == "p" || n.Data == "br" || n.Data == "div") {
			b.WriteString("\n")
		}
	}
	walk(doc)
	return strings.TrimSpace(compressSpace(b.String()))
}

func stripTagsSimple(s string) string {
	var b strings.Builder
	inTag := false
	for _, r := range s {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case !inTag:
			b.WriteRune(r)
		}
	}
	return compressSpace(b.String())
}

func compressSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

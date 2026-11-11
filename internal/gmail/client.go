package gmail

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/dharundds/jobtrack/internal/auth"
	"github.com/dharundds/jobtrack/internal/crypto"
	"github.com/dharundds/jobtrack/internal/store"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

const maxBodyExcerpt = 2000

// Message is normalized metadata from Gmail (no attachment payloads).
type Message struct {
	ID          string
	ThreadID    string
	LabelIDs    []string
	From        string
	To          string
	Subject     string
	ReceivedAt  time.Time
	Snippet     string
	BodyExcerpt string
}

type Client struct {
	svc   *gmail.Service
	log   zerolog.Logger
	quota *QuotaLimiter
}

// NewClientFromAccount decrypts the refresh token and builds a Gmail API client.
func NewClientFromAccount(
	ctx context.Context,
	oauth *auth.Service,
	tokenKey []byte,
	acc *store.Account,
	unitsPerMinute int,
	log zerolog.Logger,
) (*Client, error) {
	refresh, err := crypto.DecryptToken(tokenKey, acc.EncRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("decrypt refresh token: %w", err)
	}
	tok := &oauth2.Token{RefreshToken: refresh}
	// Token refresh must outlive per-sync cancellation (oauth2: client invalid after ctx ends).
	client := oauth.OAuthConfig().Client(context.Background(), tok)
	svc, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("gmail service: %w", err)
	}
	return &Client{
		svc:   svc,
		log:   log,
		quota: NewQuotaLimiter(unitsPerMinute),
	}, nil
}

func (c *Client) do(ctx context.Context, cost int, fn func() error) error {
	return withRetry(ctx, c.log, func() error {
		if err := c.quota.Wait(ctx, cost); err != nil {
			return err
		}
		return fn()
	})
}

// CurrentHistoryID returns the mailbox historyId from the user profile.
func (c *Client) CurrentHistoryID(ctx context.Context) (uint64, error) {
	var out uint64
	err := c.do(ctx, CostUsersGetProfile, func() error {
		profile, err := c.svc.Users.GetProfile("me").Context(ctx).Do()
		if err != nil {
			return err
		}
		out = profile.HistoryId
		return nil
	})
	return out, err
}

// ListMessageIDs lists message IDs matching query with optional since filter (full sync).
func (c *Client) ListMessageIDs(ctx context.Context, query string, since time.Time) ([]string, error) {
	q := query
	if !since.IsZero() {
		q = fmt.Sprintf("%s after:%s", query, since.Format("2006/01/02"))
	}

	var ids []string
	var pageToken string
	for {
		var batch []string
		var nextPage string
		err := c.do(ctx, CostMessagesList, func() error {
			call := c.svc.Users.Messages.List("me").Q(q).MaxResults(500)
			if pageToken != "" {
				call = call.PageToken(pageToken)
			}
			res, err := call.Context(ctx).Do()
			if err != nil {
				return err
			}
			for _, m := range res.Messages {
				if m.Id != "" {
					batch = append(batch, m.Id)
				}
			}
			nextPage = res.NextPageToken
			return nil
		})
		if err != nil {
			return nil, err
		}
		ids = append(ids, batch...)
		pageToken = nextPage
		if pageToken == "" {
			break
		}
	}
	return ids, nil
}

// ListHistory returns message IDs added since startHistoryID. historyExpired is true on 404.
func (c *Client) ListHistory(ctx context.Context, startHistoryID uint64) (ids []string, historyExpired bool, err error) {
	var pageToken string
	for {
		var res *gmail.ListHistoryResponse
		expired := false
		retryErr := c.do(ctx, CostHistoryList, func() error {
			call := c.svc.Users.History.List("me").
				StartHistoryId(startHistoryID).
				HistoryTypes("messageAdded")
			if pageToken != "" {
				call = call.PageToken(pageToken)
			}
			var e error
			res, e = call.Context(ctx).Do()
			if e != nil {
				if isHistoryExpired(e) {
					expired = true
					return nil
				}
				return e
			}
			return nil
		})
		if expired {
			return nil, true, nil
		}
		if retryErr != nil {
			return nil, false, retryErr
		}
		for _, h := range res.History {
			for _, added := range h.MessagesAdded {
				if added.Message != nil && added.Message.Id != "" {
					ids = append(ids, added.Message.Id)
				}
			}
		}
		if res.NextPageToken == "" {
			break
		}
		pageToken = res.NextPageToken
	}
	return ids, false, nil
}

// GetMessage fetches and parses a single message.
func (c *Client) GetMessage(ctx context.Context, id string) (*Message, error) {
	var gm *gmail.Message
	err := c.do(ctx, CostMessagesGet, func() error {
		var e error
		gm, e = c.svc.Users.Messages.Get("me", id).Format("full").Context(ctx).Do()
		return e
	})
	if err != nil {
		return nil, err
	}

	from := headerValue(gm.Payload.Headers, "From")
	to := headerValue(gm.Payload.Headers, "To")
	subject := headerValue(gm.Payload.Headers, "Subject")
	dateStr := headerValue(gm.Payload.Headers, "Date")
	received := parseEmailDate(dateStr)

	body, _ := extractBody(gm.Payload)
	if len(body) > maxBodyExcerpt {
		body = body[:maxBodyExcerpt]
	}

	return &Message{
		ID:          gm.Id,
		ThreadID:    gm.ThreadId,
		LabelIDs:    gm.LabelIds,
		From:        from,
		To:          to,
		Subject:     subject,
		ReceivedAt:  received,
		Snippet:     gm.Snippet,
		BodyExcerpt: body,
	}, nil
}

func headerValue(headers []*gmail.MessagePartHeader, name string) string {
	for _, h := range headers {
		if strings.EqualFold(h.Name, name) {
			return h.Value
		}
	}
	return ""
}

func parseEmailDate(raw string) time.Time {
	if raw == "" {
		return time.Now().UTC()
	}
	if t, err := mail.ParseDate(raw); err == nil {
		return t.UTC()
	}
	return time.Now().UTC()
}

// IsNotFound reports whether err is a Gmail API 404 (e.g. message deleted).
func IsNotFound(err error) bool {
	var gerr *googleapi.Error
	if errors.As(err, &gerr) && gerr.Code == 404 {
		return true
	}
	return false
}

// IsOAuthTokenError reports refresh/revocation failures that require re-connecting Gmail.
func IsOAuthTokenError(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "invalid_grant") ||
		strings.Contains(s, "cannot fetch token") ||
		strings.Contains(s, "token expired and refresh token is not set")
}

func isHistoryExpired(err error) bool {
	var gerr *googleapi.Error
	if errors.As(err, &gerr) && gerr.Code == 404 {
		return true
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "404") && strings.Contains(s, "history")
}

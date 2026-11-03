package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	sessionKeyOAuthState = "oauth_state"
	SessionKeyAccountID  = "account_id"
)

type Service struct {
	oauth *oauth2.Config
}

func NewService(clientID, clientSecret, redirectURL string) *Service {
	return &Service{
		oauth: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"https://www.googleapis.com/auth/gmail.readonly"},
			Endpoint:     google.Endpoint,
		},
	}
}

func (s *Service) OAuthConfig() *oauth2.Config {
	return s.oauth
}

func (s *Service) LoginURL(state string) string {
	return s.oauth.AuthCodeURL(state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "consent"),
	)
}

func (s *Service) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.oauth.Exchange(ctx, code)
}

func NewOAuthState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func SetOAuthState(sess sessions.Session, state string) error {
	sess.Set(sessionKeyOAuthState, state)
	return sess.Save()
}

func ValidateOAuthState(sess sessions.Session, received string) bool {
	v, ok := sess.Get(sessionKeyOAuthState).(string)
	if !ok || v == "" || received == "" {
		return false
	}
	valid := v == received
	sess.Delete(sessionKeyOAuthState)
	_ = sess.Save()
	return valid
}

func SetLoggedIn(sess sessions.Session, accountID int64) error {
	sess.Set(SessionKeyAccountID, accountID)
	return sess.Save()
}

func ClearSession(sess sessions.Session) error {
	sess.Clear()
	return sess.Save()
}

func AccountIDFromSession(sess sessions.Session) (int64, bool) {
	v := sess.Get(SessionKeyAccountID)
	switch id := v.(type) {
	case int64:
		return id, id > 0
	case int:
		return int64(id), id > 0
	case float64:
		return int64(id), id > 0
	default:
		return 0, false
	}
}

func SessionFromContext(c *gin.Context) sessions.Session {
	return sessions.Default(c)
}

func SecureCookie(r *http.Request) bool {
	host := r.Host
	if idx := strings.Index(host, ":"); idx >= 0 {
		host = host[:idx]
	}
	host = strings.ToLower(host)
	return host != "localhost" && host != "127.0.0.1"
}

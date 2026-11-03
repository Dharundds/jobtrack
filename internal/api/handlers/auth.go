package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/dharundds/jobtrack/internal/api/apperrors"
	"github.com/dharundds/jobtrack/internal/api/middleware"
	"github.com/dharundds/jobtrack/internal/auth"
	"github.com/dharundds/jobtrack/internal/config"
	"github.com/dharundds/jobtrack/internal/crypto"
	"github.com/dharundds/jobtrack/internal/privacy"
	"github.com/dharundds/jobtrack/internal/store"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type AuthHandler struct {
	Cfg   *config.Config
	Auth  *auth.Service
	Store *store.Store
	Log   zerolog.Logger
}

type authStatusResponse struct {
	Connected bool   `json:"connected"`
	Email     string `json:"email,omitempty"`
}

func (h *AuthHandler) Status(c *gin.Context) {
	acc, err := h.Store.GetAccount(c.Request.Context())
	if err != nil {
		if err == store.ErrNoAccount {
			c.JSON(http.StatusOK, authStatusResponse{Connected: false})
			return
		}
		middleware.AbortWithError(c, apperrors.Internal("failed to read account"))
		return
	}

	c.JSON(http.StatusOK, authStatusResponse{
		Connected: acc.EncRefreshToken != "",
		Email:     acc.Email,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	state, err := auth.NewOAuthState()
	if err != nil {
		middleware.AbortWithError(c, apperrors.Internal("failed to create oauth state"))
		return
	}
	sess := auth.SessionFromContext(c)
	if err := auth.SetOAuthState(sess, state); err != nil {
		middleware.AbortWithError(c, apperrors.Internal("failed to save session"))
		return
	}
	c.Redirect(http.StatusFound, h.Auth.LoginURL(state))
}

func (h *AuthHandler) Callback(c *gin.Context) {
	log := middleware.LoggerFromContext(c)
	if errParam := c.Query("error"); errParam != "" {
		log.Warn().Str("oauth_error", errParam).Msg("google oauth denied")
		c.Redirect(http.StatusFound, h.Cfg.FrontendConnectURL("auth=denied"))
		return
	}

	state := c.Query("state")
	sess := auth.SessionFromContext(c)
	if !auth.ValidateOAuthState(sess, state) {
		middleware.AbortWithError(c, apperrors.BadRequest("invalid_state", "invalid or expired OAuth state"))
		return
	}

	code := c.Query("code")
	if code == "" {
		middleware.AbortWithError(c, apperrors.BadRequest("missing_code", "authorization code missing"))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	tok, err := h.Auth.Exchange(ctx, code)
	if err != nil {
		log.Warn().Err(err).Msg("oauth token exchange failed")
		middleware.AbortWithError(c, apperrors.BadRequest("oauth_exchange", "could not complete Google sign-in"))
		return
	}
	if tok.RefreshToken == "" {
		log.Warn().Msg("oauth response missing refresh token")
		middleware.AbortWithError(c, apperrors.BadRequest("no_refresh_token", "no refresh token returned; revoke app access in Google Account and try Connect again"))
		return
	}

	email, err := fetchGmailAddress(ctx, h.Auth.OAuthConfig(), tok)
	if err != nil {
		log.Warn().Err(err).Msg("failed to resolve gmail profile")
		middleware.AbortWithError(c, apperrors.Internal("could not read Gmail profile"))
		return
	}

	enc, err := crypto.EncryptToken(h.Cfg.TokenEncKey, tok.RefreshToken)
	if err != nil {
		log.Error().Err(err).Msg("encrypt refresh token")
		middleware.AbortWithError(c, apperrors.Internal("failed to store credentials"))
		return
	}

	acc, err := h.Store.UpsertAccount(ctx, email, enc)
	if err != nil {
		log.Error().Err(err).Str("email", privacy.MaskEmail(email)).Msg("save account")
		middleware.AbortWithError(c, apperrors.Internal("failed to save account"))
		return
	}

	if err := auth.SetLoggedIn(sess, acc.ID); err != nil {
		middleware.AbortWithError(c, apperrors.Internal("failed to save session"))
		return
	}

	log.Info().Str("email", privacy.MaskEmail(email)).Int64("account_id", acc.ID).Msg("gmail connected")
	c.Redirect(http.StatusFound, h.Cfg.PostLoginRedirectURL())
}

func (h *AuthHandler) Logout(c *gin.Context) {
	if err := h.Store.DeleteAllAccounts(c.Request.Context()); err != nil {
		middleware.AbortWithError(c, apperrors.Internal("failed to disconnect"))
		return
	}
	sess := auth.SessionFromContext(c)
	if err := auth.ClearSession(sess); err != nil {
		middleware.AbortWithError(c, apperrors.Internal("failed to clear session"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func fetchGmailAddress(ctx context.Context, oauthCfg *oauth2.Config, tok *oauth2.Token) (string, error) {
	client := oauthCfg.Client(ctx, tok)
	svc, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return "", err
	}
	profile, err := svc.Users.GetProfile("me").Context(ctx).Do()
	if err != nil {
		return "", err
	}
	return profile.EmailAddress, nil
}

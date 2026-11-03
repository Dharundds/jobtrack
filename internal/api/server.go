package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/dharundds/jobtrack/internal/api/handlers"
	"github.com/dharundds/jobtrack/internal/api/middleware"
	"github.com/dharundds/jobtrack/internal/auth"
	"github.com/dharundds/jobtrack/internal/config"
	"github.com/dharundds/jobtrack/internal/logging"
	"github.com/dharundds/jobtrack/internal/store"
	syncworker "github.com/dharundds/jobtrack/internal/sync"
)

type Server struct {
	engine *gin.Engine
	cfg    *config.Config
}

func NewServer(cfg *config.Config, st *store.Store, worker *syncworker.Worker) *Server {
	if cfg.IsDev() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLog())
	r.Use(middleware.Recovery())
	r.Use(middleware.ErrorHandler())

	sessionStore := cookie.NewStore(cfg.SessionSecret)
	sessionStore.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 30,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	r.Use(sessions.Sessions("jobtrack_session", sessionStore))
	r.Use(sessionOptions())

	if cfg.IsDev() {
		r.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
			AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
			AllowCredentials: true,
		}))
	}

	authSvc := auth.NewService(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURL)
	log := logging.L()

	health := &handlers.HealthHandler{Store: st}
	authH := &handlers.AuthHandler{
		Cfg:   cfg,
		Auth:  authSvc,
		Store: st,
		Log:   log,
	}
	syncH := &handlers.SyncHandler{Worker: worker}
	appsH := &handlers.ApplicationsHandler{Store: st}
	rulesH := &handlers.RulesHandler{Store: st}
	statsH := &handlers.StatsHandler{Store: st}
	settingsH := &handlers.SettingsHandler{Store: st, Cfg: cfg}

	apiGroup := r.Group("/api")
	apiGroup.GET("/health", health.Get)
	apiGroup.GET("/auth/status", authH.Status)

	protected := apiGroup.Group("")
	protected.Use(middleware.RequireAccount(st))
	protected.POST("/sync", syncH.Trigger)
	protected.GET("/sync/status", syncH.Status)

	protected.GET("/applications", appsH.List)
	protected.GET("/applications/:id", appsH.Get)
	protected.PATCH("/applications/:id", appsH.Patch)
	protected.DELETE("/applications/:id", appsH.Delete)
	protected.POST("/applications/:id/messages/:mid/split", appsH.SplitMessage)
	protected.DELETE("/applications/:id/messages/:mid", appsH.IgnoreMessage)
	protected.POST("/applications/merge", appsH.Merge)

	protected.GET("/rules", rulesH.List)
	protected.POST("/rules", rulesH.Create)
	protected.PATCH("/rules/:id", rulesH.Patch)
	protected.DELETE("/rules/:id", rulesH.Delete)

	protected.GET("/stats", statsH.Get)
	protected.GET("/settings", settingsH.Get)
	protected.PATCH("/settings", settingsH.Patch)

	authRoutes := r.Group("/auth")
	authRoutes.GET("/login", authH.Login)
	authRoutes.GET("/callback", authH.Callback)
	authRoutes.POST("/logout", authH.Logout)

	return &Server{engine: r, cfg: cfg}
}

func sessionOptions() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := sessions.Default(c)
		sess.Options(sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 30,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   auth.SecureCookie(c.Request),
		})
		c.Next()
	}
}

func (s *Server) Engine() *gin.Engine {
	return s.engine
}

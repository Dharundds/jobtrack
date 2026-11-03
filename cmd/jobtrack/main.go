package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dharundds/jobtrack/internal/api"
	"github.com/dharundds/jobtrack/internal/auth"
	"github.com/dharundds/jobtrack/internal/config"
	"github.com/dharundds/jobtrack/internal/logging"
	"github.com/dharundds/jobtrack/internal/store"
	syncworker "github.com/dharundds/jobtrack/internal/sync"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		os.Stderr.WriteString("configuration error: " + err.Error() + "\n")
		os.Exit(1)
	}

	log := logging.Setup(cfg.LogLevel, cfg.LogFormat)
	log.Info().
		Str("app_env", cfg.AppEnv).
		Str("listen", cfg.ListenAddr()).
		Msg("starting jobtrack")

	st, err := store.Open(cfg.DBPath)
	if err != nil {
		log.Fatal().Err(err).Msg("open database")
	}
	defer st.Close()

	if err := store.RunMigrations(st.DB, "migrations"); err != nil {
		log.Fatal().Err(err).Msg("run migrations")
	}
	log.Info().Str("path", cfg.DBPath).Msg("database ready")

	authSvc := auth.NewService(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURL)
	worker := syncworker.NewWorker(cfg, st, authSvc, log)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go worker.Start(ctx)

	srv := api.NewServer(cfg, st, worker)

	go func() {
		if err := srv.Engine().Run(cfg.ListenAddr()); err != nil {
			log.Fatal().Err(err).Msg("http server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down")
	cancel()
	time.Sleep(500 * time.Millisecond)
}

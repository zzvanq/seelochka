package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/zzvanq/seelochka/internal/config"
	"github.com/zzvanq/seelochka/internal/http/handlers/url/redirect"
	"github.com/zzvanq/seelochka/internal/http/handlers/url/save"
	"github.com/zzvanq/seelochka/internal/http/middlewares/reqdata"
	"github.com/zzvanq/seelochka/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/zzvanq/seelochka/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

//	@title		Seelochka swagger
//	@version	1.0

// @license.name	Apache 2.0

func main() {
	cfg := config.MustLoad()

	stlog := setupLogger(cfg.Env)
	stlog.Info("logger has been initialized")

	// Storage
	strg, err := sqlite.New(cfg.DbPath)
	if err != nil {
		stlog.Error("failed to initialize a storage", slog.Any("error", err))
		os.Exit(1)
	}
	stlog.Info("storage is initialized")

	// Close the storage
	defer func() {
		if err := strg.Close(); err != nil {
			stlog.Error("failed to close the storage", slog.Any("error", err))
			return
		}

		stlog.Info("storage is closed")
	}()

	// Sentry
	sentrymw, err := setupSentry(cfg)
	if err != nil {
		stlog.Error("failed to setup sentry", slog.Any("error", err))
		return
	}
	defer sentry.Flush(time.Second)

	// Router
	r := chi.NewRouter()
	r.Use(reqdata.New(stlog))
	r.Use(middleware.Recoverer)
	r.Use(sentrymw.Handle)
	r.Use(middleware.RequestID)

	r.Group(func(r chi.Router) {
		r.Post("/", save.New(stlog, strg))
		r.Get("/{alias}", redirect.New(stlog, strg))
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", cfg.Http.Address))))

	// Server
	srv := &http.Server{
		Addr:         cfg.Http.BindHost,
		Handler:      r,
		WriteTimeout: cfg.Http.Timeout,
		IdleTimeout:  cfg.Http.IdleTimeout,
	}

	// Server shutdown
	done := make(chan struct{})
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)

	go func() {
		defer close(done)

		<-sigterm
		stlog.Info("shutting down the server")

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Http.ShutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			stlog.Error("failed to shutdown the server", slog.Any("error", err))
		}
	}()

	stlog.Info("server is starting", slog.String("address", cfg.Http.BindHost))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		signal.Stop(sigterm)
		stlog.Error("failed to start the server", slog.Any("error", err))
		return
	}

	<-done
	stlog.Info("server has stopped")
}

func setupSentry(cfg *config.Config) (*sentryhttp.Handler, error) {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: cfg.SentryDSN,
	}); err != nil {
		return nil, err
	}

	return sentryhttp.New(sentryhttp.Options{Repanic: true}), nil
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}

package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"seelochka/internal/configs"
	"seelochka/internal/http/handlers/urls"
	mwReqdata "seelochka/internal/http/middlewares/reqdata"
	storage "seelochka/internal/storages/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "seelochka/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

//	@title		Seelochka swagger
//	@version	1.0

// @license.name	Apache 2.0
func main() {
	cfg := configs.MustLoad()

	stlog := setupLogger(cfg.Env)
	stlog.Info("logger has been initialized")

	strg, err := storage.New(cfg.DbPath)
	if err != nil {
		stlog.Error("failed to initialize a storage", slog.Any("error", err))
		os.Exit(1)
	}
	stlog.Info("storage is initialized")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(mwReqdata.New(stlog))

	r.Post("/", urls.NewURLSave(stlog, strg))
	r.Get("/{alias}", urls.NewURLRedirect(stlog, strg))

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", cfg.Address)),
	))

	stlog.Info("server is listening", slog.String("address", cfg.Address))
	log.Fatal(http.ListenAndServe(cfg.Address, r))
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

package main

import (
	"context"
	"flag"
	"github.com/co1seam/ember-backend-media/config"
	"github.com/co1seam/ember-backend-media/internal/adapters/repository"
	"github.com/co1seam/ember-backend-media/internal/adapters/rpc"
	"github.com/co1seam/ember-backend-media/internal/core/models"
	"github.com/co1seam/ember-backend-media/internal/core/services"
	"github.com/gofiber/fiber/v2/log"
	"log/slog"
	"os"
)

func main() {
	ctx := context.Background()

	cfgFlag := flag.String("config", "", "flag to add config path")
	flag.Parse()

	cfg, err := config.New(cfgFlag)
	if err != nil {
		log.Fatal(err)
	}

	handlerOpts := &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevelChoice(cfg.App.LogLevel),
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, handlerOpts))

	slog.SetDefault(log)

	db, err := repository.NewPostgres(ctx, &cfg.Database)
	if err != nil {
		log.Error("error initializing PostgreSQL", err)
		return
	}
	defer db.Close()

	minioClient, err := repository.NewMinio(&cfg.MinIO)
	if err != nil {
		log.Error("error initializing MinIO", err)
		return
	}

	cache := repository.NewRedis(cfg.Redis.Host, cfg.Redis.Port)

	opts := &models.Options{
		Logger: log,
		Config: cfg,
	}

	repos := repository.NewRepository(db.DB, minioClient, cache, opts)
	service := services.NewService(repos, opts)
	handler := rpc.NewHandler(service, opts)

	server := rpc.NewServer()
	if err := server.Run(handler); err != nil {
		log.Error("server run error", err)
	}

	if err := db.Close(); err != nil {
		log.Error("error closing DB", err)
	}
}

func logLevelChoice(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

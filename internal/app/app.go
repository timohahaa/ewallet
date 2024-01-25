package app

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/timohahaa/ewallet/config"
	v1 "github.com/timohahaa/ewallet/internal/controllers/http/v1"
	"github.com/timohahaa/ewallet/internal/repository"
	"github.com/timohahaa/ewallet/internal/service"
	"github.com/timohahaa/ewallet/pkg/httpserver"
	"github.com/timohahaa/postgres"
)

func Run(configFilePath string) {
	// Config
	cfg, err := config.NewConfig(configFilePath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// database
	slog.Info("Initializing postgres connection...")
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxConnPoolSize(cfg.PG.ConnPoolSize))
	if err != nil {
		log.Fatalf("error connecting to postgres: %s", err)
	}

	// транспортный слой
	slog.Info("Initializing repositories...")
	walletRepo := repository.NewWalletRepo(pg)

	// слой БЛ
	slog.Info("Initializing services...")
	walletService := service.NewWalletService(walletRepo)

	// слой представления - handlers and routes
	slog.Info("Initializing handlers and routes...")
	handler := v1.NewRouter(walletService)

	slog.Info("Starting http server...", "port", cfg.Server.Port)
	server := httpserver.New(handler, httpserver.Port(cfg.Server.Port))

	// gracefull shutdown
	slog.Info("Configuring gracefull shutdown...")
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	slog.Info("Server started!", "port", cfg.Server.Port)

	<-shutdownChan

	slog.Info("Shutting down...")
	err = server.Shutdown()
	if err != nil {
		log.Fatalf("Error shutting down http server: %s", err)
	}
}

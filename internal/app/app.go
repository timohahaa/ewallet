package app

import (
	log2 "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/timohahaa/ewallet/config"
	v1 "github.com/timohahaa/ewallet/internal/controllers/http/v1"
	"github.com/timohahaa/ewallet/internal/repository"
	"github.com/timohahaa/ewallet/internal/service"
	"github.com/timohahaa/ewallet/pkg/httpserver"
	log "github.com/timohahaa/ewallet/pkg/logger"
	"github.com/timohahaa/postgres"
)

func Run(configFilePath string) {
	// Config
	cfg, err := config.NewConfig(configFilePath)
	if err != nil {
		log2.Fatalf("config error: %s", err)
	}
	logger := log.GetLogger("internal.log", cfg.Server.LogPath)
	httpLogger := log.GetLogger("requests.log", cfg.Server.LogPath)

	// database
	logger.Info("initializing postgres connection...")
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxConnPoolSize(cfg.PG.ConnPoolSize))
	if err != nil {
		logger.WithFields(logrus.Fields{"error": err}).Fatal("error connecting to postgres")
	}

	// транспортный слой
	logger.Info("initializing repositories...")
	walletRepo := repository.NewWalletRepo(pg, logger)

	// слой БЛ
	logger.Info("initializing services...")
	walletService := service.NewWalletService(walletRepo, logger)

	// слой представления - handlers and routes
	logger.Info("initializing handlers and routes...")
	handler := v1.NewRouter(walletService, httpLogger)

	logger.Infof("starting http server...")
	server := httpserver.New(handler, httpserver.Port(cfg.Server.Port))

	// gracefull shutdown
	logger.Info("configuring gracefull shutdown...")
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	logger.WithFields(logrus.Fields{"port": cfg.Server.Port}).Info("server started!")

	<-shutdownChan

	logger.Info("shutting down...")
	err = server.Shutdown()
	if err != nil {
		logger.WithFields(logrus.Fields{"error": err}).Fatal("error shutting down the server")
	}
}

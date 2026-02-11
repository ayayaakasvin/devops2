package main

import (
	"application-for-kubernetes/internal/app"
	"application-for-kubernetes/internal/config"
	"application-for-kubernetes/internal/logger"
	"application-for-kubernetes/internal/repository/postgresql"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log := logger.SetupLogger("Application")

	cfg := config.MustLoadConfig()
	log.Infof("Configs retrieved: %v", cfg)
	log.Infof("%s", time.Now().String())

	repo, err := postgresql.NewPostgreSQLConnection(cfg.DatabaseConfig)
	if err != nil {
		log.Errorf("failed to set session id: %v", err)
		os.Exit(1)
	}

	srv := app.NewServerApp(&cfg.HTTPServer, log, repo)

	// Error channel to monitor critical errors
	errChan := make(chan error, 1)
	go func() {
		errChan <- srv.Run()
	}()

	// Signal channel for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		log.Infof("Received signal: %v, shutting down gracefully", sig)
		srv.Shutdown()
	case err := <-errChan:
		if err != nil {
			log.Errorf("Server error: %v", err)
		}
	}
}

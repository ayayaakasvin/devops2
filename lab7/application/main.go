package main

import (
	"application-for-kubernetes/internal/app"
	"application-for-kubernetes/internal/config"
	aliveapp "application-for-kubernetes/internal/libs/alive-app"
	"application-for-kubernetes/internal/logger"
	"application-for-kubernetes/internal/repository/valkey"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ayayaakasvin/goroutinesupervisor"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()

	log := logger.SetupLogger("Application")
	cfg := config.MustLoadConfig()

	log.Infof("Configs retrieved: %v", cfg)
	log.Infof("%s", time.Now().String())
	log.Infof("New Image Set, logs for checking")

	cache, err := valkey.NewValkeyClient(cfg.CacheConfig)
	if err != nil {
		log.Errorf("failed to set session id: %v", err)
		os.Exit(1)
	}

	gs := goroutinesupervisor.NewSupervisor(ctx)

	srv := app.NewServerApp(&cfg.HTTPServer, log, cache)

	gs.Go("HTTP-Server", srv.Start)
	gs.Go("Server Status", aliveapp.LogAppStatus(3*time.Minute, log, ctx))
	gs.Go("Memory Stats", aliveapp.MemStat(2*time.Minute, log, ctx))

	err = gs.Wait()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	srv.Stop(shutdownCtx)
	cache.Close()

	return err
}
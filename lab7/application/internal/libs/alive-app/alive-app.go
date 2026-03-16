package aliveapp

import (
	"context"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

func LogAppStatus(trate time.Duration, log *logrus.Logger, ctx context.Context) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		ticker := time.NewTicker(trate)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				log.Info("Server is alive...")
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

func MemStat(trate time.Duration, log *logrus.Logger, ctx context.Context) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		ticker := time.NewTicker(time.Second * 15)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				var m runtime.MemStats
				runtime.ReadMemStats(&m)
				log.Infof("Alloc = %v MiB", m.Alloc/1024/1024)
				time.Sleep(1 * time.Second)
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

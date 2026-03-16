package domain

import (
	"context"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string) (any, error)
	Delete(ctx context.Context, key string) error
	Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error)

	Info(ctx context.Context) (map[string]string, error)
	Ping(ctx context.Context) error
	Close() error
}

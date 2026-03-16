package valkey

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"application-for-kubernetes/internal/config"
	"application-for-kubernetes/internal/domain"

	"github.com/redis/go-redis/v9"
)

// for storing methods of storing and retrieving session_id
type Cache struct {
	connection *redis.Client
}

func NewValkeyClient(cfg config.CacheConfig) (domain.Cache, error) {
	ctx := context.Background()
	opt, err := redis.ParseURL(cfg.URL)
	log.Printf("URL: %s", cfg.URL)
	if err != nil {
		msg := fmt.Sprintf("failed to parse Redis URL: %v", err)
		log.Println(msg)
		return nil, err
	}

	conn := redis.NewClient(opt)
	if err := conn.Ping(ctx).Err(); err != nil {
		msg := fmt.Sprintf("failed to connect to db: %v", err)
		log.Println(msg)
		return nil, err
	}

	return &Cache{
		connection: conn,
	}, nil
}

func (c *Cache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return c.connection.Set(ctx, key, value, ttl).Err()
}

func (c *Cache) Get(ctx context.Context, key string) (any, error) {
	return c.connection.Get(ctx, key).Result()
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.connection.Del(ctx, key).Err()
}

func (c *Cache) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return c.connection.Scan(ctx, cursor, match, count).Result()
}

func (c *Cache) Info(ctx context.Context) (map[string]string, error) {
	info, err := c.connection.Info(ctx).Result()
	if err != nil {
		return nil, err
	}

	return parseRedisInfo(info), nil
}

func (c *Cache) Ping(ctx context.Context) error {
	return c.connection.Ping(ctx).Err()
}

func (c *Cache) Close() error {
	return c.connection.Close()
}

func parseRedisInfo(info string) map[string]string {
	result := make(map[string]string)

	var currentSection string

	lines := strings.Split(info, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			// comment or empty line
			if strings.HasPrefix(line, "# ") {
				currentSection = strings.TrimSpace(strings.TrimPrefix(line, "#"))
			}
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if currentSection == "" {
			// rare case — key without section
			result[key] = value
		} else {
			fullKey := currentSection + ":" + key
			result[fullKey] = value
		}
	}

	return result
}
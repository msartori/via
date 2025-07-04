package redis_ds

import (
	"context"
	"strconv"
	"time"
	"via/internal/ds"
	"via/internal/secret"

	"github.com/redis/go-redis/v9"
)

type RedisDS struct {
	client *redis.Client
}

func New(cfg ds.DSConfig) ds.DS {
	if cfg.Password == "" {
		cfg.Password = secret.ReadSecret(cfg.PasswordFile)
	}
	base, _ := strconv.Atoi(cfg.Base)
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       base,
	})
	return &RedisDS{client: rdb}
}

func (r *RedisDS) Set(ctx context.Context, key string, value string, ttlSeconds int) error {
	return r.client.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *RedisDS) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisDS) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

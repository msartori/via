package ratelimit

import (
	"context"
	"time"
	"via/internal/ds"
)

type RateLimiter struct {
	ID     string
	limit  int
	window time.Duration
	ds     ds.DS
}

func New(id string, limit int, window time.Duration, ds ds.DS) *RateLimiter {
	return &RateLimiter{
		ID:     id,
		limit:  limit,
		window: window,
		ds:     ds,
	}
}

func (rl *RateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	counter, err := rl.ds.Incr(ctx, key)
	if err != nil {
		return false, err
	}
	if counter == 1 {
		// Set expiration for the key if it's the first request
		err = rl.ds.Set(ctx, key, "1", int(rl.window.Seconds()))
		if err != nil {
			return false, err
		}
	}
	if counter > int64(rl.limit) {
		// If the limit is exceeded, we can either return false or reset the counter
		return false, nil
	}
	return true, nil
}

package redis_pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"via/internal/pubsub"
	"via/internal/secret"

	"github.com/redis/go-redis/v9"
)

// RedisPubSub implements pub/sub using Redis.
type RedisPubSub struct {
	client *redis.Client
}

// NewRedisPubSub creates a new RedisPubSub.
func New(cfg pubsub.PubSubConfig) *RedisPubSub {
	if cfg.Password == "" {
		cfg.Password = secret.Get().Read(cfg.PasswordFile)
	}
	base, _ := strconv.Atoi(cfg.Base)
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       base,
	})
	return &RedisPubSub{client: rdb}
}

// Publish sends a message to a Redis channel.
func (r *RedisPubSub) Publish(ctx context.Context, channel string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}
	return r.client.Publish(ctx, channel, data).Err()
}

// Subscribe subscribes to multiple channels and returns a subscription.
func (r *RedisPubSub) Subscribe(ctx context.Context, channels ...string) (pubsub.Subscription, error) {
	ps := r.client.Subscribe(ctx, channels...)
	sub := &redisSubscription{
		ps:  ps,
		ch:  make(chan pubsub.Message),
		ctx: ctx,
	}
	go sub.listen()
	return sub, nil
}

// redisSubscription handles the actual message loop.
type redisSubscription struct {
	ps  pubSub
	ch  chan pubsub.Message
	ctx context.Context
}

type pubSub interface {
	Channel(opts ...redis.ChannelOption) <-chan *redis.Message
	Close() error
}

func (r *redisSubscription) listen() {
	defer close(r.ch)
	redisChannel := r.ps.Channel()

	for {
		select {
		case <-r.ctx.Done():
			return
		case msg, ok := <-redisChannel:
			if !ok {
				return
			}
			var payload any
			err := json.Unmarshal([]byte(msg.Payload), &payload)
			if err != nil {
				//log error
				continue
			}
			r.ch <- pubsub.Message{Channel: msg.Channel, Payload: payload}
		}
	}
}

func (r *redisSubscription) Channel() <-chan pubsub.Message {
	return r.ch
}

func (r *redisSubscription) Close() error {
	return r.ps.Close()
}

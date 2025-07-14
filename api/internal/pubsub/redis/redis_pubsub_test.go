package redis_pubsub

import (
	"context"
	"encoding/json"
	"testing"
	"time"
	"via/internal/pubsub"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"

	"github.com/stretchr/testify/assert"
)

func TestPublish_Success(t *testing.T) {
	db, mock := redismock.NewClientMock()
	rds := &RedisPubSub{client: db}

	ctx := context.Background()
	payload := map[string]string{"foo": "bar"}
	jsonData, _ := json.Marshal(payload)

	mock.ExpectPublish("test_channel", jsonData).SetVal(1)

	err := rds.Publish(ctx, "test_channel", payload)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPublish_JSONError(t *testing.T) {
	db, _ := redismock.NewClientMock()
	rds := &RedisPubSub{client: db}

	ctx := context.Background()
	payload := make(chan int) // cannot be JSON marshaled

	err := rds.Publish(ctx, "test_channel", payload)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "marshal error")
}

func TestSubscribe(t *testing.T) {
	db, _ := redismock.NewClientMock()
	rds := &RedisPubSub{client: db}

	ctx := context.Background()

	// Test that Subscribe returns a non-nil Subscription without error
	sub, err := rds.Subscribe(ctx, "ch1", "ch2")
	assert.NoError(t, err)
	assert.NotNil(t, sub)
}

type fakePubSub struct {
	ch       chan *redis.Message
	closeErr error
}

func (f *fakePubSub) Channel(opts ...redis.ChannelOption) <-chan *redis.Message {
	return f.ch
}

func (f *fakePubSub) Close() error {
	return f.closeErr
}

func TestRedisSubscription_Listen_Success(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fps := &fakePubSub{
		ch: make(chan *redis.Message),
	}

	sub := &redisSubscription{
		ps:  fps,
		ch:  make(chan pubsub.Message, 1),
		ctx: ctx,
	}

	go sub.listen()

	go func() {
		fps.ch <- &redis.Message{Channel: "ch1", Payload: `{"foo":"bar"}`}
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	select {
	case msg := <-sub.Channel():
		assert.Equal(t, "ch1", msg.Channel)
		assert.NotNil(t, msg.Payload)
	case <-time.After(time.Second):
		t.Fatal("Timeout waiting for message")
	}
}

func TestRedisSubscription_Listen_InvalidJSON(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fps := &fakePubSub{
		ch: make(chan *redis.Message),
	}

	sub := &redisSubscription{
		ps:  fps,
		ch:  make(chan pubsub.Message, 1),
		ctx: ctx,
	}

	go sub.listen()

	go func() {
		fps.ch <- &redis.Message{Channel: "ch1", Payload: "invalid_json"}
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	select {
	case _, ok := <-sub.Channel():
		if ok {
			t.Fatal("Should not receive message with invalid JSON")
		}
	case <-time.After(100 * time.Millisecond):
		// Correct behavior: invalid JSON ignored
	}
}

func TestRedisSubscription_Close(t *testing.T) {
	db, _ := redismock.NewClientMock()
	ctx := context.Background()

	ps := db.Subscribe(ctx, "ch1")
	sub := &redisSubscription{
		ps:  ps,
		ch:  make(chan pubsub.Message, 1),
		ctx: ctx,
	}

	assert.NoError(t, sub.Close())
}

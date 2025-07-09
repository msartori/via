package redis_ds

import (
	"context"
	"testing"
	"time"
	"via/internal/ds"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestNewRedisDS(t *testing.T) {
	// Override secret.ReadSecret temporarily
	originalReadSecret := readSecret
	readSecret = func(_ string) string {
		return "mockedPassword"
	}
	defer func() { readSecret = originalReadSecret }()

	cfg := ds.DSConfig{
		Host:         "localhost",
		Port:         "6379",
		Base:         "0",
		Password:     "", // triggers secret.ReadSecret
		PasswordFile: "dummy/path",
	}

	dsInstance := New(cfg)

	assert.NotNil(t, dsInstance)
	// We can’t really test the redis.Client internals without integration, but we can check type
	_, ok := dsInstance.(*RedisDS)
	assert.True(t, ok, "expected type *RedisDS")
}

func TestRedisDS(t *testing.T) {
	db, mock := redismock.NewClientMock()
	r := &RedisDS{client: db}
	ctx := context.Background()

	t.Run("Set success", func(t *testing.T) {
		mock.ExpectSet("key", "value", 5*time.Second).SetVal("OK")
		err := r.Set(ctx, "key", "value", 5)
		assert.NoError(t, err)
	})

	t.Run("Get key exists", func(t *testing.T) {
		mock.ExpectGet("key").SetVal("value")
		found, val, err := r.Get(ctx, "key")
		assert.NoError(t, err)
		assert.True(t, found)
		assert.Equal(t, "value", val)
	})

	t.Run("Get key does not exist", func(t *testing.T) {
		mock.ExpectGet("key").RedisNil()
		found, val, err := r.Get(ctx, "key")
		assert.NoError(t, err)
		assert.False(t, found)
		assert.Equal(t, "", val)
	})

	t.Run("Del success", func(t *testing.T) {
		mock.ExpectDel("key").SetVal(1)
		err := r.Del(ctx, "key")
		assert.NoError(t, err)
	})

	t.Run("Incr success", func(t *testing.T) {
		mock.ExpectIncr("key").SetVal(42)
		val, err := r.Incr(ctx, "key")
		assert.NoError(t, err)
		assert.Equal(t, int64(42), val)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

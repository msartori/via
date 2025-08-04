package pubsub

import (
	"context"
	"sync"
)

type PubSubConfig struct {
	Host         string `env:"HOST" envDefault:"localhost:6379" json:"host"`
	Password     string `env:"PASSWORD" envDefault:"" json:"-"`
	PasswordFile string `env:"PASSWORD_FILE" envDefault:"" json:"passwordFile"`
	Base         string `env:"BASE" envDefault:"0" json:"-"`
	Port         string `env:"PORT" envDefault:"6379" json:"port"`
}

// PubSub defines a generic pub/sub interface.
type PubSub interface {
	// Subscribe subscribes to multiple channels and returns a subscription.
	Subscribe(ctx context.Context, channels ...string) (Subscription, error)

	// Publish publishes a message to a channel.
	Publish(ctx context.Context, channel string, message any) error
}

// Subscription represents an active subscription.
type Subscription interface {
	Channel() <-chan Message
	Close() error
}

// Message represents a message with a generic payload.
type Message struct {
	Channel string
	Payload any
}

var (
	instance PubSub
	mutex    = &sync.RWMutex{}
)

func Get() PubSub {
	mutex.RLock()
	defer mutex.RUnlock()
	return instance
}

func Set(pubsub PubSub) {
	mutex.Lock()
	defer mutex.Unlock()
	instance = pubsub
}

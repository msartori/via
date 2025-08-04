package ds

import (
	"context"
	"sync"
)

type DSConfig struct {
	Host         string `env:"HOST" envDefault:"localhost:6379" json:"host"`
	Password     string `env:"PASSWORD" envDefault:"" json:"-"`
	PasswordFile string `env:"PASSWORD_FILE" envDefault:"" json:"passwordFile"`
	Base         string `env:"BASE" envDefault:"0" json:"-"`
	Port         string `env:"PORT" envDefault:"6379" json:"port"`
}

// DataStore is an interface that defines methods for setting, getting, and deleting data in a datastore.
type DS interface {
	Set(ctx context.Context, key string, value string, ttlSeconds int) error
	Get(ctx context.Context, key string) (bool, string, error)
	Del(ctx context.Context, key string) error
	Incr(ctx context.Context, key string) (int64, error)
}

var (
	instance DS
	mutex    = &sync.RWMutex{}
)

func Get() DS {
	mutex.RLock()
	defer mutex.RUnlock()
	return instance
}

func Set(ds DS) {
	mutex.Lock()
	defer mutex.Unlock()
	instance = ds
}

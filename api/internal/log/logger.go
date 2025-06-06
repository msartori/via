package log

import (
	"context"
	"net/http"
	"sync"
)

type Logger interface {
	Trace(ctx context.Context, pairs ...any)
	Debug(ctx context.Context, pairs ...any)
	Info(ctx context.Context, pairs ...any)
	Warn(ctx context.Context, pairs ...any)
	Error(ctx context.Context, err error, pairs ...any)
	Fatal(ctx context.Context, err error, pairs ...any)
	WithLogFieldsInRequest(r *http.Request, pairs ...any) *http.Request
	WithLogFields(ctx context.Context, pairs ...any) context.Context
}

var (
	instance Logger
	mutex    = &sync.RWMutex{}
)

func Get() Logger {
	mutex.RLock()
	defer mutex.RUnlock()
	return instance
}

func Set(logger Logger) {
	mutex.Lock()
	defer mutex.Unlock()
	instance = logger
}

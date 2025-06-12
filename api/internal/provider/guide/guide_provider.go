package guide_provider

import (
	"context"
	"sync"
	"via/internal/model"
)

type GuideProvider interface {
	GetGuide(ctx context.Context, id string) (model.Guide, error)
}

var (
	instance GuideProvider
	mutex    = &sync.RWMutex{}
)

func Get() GuideProvider {
	mutex.RLock()
	defer mutex.RUnlock()
	return instance
}

func Set(logger GuideProvider) {
	mutex.Lock()
	defer mutex.Unlock()
	instance = logger
}

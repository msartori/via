package via_guide_provider

import (
	"context"
	"sync"
	"via/internal/model"
)

type ViaGuideProvider interface {
	GetGuide(ctx context.Context, id string) (model.ViaGuide, error)
}

var (
	instance ViaGuideProvider
	mutex    = &sync.RWMutex{}
)

func Get() ViaGuideProvider {
	mutex.RLock()
	defer mutex.RUnlock()
	return instance
}

func Set(provider ViaGuideProvider) {
	mutex.Lock()
	defer mutex.Unlock()
	instance = provider
}

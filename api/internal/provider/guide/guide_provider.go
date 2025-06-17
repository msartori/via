package guide_provider

import (
	"context"
	"sync"
	"via/internal/model"
)

type GuideProvider interface {
	GetGuideByCode(ctx context.Context, code string) (model.GuideProcess, error)
	CreateGuide(ctx context.Context, guide model.ViaGuide) (int, error)
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

func Set(guideProcessProvider GuideProvider) {
	mutex.Lock()
	defer mutex.Unlock()
	instance = guideProcessProvider
}

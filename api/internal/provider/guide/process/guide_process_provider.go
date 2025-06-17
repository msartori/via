package guide_process_provider

import (
	"context"
	"sync"
	"via/internal/model"
)

type GuideProcessProvider interface {
	GetGuideProcessByCode(ctx context.Context, code string) (model.GuideProcess, error)
	CreateGuide(ctx context.Context, guide model.Guide) (int, error)
}

var (
	instance GuideProcessProvider
	mutex    = &sync.RWMutex{}
)

func Get() GuideProcessProvider {
	mutex.RLock()
	defer mutex.RUnlock()
	return instance
}

func Set(guideProcessProvider GuideProcessProvider) {
	mutex.Lock()
	defer mutex.Unlock()
	instance = guideProcessProvider
}

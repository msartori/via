package guide_provider

import (
	"context"
	"sync"
	"via/internal/model"
)

type GuideProvider interface {
	GetGuideByViaGuideId(ctx context.Context, viaGuideId string) (model.Guide, error)
	CreateGuide(ctx context.Context, guide model.ViaGuide) (int, error)
	ReinitGuide(ctx context.Context, id int) (int, error)
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

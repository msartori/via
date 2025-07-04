package guide_provider

import (
	"context"
	"sync"
	"via/internal/model"
)

type GuideProvider interface {
	GetGuideByViaGuideId(ctx context.Context, viaGuideId string) (model.Guide, error)
	CreateGuide(ctx context.Context, guide model.ViaGuide) (int, error)
	//ReinitGuide(ctx context.Context, id int) (int, error)
	GetGuidesByStatus(ctx context.Context, status []string) ([]model.Guide, error)
	UpdateGuide(ctx context.Context, guide model.Guide) error
	GetGuideById(ctx context.Context, id int) (model.Guide, error)
	GetGuideHistory(ctx context.Context, guideId int) ([]model.GuideHistory, error)
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

func Set(guideProvider GuideProvider) {
	mutex.Lock()
	defer mutex.Unlock()
	instance = guideProvider
}

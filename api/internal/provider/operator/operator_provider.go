package operator_provider

import (
	"context"
	"sync"
	"via/internal/model"
)

type OperatorProvider interface {
	GetOperators(ctx context.Context) ([]model.Operator, error)
	GetOperatorByAccount(ctx context.Context, account string) (model.Operator, error)
	GetOperatorById(ctx context.Context, id int) (model.Operator, error)
}

var (
	instance OperatorProvider
	mutex    = &sync.RWMutex{}
)

func Get() OperatorProvider {
	mutex.RLock()
	defer mutex.RUnlock()
	return instance
}

func Set(operatorProvider OperatorProvider) {
	mutex.Lock()
	defer mutex.Unlock()
	instance = operatorProvider
}

package biz_operator

import (
	"context"
	"time"
	"via/internal/model"
	operator_provider "via/internal/provider/operator"

	"github.com/patrickmn/go-cache"
)

var operatorCache = cache.New(5*time.Minute, 10*time.Minute)

const OPERATOR_SYSTEM int = 1

func GetOperatorByAccount(ctx context.Context, account string) (model.Operator, error) {
	if data, found := operatorCache.Get(account); found {
		return data.(model.Operator), nil
	}
	operator, err := operator_provider.Get().GetOperatorByAccount(ctx, account)
	if err != nil {
		return model.Operator{}, err
	}
	if operator.ID != 0 {
		operatorCache.Set(account, operator, cache.DefaultExpiration)
	}
	return operator, nil
}

func ClearCache() {
	operatorCache.Flush()
}

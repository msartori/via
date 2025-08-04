package biz_operator

import (
	"context"
	"errors"
	"testing"

	"via/internal/model"
	operator_provider "via/internal/provider/operator"
	mock_operator_provider "via/internal/provider/operator/mock"

	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetOperatorByAccount(t *testing.T) {
	ctx := context.Background()
	account := "test_account"
	expectedOperator := model.Operator{ID: 123, Name: "Test Operator"}

	t.Run("returns from cache", func(t *testing.T) {
		// Prime the cache
		operatorCache.Set(account, expectedOperator, cache.DefaultExpiration)

		// Call
		result, err := GetOperatorByAccount(ctx, account)
		assert.NoError(t, err)
		assert.Equal(t, expectedOperator, result)
	})

	t.Run("fetches from provider and caches", func(t *testing.T) {
		operatorCache.Flush()

		mockProvider := new(mock_operator_provider.MockOperatorProvider)
		mockProvider.On("GetOperatorByAccount", mock.Anything, account).
			Return(expectedOperator, nil).Once()

		operator_provider.Set(mockProvider)

		// Call
		result, err := GetOperatorByAccount(ctx, account)
		assert.NoError(t, err)
		assert.Equal(t, expectedOperator, result)

		// Ensure cached now
		cached, found := operatorCache.Get(account)
		assert.True(t, found)
		assert.Equal(t, expectedOperator, cached)

		mockProvider.AssertExpectations(t)
	})

	t.Run("provider returns error", func(t *testing.T) {
		operatorCache.Flush()

		mockProvider := new(mock_operator_provider.MockOperatorProvider)
		mockProvider.On("GetOperatorByAccount", mock.Anything, account).
			Return(model.Operator{}, errors.New("provider error")).Once()

		operator_provider.Set(mockProvider)

		// Call
		result, err := GetOperatorByAccount(ctx, account)
		assert.Error(t, err)
		assert.Empty(t, result)

		mockProvider.AssertExpectations(t)
	})

	t.Run("provider returns operator with ID 0 (not cached)", func(t *testing.T) {
		operatorCache.Flush()

		mockProvider := new(mock_operator_provider.MockOperatorProvider)
		mockProvider.On("GetOperatorByAccount", mock.Anything, account).
			Return(model.Operator{ID: 0}, nil).Once()

		operator_provider.Set(mockProvider)

		// Call
		result, err := GetOperatorByAccount(ctx, account)
		assert.NoError(t, err)
		assert.Equal(t, model.Operator{ID: 0}, result)

		_, found := operatorCache.Get(account)
		assert.False(t, found)

		mockProvider.AssertExpectations(t)
	})
}

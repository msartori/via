package ds

import (
	"context"
	"testing"
	mock_ds "via/internal/ds/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSetAndGetInstance(t *testing.T) {
	mockDs := new(mock_ds.MockDS)

	// Set mock DS instance
	Set(mockDs)

	// Retrieve it
	instance := Get()
	assert.NotNil(t, instance)

	// Test Set
	mockDs.On("Set", mock.Anything, "key", "value", 10).Return(nil).Once()
	err := instance.Set(context.Background(), "key", "value", 10)
	assert.NoError(t, err)
	mockDs.AssertExpectations(t)

	// Test Get
	mockDs.On("Get", mock.Anything, "key").Return(true, "value", nil).Once()
	found, val, err := instance.Get(context.Background(), "key")
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "value", val)
	mockDs.AssertExpectations(t)

	// Test Del
	mockDs.On("Del", mock.Anything, "key").Return(nil).Once()
	err = instance.Del(context.Background(), "key")
	assert.NoError(t, err)
	mockDs.AssertExpectations(t)

	// Test Incr
	mockDs.On("Incr", mock.Anything, "key").Return(1, nil).Once()
	count, err := instance.Incr(context.Background(), "key")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
	mockDs.AssertExpectations(t)
}

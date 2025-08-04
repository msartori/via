package pubsub_test

import (
	"context"
	"testing"
	"via/internal/pubsub"
	mock_pubsub "via/internal/pubsub/mock"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
)

func TestSetAndGetInstance(t *testing.T) {
	mockPubSub := new(mock_pubsub.MockPubSub)
	mockSubscription := new(mock_pubsub.MockSubscription)

	// Set mock PubSub instance
	pubsub.Set(mockPubSub)

	instance := pubsub.Get()
	if instance != mockPubSub {
		t.Error("expected same instance from Get after Set")
	}

	// Test Subscribe
	mockPubSub.On("Subscribe", mock.Anything, "channel").Return(mockSubscription, nil)
	subscription, err := instance.Subscribe(context.Background(), "channel")
	assert.NoError(t, err)
	assert.Equal(t, mockSubscription, subscription)

	// Test Publish
	mockPubSub.On("Publish", mock.Anything, "channel", "message").Return(nil)
	err = instance.Publish(context.Background(), "channel", "message")
	assert.NoError(t, err)

	// Test Channel
	ch := make(chan pubsub.Message)
	var recvCh <-chan pubsub.Message = ch
	mockSubscription.On("Channel").Return(recvCh)
	rch := mockSubscription.Channel()
	assert.Equal(t, recvCh, rch)

	// Test Close
	mockSubscription.On("Close").Return(nil)
	err = mockSubscription.Close()
	assert.NoError(t, err)

	mockPubSub.AssertExpectations(t)
}

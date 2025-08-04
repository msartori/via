package mock_pubsub

import (
	"context"
	"via/internal/pubsub"

	"github.com/stretchr/testify/mock"
)

type MockSubscription struct {
	mock.Mock
}

type MockPubSub struct {
	mock.Mock
}

func (m *MockSubscription) Channel() <-chan pubsub.Message {
	args := m.Called()
	return args.Get(0).(<-chan pubsub.Message)
}

func (m *MockSubscription) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPubSub) Subscribe(ctx context.Context, channels ...string) (pubsub.Subscription, error) {
	args := m.Called(append([]any{ctx}, stringSliceToInterfaceSlice(channels)...)...)
	//args := m.Called(ctx, channels)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockSubscription), args.Error(1)
}

func (m *MockPubSub) Publish(ctx context.Context, channel string, message any) error {
	args := m.Called(ctx, channel, message)
	return args.Error(0)
}

func stringSliceToInterfaceSlice(s []string) []interface{} {
	res := make([]interface{}, len(s))
	for i, v := range s {
		res[i] = v
	}
	return res
}

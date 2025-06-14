package mock_log

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Trace(ctx context.Context, keyvals ...any) {
	m.Called(ctx, keyvals)
}

func (m *MockLogger) Debug(ctx context.Context, keyvals ...any) {
	m.Called(ctx, keyvals)
}

func (m *MockLogger) Info(ctx context.Context, keyvals ...any) {
	m.Called(ctx, keyvals)
}

func (m *MockLogger) Warn(ctx context.Context, keyvals ...any) {
	m.Called(ctx, keyvals)
}
func (m *MockLogger) Error(ctx context.Context, err error, keyvals ...any) {
	m.Called(ctx, err, keyvals)
}
func (m *MockLogger) Fatal(ctx context.Context, err error, keyvals ...any) {
	m.Called(ctx, err, keyvals)
}
func (m *MockLogger) WithLogFields(ctx context.Context, keyvals ...any) context.Context {
	args := m.Called(ctx, keyvals)
	return args.Get(0).(context.Context)
}

func (m *MockLogger) WithLogFieldsInRequest(r *http.Request, keyvals ...any) *http.Request {
	args := m.Called(r, keyvals)
	return args.Get(0).(*http.Request)
}

type MockNoOpLogger struct {
}

func (m *MockNoOpLogger) Trace(ctx context.Context, keyvals ...any) {

}

func (m *MockNoOpLogger) Debug(ctx context.Context, keyvals ...any) {
}

func (m *MockNoOpLogger) Info(ctx context.Context, keyvals ...any) {
}

func (m *MockNoOpLogger) Warn(ctx context.Context, keyvals ...any) {
}
func (m *MockNoOpLogger) Error(ctx context.Context, err error, keyvals ...any) {
}
func (m *MockNoOpLogger) Fatal(ctx context.Context, err error, keyvals ...any) {

}
func (m *MockNoOpLogger) WithLogFields(ctx context.Context, keyvals ...any) context.Context {
	return ctx
}

func (m *MockNoOpLogger) WithLogFieldsInRequest(r *http.Request, keyvals ...any) *http.Request {
	return r
}

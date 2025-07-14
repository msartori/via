package sse

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"via/internal/i18n"
	"via/internal/pubsub"
	mock_pubsub "via/internal/pubsub/mock"
	"via/internal/response"
	"via/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSubscription struct {
	mock.Mock
	ch chan pubsub.Message
}

func (m *mockSubscription) Channel() <-chan pubsub.Message {
	return m.ch
}

func (m *mockSubscription) Close() error {
	return nil
}

type mockWriter struct {
	header http.Header
	events []string
	Body   *bytes.Buffer
}

func newMockWriter() *mockWriter {
	return &mockWriter{
		header: http.Header{},
		events: make([]string, 0),
		Body:   new(bytes.Buffer),
	}
}

func (w *mockWriter) Header() http.Header {
	return w.header
}

func (w *mockWriter) Write(data []byte) (int, error) {
	w.events = append(w.events, string(data))
	w.Body.Write(data)
	return len(data), nil
}

func (w *mockWriter) WriteHeader(statusCode int) {}

func TestHandleSSE_AllPaths(t *testing.T) {
	testutil.InjectNoOpLogger()
	loader := func(r *http.Request) response.Response[any] {
		return response.Response[any]{Data: "ok"}
	}
	t.Run("invalid flusher", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?lang=es", nil)
		w := newMockWriter()
		HandleSSE(w, req, loader)
		clean := strings.TrimPrefix(w.Body.String(), "data:")
		clean = strings.TrimSpace(clean)
		var resp response.Response[any]
		err := json.Unmarshal([]byte(clean), &resp)
		assert.NoError(t, err)
		assert.Equal(t, i18n.Get(req, i18n.MsgInternalServerError), resp.Message)
	})
	t.Run("subscribe error", func(t *testing.T) {
		mockPubSub := new(mock_pubsub.MockPubSub)
		pubsub.Set(mockPubSub)
		mockPubSub.On("Subscribe", mock.Anything, mock.Anything).Return(nil, errors.New("fail"))
		req := httptest.NewRequest("GET", "/?lang=es", nil)
		w := newMockWriter()
		HandleSSE(w, req, loader)
		clean := strings.TrimPrefix(w.Body.String(), "data:")
		clean = strings.TrimSpace(clean)
		var resp response.Response[any]
		err := json.Unmarshal([]byte(clean), &resp)
		assert.NoError(t, err)
		assert.Equal(t, i18n.Get(req, i18n.MsgInternalServerError), resp.Message)
	})
	t.Run("channel closed", func(t *testing.T) {
		mockPubSub := new(mock_pubsub.MockPubSub)
		pubsub.Set(mockPubSub)
		mockSubscription := new(mock_pubsub.MockSubscription)
		mockPubSub.On("Subscribe", mock.Anything, mock.Anything).Return(mockSubscription, nil)
		ch := make(chan pubsub.Message)
		mockSubscription.On("Channel").Return((<-chan pubsub.Message)(ch))
		mockSubscription.On("Close").Return(nil)
		close(ch)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/?lang=es", nil)
		HandleSSE(w, req, loader)
		body := w.Body.String()
		events := strings.Split(body, "\n\n")
		clean := strings.TrimPrefix(events[0], "data:")
		clean = strings.TrimSpace(clean)
		var resp response.Response[any]
		err := json.Unmarshal([]byte(clean), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "", resp.Message)
		assert.Equal(t, "ok", resp.Data)
		clean = strings.TrimPrefix(events[1], "data:")
		clean = strings.TrimSpace(clean)
		err = json.Unmarshal([]byte(clean), &resp)
		assert.NoError(t, err)
		assert.Equal(t, i18n.Get(req, i18n.MsgInternalServerError), resp.Message)
	})

	t.Run("client disconnect", func(t *testing.T) {
		mockPubSub := new(mock_pubsub.MockPubSub)
		pubsub.Set(mockPubSub)
		mockSubscription := new(mock_pubsub.MockSubscription)
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // cancel immediately
		mockPubSub.On("Subscribe", mock.Anything, mock.Anything).Return(mockSubscription, nil)
		ch := make(chan pubsub.Message)
		mockSubscription.On("Channel").Return((<-chan pubsub.Message)(ch))
		mockSubscription.On("Close").Return(nil)
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/?lang=es", nil).WithContext(ctx)
		HandleSSE(w, req, loader)
		clean := strings.TrimPrefix(w.Body.String(), "data:")
		clean = strings.TrimSpace(clean)
		var resp response.Response[any]
		err := json.Unmarshal([]byte(clean), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "", resp.Message)
		assert.Equal(t, "ok", resp.Data)

	})

	t.Run("message received", func(t *testing.T) {
		mockPubSub := new(mock_pubsub.MockPubSub)
		pubsub.Set(mockPubSub)
		mockSubscription := new(mock_pubsub.MockSubscription)
		ch := make(chan pubsub.Message, 1)
		mockPubSub.On("Subscribe", mock.Anything, mock.Anything).Return(mockSubscription, nil)
		mockSubscription.On("Channel").Return((<-chan pubsub.Message)(ch))
		mockSubscription.On("Close").Return(nil)
		ctx, cancel := context.WithCancel(context.Background())
		req := httptest.NewRequest("GET", "/?lang=es", nil).WithContext(ctx)
		w := httptest.NewRecorder()
		go func() {
			time.Sleep(100 * time.Millisecond)
			ch <- pubsub.Message{Payload: "event"}
			cancel()
		}()
		HandleSSE(w, req, loader)

		body := w.Body.String()
		events := strings.Split(body, "\n\n")
		clean := strings.TrimPrefix(events[0], "data:")
		clean = strings.TrimSpace(clean)
		var resp response.Response[any]
		err := json.Unmarshal([]byte(clean), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "", resp.Message)
		assert.Equal(t, "ok", resp.Data)
		clean = strings.TrimPrefix(events[1], "data:")
		clean = strings.TrimSpace(clean)
		err = json.Unmarshal([]byte(clean), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "", resp.Message)
		assert.Equal(t, "ok", resp.Data)
	})

}

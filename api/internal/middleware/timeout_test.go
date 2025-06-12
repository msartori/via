package middleware

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"via/internal/log"
	mock_log "via/internal/log/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTimeoutMiddleware_DeadlineExceeded(t *testing.T) {
	mockLogger := new(mock_log.MockLogger)
	log.Set(mockLogger)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond) // longer than timeout
		io.WriteString(w, "should not write")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// Expected log calls
	mockLogger.On("WithLogFieldsInRequest", mock.Anything, []any{"status", http.StatusGatewayTimeout}).
		Return(req).Once()
	mockLogger.On("Error", mock.Anything, context.DeadlineExceeded, []any{"msg", "request timeout"}).
		Return().Once()

	timeoutMiddleware := Timeout(10 * time.Millisecond)
	timeoutMiddleware(handler).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusGatewayTimeout, rec.Code)
	mockLogger.AssertExpectations(t)
}

func TestTimeoutMiddleware_Canceled(t *testing.T) {
	mockLogger := new(mock_log.MockLogger)
	log.Set(mockLogger)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate external cancellation
		<-r.Context().Done()
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// Create a parent context and cancel it
	ctx, cancel := context.WithCancel(req.Context())
	cancel()
	req = req.WithContext(ctx)

	mockLogger.On("WithLogFieldsInRequest", mock.Anything, []any{"status", HttpStatusClientCloseRequest}).
		Return(req).Once()
	mockLogger.On("Error", mock.Anything, context.Canceled, []any{"msg", "request canceled by client"}).
		Return().Once()

	timeoutMiddleware := Timeout(100 * time.Millisecond)
	timeoutMiddleware(handler).ServeHTTP(rec, req)

	assert.Equal(t, HttpStatusClientCloseRequest, rec.Code)
	mockLogger.AssertExpectations(t)
}

func TestTimeoutMiddleware_Success(t *testing.T) {
	mockLogger := new(mock_log.MockLogger)
	log.Set(mockLogger)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Test", "ok")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("done"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	mockLogger.On("WithLogFieldsInRequest", mock.Anything, []any{"status", http.StatusAccepted}).
		Return(req).Once()

	timeoutMiddleware := Timeout(100 * time.Millisecond)
	timeoutMiddleware(handler).ServeHTTP(rec, req)

	res := rec.Result()
	body, _ := io.ReadAll(res.Body)

	assert.Equal(t, http.StatusAccepted, res.StatusCode)
	assert.Equal(t, "ok", res.Header.Get("X-Test"))
	assert.Equal(t, "done", string(body))

	mockLogger.AssertExpectations(t)
}

// Custom context that returns unexpected error
type dummyCtx struct {
	r *http.Request
}

func (d *dummyCtx) Deadline() (time.Time, bool) { return time.Time{}, false }
func (d *dummyCtx) Done() <-chan struct{}       { ch := make(chan struct{}); close(ch); return ch }
func (d *dummyCtx) Err() error                  { return errors.New("unexpected error") }
func (d *dummyCtx) Value(key any) any           { return d.r.Context().Value(key) }

func TestTimeoutMiddleware_DefaultErrorCase(t *testing.T) {
	// Arrange
	mockLogger := new(mock_log.MockLogger)
	log.Set(mockLogger) // You must have this method available in your logger

	mockLogger.On("WithLogFieldsInRequest", mock.Anything, []any{"status", http.StatusServiceUnavailable}).Return(&http.Request{})
	mockLogger.On("Error", mock.Anything, mock.Anything, []any{"msg", "unexpected context error"})

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("should be discarded"))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	// Override context with one that is "done" and returns an unexpected error
	dummy := &dummyCtx{req}
	req = req.WithContext(dummy)

	w := httptest.NewRecorder()
	timeoutMiddleware := Timeout(5 * time.Second)
	timeoutMiddleware(handler).ServeHTTP(w, req)

	// Assert
	require.Equal(t, http.StatusServiceUnavailable, w.Code)
	mockLogger.AssertExpectations(t)
}

func TestTimeoutWriter_WriteAfterTimeout(t *testing.T) {
	rec := httptest.NewRecorder()
	tw := &timeoutWriter{ResponseWriter: rec, header: http.Header{}}
	tw.markTimedOut()
	n, err := tw.Write([]byte("late"))
	assert.Equal(t, 0, n)
	assert.Equal(t, "write after timeout", err.Error())
}

func TestTimeoutWriter_Copy_SkipsWhenTimedOut(t *testing.T) {
	rr := httptest.NewRecorder()
	tw := &timeoutWriter{
		ResponseWriter: rr,
		header:         make(http.Header),
		timedOut:       true, // <-- Simulate timeout
		written:        true, // written true so only timedOut triggers
	}
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	tw.copy(req)

	if rr.Body.Len() != 0 {
		t.Error("expected no response written when timed out")
	}
}

func TestTimeoutWriter_Copy_SkipsWhenNotWritten(t *testing.T) {
	rr := httptest.NewRecorder()
	tw := &timeoutWriter{
		ResponseWriter: rr,
		header:         make(http.Header),
		timedOut:       false,
		written:        false, // <-- Simulate nothing was written
	}
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	tw.copy(req)

	if rr.Body.Len() != 0 {
		t.Error("expected no response written when nothing was written")
	}
}

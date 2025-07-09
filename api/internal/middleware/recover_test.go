package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"via/internal/log"
	mock_log "via/internal/log/mock"
	"via/internal/middleware"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRecoverMiddleware(t *testing.T) {
	// Mock the logger
	mockLog := new(mock_log.MockLogger)
	mockLog.On("Error", mock.Anything, mock.Anything, mock.Anything).Return()
	log.Set(mockLog)

	// Handler with panic
	handler := middleware.Recover(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("something went wrong")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	// Request Execution
	handler.ServeHTTP(rec, req)

	// Response verification
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Equal(t, "internal server error\n", rec.Body.String())

	// Verify that the logger was called with the expected message
	mockLog.AssertCalled(t, "Error", mock.Anything, mock.Anything, mock.Anything)

	// Check the arguments passed to the logger
	args := mockLog.Calls[0].Arguments
	errArg, ok := args.Get(1).(error)
	assert.True(t, ok)
	assert.Contains(t, errArg.Error(), "something went wrong")
	keyvals := args.Get(2).([]any)
	assert.Contains(t, keyvals, "msg")
	assert.Contains(t, keyvals, "recovering from panic")
}

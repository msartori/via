package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"via/internal/middleware"
	"via/internal/testutil"

	"github.com/stretchr/testify/assert"
)

func TestLogHandlerExecution(t *testing.T) {
	// Inject no-op logger to avoid real logging during test
	testutil.InjectNoOpLogger()

	var handlerCalled bool

	// Create a dummy handler to wrap
	handler := func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusTeapot)
	}

	// Wrap with middleware
	wrapped := middleware.LogHandlerExecution("TestHandler", handler)

	// Perform the test request
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	wrapped(rr, req)

	// Assert handler was called and response code is from handler
	assert.True(t, handlerCalled, "handler should have been called")
	assert.Equal(t, http.StatusTeapot, rr.Code, "unexpected status code")
}

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"via/internal/log"
	mock_log "via/internal/log/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRequestMiddleware_GeneratesRequestID(t *testing.T) {
	// Setup
	mockLogger := new(mock_log.MockLogger)
	log.Set(mockLogger) // Asume que tienes una función SetMock para sustituir el singleton
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	mockLogger.On("WithLogFieldsInRequest", mock.Anything, mock.Anything).Return(req)
	mockLogger.On("Info", mock.Anything, mock.Anything).Twice()

	var receivedReqID string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedReqID = r.Context().Value(RequestIDKey).(string)
		assert.NotEmpty(t, receivedReqID)
		w.WriteHeader(http.StatusOK)
	})

	rec := httptest.NewRecorder()

	// Execute
	mw := Request(handler)
	mw.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotEmpty(t, rec.Header().Get(RequestIDHeader))
	assert.Equal(t, rec.Header().Get(RequestIDHeader), receivedReqID)
	mockLogger.AssertExpectations(t)
}

func TestRequestMiddleware_UsesExistingRequestID(t *testing.T) {
	// Setup
	mockLogger := new(mock_log.MockLogger)
	log.Set(mockLogger)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	mockLogger.On("WithLogFieldsInRequest", mock.Anything, mock.Anything).Return(req)
	mockLogger.On("Info", mock.Anything, mock.Anything).Twice()

	existingID := "test-id-123"
	var receivedReqID string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedReqID = r.Context().Value(RequestIDKey).(string)
		assert.Equal(t, existingID, receivedReqID)
		w.WriteHeader(http.StatusOK)
	})

	req.Header.Set(RequestIDHeader, existingID)
	rec := httptest.NewRecorder()

	// Execute
	mw := Request(handler)
	mw.ServeHTTP(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, existingID, rec.Header().Get(RequestIDHeader))
	assert.Equal(t, existingID, receivedReqID)
	mockLogger.AssertExpectations(t)
}

package middleware_test

import (
	"fmt"
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
	mockLog.On("Info", mock.Anything, mock.Anything).Return()
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

	// Verificar que se llamó Info con algún contexto y algún par de argumentos
	mockLog.AssertCalled(t, "Info", mock.Anything, mock.Anything)

	// También podés verificar parcialmente los argumentos:
	// Extraer los argumentos reales para validarlos
	args := mockLog.Calls[0].Arguments
	pairs, ok := args.Get(1).([]any)
	assert.True(t, ok, "expected second arg to be []any")

	// Validar que contiene el mensaje y el error
	str := fmt.Sprint(pairs)
	assert.Contains(t, str, "recovering from panic")
	assert.Contains(t, str, "something went wrong")
}

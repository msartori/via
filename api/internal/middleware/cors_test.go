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

// dummy function that will be invoked by the middleware
func testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	w.Write([]byte("OK"))
}

func TestCORSMiddleware(t *testing.T) {
	mockLog := new(mock_log.MockLogger)

	// logger mock setup
	log.Set(mockLog)

	// middleware configuration
	cfg := CORSCfg{
		Origins: "http://localhost",
		Methods: "GET,POST,OPTIONS",
		Headers: "Content-Type,Authorization",
	}

	// handler setup with CORS middleware
	handler := CORS(cfg)

	//handler := middleware.CORS(http.HandlerFunc(testHandler), cfg)

	t.Run("regular request passes with correct headers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/resource", nil)
		rec := httptest.NewRecorder()

		mockLog.On("WithLogFieldsInRequest", req, mock.Anything).Return(req)
		mockLog.On("Info", req.Context(), mock.Anything)

		handler(http.HandlerFunc(testHandler)).ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusTeapot, res.StatusCode)
		assert.Equal(t, cfg.Origins, res.Header.Get("Access-Control-Allow-Origin"))
		assert.Equal(t, cfg.Methods, res.Header.Get("Access-Control-Allow-Methods"))
		assert.Equal(t, cfg.Headers, res.Header.Get("Access-Control-Allow-Headers"))

		mockLog.AssertCalled(t, "WithLogFieldsInRequest", req, mock.Anything)
		mockLog.AssertCalled(t, "Info", req.Context(), mock.Anything)
	})

	t.Run("OPTIONS request returns 200 OK with headers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/resource", nil)
		rec := httptest.NewRecorder()

		mockLog.On("WithLogFieldsInRequest", req, mock.Anything).Return(req)
		mockLog.On("Info", req.Context(), mock.Anything)

		handler(http.HandlerFunc(testHandler)).ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, cfg.Origins, res.Header.Get("Access-Control-Allow-Origin"))
		assert.Equal(t, cfg.Methods, res.Header.Get("Access-Control-Allow-Methods"))
		assert.Equal(t, cfg.Headers, res.Header.Get("Access-Control-Allow-Headers"))
	})

	t.Run("handles empty config gracefully", func(t *testing.T) {
		emptyCfg := CORSCfg{}
		handlerEmpty := CORS(emptyCfg)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()

		mockLog.On("WithLogFieldsInRequest", req, mock.Anything).Return(req)
		mockLog.On("Info", req.Context(), mock.Anything)

		handlerEmpty(http.HandlerFunc(testHandler)).ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, "", res.Header.Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "", res.Header.Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "", res.Header.Get("Access-Control-Allow-Headers"))
	})
}

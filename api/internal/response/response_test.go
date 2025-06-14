package response

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"via/internal/global"

	"github.com/stretchr/testify/assert"
)

func TestWriteJSON_WithRequestID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := context.WithValue(req.Context(), global.RequestIDKey, "1234")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	res := Response{
		Data:    map[string]string{"foo": "bar"},
		Message: "ok",
	}
	WriteJSON(rr, req, res, http.StatusOK)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var got Response
	body, _ := io.ReadAll(rr.Body)
	err := json.Unmarshal(body, &got)
	assert.NoError(t, err)
	assert.Equal(t, "ok", got.Message)
	assert.Equal(t, "1234", got.RequestID)
	assert.Equal(t, map[string]interface{}{"foo": "bar"}, got.Data)
}

func TestWriteJSON_NoRequestID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	res := Response{Message: "no ID"}
	WriteJSON(rr, req, res, http.StatusAccepted)

	assert.Equal(t, http.StatusAccepted, rr.Code)
	var got Response
	json.NewDecoder(rr.Body).Decode(&got)
	assert.Equal(t, "no ID", got.Message)
	assert.Empty(t, got.RequestID)
}

func TestWriteJSON_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
	rr := httptest.NewRecorder()

	res := Response{Message: "should not be written"}
	WriteJSON(rr, req, res, http.StatusOK)

	assert.Equal(t, "", rr.Body.String())
	assert.Equal(t, "", rr.Header().Get("Content-Type"))
}

func TestWriteJSONError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := context.WithValue(req.Context(), global.RequestIDKey, "err-id")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	res := Response{Message: "error"}
	WriteJSONError(rr, req, res, http.StatusBadRequest)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var got Response
	json.NewDecoder(rr.Body).Decode(&got)
	assert.Equal(t, "error", got.Message)
	assert.Equal(t, "err-id", got.RequestID)
}

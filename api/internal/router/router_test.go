// internal/router/router_test.go
package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"via/internal/config"
	"via/internal/log"
	mock_log "via/internal/log/mock"
	"via/internal/model"
	guide_provider "via/internal/provider/guide"
	mock_guide_provider "via/internal/provider/guide/mock"
	"via/internal/router"

	"github.com/stretchr/testify/mock"
)

func TestRouter_New(t *testing.T) {
	mockLog := new(mock_log.MockNoOpLogger)
	// logger mock setup
	log.Set(mockLog)
	cfg := config.Config{
		Application: config.Application{
			RequestTimeout: 5,
		},
		CORS: config.CORS{
			Origins: "*",
		},
		Bussiness: config.Bussiness{
			ViaBranch:       "001",
			PendingStatus:   "PENDING",
			DeliveredStatus: "DELIVERED",
			WithdrawStatus:  "WITHDRAW",
		},
	}

	h := router.New(cfg)
	mockProvider := new(mock_guide_provider.MockGuideProvider)
	mockProvider.On("GetGuide", mock.Anything, "123456789012").Return(model.Guide{ID: "123456789012"}, nil)
	guide_provider.Set(mockProvider)
	r := httptest.NewRequest(http.MethodGet, "/guide/123456789012", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	if w.Code == 400 || w.Code == 500 {
		t.Errorf("expected router to handle request gracefully, got status %d", w.Code)
	}
}

// internal/router/router_test.go
package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"via/internal/config"
	"via/internal/log"
	mock_log "via/internal/log/mock"
	"via/internal/middleware"
	"via/internal/model"
	guide_provider "via/internal/provider/guide"
	mock_guide_provider "via/internal/provider/guide/mock"
	via_guide_provider "via/internal/provider/via/guide"
	mock_via_guide_provider "via/internal/provider/via/guide/mock"
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
		CORS: middleware.CORSCfg{
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
	mockViaGuideProvider := new(mock_via_guide_provider.MockViaGuideProvider)
	mockViaGuideProvider.On("GetGuide", mock.Anything, "123456789012").Return(model.ViaGuide{ID: "123456789012"}, nil)
	via_guide_provider.Set(mockViaGuideProvider)

	mockGuideProvider := new(mock_guide_provider.MockGuideProvider)
	guide_provider.Set(mockGuideProvider)
	r := httptest.NewRequest(http.MethodGet, "/guide/123456789012", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	if w.Code == 400 || w.Code == 500 {
		t.Errorf("expected router to handle request gracefully, got status %d", w.Code)
	}
}

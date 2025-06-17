package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"via/internal/config"
	"via/internal/i18n"
	"via/internal/log"
	mock_log "via/internal/log/mock"
	"via/internal/model"
	via_guide_provider "via/internal/provider/via/guide"
	mock_via_guide_provider "via/internal/provider/via/guide/mock"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Helper para crear requests con parámetros
func newRequestWithID(id string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/guides/"+id, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", id)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	req.Header.Set("Accept-Language", "es")
	return req
}

func TestGetGuideHandler(t *testing.T) {
	mockLog := new(mock_log.MockNoOpLogger)
	// logger mock setup
	log.Set(mockLog)
	mockProvider := new(mock_via_guide_provider.MockViaGuideProvider)
	via_guide_provider.Set(mockProvider)

	biz := config.Bussiness{
		ViaBranch:       "001",
		PendingStatus:   "IN_TRANSIT,WAITING",
		DeliveredStatus: "DELIVERED",
		WithdrawStatus:  "READY",
	}

	tests := []struct {
		name                 string
		guideID              string
		mockReturn           model.ViaGuide
		mockError            error
		expectedStatus       int
		expectedMsg          string
		setupHeaders         func(*http.Request)
		providerCallExpected bool
	}{
		{
			name:           "missing id",
			guideID:        "",
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    i18n.Get(newRequestWithID(""), i18n.MsgGuideRequired),
		},
		{
			name:           "invalid id",
			guideID:        "abc123",
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    i18n.Get(newRequestWithID("abc123"), i18n.MsgGuideInvalid),
		},
		{
			name:                 "provider error",
			guideID:              "123456789012",
			mockError:            errors.New("fail"),
			expectedStatus:       http.StatusInternalServerError,
			expectedMsg:          i18n.Get(newRequestWithID("123456789012"), i18n.MsgInternalServerError),
			providerCallExpected: true,
		},
		{
			name:                 "guide not found (empty ID)",
			guideID:              "123456789013",
			mockReturn:           model.ViaGuide{},
			expectedStatus:       http.StatusNotFound,
			expectedMsg:          i18n.Get(newRequestWithID("123456789013"), i18n.MsgGuideNotFound),
			providerCallExpected: true,
		},
		{
			name:    "guide from other branch",
			guideID: "123456789014",
			mockReturn: model.ViaGuide{
				ID: "123456789014",
				Destination: model.ViaDestination{
					ID:          "999",
					Description: "Sucursal X",
				},
			},
			expectedStatus:       http.StatusOK,
			expectedMsg:          i18n.Get(newRequestWithID("123456789014"), i18n.MsgOtherBranch, "Sucursal X"),
			providerCallExpected: true,
		},
		{
			name:    "guide in transit",
			guideID: "123456789015",
			mockReturn: model.ViaGuide{
				ID: "123456789015",
				Destination: model.ViaDestination{
					ID: "001",
				},
				Status: "IN_TRANSIT",
			},
			expectedStatus:       http.StatusOK,
			expectedMsg:          i18n.Get(newRequestWithID("123456789015"), i18n.MsgInTransit),
			providerCallExpected: true,
		},
		{
			name:    "guide delivered",
			guideID: "123456789016",
			mockReturn: model.ViaGuide{
				ID: "123456789016",
				Destination: model.ViaDestination{
					ID: "001",
				},
				Status: "DELIVERED",
			},
			expectedStatus:       http.StatusOK,
			expectedMsg:          i18n.Get(newRequestWithID("123456789016"), i18n.MsgDelivered),
			providerCallExpected: true,
		},
		{
			name:    "guide withdraw",
			guideID: "123456789017",
			mockReturn: model.ViaGuide{
				ID: "123456789017",
				Destination: model.ViaDestination{
					ID: "001",
				},
				Status: "READY",
			},
			expectedStatus:       http.StatusOK,
			expectedMsg:          i18n.Get(newRequestWithID("123456789017"), i18n.MsgWithdrawAvailable),
			providerCallExpected: true,
		},
		{
			name:    "guide not available",
			guideID: "123456789018",
			mockReturn: model.ViaGuide{
				ID: "123456789018",
				Destination: model.ViaDestination{
					ID: "001",
				},
				Status: "UNKNOWN",
			},
			expectedStatus:       http.StatusOK,
			expectedMsg:          i18n.Get(newRequestWithID("123456789018"), i18n.MsgNotAvailable),
			providerCallExpected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.providerCallExpected {
				mockProvider.On("GetGuide", mock.Anything, tt.guideID).Return(tt.mockReturn, tt.mockError)
			}

			req := newRequestWithID(tt.guideID)
			rr := httptest.NewRecorder()
			GetGuide(biz).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			body := rr.Body.String()
			assert.Contains(t, body, tt.expectedMsg)
		})
	}
}

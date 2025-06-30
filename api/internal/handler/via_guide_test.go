package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	biz_config "via/internal/biz/config"
	biz_guide_status "via/internal/biz/guide/status"
	"via/internal/global"
	"via/internal/i18n"
	"via/internal/log"
	mock_log "via/internal/log/mock"
	"via/internal/model"
	guide_provider "via/internal/provider/guide"
	mock_guide_provider "via/internal/provider/guide/mock"
	via_guide_provider "via/internal/provider/via/guide"
	mock_via_guide_provider "via/internal/provider/via/guide/mock"

	"via/internal/response"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Helper to create requests with parmeters
func newRequestWithID(id string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "//guide-to-withdraw/"+id, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("viaGuideId", id)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	ctx := context.WithValue(req.Context(), global.RequestIDKey, id+"0001")
	req = req.WithContext(ctx)
	req.Header.Set("Accept-Language", "es")
	return req
}

func TestGetGuideHandler(t *testing.T) {
	mockLog := new(mock_log.MockNoOpLogger)
	// logger mock setup
	log.Set(mockLog)
	mockViaGuideProvider := new(mock_via_guide_provider.MockViaGuideProvider)
	via_guide_provider.Set(mockViaGuideProvider)
	mockGuideProvider := new(mock_guide_provider.MockGuideProvider)
	guide_provider.Set(mockGuideProvider)

	biz := biz_config.Bussiness{
		ViaBranch:       "001",
		PendingStatus:   "IN_TRANSIT,WAITING",
		DeliveredStatus: "DELIVERED",
		WithdrawStatus:  "READY",
		HomeDelivery:    "CD06",
	}

	tests := []struct {
		name                         string
		guideID                      string
		guideViaProviderCallExpected bool
		guideViaMockReturn           model.ViaGuide
		guideViaMockError            error
		guideProviderCallExpected    bool
		guideMockReturn              model.Guide
		guideMockError               error
		expectedStatus               int
		expectedResponse             response.Response[GetGuideToWithdrawOutput]
	}{
		{
			name:           "missing id",
			guideID:        "",
			expectedStatus: http.StatusBadRequest,
			expectedResponse: response.Response[GetGuideToWithdrawOutput]{RequestID: "0001",
				Message: i18n.Get(newRequestWithID(""), i18n.MsgGuideRequired)},
		},

		{
			name:           "invalid id",
			guideID:        "abc123",
			expectedStatus: http.StatusBadRequest,
			expectedResponse: response.Response[GetGuideToWithdrawOutput]{RequestID: "abc1230001",
				Message: i18n.Get(newRequestWithID("abc123"), i18n.MsgGuideInvalid)},
		},

		{
			name:              "provider error",
			guideID:           "123456789012",
			guideViaMockError: errors.New("fail"),
			expectedStatus:    http.StatusInternalServerError,
			expectedResponse: response.Response[GetGuideToWithdrawOutput]{RequestID: "1234567890120001",
				Message: i18n.Get(newRequestWithID("123456789012"), i18n.MsgInternalServerError)},
			guideViaProviderCallExpected: true,
		},

		{
			name:               "guide not found (empty ID)",
			guideID:            "123456789013",
			guideViaMockReturn: model.ViaGuide{},
			expectedStatus:     http.StatusOK,
			expectedResponse: response.Response[GetGuideToWithdrawOutput]{RequestID: "1234567890130001",
				Data: GetGuideToWithdrawOutput{EnabledToWithdraw: false,
					WithdrawMessage: getWithDrawMessage(newRequestWithID("123456789013"), notFound, "123456789013")}},
			guideViaProviderCallExpected: true,
		},
		{
			name:    "guide from other branch",
			guideID: "123456789014",
			guideViaMockReturn: model.ViaGuide{
				ID: "123456789014",
				Destination: model.ViaDestination{
					ID:          "999",
					Description: "Sucursal X",
				},
			},
			expectedStatus: http.StatusOK,
			expectedResponse: response.Response[GetGuideToWithdrawOutput]{RequestID: "1234567890140001",
				Data: GetGuideToWithdrawOutput{EnabledToWithdraw: false,
					WithdrawMessage: getWithDrawMessage(newRequestWithID("123456789014"), wrongBranch, "123456789014", "Sucursal X")}},
			guideViaProviderCallExpected: true,
		},
		{
			name:    "guide in transit",
			guideID: "123456789015",
			guideViaMockReturn: model.ViaGuide{
				ID: "123456789015",
				Destination: model.ViaDestination{
					ID: "001",
				},
				Status: "IN_TRANSIT",
			},
			expectedStatus: http.StatusOK,
			expectedResponse: response.Response[GetGuideToWithdrawOutput]{RequestID: "1234567890150001",
				Data: GetGuideToWithdrawOutput{EnabledToWithdraw: false,
					WithdrawMessage: getWithDrawMessage(newRequestWithID("123456789015"), pending, "123456789015")}},
			guideViaProviderCallExpected: true,
		},
		{
			name:    "guide delivered",
			guideID: "123456789016",
			guideViaMockReturn: model.ViaGuide{
				ID: "123456789016",
				Destination: model.ViaDestination{
					ID: "001",
				},
				Status: "DELIVERED",
			},
			expectedStatus: http.StatusOK,
			expectedResponse: response.Response[GetGuideToWithdrawOutput]{RequestID: "1234567890160001",
				Data: GetGuideToWithdrawOutput{EnabledToWithdraw: false,
					WithdrawMessage: getWithDrawMessage(newRequestWithID("123456789016"), delivered, "123456789016")}},
			guideViaProviderCallExpected: true,
		},
		{
			name:                         "guide withdraw",
			guideID:                      "123456789017",
			guideViaProviderCallExpected: true,
			guideViaMockReturn: model.ViaGuide{
				ID: "123456789017",
				Destination: model.ViaDestination{
					ID: "001",
				},
				Status: "READY",
			},
			guideProviderCallExpected: true,
			guideMockReturn:           model.Guide{},
			expectedStatus:            http.StatusOK,
			expectedResponse: response.Response[GetGuideToWithdrawOutput]{RequestID: "1234567890170001",
				Data: GetGuideToWithdrawOutput{EnabledToWithdraw: true,
					WithdrawMessage: getWithDrawMessage(newRequestWithID("123456789017"), enabledToWithdraw, "123456789017")}},
		},
		{
			name:                         "guide not available",
			guideID:                      "123456789018",
			guideViaProviderCallExpected: true,
			guideViaMockReturn: model.ViaGuide{
				ID: "123456789018",
				Destination: model.ViaDestination{
					ID: "001",
				},
				Status: "UNKNOWN",
			},
			guideProviderCallExpected: false,
			expectedStatus:            http.StatusOK,
			expectedResponse: response.Response[GetGuideToWithdrawOutput]{RequestID: "1234567890180001",
				Data: GetGuideToWithdrawOutput{EnabledToWithdraw: false,
					WithdrawMessage: getWithDrawMessage(newRequestWithID("123456789018"), notAvailable, "123456789018")}},
		},
		{
			name:                         "guide with delivered status",
			guideID:                      "123456789020",
			guideViaProviderCallExpected: true,
			guideViaMockReturn: model.ViaGuide{
				ID: "123456789020",
				Destination: model.ViaDestination{
					ID: "001",
				},
				Status: "READY",
			},
			guideProviderCallExpected: true,
			guideMockReturn: model.Guide{
				ID:     1,
				Status: biz_guide_status.DELIVERED,
			},
			expectedStatus: http.StatusOK,
			expectedResponse: response.Response[GetGuideToWithdrawOutput]{
				RequestID: "1234567890200001",
				Data: GetGuideToWithdrawOutput{
					EnabledToWithdraw: false,
					WithdrawMessage:   getWithDrawMessage(newRequestWithID("123456789020"), delivered, "123456789020"),
				},
			},
		},
		{
			name:                         "guide in process",
			guideID:                      "123456789021",
			guideViaProviderCallExpected: true,
			guideViaMockReturn: model.ViaGuide{
				ID: "123456789021",
				Destination: model.ViaDestination{
					ID: "001",
				},
				Status: "READY",
			},
			guideProviderCallExpected: true,
			guideMockReturn: model.Guide{
				ID:     2,
				Status: biz_guide_status.INITIAL,
			},
			expectedStatus: http.StatusOK,
			expectedResponse: response.Response[GetGuideToWithdrawOutput]{
				RequestID: "1234567890210001",
				Data: GetGuideToWithdrawOutput{
					EnabledToWithdraw: false,
					WithdrawMessage:   getWithDrawMessage(newRequestWithID("123456789021"), alreadyInProcess, "123456789021"),
				},
			},
		},
		{
			name:                         "guide enabled to withdraw",
			guideID:                      "123456789022",
			guideViaProviderCallExpected: true,
			guideViaMockReturn: model.ViaGuide{
				ID: "123456789022",
				Destination: model.ViaDestination{
					ID: "001",
				},
				Status: "READY",
			},
			guideProviderCallExpected: true,
			guideMockReturn: model.Guide{
				ID:     3,
				Status: biz_guide_status.ON_HOLD,
			},
			expectedStatus: http.StatusOK,
			expectedResponse: response.Response[GetGuideToWithdrawOutput]{
				RequestID: "1234567890220001",
				Data: GetGuideToWithdrawOutput{
					EnabledToWithdraw: true,
					WithdrawMessage:   getWithDrawMessage(newRequestWithID("123456789022"), enabledToWithdraw, "123456789022"),
				},
			},
		},
		{
			name:                         "local provider error",
			guideID:                      "123456789023",
			guideViaProviderCallExpected: true,
			guideViaMockReturn: model.ViaGuide{
				ID: "123456789023",
				Destination: model.ViaDestination{
					ID: "001",
				},
				Status: "READY",
			},
			guideProviderCallExpected: true,
			guideMockError:            errors.New("fail"),
			expectedStatus:            http.StatusInternalServerError,
			expectedResponse: response.Response[GetGuideToWithdrawOutput]{RequestID: "1234567890230001",
				Message: i18n.Get(newRequestWithID("123456789023"), i18n.MsgInternalServerError)},
		},
		{
			name:                         "guide home delivery",
			guideID:                      "123456789024",
			guideViaProviderCallExpected: true,
			guideViaMockReturn: model.ViaGuide{
				ID: "123456789024",
				Destination: model.ViaDestination{
					ID: "001",
				},
				Status: "READY",
				Route:  "CD06",
			},
			expectedStatus: http.StatusOK,
			expectedResponse: response.Response[GetGuideToWithdrawOutput]{
				RequestID: "1234567890240001",
				Data: GetGuideToWithdrawOutput{
					EnabledToWithdraw: false,
					WithdrawMessage:   getWithDrawMessage(newRequestWithID("123456789024"), homeDelivery, "123456789024"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.guideViaProviderCallExpected {
				mockViaGuideProvider.On("GetGuide", mock.Anything, tt.guideID).Return(tt.guideViaMockReturn, tt.guideViaMockError)
			}
			if tt.guideProviderCallExpected {
				mockGuideProvider.On("GetGuideByViaGuideId", mock.Anything, tt.guideID).Return(tt.guideMockReturn, tt.guideMockError)
			}

			req := newRequestWithID(tt.guideID)
			rr := httptest.NewRecorder()
			GetGuideToWithdraw(biz).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			res := response.Response[GetGuideToWithdrawOutput]{}
			if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
				t.Errorf("unable to decode response:%v", err)
			}
			assert.Equal(t, tt.expectedResponse, res)
		})
	}
}

func TestGetWithDrawMessage(t *testing.T) {
	tests := []struct {
		name           string
		acceptLang     string
		key            string
		args           []interface{}
		expectedOutput string
	}{
		{
			name:           "language and key with args",
			acceptLang:     "es",
			key:            notFound,
			args:           []interface{}{"ABC123"},
			expectedOutput: "El código de guía ingresado [ABC123] no se encuentra, por favor verifique que sea correcto.",
		},
		{
			name:           "language and key without args (no placeholder)",
			acceptLang:     "es",
			key:            alreadyInProcess,
			args:           nil,
			expectedOutput: "El código de guía ingresado [%s] es correcto, el envío esta en proceso de entrega, aguarde y será atendido.",
		},
		{
			name:           "language not set, fallback to es",
			acceptLang:     "",
			key:            delivered,
			args:           []interface{}{"XYZ789"},
			expectedOutput: "El código de guía ingresado [XYZ789] es correcto, pero el envío ya ha sido retirado.",
		},
		{
			name:           "key not found in language, returns key",
			acceptLang:     "es",
			key:            "non_existent_key",
			args:           nil,
			expectedOutput: "non_existent_key",
		},
		{
			name:           "lang not supported, fallback to key",
			acceptLang:     "fr",
			key:            notFound,
			args:           nil,
			expectedOutput: notFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.acceptLang != "" {
				req.Header.Set("Accept-Language", tt.acceptLang)
			}

			got := getWithDrawMessage(req, tt.key, tt.args...)
			assert.Equal(t, tt.expectedOutput, got)
		})
	}
}

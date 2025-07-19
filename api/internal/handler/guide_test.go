package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	biz_config "via/internal/biz/config"
	biz_guide_status "via/internal/biz/guide/status"
	biz_operator "via/internal/biz/operator"
	"via/internal/i18n"
	"via/internal/log"
	mock_log "via/internal/log/mock"
	"via/internal/middleware"
	"via/internal/model"
	guide_provider "via/internal/provider/guide"
	mock_guide_provider "via/internal/provider/guide/mock"
	via_guide_provider "via/internal/provider/via/guide"
	mock_via_guide_provider "via/internal/provider/via/guide/mock"
	"via/internal/pubsub"
	mock_pubsub "via/internal/pubsub/mock"
	"via/internal/response"
	"via/internal/testutil"
)

func assertJSONErrorResponse(t *testing.T, req *http.Request, w *httptest.ResponseRecorder, expectedStatus int, expectedMsg string) {
	t.Helper()
	var resp response.Response[any]
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, i18n.Get(req, expectedMsg), resp.Message)
	assert.Equal(t, expectedStatus, w.Code)
}

func newGuideRequest(method, path, guideId string, body []byte, operatorId any) *http.Request {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	rctx := chi.NewRouteContext()
	if guideId != "" {
		rctx.URLParams.Add("guideId", guideId)
	}
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	if operatorId != nil {
		req = req.WithContext(context.WithValue(req.Context(), middleware.OperatorIDKey, operatorId))
	}
	return req
}

func newOperatorRequest(method, path string, body []byte, operatorID any) *http.Request {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	if operatorID != nil {
		req = req.WithContext(context.WithValue(req.Context(), middleware.OperatorIDKey, operatorID))
	}
	return req
}

func TestCreateGuideToWithdraw(t *testing.T) {
	biz := biz_config.BussinessCfg{ViaBranch: "123", WithdrawStatus: "CCR", DeliveredStatus: "ENT", PendingStatus: "ASP,CPO,ORI,PTE"}
	mockVia := new(mock_via_guide_provider.MockViaGuideProvider)
	via_guide_provider.Set(mockVia)
	mockGuide := new(mock_guide_provider.MockGuideProvider)
	guide_provider.Set(mockGuide)
	log.Set(new(mock_log.MockNoOpLogger))
	mockPubSub := new(mock_pubsub.MockPubSub)
	pubsub.Set(mockPubSub)

	resetMocks := func() {
		mockVia.ExpectedCalls = nil
		mockGuide.ExpectedCalls = nil
	}

	makeRequest := func(viaGuideID string) *httptest.ResponseRecorder {
		body := bytes.NewBufferString(fmt.Sprintf(`{"viaGuideId":"%s"}`, viaGuideID))
		req := httptest.NewRequest(http.MethodPost, "/", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		CreateGuideToWidthdraw(biz).ServeHTTP(rec, req)
		return rec
	}

	t.Run("invalid JSON body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("{invalid"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		CreateGuideToWidthdraw(biz).ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid ViaGuideId", func(t *testing.T) {
		rec := makeRequest("")
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("error fetching via guide", func(t *testing.T) {
		resetMocks()
		mockVia.On("GetGuide", mock.Anything, "123456789012").Return(model.ViaGuide{}, errors.New("error"))
		rec := makeRequest("123456789012")
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("via guide not found", func(t *testing.T) {
		resetMocks()
		mockVia.On("GetGuide", mock.Anything, "123456789012").Return(model.ViaGuide{}, nil)
		rec := makeRequest("123456789012")
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("invalid via guide to withdraw", func(t *testing.T) {
		resetMocks()
		mockVia.On("GetGuide", mock.Anything, "123456789012").
			Return(model.ViaGuide{ID: "123456789012", Status: biz.DeliveredStatus}, nil)
		rec := makeRequest("123456789012")
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("invalid via guide to withdraw 2", func(t *testing.T) {
		resetMocks()
		mockVia.On("GetGuide", mock.Anything, "123456789012").
			Return(model.ViaGuide{ID: "123456789012", Status: "CPO", Destination: model.ViaDestination{ID: "123"}}, nil)
		rec := makeRequest("123456789012")
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("failed to fetch guide", func(t *testing.T) {
		resetMocks()
		mockVia.On("GetGuide", mock.Anything, "123456789012").Return(
			model.ViaGuide{ID: "123456789012", Status: "VALID", Destination: model.ViaDestination{ID: biz.ViaBranch}}, nil)
		mockGuide.On("GetGuideByViaGuideId", mock.Anything, "123456789012").
			Return(model.Guide{}, errors.New("error fetching guide"))
		rec := makeRequest("123456789012")
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("existing guide is not valid to create a new one", func(t *testing.T) {
		resetMocks()
		mockVia.On("GetGuide", mock.Anything, "123456789012").Return(
			model.ViaGuide{ID: "123456789012", Status: "VALID", Destination: model.ViaDestination{ID: biz.ViaBranch}}, nil)
		mockGuide.On("GetGuideByViaGuideId", mock.Anything, "123456789012").
			Return(model.Guide{ID: 10, Status: biz_guide_status.PAID}, nil)
		rec := makeRequest("123456789012")
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("re-init existing guide", func(t *testing.T) {
		resetMocks()
		mockVia.On("GetGuide", mock.Anything, "123456789012").Return(
			model.ViaGuide{ID: "123456789012", Status: "VALID", Destination: model.ViaDestination{ID: biz.ViaBranch}}, nil)
		mockGuide.On("GetGuideByViaGuideId", mock.Anything, "123456789012").
			Return(model.Guide{ID: 10, Status: biz_guide_status.ON_HOLD}, nil)
		mockGuide.On("UpdateGuide", mock.Anything, mock.Anything).Return(nil)
		mockPubSub.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		rec := makeRequest("123456789012")
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("invalid status to create guide", func(t *testing.T) {
		resetMocks()
		mockVia.On("GetGuide", mock.Anything, "123456789012").Return(
			model.ViaGuide{ID: "123456789012", Status: biz.WithdrawStatus}, nil)
		mockGuide.On("GetGuideByViaGuideId", mock.Anything, "123456789012").
			Return(model.Guide{ID: 10, Status: biz_guide_status.INITIAL}, nil)
		rec := makeRequest("123456789012")
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("success create new guide", func(t *testing.T) {
		resetMocks()
		viaGuideId := "123456789013"
		mockVia.On("GetGuide", mock.Anything, viaGuideId).Return(
			model.ViaGuide{ID: viaGuideId, Status: biz.WithdrawStatus, Destination: model.ViaDestination{ID: biz.ViaBranch}}, nil)
		mockGuide.On("GetGuideByViaGuideId", mock.Anything, viaGuideId).Return(model.Guide{}, nil)
		mockGuide.On("CreateGuide", mock.Anything, mock.Anything).Return(200, nil)
		mockPubSub.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		rec := makeRequest(viaGuideId)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("error creating guide", func(t *testing.T) {
		resetMocks()
		viaGuideId := "123456789014"
		mockVia.On("GetGuide", mock.Anything, viaGuideId).Return(
			model.ViaGuide{ID: viaGuideId, Status: biz.WithdrawStatus, Destination: model.ViaDestination{ID: biz.ViaBranch}}, nil)
		mockGuide.On("GetGuideByViaGuideId", mock.Anything, viaGuideId).Return(model.Guide{}, nil)
		mockGuide.On("CreateGuide", mock.Anything, mock.Anything).Return(0, errors.New("error"))
		rec := makeRequest(viaGuideId)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestGetOperatorGuide(t *testing.T) {
	testutil.InjectNoOpLogger()

	t.Run("success", func(t *testing.T) {
		mockGuideProvider := new(mock_guide_provider.MockGuideProvider)
		guide_provider.Set(mockGuideProvider)
		mockPubSub := new(mock_pubsub.MockPubSub)
		pubsub.Set(mockPubSub)
		guides := []model.Guide{
			{ID: 2, Recipient: "John", Status: "INIT", Operator: model.Operator{ID: 42}, ViaGuideID: "V123", Payment: "PREPAID"},
		}

		req := newOperatorRequest(http.MethodGet, "/guide/operator", nil, 42)

		operatorGuide := model.OperatorGuide{
			GuideId:    guides[0].ID,
			Recipient:  guides[0].Recipient,
			Status:     biz_guide_status.GetStatusDescription(response.GetLanguage(req), guides[0].Status),
			Operator:   guides[0].Operator,
			Selectable: guides[0].Operator.ID == biz_operator.OPERATOR_SYSTEM || guides[0].Operator.ID == 42,
			ViaGuideId: guides[0].ViaGuideID,
			Payment:    biz_config.GetPaymentDescription(response.GetLanguage(req), guides[0].Payment),
			LastChange: guides[0].UpdatedAt,
		}

		mockGuideProvider.On("GetGuidesByStatus", mock.Anything, mock.Anything).
			Return(guides, nil).Once()

		mockPubSub.On("Subscribe", mock.Anything, []any{})

		resp := GetOperatorGuide(req)
		data, ok := resp.Data.(GetOperatorGuideOutput)
		assert.True(t, ok, "expected type GetOperatorGuideOutput")
		assert.Equal(t, []model.OperatorGuide{operatorGuide}, data.OperatorGuides)
		assert.Equal(t, http.StatusOK, resp.HttpStatus)
		mockGuideProvider.AssertExpectations(t)
	})

	t.Run("invalid operator ID", func(t *testing.T) {
		req := newOperatorRequest(http.MethodGet, "/guide/operator", nil, nil)

		resp := GetOperatorGuide(req)

		assert.Equal(t, http.StatusUnauthorized, resp.HttpStatus)
	})

	t.Run("error fetching guides", func(t *testing.T) {
		mockGuideProvider := new(mock_guide_provider.MockGuideProvider)
		guide_provider.Set(mockGuideProvider)
		t.Cleanup(func() { guide_provider.Set(nil) })

		mockGuideProvider.On("GetGuidesByStatus", mock.Anything, mock.Anything).
			Return([]model.Guide{}, errors.New("db error")).Once()

		req := newOperatorRequest(http.MethodGet, "/guide/operator", nil, 42)

		resp := GetOperatorGuide(req)

		assert.Equal(t, http.StatusInternalServerError, resp.HttpStatus)
	})
}

func TestAssignGuideToOperator(t *testing.T) {
	testutil.InjectNoOpLogger()

	t.Run("success", func(t *testing.T) {
		mockGuideProvider := new(mock_guide_provider.MockGuideProvider)
		guide_provider.Set(mockGuideProvider)
		mockPubSub := new(mock_pubsub.MockPubSub)
		pubsub.Set(mockPubSub)

		mockGuideProvider.On("UpdateGuide", mock.Anything, mock.Anything).
			Return(nil).Once()

		mockPubSub.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		req := newGuideRequest(http.MethodPost, "/guide/assign/123", "123", nil, 42)
		w := httptest.NewRecorder()

		AssignGuideToOperator().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockGuideProvider.AssertExpectations(t)
	})

	t.Run("missing guide ID", func(t *testing.T) {
		req := newGuideRequest(http.MethodPost, "/guide/assign/", "", nil, 42)
		w := httptest.NewRecorder()

		AssignGuideToOperator().ServeHTTP(w, req)

		assertJSONErrorResponse(t, req, w, http.StatusBadRequest, i18n.MsgGuideRequired)
	})

	t.Run("invalid guide ID", func(t *testing.T) {
		req := newGuideRequest(http.MethodPost, "/guide/assign/abc", "abc", nil, 42)
		w := httptest.NewRecorder()

		AssignGuideToOperator().ServeHTTP(w, req)

		assertJSONErrorResponse(t, req, w, http.StatusBadRequest, i18n.MsgGuideInvalid)
	})

	t.Run("missing operator ID", func(t *testing.T) {
		req := newGuideRequest(http.MethodPost, "/guide/assign/123", "123", nil, nil)
		w := httptest.NewRecorder()

		AssignGuideToOperator().ServeHTTP(w, req)

		assertJSONErrorResponse(t, req, w, http.StatusUnauthorized, i18n.MsgOperatorInvalid)
	})

	t.Run("update guide fails", func(t *testing.T) {
		mockGuideProvider := new(mock_guide_provider.MockGuideProvider)
		guide_provider.Set(mockGuideProvider)
		t.Cleanup(func() { guide_provider.Set(nil) })

		mockGuideProvider.On("UpdateGuide", mock.Anything, mock.Anything).
			Return(errors.New("db error")).Once()

		req := newGuideRequest(http.MethodPost, "/guide/assign/123", "123", nil, 42)
		w := httptest.NewRecorder()

		AssignGuideToOperator().ServeHTTP(w, req)

		assertJSONErrorResponse(t, req, w, http.StatusInternalServerError, i18n.MsgInternalServerError)
		mockGuideProvider.AssertExpectations(t)
	})
}

func TestGetGuideStatusOptions(t *testing.T) {
	testutil.InjectNoOpLogger()

	t.Run("success", func(t *testing.T) {
		mockGuideProvider := new(mock_guide_provider.MockGuideProvider)
		guide_provider.Set(mockGuideProvider)
		t.Cleanup(func() { guide_provider.Set(nil) })

		mockGuideProvider.On("GetGuideById", mock.Anything, 123).
			Return(model.Guide{
				ID:         123,
				ViaGuideID: "V123",
				Status:     biz_guide_status.INITIAL,
				Payment:    biz_config.PAID_ON_DESTINATION,
			}, nil).Once()

		mockGuideProvider.On("GetGuideHistory", mock.Anything, 123).
			Return([]model.GuideHistory{{Status: "INIT"}}, nil).Once()

		req := newGuideRequest(http.MethodGet, "/guide/status/options/123", "123", nil, nil)
		w := httptest.NewRecorder()

		GetGuideStatusOptions().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockGuideProvider.AssertExpectations(t)
	})

	t.Run("missing guide ID", func(t *testing.T) {
		req := newGuideRequest(http.MethodGet, "/guide/status/options/", "", nil, nil)
		w := httptest.NewRecorder()

		GetGuideStatusOptions().ServeHTTP(w, req)

		assertJSONErrorResponse(t, req, w, http.StatusBadRequest, i18n.MsgGuideRequired)
	})

	t.Run("invalid guide ID", func(t *testing.T) {
		req := newGuideRequest(http.MethodGet, "/guide/status/options/abc", "abc", nil, nil)
		w := httptest.NewRecorder()

		GetGuideStatusOptions().ServeHTTP(w, req)

		assertJSONErrorResponse(t, req, w, http.StatusBadRequest, i18n.MsgGuideInvalid)
	})

	t.Run("error fetching guide", func(t *testing.T) {
		mockGuideProvider := new(mock_guide_provider.MockGuideProvider)
		guide_provider.Set(mockGuideProvider)
		t.Cleanup(func() { guide_provider.Set(nil) })

		mockGuideProvider.On("GetGuideById", mock.Anything, 123).
			Return(model.Guide{}, errors.New("db error")).Once()

		req := newGuideRequest(http.MethodGet, "/guide/status/options/123", "123", nil, nil)
		w := httptest.NewRecorder()

		GetGuideStatusOptions().ServeHTTP(w, req)

		assertJSONErrorResponse(t, req, w, http.StatusInternalServerError, i18n.MsgInternalServerError)
	})

	t.Run("guide not found", func(t *testing.T) {
		mockGuideProvider := new(mock_guide_provider.MockGuideProvider)
		guide_provider.Set(mockGuideProvider)
		t.Cleanup(func() { guide_provider.Set(nil) })

		mockGuideProvider.On("GetGuideById", mock.Anything, 123).
			Return(model.Guide{ID: 123, ViaGuideID: ""}, nil).Once()

		req := newGuideRequest(http.MethodGet, "/guide/status/options/123", "123", nil, nil)
		w := httptest.NewRecorder()

		GetGuideStatusOptions().ServeHTTP(w, req)

		assertJSONErrorResponse(t, req, w, http.StatusNotFound, i18n.MsgGuideNotFound)
	})

	t.Run("error fetching guide history", func(t *testing.T) {
		mockGuideProvider := new(mock_guide_provider.MockGuideProvider)
		guide_provider.Set(mockGuideProvider)
		t.Cleanup(func() { guide_provider.Set(nil) })

		mockGuideProvider.On("GetGuideById", mock.Anything, 123).
			Return(model.Guide{ID: 123, ViaGuideID: "V123"}, nil).Once()

		mockGuideProvider.On("GetGuideHistory", mock.Anything, 123).
			Return([]model.GuideHistory{}, errors.New("db error")).Once()

		req := newGuideRequest(http.MethodGet, "/guide/status/options/123", "123", nil, nil)
		w := httptest.NewRecorder()

		GetGuideStatusOptions().ServeHTTP(w, req)

		assertJSONErrorResponse(t, req, w, http.StatusInternalServerError, i18n.MsgInternalServerError)
	})
}

func TestUpdateGuideStatus(t *testing.T) {
	testutil.InjectNoOpLogger()
	mockGuideProvider := new(mock_guide_provider.MockGuideProvider)
	guide_provider.Set(mockGuideProvider)
	defer guide_provider.Set(nil)

	t.Run("missing guide ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := newGuideRequest(http.MethodPost, "/guide/update/", "", nil, 42)

		UpdateGuideStatus().ServeHTTP(w, req)
		assertJSONErrorResponse(t, req, w, http.StatusBadRequest, i18n.MsgGuideRequired)
		mockGuideProvider.AssertExpectations(t)
	})

	t.Run("invalid guide ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := newGuideRequest(http.MethodPost, "/guide/update/abc", "abc", nil, 42)

		UpdateGuideStatus().ServeHTTP(w, req)
		assertJSONErrorResponse(t, req, w, http.StatusBadRequest, i18n.MsgGuideInvalid)
		mockGuideProvider.AssertExpectations(t)
	})

	t.Run("missing operator ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := newGuideRequest(http.MethodPost, "/guide/update/123", "123", nil, nil)

		UpdateGuideStatus().ServeHTTP(w, req)
		assertJSONErrorResponse(t, req, w, http.StatusUnauthorized, i18n.MsgOperatorInvalid)
		mockGuideProvider.AssertExpectations(t)
	})

	t.Run("invalid JSON body", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := newGuideRequest(http.MethodPost, "/guide/update/123", "123", []byte("{invalid"), 42)

		UpdateGuideStatus().ServeHTTP(w, req)
		assertJSONErrorResponse(t, req, w, http.StatusBadRequest, i18n.MsgBadRequest)
		mockGuideProvider.AssertExpectations(t)
	})

	t.Run("failed DB update", func(t *testing.T) {
		mockGuideProvider.On("UpdateGuide", mock.Anything, mock.Anything).Return(errors.New("db error")).Once()

		body, _ := json.Marshal(UpdateGuideStatusInput{Status: "DELIVERED"})
		w := httptest.NewRecorder()
		req := newGuideRequest(http.MethodPost, "/guide/update/123", "123", body, 42)

		UpdateGuideStatus().ServeHTTP(w, req)
		assertJSONErrorResponse(t, req, w, http.StatusInternalServerError, i18n.MsgInternalServerError)
		mockGuideProvider.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		mockGuideProvider.On("UpdateGuide", mock.Anything, mock.Anything).Return(nil).Once()

		body, _ := json.Marshal(UpdateGuideStatusInput{Status: "DELIVERED"})
		w := httptest.NewRecorder()
		req := newGuideRequest(http.MethodPost, "/guide/update/123", "123", body, 42)

		UpdateGuideStatus().ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		mockGuideProvider.AssertExpectations(t)
	})
}

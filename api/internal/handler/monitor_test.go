package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	biz_guide_status "via/internal/biz/guide/status"
	"via/internal/model"
	guide_provider "via/internal/provider/guide"
	mock_guide_provider "via/internal/provider/guide/mock"
	"via/internal/response"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetMonitorEvents(t *testing.T) {
	mockGuideProvider := new(mock_guide_provider.MockGuideProvider)
	guide_provider.Set(mockGuideProvider)
	defer guide_provider.Set(nil) // Restore after test

	t.Run("success", func(t *testing.T) {
		mockGuideProvider.On("GetGuidesByStatus", mock.Anything, biz_guide_status.GetMonitorStatus()).
			Return([]model.Guide{
				{
					ViaGuideID: "V123",
					Recipient:  "John Doe",
					Status:     biz_guide_status.ON_HOLD,
				},
			}, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/monitor/events", nil)
		w := httptest.NewRecorder()

		GetMonitorEvents().ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp response.Response[GetMonitorEventOutput]
		err := json.NewDecoder(w.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.Len(t, resp.Data.Events, 1)
		assert.Equal(t, "V123", resp.Data.Events[0].GuideId)
		assert.Equal(t, "John Doe", resp.Data.Events[0].Recipient)
		assert.NotEmpty(t, resp.Data.Events[0].Status)
	})

	t.Run("error fetching guides", func(t *testing.T) {
		mockGuideProvider.On("GetGuidesByStatus", mock.Anything, biz_guide_status.GetMonitorStatus()).
			Return([]model.Guide{}, errors.New("db error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/monitor/events", nil)
		w := httptest.NewRecorder()

		GetMonitorEvents().ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

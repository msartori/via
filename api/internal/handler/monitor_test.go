package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	biz_guide_status "via/internal/biz/guide/status"
	"via/internal/i18n"
	"via/internal/model"
	guide_provider "via/internal/provider/guide"
	mock_guide_provider "via/internal/provider/guide/mock"
	"via/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetMonitorEvents(t *testing.T) {
	testutil.InjectNoOpLogger()

	t.Run("success", func(t *testing.T) {
		mockGuideProvider := new(mock_guide_provider.MockGuideProvider)
		guide_provider.Set(mockGuideProvider)
		mockGuideProvider.On("GetGuidesByStatus", mock.Anything, biz_guide_status.GetMonitorStatus()).
			Return([]model.Guide{
				{
					ViaGuideID: "V123",
					Recipient:  "John Doe",
					Status:     biz_guide_status.ON_HOLD,
				},
			}, nil).Twice()

		req := httptest.NewRequest(http.MethodGet, "/monitor/events", nil)

		res := GetMonitorEvents(req)

		data := res.Data.(GetMonitorEventOutput)

		assert.Len(t, data.Events, 1)
		assert.Equal(t, "V123", data.Events[0].GuideId)
		assert.Equal(t, "John Doe", data.Events[0].Recipient)
		assert.NotEmpty(t, data.Events[0].Status)

	})

	t.Run("error fetching guides", func(t *testing.T) {
		mockGuideProvider := new(mock_guide_provider.MockGuideProvider)
		guide_provider.Set(mockGuideProvider)
		mockGuideProvider.On("GetGuidesByStatus", mock.Anything, biz_guide_status.GetMonitorStatus()).
			Return([]model.Guide{}, errors.New("db error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/monitor/events", nil)

		res := GetMonitorEvents(req)

		assert.Equal(t, i18n.Get(req, i18n.MsgInternalServerError), res.Message)
	})
}

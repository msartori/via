package handler

import (
	"math/rand"
	"net/http"
	"time"
	biz_guide_status "via/internal/biz/guide/status"
	"via/internal/model"
	guide_provider "via/internal/provider/guide"
	"via/internal/response"
)

type GetMonitorEventOutput struct {
	Events []model.MonitorEvent `json:"events"`
}

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func GetMonitorEvents() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := response.Response[GetMonitorEventOutput]{}
		monitorEvents := []model.MonitorEvent{}

		guides, err := guide_provider.Get().GetGuidesByStatus(r.Context(), biz_guide_status.GetMonitorStatus())
		if isFailedToFetchGuide(w, r, err) {
			return
		}
		for _, guide := range guides {
			monitorEvents = append(monitorEvents,
				model.MonitorEvent{GuideId: guide.ViaGuideID,
					Recipient: guide.Recipient,
					Status:    model.GetMonitorEventStatusDescriptionByGuideStatus(response.GetLanguage(r), guide.Status),
					Highlight: model.GetHighlightByGuideStatus(guide.Status)})
		}
		res.Data = GetMonitorEventOutput{Events: monitorEvents}
		response.WriteJSON(w, r, res, http.StatusOK)
	})
}

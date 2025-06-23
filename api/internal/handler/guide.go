package handler

import (
	"net/http"
	biz_guide_status "via/internal/biz/guide/status"
	"via/internal/config"
	"via/internal/i18n"
	"via/internal/log"
	"via/internal/model"
	guide_provider "via/internal/provider/guide"
	via_guide_provider "via/internal/provider/via/guide"
	response "via/internal/response"

	"github.com/go-chi/chi/v5"
)

type GetGuideByViaGuideIdOutput struct {
	Guide model.Guide `json:"guide"`
}

func GetGuideByViaGuideId() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		res := response.Response[GetGuideByViaGuideIdOutput]{}
		viaGuideId := chi.URLParam(r, "viaGuideId")

		if ok := isValidViaGuideId(w, r, viaGuideId); !ok {
			return
		}

		logger := log.Get()
		logger.WithLogFieldsInRequest(r, "via_guide_id", viaGuideId)

		guide, err := guide_provider.Get().GetGuideByViaGuideId(r.Context(), viaGuideId)

		if ok := isFailedToFetchGuide(w, r, err); ok {
			return
		}

		if ok := isGuideNotFound(w, r, guide.ViaGuideID); !ok {
			return
		}

		res.Data = GetGuideByViaGuideIdOutput{guide}
		response.WriteJSON(w, r, res, http.StatusOK)
	})
}

type CreateGuideToWidthdrawInput struct {
	ViaGuideId string `json:"viaGuideId"`
}

type CreateGuideToWidthdrawOutput struct {
	WithdrawMessage string `json:"withdrawMessage"`
}

func CreateGuideToWidthdraw(biz config.Bussiness) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := response.Response[*CreateGuideToWidthdrawOutput]{}
		var input CreateGuideToWidthdrawInput

		if ok := getJsonBody(w, r, &input); !ok {
			return
		}

		if ok := isValidViaGuideId(w, r, input.ViaGuideId); !ok {
			return
		}
		logger := log.Get()
		logger.WithLogFieldsInRequest(r, "via_guide_id", input.ViaGuideId)

		viaGuide, err := via_guide_provider.Get().GetGuide(r.Context(), input.ViaGuideId)

		if ok := isFailedToFetchGuide(w, r, err); ok {
			return
		}

		data := &CreateGuideToWidthdrawOutput{}
		res.Data = data

		if ok := isGuideNotFound(w, r, viaGuide.ID); !ok {
			return
		}

		if ok := isInvalidViaGuideToWithdraw(viaGuide, biz); !ok {
			logger.Warn(r.Context(), "msg", "not able to create a new guide to process", "via_guide", viaGuide)
			res.Message = i18n.Get(r, i18n.MsgGuideInvalid)
			response.WriteJSON(w, r, res, http.StatusBadRequest)
			return
		}

		guide, err := guide_provider.Get().GetGuideByViaGuideId(r.Context(), viaGuide.ID)

		if ok := isFailedToFetchGuide(w, r, err); ok {
			return
		}

		if guide.ID != 0 {
			logger.WithLogFieldsInRequest(r, "guide_id", guide.ID)
			if biz_guide_status.IsAbleToReInit(guide.Status) {
				guide_provider.Get().ReinitGuide(r.Context(), guide.ID)
				logger.Info(r.Context(), "msg", "guide re-init")
				data.WithdrawMessage = getWithDrawMessage(r, inProcess, viaGuide.ID)
				response.WriteJSON(w, r, res, http.StatusOK)
				return
			}

			if !biz_guide_status.IsValidToCreateForWithdraw(guide.Status) {
				logger.Warn(r.Context(), "msg", "not able to create a new guide to process",
					"via_guide", viaGuide, "guide", guide)
				res.Message = i18n.Get(r, i18n.MsgGuideInvalid)
				response.WriteJSON(w, r, res, http.StatusBadRequest)
				return
			}
		}

		id, err := guide_provider.Get().CreateGuide(r.Context(), viaGuide)

		if err != nil {
			logger.Error(r.Context(), err, "msg", "failed to create guide")
			res.Message = i18n.Get(r, i18n.MsgInternalServerError)
			status := http.StatusInternalServerError
			response.WriteJSON(w, r, res, status)
			return
		}
		logger.WithLogFieldsInRequest(r, "guide_id", id)
		logger.Info(r.Context(), "msg", "guide to withdraw created")
		data.WithdrawMessage = getWithDrawMessage(r, inProcess, viaGuide.ID)
		response.WriteJSON(w, r, res, http.StatusOK)
	})
}

type GetOperatorGuideOutput struct {
	OperatorGuides []model.OperatorGuide `json:"operatorGuides"`
}

func GetOperatorGuide() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := response.Response[GetOperatorGuideOutput]{}
		operatorGuides := []model.OperatorGuide{}
		guides, err := guide_provider.Get().GetGuidesByStatus(r.Context(), biz_guide_status.GetOperatorStatus())
		if ok := isFailedToFetchGuide(w, r, err); ok {
			return
		}
		for _, guide := range guides {
			operatorGuides = append(operatorGuides,
				model.OperatorGuide{
					GuideId:   guide.ID,
					Recipient: guide.Recipient,
					Status:    biz_guide_status.Description(GetLanguage(r), guide.Status),
					Operator:  guide.Operator,
				})
		}
		logger := log.Get()
		logger.Info(r.Context(), "msg", "returning operator guides")
		res.Data = GetOperatorGuideOutput{OperatorGuides: operatorGuides}
		response.WriteJSON(w, r, res, http.StatusOK)
	})
}

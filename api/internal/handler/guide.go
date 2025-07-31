package handler

import (
	"fmt"
	"net/http"
	biz_config "via/internal/biz/config"
	biz_guide_status "via/internal/biz/guide/status"
	biz_operator "via/internal/biz/operator"
	"via/internal/global"
	"via/internal/i18n"
	"via/internal/log"
	"via/internal/middleware"
	"via/internal/model"
	guide_provider "via/internal/provider/guide"
	via_guide_provider "via/internal/provider/via/guide"
	"via/internal/pubsub"
	response "via/internal/response"

	"github.com/go-chi/chi/v5"
)

type CreateGuideToWidthdrawInput struct {
	ViaGuideId string `json:"viaGuideId"`
}

type CreateGuideToWidthdrawOutput struct {
	WithdrawMessage string `json:"withdrawMessage"`
}

func CreateGuideToWidthdraw(biz biz_config.BussinessCfg) http.Handler {
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
				guide_provider.Get().UpdateGuide(r.Context(), model.Guide{ID: guide.ID, Status: biz_guide_status.INITIAL})
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
		err = pubsub.Get().Publish(r.Context(), global.NewGuideChannel, fmt.Sprintf("{\"guide_id\":\"%d\"}", id))
		if err != nil {
			logger.Error(r.Context(), err, "msg", "unable to publish event", "channel", global.NewGuideChannel)
		}
		data.WithdrawMessage = getWithDrawMessage(r, inProcess, viaGuide.ID)
		response.WriteJSON(w, r, res, http.StatusOK)
	})
}

type GetOperatorGuideOutput struct {
	OperatorGuides []model.OperatorGuide `json:"operatorGuides"`
}

func GetOperatorGuide(r *http.Request) response.Response[any] {
	res := response.Response[any]{}
	operatorGuides := []model.OperatorGuide{}

	operatorIdCtx := r.Context().Value(middleware.OperatorIDKey)
	operatorId, ok := operatorIdCtx.(int)
	if !ok {
		log.Get().Warn(r.Context(), "msg", "missing operator id in context", "operator_id", operatorIdCtx)
		res.Message = i18n.Get(r, i18n.MsgOperatorInvalid)
		res.HttpStatus = http.StatusUnauthorized
		return res
	}

	logger := log.Get()
	logger.WithLogFieldsInRequest(r, "operator_id", operatorId)

	guides, err := guide_provider.Get().GetGuidesByStatus(r.Context(), biz_guide_status.GetOperatorStatus())

	if err != nil {
		log.Get().Error(r.Context(), err, "msg", "failed to fetch guide")
		res.Message = i18n.Get(r, i18n.MsgInternalServerError)
		res.HttpStatus = http.StatusInternalServerError
		return res
	}

	for _, guide := range guides {
		operatorGuides = append(operatorGuides,
			model.OperatorGuide{
				GuideId:    guide.ID,
				Recipient:  guide.Recipient,
				Status:     biz_guide_status.GetStatusDescription(response.GetLanguage(r), guide.Status),
				Operator:   guide.Operator,
				Selectable: guide.Operator.ID == biz_operator.OPERATOR_SYSTEM || guide.Operator.ID == operatorId,
				ViaGuideId: guide.ViaGuideID,
				Payment:    biz_config.GetPaymentDescription(response.GetLanguage(r), guide.Payment),
				LastChange: guide.UpdatedAt,
			})
	}
	logger.Info(r.Context(), "msg", "returning operator guides")
	res.Data = GetOperatorGuideOutput{OperatorGuides: operatorGuides}
	res.HttpStatus = http.StatusOK
	return res
}

type AssignGuideToOperatorOutput struct {
	Guide model.Guide `json:"guide"`
}

func AssignGuideToOperator() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := response.Response[AssignGuideToOperatorOutput]{}
		valid, guideId := isValidGuideId(w, r, chi.URLParam(r, "guideId"))
		if !valid {
			return
		}
		logger := log.Get()
		logger.WithLogFieldsInRequest(r, "guide_id", guideId)
		operatorId := 0
		if operatorId = isValidOperatorId(w, r); operatorId == 0 {
			return
		}
		logger.WithLogFieldsInRequest(r, "operator_id", operatorId)
		err := guide_provider.Get().UpdateGuide(r.Context(), model.Guide{ID: guideId, Operator: model.Operator{ID: operatorId}})
		if err != nil {
			log.Get().Error(r.Context(), err, "msg", "failed assigning operator to guide")
			res.Message = i18n.Get(r, i18n.MsgInternalServerError)
			status := http.StatusInternalServerError
			response.WriteJSON(w, r, res, status)
			return
		}
		logger.Info(r.Context(), "msg", "operator assigned to guide")
		pubsub.Get().Publish(r.Context(), global.GuideAssignmentChannel, fmt.Sprintf("{\"guide_id\":\"%d\"}", guideId))
		response.WriteJSON(w, r, res, http.StatusOK)
	})
}

type GetGuideStatusOptionsOutput struct {
	StatusOptions []model.GenericIdDesc `json:"statusOption"`
}

func GetGuideStatusOptions() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := response.Response[GetGuideStatusOptionsOutput]{}
		valid, guideId := isValidGuideId(w, r, chi.URLParam(r, "guideId"))
		if !valid {
			return
		}
		logger := log.Get()
		logger.WithLogFieldsInRequest(r, "guide_id", guideId)

		guide, err := guide_provider.Get().GetGuideById(r.Context(), guideId)
		if ok := isFailedToFetchGuide(w, r, err); ok {
			return
		}
		if ok := isGuideNotFound(w, r, guide.ViaGuideID); !ok {
			return
		}

		guideHistory, err := guide_provider.Get().GetGuideHistory(r.Context(), guide.ID)
		if ok := isFailedToFetchGuide(w, r, err); ok {
			return
		}

		statusHistory := make([]string, len(guideHistory))
		for i, h := range guideHistory {
			statusHistory[i] = h.Status
		}

		nextStatus := biz_guide_status.GetNextStatus(guide.Status, statusHistory, guide.Payment)

		statusOptions := []model.GenericIdDesc{}

		for _, nStatus := range nextStatus {
			statusOptions = append(statusOptions,
				model.GenericIdDesc{ID: nStatus,
					Description: biz_guide_status.GetStatusDescription(response.GetLanguage(r), nStatus),
					Extra:       biz_guide_status.GetStatusState(nStatus)})
		}
		res.Data = GetGuideStatusOptionsOutput{StatusOptions: statusOptions}
		response.WriteJSON(w, r, res, http.StatusOK)
	})
}

type UpdateGuideStatusInput struct {
	Status string `json:"status"`
}

func UpdateGuideStatus() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := response.Response[any]{}
		valid, guideId := isValidGuideId(w, r, chi.URLParam(r, "guideId"))
		if !valid {
			return
		}
		logger := log.Get()
		logger.WithLogFieldsInRequest(r, "guide_id", guideId)

		operatorId := 0
		if operatorId = isValidOperatorId(w, r); operatorId == 0 {
			return
		}
		logger.WithLogFieldsInRequest(r, "operator_id", operatorId)

		var input UpdateGuideStatusInput
		if ok := getJsonBody(w, r, &input); !ok {
			return
		}

		logger.WithLogFieldsInRequest(r, "status", input.Status)
		err := guide_provider.Get().UpdateGuide(r.Context(), model.Guide{ID: guideId, Status: input.Status})
		if err != nil {
			log.Get().Error(r.Context(), err, "msg", "failed updating guide status")
			res.Message = i18n.Get(r, i18n.MsgInternalServerError)
			status := http.StatusInternalServerError
			response.WriteJSON(w, r, res, status)
			return
		}
		logger.Info(r.Context(), "msg", "guide status updated")
		pubsub.Get().Publish(r.Context(), global.GuideStatusChangeChannel, fmt.Sprintf("{\"guide_id\":\"%d\"}", guideId))
		response.WriteJSON(w, r, res, http.StatusOK)
	})
}

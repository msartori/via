package guide_ent_provider

import (
	"context"
	"fmt"
	biz_guide_status "via/internal/biz/guide/status"
	biz_operator "via/internal/biz/operator"
	ent_client "via/internal/client/ent"
	"via/internal/ent"
	"via/internal/ent/guide"
	"via/internal/ent/guidehistory"
	"via/internal/log"
	"via/internal/model"

	"entgo.io/ent/dialect/sql"
)

type GuideEntProvider struct {
	client *ent.Client
}

func New() GuideEntProvider {
	return GuideEntProvider{ent_client.Get()}
}

func fromEntGuide(guide ent.Guide) model.Guide {
	operator := model.Operator{ID: guide.OperatorID}
	if guide.Edges.Operator != nil {
		operator = model.Operator{
			ID:      guide.Edges.Operator.ID,
			Account: guide.Edges.Operator.Account,
			Name:    guide.Edges.Operator.Name,
			Enabled: guide.Edges.Operator.Enabled}
	}
	return model.Guide{
		ID:         guide.ID,
		ViaGuideID: guide.ViaGuideID,
		Recipient:  guide.Recipient,
		Payment:    guide.Payment,
		Status:     guide.Status,
		Operator:   operator,
		CreatedAt:  guide.CreatedAt,
		UpdatedAt:  guide.UpdatedAt,
	}
}

func fromEntGuideHistory(guideHistory ent.GuideHistory) model.GuideHistory {
	return model.GuideHistory{
		ID:        guideHistory.ID,
		Status:    guideHistory.Status,
		Timestamp: guideHistory.CreatedAt,
	}
}

func (p GuideEntProvider) GetGuideByViaGuideId(ctx context.Context, viaGuideId string) (model.Guide, error) {
	guide, err := p.client.Guide.
		Query().
		Where(guide.ViaGuideID(viaGuideId)).
		Order(guide.ByCreatedAt(sql.OrderDesc())).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return model.Guide{}, nil

		}
		log.Get().Error(ctx, err, "msg", "failed querying guide by viaGuideId")
		return model.Guide{}, fmt.Errorf("failed querying guide by viaGuideId: %w", err)

	}
	return fromEntGuide(*guide), err
}

func (p GuideEntProvider) CreateGuide(ctx context.Context, viaGuide model.ViaGuide) (int, error) {
	gp, err := p.client.Guide.
		Create().
		SetViaGuideID(viaGuide.ID).
		SetStatus(biz_guide_status.INITIAL).
		SetRecipient(viaGuide.Recipient).
		SetOperatorID(biz_operator.OPERATOR_SYSTEM).
		SetPayment(viaGuide.Payment).
		Save(ctx)
	if err != nil {
		log.Get().Error(ctx, err, "msg", "failed creating Guide")
		return 0, fmt.Errorf("failed creating Guide: %w", err)
	}
	return gp.ID, nil
}

func (p GuideEntProvider) GetGuidesByStatus(ctx context.Context, status []string) ([]model.Guide, error) {
	guides := []model.Guide{}
	gp, err := p.client.Guide.
		Query().
		Where(guide.StatusIn(status...)).
		WithOperator().
		Order(guide.ByCreatedAt(sql.OrderAsc())).
		All(ctx)
	if err != nil {
		log.Get().Error(ctx, err, "msg", "failed getting Guides by status")
		return guides, fmt.Errorf("failed getting Monitor Guides: %w", err)
	}
	for _, guide := range gp {
		guides = append(guides, fromEntGuide(*guide))
	}
	return guides, nil
}

func (p GuideEntProvider) UpdateGuide(ctx context.Context, guide model.Guide) error {
	guideUpdateOne := p.client.Guide.UpdateOneID(guide.ID)
	if guide.Operator.ID != 0 {
		guideUpdateOne.SetOperatorID(guide.Operator.ID)
	}
	if guide.Status != "" {
		guideUpdateOne.SetStatus(guide.Status)
	}
	guideUpdated, err := guideUpdateOne.Save(ctx)
	if err != nil {
		log.Get().Error(ctx, err, "msg", "error updating guide", "guide_id", guide.ID)
	} else {
		log.Get().Info(ctx, "msg", "guide updated", "guide_id", guideUpdated.ID)
	}
	return err
}

func (p GuideEntProvider) GetGuideById(ctx context.Context, id int) (model.Guide, error) {
	guide, err := p.client.Guide.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return model.Guide{}, nil

		}
		log.Get().Error(ctx, err, "msg", "failed querying guide by id")
		return model.Guide{}, fmt.Errorf("failed querying guide by id: %w", err)

	}
	return fromEntGuide(*guide), err
}

func (p GuideEntProvider) GetGuideHistory(ctx context.Context, guideId int) ([]model.GuideHistory, error) {
	guideHistory := []model.GuideHistory{}
	ghs, err := p.client.GuideHistory.Query().
		Where(guidehistory.GuideID(guideId)).
		Order(guidehistory.ByCreatedAt(sql.OrderAsc())).
		All(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return []model.GuideHistory{}, nil

		}
		log.Get().Error(ctx, err, "msg", "failed querying guide history by guide id")
		return []model.GuideHistory{}, fmt.Errorf("failed querying guide history by guide id: %w", err)
	}
	for _, gh := range ghs {
		guideHistory = append(guideHistory, fromEntGuideHistory(*gh))
	}
	return guideHistory, nil
}

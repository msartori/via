package guide_ent_provider

import (
	"context"
	"fmt"
	biz_operator "via/internal/biz"
	biz_guide_status "via/internal/biz/guide/status"
	ent_client "via/internal/client/ent"
	"via/internal/ent"
	"via/internal/ent/guide"
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
			Name:    guide.Edges.Operator.Account,
			Enabled: guide.Edges.Operator.Enabled}
	}
	return model.Guide{
		ID:         guide.ID,
		ViaGuideID: guide.ViaGuideID,
		Recipient:  guide.Recipient,
		Status:     guide.Status,
		Operator:   operator,
		CreatedAt:  guide.CreatedAt,
		UpdatedAt:  guide.UpdatedAt,
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
		Save(ctx)
	if err != nil {
		log.Get().Error(ctx, err, "msg", "failed creating Guide")
		return 0, fmt.Errorf("failed creating Guide: %w", err)
	}
	return gp.ID, nil
}

func (p GuideEntProvider) ReinitGuide(ctx context.Context, id int) (int, error) {
	gp, err := p.client.Guide.
		UpdateOneID(id).
		SetStatus(biz_guide_status.INITIAL).
		Save(ctx)
	if err != nil {
		log.Get().Error(ctx, err, "msg", "failed updating Guide")
		return 0, fmt.Errorf("failed updating Guide: %w", err)
	}
	return gp.ID, nil
}

func (p GuideEntProvider) GetGuidesByStatus(ctx context.Context, status []string) ([]model.Guide, error) {
	guides := []model.Guide{}
	gp, err := p.client.Guide.
		Query().
		Where(guide.StatusIn(status...)).
		WithOperator().
		Order(guide.ByUpdatedAt(sql.OrderAsc())).
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

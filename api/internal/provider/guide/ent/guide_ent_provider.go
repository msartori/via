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
	return model.FromEntGuide(*guide), err
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

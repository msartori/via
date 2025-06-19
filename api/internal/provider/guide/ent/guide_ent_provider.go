package guide_ent_provider

import (
	"context"
	"fmt"
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
		SetStatus("Initial").
		SetRecipient(viaGuide.Recipient).
		SetOperatorID(1).
		Save(ctx)
	if err != nil {
		log.Get().Error(ctx, err, "msg", "failed creating Guide")
		return 0, fmt.Errorf("failed creating Guide: %w", err)
	}
	return gp.ID, nil
}

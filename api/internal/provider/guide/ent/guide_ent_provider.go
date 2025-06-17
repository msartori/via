package guide_ent_provider

import (
	"context"
	"fmt"
	ent_client "via/internal/client/ent"
	"via/internal/ent"
	"via/internal/ent/guide"
	"via/internal/log"
	"via/internal/model"
)

type GuideEntProvider struct {
	client *ent.Client
}

func New() GuideEntProvider {
	return GuideEntProvider{ent_client.Get()}
}

func (p GuideEntProvider) GetGuideByCode(ctx context.Context, code string) (model.GuideProcess, error) {
	logger := log.Get()
	logger.Info(ctx, "msg", "guide_ent_process_provider.GetGuideProcessByCode_start")
	defer logger.Info(ctx, "msg", "guide_ent_process_provider.GetGuideProcessByCode_end")
	gp, err := p.client.Guide.
		Query().
		Where(guide.Code(code)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return model.GuideProcess{}, nil

		}
		log.Get().Error(ctx, err, "msg", "failed querying guide process by code")
		return model.GuideProcess{}, fmt.Errorf("failed querying guide process by code: %w", err)

	}
	return model.GuideProcess{
		ID:        gp.ID,
		Code:      gp.Code,
		Recipient: gp.Recipient,
		Status:    gp.Status,
		CreatedAt: gp.CreatedAt,
		UpdatedAt: gp.UpdatedAt,
	}, err
}

func (p GuideEntProvider) CreateGuide(ctx context.Context, guide model.ViaGuide) (int, error) {
	logger := log.Get()
	logger.Info(ctx, "msg", "guide_ent_process_provider.CreateGuide_start")
	defer logger.Info(ctx, "msg", "guide_ent_process_provider.CreateGuide_end")
	gp, err := p.client.Guide.
		Create().
		SetCode(guide.ID).
		SetStatus("Initial").
		SetRecipient(guide.Recipient).
		Save(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed creating user: %w", err)
	}
	return gp.ID, nil
}

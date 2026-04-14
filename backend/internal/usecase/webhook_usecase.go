package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
	domainUC "github.com/org/itsm/internal/domain/usecase"
)

type webhookUseCase struct {
	webhookRepo repository.WebhookRepository
	dispatcher  WebhookDispatcher
}

// WebhookDispatcher is an interface for the infrastructure dispatcher
type WebhookDispatcher interface {
	Dispatch(ctx context.Context, event entity.WebhookEvent, payload interface{}, configs []*entity.WebhookConfig) error
}

func NewWebhookUseCase(webhookRepo repository.WebhookRepository, dispatcher WebhookDispatcher) domainUC.WebhookUseCase {
	return &webhookUseCase{webhookRepo: webhookRepo, dispatcher: dispatcher}
}

func (uc *webhookUseCase) CreateWebhookConfig(ctx context.Context, req domainUC.CreateWebhookRequest) (*entity.WebhookConfig, error) {
	now := time.Now().UTC()
	config := &entity.WebhookConfig{
		ID:        uuid.New(),
		URL:       req.URL,
		Events:    req.Events,
		SecretKey: req.Secret,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := uc.webhookRepo.Create(ctx, config); err != nil {
		return nil, err
	}
	return config, nil
}

func (uc *webhookUseCase) UpdateWebhookConfig(ctx context.Context, id uuid.UUID, req domainUC.CreateWebhookRequest) error {
	config, err := uc.webhookRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	config.URL = req.URL
	config.Events = req.Events
	config.SecretKey = req.Secret
	config.UpdatedAt = time.Now().UTC()
	return uc.webhookRepo.Update(ctx, config)
}

func (uc *webhookUseCase) DeleteWebhookConfig(ctx context.Context, id uuid.UUID) error {
	return uc.webhookRepo.Delete(ctx, id)
}

func (uc *webhookUseCase) ListWebhookConfigs(ctx context.Context) ([]*entity.WebhookConfig, error) {
	return uc.webhookRepo.FindAll(ctx)
}

func (uc *webhookUseCase) Dispatch(ctx context.Context, event entity.WebhookEvent, payload interface{}) error {
	configs, err := uc.webhookRepo.FindAll(ctx)
	if err != nil {
		return err
	}
	// Filter active configs that subscribe to this event
	var matching []*entity.WebhookConfig
	for _, c := range configs {
		if !c.IsActive {
			continue
		}
		for _, e := range c.Events {
			if e == event {
				matching = append(matching, c)
				break
			}
		}
	}
	if len(matching) == 0 {
		return nil
	}
	return uc.dispatcher.Dispatch(ctx, event, payload, matching)
}

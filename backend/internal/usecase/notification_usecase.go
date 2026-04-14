package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/org/itsm/internal/domain/entity"
	"github.com/org/itsm/internal/domain/repository"
	domainUC "github.com/org/itsm/internal/domain/usecase"
)

type notificationUseCase struct {
	notifRepo repository.NotificationRepository
}

func NewNotificationUseCase(notifRepo repository.NotificationRepository) domainUC.NotificationUseCase {
	return &notificationUseCase{notifRepo: notifRepo}
}

func (uc *notificationUseCase) GetNotifications(ctx context.Context, userID uuid.UUID) ([]*entity.Notification, error) {
	return uc.notifRepo.FindByUserID(ctx, userID)
}

func (uc *notificationUseCase) MarkAsRead(ctx context.Context, notifID uuid.UUID, userID uuid.UUID) error {
	return uc.notifRepo.MarkAsRead(ctx, notifID, userID)
}

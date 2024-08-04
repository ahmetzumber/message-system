package repository

import (
	"context"
	"message-system/app/domain"
)

type MessageRepository interface {
	GetUnsentMessages(ctx context.Context, limit int) ([]domain.Message, error)
	GetSentMessages(ctx context.Context) ([]domain.Message, error)
	MarkAsSent(ctx context.Context, id int) error
}

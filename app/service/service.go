package service

import (
	"context"
	"log"
	"message-system/app/constants"
	"message-system/app/domain"
	"message-system/app/types"
	"time"
)

type CacheService interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}

type WebhookClient interface {
	SendMessage(ctx context.Context, message *domain.Message) (*types.MessageResponse, error)
}

type MessageRepository interface {
	GetUnsentMessages(ctx context.Context, limit int) ([]domain.Message, error)
	GetSentMessages(ctx context.Context) ([]domain.Message, error)
	MarkAsSent(ctx context.Context, id int) error
}

type Service struct {
	cache      CacheService
	client     WebhookClient
	repository MessageRepository
}

func NewService(cache CacheService,
	client WebhookClient,
	repository MessageRepository) *Service {
	return &Service{
		cache:      cache,
		client:     client,
		repository: repository,
	}
}

func (s *Service) StartSending() {
	ctx := context.Background()
	messages, err := s.repository.GetUnsentMessages(ctx, constants.MessageLimit)
	if err != nil {
		log.Println(err)
		return
	}

	for _, message := range messages {
		startTime := time.Now()
		response, err := s.client.SendMessage(ctx, &message)
		if err != nil {
			log.Println(constants.ErrorSendingMessages)
			return
		}

		err = s.repository.MarkAsSent(ctx, message.ID)
		if err != nil {
			log.Println(constants.ErrorMarkingMessages)
			return
		}

		err = s.cache.Set(ctx, response.Message, startTime.GoString(), time.Hour)
		if err != nil {
			log.Println(constants.ErrorCachingMessages)
			return
		}
	}
}

func (s *Service) GetSentMessages(ctx context.Context) ([]types.Message, error) {
	messages, err := s.repository.GetSentMessages(ctx)
	if err != nil {
		return nil, err
	}
	result := convertMessagesToResponse(messages)
	if len(result) == 0 {
		return []types.Message{}, nil
	}
	return result, nil
}

func convertMessagesToResponse(messages []domain.Message) []types.Message {
	var result []types.Message
	for _, message := range messages {
		result = append(result, types.Message{
			ID:      message.ID,
			To:      message.To,
			Content: message.Content,
			IsSent:  message.IsSent,
		})
	}
	return result
}

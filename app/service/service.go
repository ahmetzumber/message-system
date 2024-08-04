package service

import (
	"context"
	"log"
	"message-system/app/cache/redis"
	"message-system/app/client"
	"message-system/app/constants"
	"message-system/app/domain"
	"message-system/app/repository/mongodb"
	"message-system/app/types"
	"time"
)

type Service struct {
	cache      *redis.Cache
	client     *client.WebhookClient
	repository *mongodb.MessageRepository
}

func NewService(cache *redis.Cache,
	client *client.WebhookClient,
	repository *mongodb.MessageRepository) *Service {
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
	}

	for _, message := range messages {
		startTime := time.Now()
		response, err := s.client.SendMessage(ctx, &message)
		if err != nil {
			log.Println(constants.ErrorSendingMessages)
		}

		err = s.repository.MarkAsSent(ctx, message.ID)
		if err != nil {
			log.Println(constants.ErrorMarkingMessages)
		}

		err = s.cache.Set(ctx, response.Message, startTime.GoString(), time.Hour)
		if err != nil {
			log.Println(constants.ErrorCachingMessages)
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

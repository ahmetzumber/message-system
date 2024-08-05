//go:build unit
// +build unit

package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"message-system/app/constants"
	"message-system/app/domain"
	"message-system/app/types"
	"message-system/mocks"
	"testing"
	"time"
)

func TestService_StartSending(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockCacheService := mocks.NewMockCacheService(mockController)
	mockWebhookClient := mocks.NewMockWebhookClient(mockController)
	mockMessageRepository := mocks.NewMockMessageRepository(mockController)

	service := NewService(mockCacheService, mockWebhookClient, mockMessageRepository)
	ctx := context.Background()

	t.Run("successful message sending", func(t *testing.T) {
		messages := []domain.Message{
			{ID: 1, Content: "Message 1", To: "+90458748", IsSent: false},
			{ID: 2, Content: "Message 2", To: "+902131542", IsSent: false},
		}

		mockMessageRepository.
			EXPECT().
			GetUnsentMessages(ctx, constants.MessageLimit).
			Return(messages, nil)

		mockWebhookClient.
			EXPECT().
			SendMessage(ctx, &messages[0]).
			Return(&types.MessageResponse{Message: "Message 1"}, nil)

		mockWebhookClient.
			EXPECT().
			SendMessage(ctx, &messages[1]).
			Return(&types.MessageResponse{Message: "Message 2"}, nil)

		mockMessageRepository.
			EXPECT().
			MarkAsSent(ctx, messages[0].ID).
			Return(nil)

		mockMessageRepository.
			EXPECT().
			MarkAsSent(ctx, messages[1].ID).
			Return(nil)

		mockCacheService.
			EXPECT().
			Set(ctx, "Message 1", gomock.Any(), time.Hour).
			Return(nil)

		mockCacheService.
			EXPECT().
			Set(ctx, "Message 2", gomock.Any(), time.Hour).
			Return(nil)

		service.StartSending()
	})
	t.Run("error on getting unsent messages", func(t *testing.T) {
		mockMessageRepository.
			EXPECT().
			GetUnsentMessages(ctx, constants.MessageLimit).
			Return(nil, errors.New("db error"))

		service.StartSending()
	})
	t.Run("error on sending message", func(t *testing.T) {
		messages := []domain.Message{
			{ID: 1, Content: "Message 1"},
		}

		mockMessageRepository.
			EXPECT().
			GetUnsentMessages(ctx, constants.MessageLimit).
			Return(messages, nil)

		mockWebhookClient.
			EXPECT().
			SendMessage(ctx, &messages[0]).
			Return(nil, errors.New("send error"))

		mockMessageRepository.
			EXPECT().
			MarkAsSent(gomock.Any(), gomock.Any()).
			Times(0)

		service.StartSending()
	})
	t.Run("error on marking message as sent", func(t *testing.T) {
		messages := []domain.Message{
			{ID: 1, Content: "Message 1"},
		}

		mockMessageRepository.
			EXPECT().
			GetUnsentMessages(ctx, constants.MessageLimit).
			Return(messages, nil)

		mockWebhookClient.
			EXPECT().
			SendMessage(ctx, &messages[0]).
			Return(&types.MessageResponse{Message: "Message 1"}, nil)

		mockMessageRepository.
			EXPECT().
			MarkAsSent(ctx, messages[0].ID).
			Return(errors.New("mark error"))

		service.StartSending()
	})
	t.Run("error on setting cache", func(t *testing.T) {
		messages := []domain.Message{
			{ID: 1, Content: "Message 1"},
		}

		mockMessageRepository.
			EXPECT().
			GetUnsentMessages(ctx, constants.MessageLimit).
			Return(messages, nil)

		mockWebhookClient.
			EXPECT().
			SendMessage(ctx, &messages[0]).
			Return(&types.MessageResponse{Message: "Message 1"}, nil)

		mockMessageRepository.
			EXPECT().
			MarkAsSent(ctx, messages[0].ID).
			Return(nil)

		mockCacheService.
			EXPECT().
			Set(ctx, "Message 1", gomock.Any(), time.Hour).
			Return(errors.New("cache error"))

		service.StartSending()
	})
}

func TestService_GetSentMessages(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockMessageRepository := mocks.NewMockMessageRepository(mockController)
	service := NewService(nil, nil, mockMessageRepository)
	ctx := context.Background()

	t.Run("successful scenario", func(t *testing.T) {
		messages := []domain.Message{
			{ID: 1, To: "+90458748", Content: "deneme1", IsSent: true},
			{ID: 2, To: "+902131542", Content: "deneme2", IsSent: true},
		}

		mockMessageRepository.
			EXPECT().
			GetSentMessages(ctx).
			Return(messages, nil)

		result, err := service.GetSentMessages(ctx)
		assert.NoError(t, err)
		assert.Len(t, result, len(messages))
		assert.Equal(t, "+90458748", result[0].To)
		assert.Equal(t, "+902131542", result[1].To)

	})

	t.Run("error scenario", func(t *testing.T) {
		mockMessageRepository.
			EXPECT().
			GetSentMessages(ctx).
			Return(nil, errors.New("db error"))

		result, err := service.GetSentMessages(ctx)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

//go:build unit
// +build unit

package handler

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"message-system/app/constants"
	"message-system/app/types"
	"message-system/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_StartStopSending(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	mockService := mocks.NewMockService(mockController)
	handler := NewHandler(mockService)

	t.Run("start-stop scheduler test", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/start-stop", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		assert.Nil(t, scheduler)

		if assert.NoError(t, handler.StartStopSending(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), constants.SendingStarted)
			assert.NotNil(t, scheduler)
		}

		rec = httptest.NewRecorder()
		c = echo.New().NewContext(req, rec)
		if assert.NoError(t, handler.StartStopSending(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), constants.SendingStopped)
			assert.Nil(t, scheduler)
		}
	})
}

func Test_GetSentMessages(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	mockService := mocks.NewMockService(mockController)
	handler := NewHandler(mockService)

	t.Run("success scenario", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/sent-messages", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		mockService.
			EXPECT().
			GetSentMessages(c.Request().Context()).
			Return([]types.Message{}, nil)

		err := handler.GetSentMessages(c)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("error scenario", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/sent-messages", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		mockService.
			EXPECT().
			GetSentMessages(c.Request().Context()).
			Return(nil, errors.New("esrror"))

		err := handler.GetSentMessages(c)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

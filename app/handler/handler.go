package handler

import (
	"message-system/app/constants"
	"message-system/app/service"
	"message-system/app/types"
	"net/http"
	"sync"

	"github.com/jasonlvhit/gocron"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Message System API
// @version 1.0
// @description This is a sample server for a message system.
// @host localhost:8080
// @BasePath /api
var (
	scheduler *gocron.Scheduler
	mutex     sync.Mutex
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func RegisterRoutes(handler *Handler) *echo.Echo {
	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/start-stop", handler.StartStopSending)
	e.GET("/sent-messages", handler.GetSentMessages)
	return e
}

// StartStopSending handles starting and stopping the message sending scheduler
// @Summary Start or stop the message sending scheduler
// @Description Starts a scheduler to send messages every 10 seconds, or stops it if it's already running
// @Tags Scheduler
// @Produce json
// @Success 200 {object} types.Response
// @Router /start-stop [get]
func (h *Handler) StartStopSending(c echo.Context) error {

	if scheduler != nil {
		scheduler.Clear()
		scheduler = nil

		return c.JSON(http.StatusOK, types.Response{
			Message: constants.SendingStopped,
		})
	}

	scheduler = gocron.NewScheduler()
	scheduler.Every(10).
		Second().
		Do(h.service.StartSending)
	scheduler.Start()

	return c.JSON(http.StatusOK, types.Response{
		Message: constants.SendingStarted,
	})
}

// GetSentMessages retrieves all sent messages
// @Summary Retrieve sent messages
// @Description Retrieves all messages that have been sent
// @Tags Messages
// @Produce json
// @Success 200 {object} types.Response
// @Failure 500 {object} types.Response
// @Router /sent-messages [get]
func (h *Handler) GetSentMessages(c echo.Context) error {
	ctx := c.Request().Context()
	messages, err := h.service.GetSentMessages(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, types.Response{Message: constants.ErrorGettingMessages})
	}
	return c.JSON(http.StatusOK, types.Response{
		Message: constants.MessagesRetrieved,
		Data:    messages,
	})
}

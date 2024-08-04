package client

import (
	"bytes"
	"context"
	"encoding/json"
	"message-system/app/domain"
	"message-system/app/types"
	"message-system/config"
	"net/http"
)

type WebhookClient struct {
	BaseURL    string
	httpClient *http.Client
}

func NewWebhookClient(config *config.WebhookConfig) *WebhookClient {
	return &WebhookClient{
		BaseURL:    config.BaseURL,
		httpClient: &http.Client{},
	}
}

func (c *WebhookClient) SendMessage(ctx context.Context, message *domain.Message) (*types.MessageResponse, error) {
	reqBody, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var responseBody *types.MessageResponse
	err = json.NewDecoder(res.Body).Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, err
}

package types

type Response struct {
	Message string    `json:"message"`
	Data    []Message `json:"data,omitempty"`
}

type MessageResponse struct {
	Message   string `json:"message"`
	MessageID string `json:"messageId"`
}

type Message struct {
	ID      int    `json:"id"`
	To      string `json:"to"`
	Content string `json:"content"`
	IsSent  bool   `json:"isSent"`
}

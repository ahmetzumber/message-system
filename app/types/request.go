package types

type MessageRequest struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

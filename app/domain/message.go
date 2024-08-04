package domain

type Message struct {
	ID      int    `json:"id",bson:"id"`
	To      string `json:"to",bson:"to"`
	Content string `json:"content",bson:"content"`
	IsSent  bool   `json:"isSent",bson:"isSent"`
}

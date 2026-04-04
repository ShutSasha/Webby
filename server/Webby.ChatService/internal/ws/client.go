package ws

import "github.com/google/uuid"

type Client struct {
	UserID uuid.UUID
	ChatID uuid.UUID
	send   chan *BroadcastMessage
}

func NewClient(userId, chatId uuid.UUID) *Client {
	return &Client{
		UserID: userId,
		ChatID: chatId,
		send:   make(chan *BroadcastMessage, 256),
	}
}

func (c *Client) Send() <-chan *BroadcastMessage {
	return c.send
}

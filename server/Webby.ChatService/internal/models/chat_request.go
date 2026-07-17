package models

// CreateChatRequest contains payload data for creating a chat.
type CreateChatRequest struct {
	RoomId *string `json:"roomId"`
}

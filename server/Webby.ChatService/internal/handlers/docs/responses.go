package docs

// ErrorResponse is the standard error response shape.
type ErrorResponse struct {
	Success bool              `json:"success" example:"false"`
	Message string            `json:"message" example:"An error occurred"`
	Errors  map[string]string `json:"errors,omitempty"`
	Data    *string           `json:"data" example:"null"`
}

// SuccessResponse is a generic success wrapper (no typed data).
type SuccessResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Operation completed successfully"`
}

// --- Chat ---

type CreateRequest struct {
	RoomId *string `json:"roomId" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type ChatResponse struct {
	Id        string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	RoomId    *string `json:"roomId,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	CreatedAt string  `json:"createdAt" example:"2026-04-04T12:00:00Z"`
}

type ChatApiResponse struct {
	Success bool          `json:"success" example:"true"`
	Message string        `json:"message" example:"Chat created"`
	Data    *ChatResponse `json:"data,omitempty"`
}

// --- Message ---

type MessageResponse struct {
	Id        string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	SenderId  string  `json:"senderId" example:"550e8400-e29b-41d4-a716-446655440000"`
	ChatId    string  `json:"chatId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Content   string  `json:"content" example:"Hello, world!"`
	IsEdited  bool    `json:"isEdited" example:"false"`
	EditedAt  *string `json:"editedAt,omitempty" example:"2026-04-04T12:00:00Z"`
	CreatedAt string  `json:"createdAt" example:"2026-04-04T12:00:00Z"`
}

type MessageApiResponse struct {
	Success bool             `json:"success" example:"true"`
	Message string           `json:"message" example:"Message sent"`
	Data    *MessageResponse `json:"data,omitempty"`
}

type MessageListPagData struct {
	Items      []MessageResponse `json:"items"`
	Page       int               `json:"page" example:"1"`
	PageSize   int               `json:"pageSize" example:"50"`
	TotalCount int               `json:"totalCount" example:"100"`
}

type MessageListApiResponse struct {
	Success bool                `json:"success" example:"true"`
	Message string              `json:"message" example:"Messages retrieved"`
	Data    *MessageListPagData `json:"data,omitempty"`
}

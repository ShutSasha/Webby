package docs

import "github.com/google/uuid"

type ApiResponse[T any] struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Operation completed successfully"`
	Data    *T     `json:"data,omitempty"`
	Errors  string `json:"errors" example:"null" extensions:"x-nullable"`
}

type PaginatedResponse[T any] struct {
	Items []T `json:"items"`
	Page  int `json:"page" example:"1"`
	Limit int `json:"limit" example:"10"`
	Total int `json:"total" example:"100"`
}

type Error400Response struct {
	Success bool              `json:"success" example:"false"`
	Message string            `json:"message" example:"Validation failed"`
	Errors  map[string]string `json:"errors" example:"name:Name must be at least 2 characters"`
	Data    string            `json:"data" example:"null" extensions:"x-nullable"`
}

type Error401Response struct {
	Success bool              `json:"success" example:"false"`
	Message string            `json:"message" example:"Unauthorized"`
	Errors  map[string]string `json:"errors" example:"message:Missing or invalid token"`
	Data    string            `json:"data" example:"null" extensions:"x-nullable"`
}

type Error403Response struct {
	Success bool              `json:"success" example:"false"`
	Message string            `json:"message" example:"Forbidden"`
	Errors  map[string]string `json:"errors" example:"message:Admin only"`
	Data    string            `json:"data" example:"null" extensions:"x-nullable"`
}

type Error404Response struct {
	Success bool              `json:"success" example:"false"`
	Message string            `json:"message" example:"Not found"`
	Errors  map[string]string `json:"errors" example:"message:Object not found"`
	Data    string            `json:"data" example:"null" extensions:"x-nullable"`
}

type Error409Response struct {
	Success bool              `json:"success" example:"false"`
	Message string            `json:"message" example:"Conflict"`
	Errors  map[string]string `json:"errors" example:"message:Object already exists"`
	Data    string            `json:"data" example:"null" extensions:"x-nullable"`
}

type Error500Response struct {
	Success bool              `json:"success" example:"false"`
	Message string            `json:"message" example:"Internal server error"`
	Errors  map[string]string `json:"errors" example:"message:Database connection lost"`
	Data    string            `json:"data" example:"null" extensions:"x-nullable"`
}

type CategoryResponse struct {
	Id   uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name string    `json:"name" example:"Gaming"`
}

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50" example:"Gaming"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50" example:"Education"`
}

type RoomResponse struct {
	Id         uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name       string    `json:"name" example:"Room Name"`
	CategoryId uuid.UUID `json:"categoryId" example:"123e4567-e89b-12d3-a456-426614174000"`
	IsPrivate  bool      `json:"isPrivate" example:"false"`
	HostId     uuid.UUID `json:"hostId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Thumbnail  string    `json:"thumbnail" example:"https://example.com/thumbnail.jpg"`
	Token      string    `json:"token" example:"abc123def456ghij"`
}

type RoomListItem struct {
	Id         uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name       string    `json:"name" example:"Room Name"`
	CategoryId uuid.UUID `json:"categoryId" example:"123e4567-e89b-12d3-a456-426614174000"`
	IsPrivate  bool      `json:"isPrivate" example:"false"`
	HostId     uuid.UUID `json:"hostId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Thumbnail  string    `json:"thumbnail" example:"https://example.com/thumbnail.jpg"`
	Token      string    `json:"token" example:"abc123def456ghij"`
}

type CreateRoomRequest struct {
	Name       string    `json:"name" validate:"required,min=2,max=100" example:"My Room"`
	CategoryId uuid.UUID `json:"categoryId" validate:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	IsPrivate  bool      `json:"isPrivate" example:"false"`
}

type RoomMemberItem struct {
	UserId     uuid.UUID `json:"userId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username   string    `json:"username" example:"JohnDoe"`
	AvatarUrl  string    `json:"avatarUrl" example:"https://example.com/avatar.jpg"`
	RoomPoints int       `json:"roomPoints" example:"100"`
}

type AddQueueItemRequest struct {
	EntityId   string `json:"entityId" validate:"required,uuid4" example:"550e8400-e29b-41d4-a716-446655440000"`
	EntityType string `json:"entityType" validate:"required,oneof=video playlist" example:"video"`
}

type QueueItemResponse struct {
	Id         uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	EntityId   uuid.UUID `json:"entityId" example:"550e8400-e29b-41d4-a716-446655440000"`
	EntityType string    `json:"entityType" example:"video"`
	IsActive   bool      `json:"isActive" example:"false"`
}

type QueueVideoChild struct {
	Id         uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title      string    `json:"title" example:"Episode 1: Home Alone"`
	Thumbnail  string    `json:"thumbnail" example:"https://example.com/thumb.jpg"`
	VideoUrl   string    `json:"videoUrl" example:"https://example.com/video.mp4"`
	PreviewUrl string    `json:"previewUrl" example:"https://example.com/preview.mp4"`
}

type QueueItemDetailResponse struct {
	Id         uuid.UUID         `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	EntityId   uuid.UUID         `json:"entityId" example:"550e8400-e29b-41d4-a716-446655440000"`
	EntityType string            `json:"entityType" example:"video"`
	Title      string            `json:"title" example:"Fears to Fathom: Ironbark Lookout"`
	Thumbnail  string            `json:"thumbnail" example:"https://example.com/thumb.jpg"`
	VideoUrl   string            `json:"videoUrl" example:"https://example.com/video.mp4"`
	PreviewUrl string            `json:"previewUrl" example:"https://example.com/preview.mp4"`
	IsActive   bool              `json:"isActive" example:"true"`
	IsFolder   bool              `json:"isFolder" example:"false"`
	Children   []QueueVideoChild `json:"children,omitempty"`
}

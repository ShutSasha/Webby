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

// --- Room ---

type RoomResponse struct {
	Id         string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name       string `json:"name" example:"Room Name"`
	CategoryId string `json:"categoryId" example:"123e4567-e89b-12d3-a456-426614174000"`
	IsPrivate  bool   `json:"isPrivate" example:"false"`
	HostId     string `json:"hostId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Thumbnail  string `json:"thumbnail" example:"https://example.com/thumbnail.jpg"`
	Token      string `json:"token" example:"abc123def456ghij"`
}

type RoomApiResponse struct {
	Success bool          `json:"success" example:"true"`
	Message string        `json:"message" example:"Room created"`
	Data    *RoomResponse `json:"data,omitempty"`
}

type RoomListItem struct {
	Id         string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name       string `json:"name" example:"Room Name"`
	CategoryId string `json:"categoryId" example:"123e4567-e89b-12d3-a456-426614174000"`
	IsPrivate  bool   `json:"isPrivate" example:"false"`
	HostId     string `json:"hostId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Thumbnail  string `json:"thumbnail" example:"https://example.com/thumbnail.jpg"`
	Token      string `json:"token" example:"abc123def456ghij"`
}

type RoomListApiResponse struct {
	Success bool             `json:"success" example:"true"`
	Message string           `json:"message" example:"Rooms retrieved"`
	Data    *RoomListPagData `json:"data,omitempty"`
}

type RoomListPagData struct {
	Items []RoomListItem `json:"items"`
	Page  int            `json:"page" example:"1"`
	Limit int            `json:"pageSize" example:"10"`
	Total int            `json:"totalCount" example:"100"`
}

// --- Category ---

type CategoryResponse struct {
	Id   string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name string `json:"name" example:"Gaming"`
}

type CategoryApiResponse struct {
	Success bool              `json:"success" example:"true"`
	Message string            `json:"message" example:"Category created"`
	Data    *CategoryResponse `json:"data,omitempty"`
}

type CategoryListApiResponse struct {
	Success bool                 `json:"success" example:"true"`
	Message string               `json:"message" example:"Categories retrieved"`
	Data    *CategoryListPagData `json:"data,omitempty"`
}

type CategoryListPagData struct {
	Items []CategoryResponse `json:"items"`
	Page  int                `json:"page" example:"1"`
	Limit int                `json:"pageSize" example:"10"`
	Total int                `json:"totalCount" example:"100"`
}

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50" example:"Gaming"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50" example:"Education"`
}

// --- Room Members ---

type RoomMemberItem struct {
	UserId     string `json:"userId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username   string `json:"username" example:"JohnDoe"`
	AvatarUrl  string `json:"avatarUrl" example:"https://example.com/avatar.jpg"`
	RoomPoints int    `json:"roomPoints" example:"100"`
}

type MemberListApiResponse struct {
	Success bool               `json:"success" example:"true"`
	Message string             `json:"message" example:"Members retrieved"`
	Data    *MemberListPagData `json:"data,omitempty"`
}

type MemberListPagData struct {
	Items []RoomMemberItem `json:"items"`
	Page  int              `json:"page" example:"1"`
	Limit int              `json:"pageSize" example:"10"`
	Total int              `json:"totalCount" example:"100"`
}

type AddMemberRequest struct {
	UserId string `json:"userId" validate:"required,uuid4" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type GetByTokenRequest struct {
	Token string `json:"token" validate:"required,notblank" example:"abc123def456ghij"`
}

// --- Queue ---

type AddQueueItemRequest struct {
	EntityId   string `json:"entityId" validate:"required,uuid4" example:"550e8400-e29b-41d4-a716-446655440000"`
	EntityType string `json:"entityType" validate:"required,oneof=video playlist" example:"video"`
}

type QueueItemResponse struct {
	Id         string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	EntityId   string `json:"entityId" example:"550e8400-e29b-41d4-a716-446655440000"`
	EntityType string `json:"entityType" example:"video"`
	IsActive   bool   `json:"isActive" example:"false"`
}

type QueueItemApiResponse struct {
	Success bool               `json:"success" example:"true"`
	Message string             `json:"message" example:"Item added to queue"`
	Data    *QueueItemResponse `json:"data,omitempty"`
}

type QueueVideoChild struct {
	Id        string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title     string `json:"title" example:"Episode 1: Home Alone"`
	Thumbnail string `json:"thumbnail" example:"https://example.com/thumb.jpg"`
	VideoUrl  string `json:"videoUrl" example:"https://example.com/video.mp4"`
}

type QueueItemDetailResponse struct {
	Id            string            `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	EntityId      string            `json:"entityId" example:"550e8400-e29b-41d4-a716-446655440000"`
	EntityType    string            `json:"entityType" example:"video"`
	Title         string            `json:"title" example:"Fears to Fathom: Ironbark Lookout"`
	Thumbnail     string            `json:"thumbnail" example:"https://example.com/thumb.jpg"`
	VideoUrl      string            `json:"videoUrl" example:"https://example.com/video.mp4"`
	IsActive      bool              `json:"isActive" example:"true"`
	IsFolder      bool              `json:"isFolder" example:"false"`
	TotalChildren int               `json:"totalChildren" example:"5"`
	Children      []QueueVideoChild `json:"children,omitempty"`
}

type QueueListApiResponse struct {
	Success bool              `json:"success" example:"true"`
	Message string            `json:"message" example:"Queue retrieved"`
	Data    *QueueListPagData `json:"data,omitempty"`
}

type QueueListPagData struct {
	Items []QueueItemDetailResponse `json:"items"`
	Page  int                       `json:"page" example:"1"`
	Limit int                       `json:"pageSize" example:"10"`
	Total int                       `json:"totalCount" example:"100"`
}

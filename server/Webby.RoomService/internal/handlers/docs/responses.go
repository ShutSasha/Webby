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
	Id           string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name         string `json:"name" example:"Room Name"`
	CategoryName string `json:"categoryName" example:"Gaming"`
	IsPrivate    bool   `json:"isPrivate" example:"false"`
	HostId       string `json:"hostId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Thumbnail    string `json:"thumbnail" example:"https://example.com/thumbnail.jpg"`
	ChatId       string `json:"chatId,omitempty" example:"550e8400-e29b-41d4-a716-446655440000" description:"Chat ID associated with this room (from ChatService)"`
}

type RoomApiResponse struct {
	Success bool          `json:"success" example:"true"`
	Message string        `json:"message" example:"Room created"`
	Data    *RoomResponse `json:"data,omitempty"`
}

type RoomListItem struct {
	Id            string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name          string `json:"name" example:"Room Name"`
	CategoryName  string `json:"categoryName" example:"Gaming"`
	IsPrivate     bool   `json:"isPrivate" example:"false"`
	HostId        string `json:"hostId" example:"550e8400-e29b-41d4-a716-446655440000"`
	HostUsername  string `json:"hostUsername" example:"JohnDoe"`
	HostAvatarUrl string `json:"hostAvatarUrl" example:"https://example.com/avatar.jpg"`
	Thumbnail     string `json:"thumbnail" example:"https://example.com/thumbnail.jpg"`
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
	Items []string `json:"items" example:"Gaming,Education,Music"`
	Page  int      `json:"page" example:"1"`
	Limit int      `json:"pageSize" example:"10"`
	Total int      `json:"totalCount" example:"100"`
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

type AddMembersRequest struct {
	UserIds []string `json:"userIds" validate:"required,min=1,dive,uuid"`
}

// --- Queue ---

type AddQueueItemRequest struct {
	EntityId   string `json:"entityId" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	EntityType string `json:"entityType" validate:"required,oneof=video playlist" example:"video"`
}

type QueueItemResponse struct {
	Id         string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	EntityId   string `json:"entityId" example:"550e8400-e29b-41d4-a716-446655440000"`
	EntityType string `json:"entityType" example:"video" description:"Type of media entity: 'video' or 'playlist'"`
	IsActive   bool   `json:"isActive" example:"false" description:"Whether this item is currently playing"`
	Position   int    `json:"position" example:"3" description:"Ordinal position in the room queue (ascending)"`
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
	EntityType    string            `json:"entityType" example:"video" description:"Type of media entity: 'video' or 'playlist'"`
	Title         string            `json:"title" example:"Fears to Fathom: Ironbark Lookout"`
	Thumbnail     string            `json:"thumbnail" example:"https://example.com/thumb.jpg"`
	VideoUrl      string            `json:"videoUrl" example:"https://example.com/video.mp4" description:"Direct video URL (empty for playlists)"`
	IsActive      bool              `json:"isActive" example:"true" description:"Whether this item is currently playing"`
	IsFolder      bool              `json:"isFolder" example:"false" description:"True if item is a playlist containing child videos"`
	Position      int               `json:"position" example:"1" description:"Ordinal position in the room queue (ascending)"`
	TotalChildren int               `json:"totalChildren" example:"5" description:"Total number of child videos (0 for single videos)"`
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

type CreateVoteRequest struct {
	Type            string              `json:"type" validate:"required,oneof=poll next_video" example:"poll" description:"Vote type: 'poll' — text poll created by host; 'next_video' — vote for the next queue item to play"`
	VoteText        string              `json:"voteText" validate:"required,min=1,max=500" example:"What should we watch next?" description:"Question or description shown to voters"`
	DurationSeconds int                 `json:"durationSeconds" validate:"required,min=60,max=604800" example:"3600" description:"Vote duration in seconds (60 = 1 min, 604800 = 7 days). After expiry the vote is closed and winner is computed"`
	Choices         []CreateChoiceInput `json:"choices" validate:"required,min=2,max=20" description:"List of choices (2–20). For 'next_video' type each choice should reference a queueItemId"`
}

type CreateChoiceInput struct {
	Name        string  `json:"name" example:"Option A" description:"Display label for the choice"`
	IsCorrect   bool    `json:"isCorrect" example:"false" description:"Mark as the correct answer (only meaningful for 'poll' type)"`
	QueueItemId *string `json:"queueItemId" example:"550e8400-e29b-41d4-a716-446655440000" description:"Queue item UUID this choice refers to (required for 'next_video' type, null for 'poll')"`
}

type CastVoteRequest struct {
	ChoiceId string `json:"choiceId" validate:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000" description:"UUID of the choice to vote for"`
}

type VoteChoiceDetailResponse struct {
	Id          string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string  `json:"name" example:"Option A" description:"Display label of this choice"`
	Votes       int     `json:"votes" example:"5" description:"Number of users who selected this choice"`
	Percentage  float64 `json:"percentage" example:"62.5" description:"Rounded percentage of total votes (0–100)"`
	IsCorrect   bool    `json:"isCorrect" example:"false" description:"Whether this choice is marked as the correct answer (poll type only)"`
	QueueItemId *string `json:"queueItemId" example:"550e8400-e29b-41d4-a716-446655440000" description:"Associated queue item UUID (next_video type only, null for poll)"`
}

type VoteDetailResponse struct {
	Id                string                     `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	RoomId            string                     `json:"roomId" example:"550e8400-e29b-41d4-a716-446655440000"`
	Type              string                     `json:"type" example:"poll" description:"Vote type: 'poll' (text poll) or 'next_video' (queue item selection). For 'next_video', the winner is moved to the top of the queue"`
	VoteText          string                     `json:"voteText" example:"What should we watch next?" description:"Question displayed to room members"`
	CreatedAt         string                     `json:"createdAt" example:"2026-03-30T12:00:00Z"`
	DurationSeconds   int                        `json:"durationSeconds" example:"3600" description:"Duration of the vote in seconds"`
	ExpiresAt         string                     `json:"expiresAt" example:"2026-03-30T13:00:00Z" description:"Computed expiry time (createdAt + durationSeconds)"`
	IsExpired         bool                       `json:"isExpired" example:"false" description:"Whether the vote has passed its expiry time"`
	TotalVotes        int                        `json:"totalVotes" example:"8" description:"Sum of all votes across all choices"`
	UserVotedChoiceId *string                    `json:"userVotedChoiceId" example:"550e8400-e29b-41d4-a716-446655440000" description:"Choice UUID the current user voted for (null if not voted)"`
	WinnerId          *string                    `json:"winnerId" example:"550e8400-e29b-41d4-a716-446655440000" description:"Choice UUID of the winner (only set for expired 'next_video' votes). In a tie the host's vote counts double; if still tied the first candidate wins"`
	Choices           []VoteChoiceDetailResponse `json:"choices"`
}

type VoteApiResponse struct {
	Success bool                `json:"success" example:"true"`
	Message string              `json:"message" example:"Vote created"`
	Data    *VoteDetailResponse `json:"data,omitempty"`
}

type VoteListApiResponse struct {
	Success bool                  `json:"success" example:"true"`
	Message string                `json:"message" example:"Votes retrieved"`
	Data    *[]VoteDetailResponse `json:"data,omitempty"`
}

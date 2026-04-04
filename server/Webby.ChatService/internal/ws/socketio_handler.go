package ws

import (
	"context"
	"errors"
	"log/slog"
	"time"
	"webby-chat/internal/models"
	"webby-chat/internal/services"
	"webby-chat/pkg/http/middleware/auth"

	"github.com/google/uuid"
	socketio "github.com/googollee/go-socket.io"
)

type MessageResponse struct {
	Id        uuid.UUID  `json:"id"`
	SenderId  uuid.UUID  `json:"senderId"`
	ChatId    uuid.UUID  `json:"chatId"`
	Content   string     `json:"content"`
	IsEdited  bool       `json:"isEdited"`
	EditedAt  *time.Time `json:"editedAt,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
}

type DeletedMessageResponse struct {
	Id     uuid.UUID `json:"id"`
	ChatId uuid.UUID `json:"chatId"`
}

type PaginatedMessages struct {
	Items      []MessageResponse `json:"items"`
	Page       int               `json:"page"`
	PageSize   int               `json:"pageSize"`
	TotalCount int64             `json:"totalCount"`
}

func toMessageResponse(msg *models.Message) MessageResponse {
	return MessageResponse{
		Id:        msg.Id,
		SenderId:  msg.SenderId,
		ChatId:    msg.ChatId,
		Content:   msg.Content,
		IsEdited:  msg.IsEdited,
		EditedAt:  msg.EditedAt,
		CreatedAt: msg.CreatedAt,
	}
}

func SetupSocketIO(
	logger *slog.Logger,
	jwtSecret []byte,
	hubManager *HubManager,
	chatService *services.ChatService,
	messageService *services.MessageService,
) *socketio.Server {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		u := s.URL()
		token := u.Query().Get("token")
		if token == "" {
			return errors.New("authorization failed") 
		}

		userIdStr, err := auth.ParseUserIDFromToken(jwtSecret, token)
		if err != nil {
			logger.Warn("socket.io connection rejected: invalid token", slog.String("error", err.Error()))
			s.Close()
			return nil
		}

		userId, err := uuid.Parse(userIdStr)
		if err != nil {
			logger.Warn("socket.io connection rejected: invalid userId", slog.String("error", err.Error()))
			s.Close()
			return nil
		}

		s.SetContext(userId)
		logger.Info("socket.io client connected", slog.String("userId", userId.String()))
		return nil
	})

	server.OnEvent("/", "chat:join", func(s socketio.Conn, data map[string]string) string {
		userId, ok := s.Context().(uuid.UUID)
		if !ok {
			return `{"error":"unauthorized"}`
		}

		chatIdStr, ok := data["chatId"]
		if !ok || chatIdStr == "" {
			return `{"error":"chatId is required"}`
		}
		chatId, err := uuid.Parse(chatIdStr)
		if err != nil {
			return `{"error":"invalid chatId"}`
		}

		isMember, err := chatService.IsMember(context.Background(), chatId, userId)
		if err != nil || !isMember {
			return `{"error":"not a member of this chat"}`
		}

		hub := hubManager.GetOrCreateHub(chatId)
		client := NewClient(userId, chatId)
		hub.Register(client)

		s.Join(chatId.String())

		s.SetContext(clientContext{UserID: userId, ChatID: chatId, Client: client, Hub: hub})

		logger.Info("user joined chat",
			slog.String("userId", userId.String()),
			slog.String("chatId", chatId.String()),
		)

		return `{"success":true}`
	})

	server.OnEvent("/", "messages:list", func(s socketio.Conn, data map[string]any) string {
		ctx, ok := s.Context().(clientContext)
		if !ok {
			return `{"error":"unauthorized"}`
		}

		chatIdStr, _ := data["chatId"].(string)
		chatId, err := uuid.Parse(chatIdStr)
		if err != nil {
			return `{"error":"invalid chatId"}`
		}

		page := intFromAny(data["page"], 1)
		limit := intFromAny(data["limit"], 50)

		messages, total, err := messageService.List(context.Background(), chatId, ctx.UserID, page, limit)
		if err != nil {
			logger.Error("messages:list error", slog.String("error", err.Error()))
			return `{"error":"failed to list messages"}`
		}

		items := make([]MessageResponse, len(messages))
		for i, msg := range messages {
			items[i] = toMessageResponse(&msg)
		}

		resp := PaginatedMessages{
			Items:      items,
			Page:       page,
			PageSize:   limit,
			TotalCount: total,
		}

		s.Emit("messages:list:response", resp)
		return ""
	})

	server.OnEvent("/", "messages:send", func(s socketio.Conn, data map[string]string) string {
		ctx, ok := s.Context().(clientContext)
		if !ok {
			return `{"error":"unauthorized"}`
		}

		chatIdStr, _ := data["chatId"]
		chatId, err := uuid.Parse(chatIdStr)
		if err != nil {
			return `{"error":"invalid chatId"}`
		}

		content, _ := data["content"]

		msg, err := messageService.Send(context.Background(), chatId, ctx.UserID, content)
		if err != nil {
			logger.Error("messages:send error", slog.String("error", err.Error()))
			return `{"error":"` + err.Error() + `"}`
		}

		resp := toMessageResponse(msg)

		server.BroadcastToRoom("/", chatId.String(), "messages:new", resp)

		return ""
	})

	server.OnEvent("/", "messages:edit", func(s socketio.Conn, data map[string]string) string {
		ctx, ok := s.Context().(clientContext)
		if !ok {
			return `{"error":"unauthorized"}`
		}

		messageIdStr, _ := data["messageId"]
		messageId, err := uuid.Parse(messageIdStr)
		if err != nil {
			return `{"error":"invalid messageId"}`
		}

		content, _ := data["content"]

		msg, err := messageService.Edit(context.Background(), messageId, ctx.UserID, content)
		if err != nil {
			logger.Error("messages:edit error", slog.String("error", err.Error()))
			return `{"error":"` + err.Error() + `"}`
		}

		resp := toMessageResponse(msg)

		server.BroadcastToRoom("/", msg.ChatId.String(), "messages:edited", resp)

		return ""
	})

	server.OnEvent("/", "messages:delete", func(s socketio.Conn, data map[string]string) string {
		ctx, ok := s.Context().(clientContext)
		if !ok {
			return `{"error":"unauthorized"}`
		}

		messageIdStr, _ := data["messageId"]

		messageId, err := uuid.Parse(messageIdStr)
		if err != nil {
			return `{"error":"invalid messageId"}`
		}

		msg, err := messageService.Delete(context.Background(), messageId, ctx.UserID)
		if err != nil {
			logger.Error("messages:delete error", slog.String("error", err.Error()))
			return `{"error":"` + err.Error() + `"}`
		}

		resp := DeletedMessageResponse{
			Id:     messageId,
			ChatId: msg.ChatId,
		}

		server.BroadcastToRoom("/", msg.ChatId.String(), "messages:deleted", resp)

		return ""
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		if ctx, ok := s.Context().(clientContext); ok {
			ctx.Hub.Unregister(ctx.Client)
			logger.Info("socket.io client disconnected",
				slog.String("userId", ctx.UserID.String()),
				slog.String("chatId", ctx.ChatID.String()),
				slog.String("reason", reason),
			)
		}
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		logger.Error("socket.io error", slog.String("error", e.Error()))
	})

	return server
}

type clientContext struct {
	UserID uuid.UUID
	ChatID uuid.UUID
	Client *Client
	Hub    *Hub
}

func intFromAny(v any, fallback int) int {
	switch val := v.(type) {
	case float64:
		return int(val)
	case int:
		return val
	case int64:
		return int(val)
	default:
		return fallback
	}
}

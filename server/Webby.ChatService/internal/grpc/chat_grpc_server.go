package grpc

import (
	"context"
	"encoding/json"
	"log/slog"
	"webby/chat-service/internal/grpc/chatpb"
	"webby/chat-service/internal/models"
	"webby/chat-service/internal/services"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChatService interface {
	Create(ctx context.Context, roomID *uuid.UUID) (*models.Chat, error)
	GetByRoomID(ctx context.Context, roomId uuid.UUID) (*models.Chat, error)
	GetChatIDByRoomID(ctx context.Context, roomID, userID uuid.UUID) (uuid.UUID, error)
	EnsureMember(ctx context.Context, chatId, userId uuid.UUID) error
	GetByID(ctx context.Context, chatID uuid.UUID) (*models.Chat, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, channel string, payload any) error
}

type RedisPublisher struct {
	client *redis.Client
}

func NewRedisPublisher(client *redis.Client) *RedisPublisher {
	return &RedisPublisher{client: client}
}

func (p *RedisPublisher) Publish(ctx context.Context, channel string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return p.client.Publish(ctx, channel, string(data)).Err()
}

type ChatGrpcServer struct {
	chatpb.UnimplementedChatGrpcServiceServer
	chatService    ChatService
	messageService *services.MessageService
	publisher      EventPublisher
	logger         *slog.Logger
}

func NewChatGrpcServer(chatService ChatService, messageService *services.MessageService, publisher EventPublisher, logger *slog.Logger) *ChatGrpcServer {
	return &ChatGrpcServer{
		chatService:    chatService,
		messageService: messageService,
		publisher:      publisher,
		logger:         logger,
	}
}

func (s *ChatGrpcServer) CreateChat(ctx context.Context, req *chatpb.CreateChatRequest) (*chatpb.ChatResponse, error) {
	roomID, err := uuid.Parse(req.RoomID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "wrong roomID")
	}

	chat, err := s.chatService.Create(ctx, &roomID)
	if err != nil {
		s.logger.Error("gRPC CreateChat failed", slog.String("error", err.Error()))
		return nil, status.Errorf(codes.Internal, "failed to create chat: %v", err)
	}

	return &chatpb.ChatResponse{
		Id:     chat.ID.String(),
		RoomID: chat.RoomID.String(),
	}, nil
}

func (s *ChatGrpcServer) GetChatByRoomID(
	ctx context.Context,
	req *chatpb.GetChatByRoomIdRequest,
) (*chatpb.ChatResponse, error) {
	const op = "grpc.ChatGrpcServer.GetChatByRoomId"

	roomID, err := uuid.Parse(req.GetRoomID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s: invalid room_id: %v", op, err)
	}

	chat, err := s.chatService.GetByRoomID(ctx, roomID)
	if err != nil {
		s.logger.Error("gRPC execution failed",
			slog.String("op", op),
			slog.String("roomID", roomID.String()),
			slog.String("error", err.Error()),
		)

		return nil, status.Errorf(codes.NotFound, "%s: chat not found: %v", op, err)
	}

	return &chatpb.ChatResponse{
		Id:     chat.ID.String(),
		RoomID: chat.RoomID.String(),
	}, nil
}

func (s *ChatGrpcServer) AddChatMember(ctx context.Context, req *chatpb.AddChatMemberRequest) (*chatpb.AddChatMemberResponse, error) {
	chatID, err := uuid.Parse(req.GetChatID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid chat_id: %v", err)
	}

	userID, err := uuid.Parse(req.GetUserID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	if err := s.chatService.EnsureMember(ctx, chatID, userID); err != nil {
		s.logger.Error("gRPC AddChatMember failed",
			slog.String("chatID", chatID.String()),
			slog.String("userID", userID.String()),
			slog.String("error", err.Error()),
		)
		return nil, status.Errorf(codes.Internal, "failed to add chat member: %v", err)
	}

	return &chatpb.AddChatMemberResponse{}, nil
}

func (s *ChatGrpcServer) SaveMessage(
	ctx context.Context,
	req *chatpb.SaveMessageRequest,
) (*chatpb.SaveMessageResponse, error) {
	chatID, err := uuid.Parse(req.GetChatID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid chat_id: %v", err)
	}

	userID, err := uuid.Parse(req.GetChatID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	chat, err := s.chatService.GetByID(ctx, chatID)
	if err != nil {
		s.logger.Error("gRPC SaveMessage failed to resolve chat",
			slog.String("chatID", chatID.String()),
			slog.String("error", err.Error()),
		)
		return nil, status.Errorf(codes.Internal, "failed to resolve chat: %v", err)
	}

	msg, err := s.messageService.Send(ctx, chat.ID, userID, req.GetContent())
	if err != nil {
		s.logger.Error("gRPC SaveMessage failed to save message",
			slog.String("chatID", chatID.String()),
			slog.String("userId", userID.String()),
			slog.String("error", err.Error()),
		)
		return nil, status.Errorf(codes.Internal, "failed to save message: %v", err)
	}

	// TODO: MOVE TO SERVICE LAYER
	envelope := struct {
		Type    string `json:"type"`
		Payload any    `json:"payload"`
	}{
		Type: "MESSAGE_CREATED",
		Payload: map[string]any{
			"id":        msg.ID.String(),
			"senderId":  msg.SenderID.String(),
			"chatId":    msg.ChatID.String(),
			"content":   msg.Content,
			"isEdited":  msg.IsEdited,
			"editedAt":  msg.EditedAt,
			"createdAt": msg.CreatedAt,
		},
	}

	if err := s.publisher.Publish(ctx, "chat:"+chatID.String(), envelope); err != nil {
		s.logger.Error("gRPC SaveMessage failed to publish redis event",
			slog.String("chatID", chatID.String()),
			slog.String("messageID", msg.ID.String()),
			slog.String("error", err.Error()),
		)
		return nil, status.Errorf(codes.Internal, "failed to publish event: %v", err)
	}

	return &chatpb.SaveMessageResponse{MessageID: msg.ID.String()}, nil
}

func (s *ChatGrpcServer) GetChatIDByRoomID(
	ctx context.Context,
	req *chatpb.GetChatIDByRoomIDRequest,
) (*chatpb.ChatIDByRoomIDResponse, error) {
	const op = "grpc.ChatGrpcServer.GetChatIDByRoomID"

	roomID, err := uuid.Parse(req.GetRoomID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid room_id: %v", err)
	}

	userID, err := uuid.Parse(req.GetUserID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	id, err := s.chatService.GetChatIDByRoomID(ctx, roomID, userID)
	if err != nil {
		s.logger.Error("gRPC GetChatIDByRoomID failed",
			slog.String("roomID", roomID.String()),
			slog.String("userID", userID.String()),
			slog.String("error", err.Error()),
		)
		return nil, status.Errorf(codes.Internal, "failed to get chat id: %v", err)
	}

	return &chatpb.ChatIDByRoomIDResponse{ChatID: id.String()}, nil
}

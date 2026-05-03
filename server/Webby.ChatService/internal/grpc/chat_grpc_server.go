package grpcserver

import (
	"context"
	"encoding/json"
	"log/slog"
	"webby-chat/internal/grpc/chatpb"
	"webby-chat/internal/services"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
	chatService    *services.ChatService
	messageService *services.MessageService
	publisher      EventPublisher
	logger         *slog.Logger
}

func NewChatGrpcServer(chatService *services.ChatService, messageService *services.MessageService, publisher EventPublisher, logger *slog.Logger) *ChatGrpcServer {
	return &ChatGrpcServer{
		chatService:    chatService,
		messageService: messageService,
		publisher:      publisher,
		logger:         logger,
	}
}

func (s *ChatGrpcServer) CreateChat(ctx context.Context, req *chatpb.CreateChatRequest) (*chatpb.ChatResponse, error) {
	roomId, err := uuid.Parse(req.GetRoomId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid room_id: %v", err)
	}

	chat, err := s.chatService.CreateForRoom(ctx, roomId)
	if err != nil {
		s.logger.Error("gRPC CreateChat failed", slog.String("error", err.Error()))
		return nil, status.Errorf(codes.Internal, "failed to create chat: %v", err)
	}

	roomIdStr := ""
	if chat.RoomId != nil {
		roomIdStr = chat.RoomId.String()
	}

	return &chatpb.ChatResponse{
		Id:     chat.Id.String(),
		RoomId: roomIdStr,
	}, nil
}

func (s *ChatGrpcServer) GetChatByRoomId(ctx context.Context, req *chatpb.GetChatByRoomIdRequest) (*chatpb.ChatResponse, error) {
	roomId, err := uuid.Parse(req.GetRoomId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid room_id: %v", err)
	}

	chat, err := s.chatService.GetByRoomId(ctx, roomId)
	if err != nil {
		s.logger.Error("gRPC GetChatByRoomId failed",
			slog.String("roomId", roomId.String()),
			slog.String("error", err.Error()),
		)
		return nil, status.Errorf(codes.NotFound, "chat not found for room %s", roomId)
	}

	roomIdStr := ""
	if chat.RoomId != nil {
		roomIdStr = chat.RoomId.String()
	}

	return &chatpb.ChatResponse{
		Id:     chat.Id.String(),
		RoomId: roomIdStr,
	}, nil
}

func (s *ChatGrpcServer) AddChatMember(ctx context.Context, req *chatpb.AddChatMemberRequest) (*chatpb.AddChatMemberResponse, error) {
	chatId, err := uuid.Parse(req.GetChatId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid chat_id: %v", err)
	}

	userId, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	if err := s.chatService.EnsureMember(ctx, chatId, userId); err != nil {
		s.logger.Error("gRPC AddChatMember failed",
			slog.String("chatId", chatId.String()),
			slog.String("userId", userId.String()),
			slog.String("error", err.Error()),
		)
		return nil, status.Errorf(codes.Internal, "failed to add chat member: %v", err)
	}

	return &chatpb.AddChatMemberResponse{}, nil
}

func (s *ChatGrpcServer) SaveMessage(ctx context.Context, req *chatpb.SaveMessageRequest) (*chatpb.SaveMessageResponse, error) {
	roomId, err := uuid.Parse(req.GetRoomId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid room_id: %v", err)
	}

	userId, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	chat, err := s.chatService.GetOrCreateByRoomId(ctx, roomId)
	if err != nil {
		s.logger.Error("gRPC SaveMessage failed to resolve chat",
			slog.String("roomId", roomId.String()),
			slog.String("error", err.Error()),
		)
		return nil, status.Errorf(codes.Internal, "failed to resolve chat: %v", err)
	}

	msg, err := s.messageService.Send(ctx, chat.Id, userId, req.GetContent())
	if err != nil {
		s.logger.Error("gRPC SaveMessage failed to save message",
			slog.String("roomId", roomId.String()),
			slog.String("userId", userId.String()),
			slog.String("error", err.Error()),
		)
		return nil, status.Errorf(codes.Internal, "failed to save message: %v", err)
	}

	envelope := struct {
		Type    string `json:"type"`
		Payload any    `json:"payload"`
	}{
		Type: "MESSAGE_CREATED",
		Payload: map[string]any{
			"id":        msg.Id.String(),
			"senderId":  msg.SenderId.String(),
			"chatId":    msg.ChatId.String(),
			"content":   msg.Content,
			"isEdited":  msg.IsEdited,
			"editedAt":  msg.EditedAt,
			"createdAt": msg.CreatedAt,
		},
	}

	if err := s.publisher.Publish(ctx, "room:"+roomId.String(), envelope); err != nil {
		s.logger.Error("gRPC SaveMessage failed to publish redis event",
			slog.String("roomId", roomId.String()),
			slog.String("messageId", msg.Id.String()),
			slog.String("error", err.Error()),
		)
		return nil, status.Errorf(codes.Internal, "failed to publish event: %v", err)
	}

	return &chatpb.SaveMessageResponse{MessageId: msg.Id.String()}, nil
}

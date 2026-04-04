package grpcserver

import (
	"context"
	"log/slog"
	"webby-chat/internal/grpc/chatpb"
	"webby-chat/internal/services"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChatGrpcServer struct {
	chatpb.UnimplementedChatGrpcServiceServer
	chatService *services.ChatService
	logger      *slog.Logger
}

func NewChatGrpcServer(chatService *services.ChatService, logger *slog.Logger) *ChatGrpcServer {
	return &ChatGrpcServer{
		chatService: chatService,
		logger:      logger,
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

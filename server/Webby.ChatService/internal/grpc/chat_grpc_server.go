package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"webby/chat-service/internal/grpc/chatpb"
	"webby/chat-service/internal/models"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type chatService interface {
	CreateForRoom(ctx context.Context, roomID uuid.UUID) (*models.Chat, error)
	GetByRoomID(ctx context.Context, roomID uuid.UUID) (*models.Chat, error)
	GetChatIDByRoomID(ctx context.Context, roomID uuid.UUID) (uuid.UUID, error)
	RoomIDByChatIDBatch(ctx context.Context, chatIDs []string) (map[string]uuid.UUID, error)
	IsChatRelatedToRoom(ctx context.Context, chatID uuid.UUID) (bool, error)
}

type chatMemberEnsurer interface {
	EnsureMember(ctx context.Context, chatId, userId uuid.UUID) error
}

type chatGrpcServer struct {
	chatpb.UnimplementedChatGrpcServiceServer
	chatService       chatService
	chatMemberEnsurer chatMemberEnsurer
	logger            *slog.Logger
}

func NewChatGrpcServer(chatService chatService, chatMemberEnsurer chatMemberEnsurer, logger *slog.Logger) *chatGrpcServer {
	return &chatGrpcServer{
		chatService:       chatService,
		chatMemberEnsurer: chatMemberEnsurer,
		logger:            logger,
	}
}

func (s *chatGrpcServer) CreateChat(ctx context.Context, req *chatpb.CreateChatRequest) (*chatpb.ChatResponse, error) {
	const op = "grpc.chatGrpcServer.CreateChat"
	log := s.logger.With("op", op)

	roomID, err := uuid.Parse(req.RoomID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s: wrong roomID", op)
	}

	chat, err := s.chatService.CreateForRoom(ctx, roomID)
	if err != nil {
		log.Error("gRPC CreateChat failed", slog.String("error", err.Error()))
		return nil, status.Errorf(codes.Internal, "%s: %v", op, err)
	}

	return &chatpb.ChatResponse{
		Id:     chat.ID.String(),
		RoomID: chat.RoomID.String(),
	}, nil
}

func (s *chatGrpcServer) GetChatByRoomID(
	ctx context.Context,
	req *chatpb.GetChatByRoomIdRequest,
) (*chatpb.ChatResponse, error) {
	const op = "grpc.ChatGrpcServer.GetChatByRoomID"
	log := s.logger.With("op", op)

	roomID, err := uuid.Parse(req.GetRoomID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s: invalid room_id: %v", op, err)
	}

	chat, err := s.chatService.GetByRoomID(ctx, roomID)
	if err != nil {
		log.Error("gRPC execution failed",
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

func (s *chatGrpcServer) AddChatMember(ctx context.Context, req *chatpb.AddChatMemberRequest) (*chatpb.AddChatMemberResponse, error) {
	const op = "grpc.chatGrpcServer.AddChatMember"
	log := s.logger.With("op", op)

	chatID, err := uuid.Parse(req.GetChatID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s: invalid chat_id: %v", op, err)
	}

	userID, err := uuid.Parse(req.GetUserID())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s: invalid user_id: %v", op, err)
	}

	if err := s.chatMemberEnsurer.EnsureMember(ctx, chatID, userID); err != nil {
		log.Error("gRPC AddChatMember failed",
			slog.String("chatID", chatID.String()),
			slog.String("userID", userID.String()),
			slog.String("error", err.Error()),
		)
		return nil, status.Errorf(codes.Internal, "%s: %v", op, err)
	}

	return &chatpb.AddChatMemberResponse{}, nil
}

func (s *chatGrpcServer) GetChatIDByRoomID(ctx context.Context, req *chatpb.GetChatIDByRoomIDRequest) (*chatpb.ChatIDByRoomIDResponse, error) {
	const op = "grpc.ChatGrpcServer.GetChatIDByRoomID"
	log := s.logger.With("op", op)

	roomID, err := uuid.Parse(req.GetRoomID())
	if err != nil {
		log.Error("gRPC GetChatIDByRoomID failed",
			slog.String("roomID", roomID.String()),
			slog.String("error", err.Error()),
		)
		return nil, status.Errorf(codes.InvalidArgument, "%s: invalid room_id: %v", op, err)
	}

	id, err := s.chatService.GetChatIDByRoomID(ctx, roomID)
	if err != nil {
		log.Error("gRPC GetChatIDByRoomID failed",
			slog.String("roomID", roomID.String()),
			slog.String("error", err.Error()),
		)
		return nil, status.Errorf(codes.Internal, "%s: %v", op, err)
	}

	return &chatpb.ChatIDByRoomIDResponse{ChatID: id.String()}, nil
}

func (s *chatGrpcServer) RoomIDByChatIDBatch(ctx context.Context, req *chatpb.RoomIDByChatIDBatchRequest) (*chatpb.RoomIDByChatIDBatchResponse, error) {
	const op = "grpc.ChatGrpcServer.RoomIDByChatIDBatch"

	chatIDroomIDMap, err := s.chatService.RoomIDByChatIDBatch(ctx, req.ChatIDs)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s: %v", op, err)
	}

	result := make(map[string]string)
	for chatID, roomID := range chatIDroomIDMap {
		result[chatID] = roomID.String()
	}

	return &chatpb.RoomIDByChatIDBatchResponse{
		ChatIDroomIDMap: result,
	}, nil
}

func (s *chatGrpcServer) IsChatRelatedToRoom(ctx context.Context, req *chatpb.IsChatRelatedToRoomRequest) (*chatpb.IsChatRelatedToRoomResponse, error) {
	const op = "grpc.chatGrpcServer.IsChatRelatedToRoom"

	chatID, err := uuid.Parse(req.ChatID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	isRelated, err := s.chatService.IsChatRelatedToRoom(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &chatpb.IsChatRelatedToRoomResponse{
		IsRelated: isRelated,
	}, nil
}

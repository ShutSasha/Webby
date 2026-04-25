package grpc

import (
	"context"
	"webby-vote-service/internal/grpc/votepb"
	"webby-vote-service/internal/models"

	"github.com/google/uuid"
)

type VoteListService interface {
	ListByRoom(ctx context.Context, roomId uuid.UUID) ([]models.Vote, error)
	GetChoicesByVoteId(
		ctx context.Context, voteId uuid.UUID,
	) ([]models.VoteChoice, error)
}

type VoteServer struct {
	votepb.UnimplementedVoteGrpcServiceServer
	service VoteListService
}

func NewVoteServer(service VoteListService) *VoteServer {
	return &VoteServer{service: service}
}

func (s *VoteServer) GetActiveVotes(
	ctx context.Context, req *votepb.GetActiveVotesRequest,
) (*votepb.GetActiveVotesResponse, error) {
	roomId, err := uuid.Parse(req.GetRoomId())
	if err != nil {
		return nil, err
	}

	votes, err := s.service.ListByRoom(ctx, roomId)
	if err != nil {
		return nil, err
	}

	result := make([]*votepb.VoteInfo, 0, len(votes))
	for _, v := range votes {
		choices, err := s.service.GetChoicesByVoteId(ctx, v.Id)
		if err != nil {
			return nil, err
		}

		totalVotes := 0
		for _, c := range choices {
			totalVotes += c.Votes
		}

		result = append(result, &votepb.VoteInfo{
			Id:              v.Id.String(),
			RoomId:          v.RoomId.String(),
			Type:            v.Type,
			VoteText:        v.VoteText,
			DurationSeconds: int32(v.DurationSeconds),
			TotalVotes:      int32(totalVotes),
		})
	}

	return &votepb.GetActiveVotesResponse{Votes: result}, nil
}

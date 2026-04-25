package services

import (
	"context"
	"fmt"
	"math"
	"time"
	"webby-vote-service/internal/apperrors"
	"webby-vote-service/internal/models"

	"github.com/google/uuid"
)

type VoteRepository interface {
	CreateVote(ctx context.Context, vote *models.Vote) (uuid.UUID, error)
	CreateChoice(
		ctx context.Context, choice *models.VoteChoice,
	) (uuid.UUID, error)
	DeleteVote(ctx context.Context, id uuid.UUID) error
	GetVoteById(ctx context.Context, id uuid.UUID) (*models.Vote, error)
	ListByRoom(
		ctx context.Context, roomId uuid.UUID,
	) ([]models.Vote, error)
	GetChoicesByVoteId(
		ctx context.Context, voteId uuid.UUID,
	) ([]models.VoteChoice, error)
	GetChoiceById(
		ctx context.Context, choiceId uuid.UUID,
	) (*models.VoteChoice, error)
	CastVote(ctx context.Context, choiceId, userId uuid.UUID) error
	RemoveUserVote(
		ctx context.Context, voteId, userId uuid.UUID,
	) error
	GetUserVoteForVote(
		ctx context.Context, voteId, userId uuid.UUID,
	) (*uuid.UUID, error)
}

type MemberChecker interface {
	Exists(ctx context.Context, roomId, userId uuid.UUID) (bool, error)
}

type RoomHostGetter interface {
	GetRoomHost(ctx context.Context, roomId uuid.UUID) (uuid.UUID, error)
}

type QueueItemMover interface {
	MoveToTop(ctx context.Context, id uuid.UUID) error
}

type CreateChoiceInput struct {
	Name        string
	IsCorrect   bool
	QueueItemId *uuid.UUID
}

type VoteChoiceDetail struct {
	Id          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Votes       int        `json:"votes"`
	Percentage  float64    `json:"percentage"`
	IsCorrect   bool       `json:"isCorrect"`
	QueueItemId *uuid.UUID `json:"queueItemId"`
}

type VoteDetail struct {
	Id                uuid.UUID          `json:"id"`
	RoomId            uuid.UUID          `json:"roomId"`
	Type              string             `json:"type"`
	VoteText          string             `json:"voteText"`
	CreatedAt         time.Time          `json:"createdAt"`
	DurationSeconds   int                `json:"durationSeconds"`
	ExpiresAt         time.Time          `json:"expiresAt"`
	IsExpired         bool               `json:"isExpired"`
	TotalVotes        int                `json:"totalVotes"`
	UserVotedChoiceId *uuid.UUID         `json:"userVotedChoiceId"`
	WinnerId          *uuid.UUID         `json:"winnerId"`
	Choices           []VoteChoiceDetail `json:"choices"`
}

type Service struct {
	voteRepo       VoteRepository
	memberChecker  MemberChecker
	roomHostGetter RoomHostGetter
	queueItemMover QueueItemMover
}

func New(
	voteRepo VoteRepository,
	memberChecker MemberChecker,
	roomHostGetter RoomHostGetter,
	queueItemMover QueueItemMover,
) *Service {
	return &Service{
		voteRepo:       voteRepo,
		memberChecker:  memberChecker,
		roomHostGetter: roomHostGetter,
		queueItemMover: queueItemMover,
	}
}

func (s *Service) CreateVote(
	ctx context.Context,
	roomId, userId uuid.UUID,
	voteType, voteText string,
	durationSeconds int,
	choices []CreateChoiceInput,
) (*VoteDetail, error) {
	hostId, err := s.roomHostGetter.GetRoomHost(ctx, roomId)
	if err != nil {
		return nil, err
	}
	if hostId != userId {
		return nil, apperrors.ErrForbidden
	}

	if voteType != "poll" && voteType != "next_video" {
		return nil, fmt.Errorf(
			"%w: type must be 'poll' or 'next_video'",
			apperrors.ErrInvalidInput,
		)
	}
	if len(choices) < 2 {
		return nil, fmt.Errorf(
			"%w: at least 2 choices required",
			apperrors.ErrInvalidInput,
		)
	}

	vote := &models.Vote{
		RoomId:          roomId,
		Type:            voteType,
		VoteText:        voteText,
		DurationSeconds: durationSeconds,
	}

	_, err = s.voteRepo.CreateVote(ctx, vote)
	if err != nil {
		return nil, err
	}

	for _, ci := range choices {
		choice := &models.VoteChoice{
			VoteId:      vote.Id,
			Name:        ci.Name,
			IsCorrect:   ci.IsCorrect,
			QueueItemId: ci.QueueItemId,
		}
		_, err := s.voteRepo.CreateChoice(ctx, choice)
		if err != nil {
			return nil, err
		}
	}

	return s.enrichVote(ctx, vote, userId)
}

func (s *Service) ListVotes(
	ctx context.Context, roomId, userId uuid.UUID,
) ([]VoteDetail, error) {
	exists, err := s.memberChecker.Exists(ctx, roomId, userId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, apperrors.ErrForbidden
	}

	votes, err := s.voteRepo.ListByRoom(ctx, roomId)
	if err != nil {
		return nil, err
	}

	details := make([]VoteDetail, 0, len(votes))
	for _, v := range votes {
		d, err := s.enrichVote(ctx, &v, userId)
		if err != nil {
			return nil, err
		}
		details = append(details, *d)
	}

	return details, nil
}

func (s *Service) GetVote(
	ctx context.Context, voteId, userId uuid.UUID,
) (*VoteDetail, error) {
	vote, err := s.voteRepo.GetVoteById(ctx, voteId)
	if err != nil {
		return nil, err
	}

	exists, err := s.memberChecker.Exists(ctx, vote.RoomId, userId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, apperrors.ErrForbidden
	}

	return s.enrichVote(ctx, vote, userId)
}

func (s *Service) CastVote(
	ctx context.Context, voteId, choiceId, userId uuid.UUID,
) (*VoteDetail, error) {
	vote, err := s.voteRepo.GetVoteById(ctx, voteId)
	if err != nil {
		return nil, err
	}

	exists, err := s.memberChecker.Exists(ctx, vote.RoomId, userId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, apperrors.ErrForbidden
	}

	expiresAt := vote.CreatedAt.Add(
		time.Duration(vote.DurationSeconds) * time.Second,
	)
	if time.Now().After(expiresAt) {
		return nil, fmt.Errorf(
			"%w: vote has expired", apperrors.ErrInvalidInput,
		)
	}

	choice, err := s.voteRepo.GetChoiceById(ctx, choiceId)
	if err != nil {
		return nil, err
	}
	if choice.VoteId != voteId {
		return nil, fmt.Errorf(
			"%w: choice does not belong to this vote",
			apperrors.ErrInvalidInput,
		)
	}

	existing, err := s.voteRepo.GetUserVoteForVote(ctx, voteId, userId)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf(
			"%w: user has already voted", apperrors.ErrConflict,
		)
	}

	if err := s.voteRepo.CastVote(ctx, choiceId, userId); err != nil {
		return nil, err
	}

	return s.enrichVote(ctx, vote, userId)
}

func (s *Service) RemoveVote(
	ctx context.Context, voteId, userId uuid.UUID,
) (*VoteDetail, error) {
	vote, err := s.voteRepo.GetVoteById(ctx, voteId)
	if err != nil {
		return nil, err
	}

	exists, err := s.memberChecker.Exists(ctx, vote.RoomId, userId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, apperrors.ErrForbidden
	}

	expiresAt := vote.CreatedAt.Add(
		time.Duration(vote.DurationSeconds) * time.Second,
	)
	if time.Now().After(expiresAt) {
		return nil, fmt.Errorf(
			"%w: vote has expired", apperrors.ErrInvalidInput,
		)
	}

	if err := s.voteRepo.RemoveUserVote(ctx, voteId, userId); err != nil {
		return nil, err
	}

	return s.enrichVote(ctx, vote, userId)
}

func (s *Service) DeleteVote(
	ctx context.Context, voteId, userId uuid.UUID,
) error {
	vote, err := s.voteRepo.GetVoteById(ctx, voteId)
	if err != nil {
		return err
	}

	hostId, err := s.roomHostGetter.GetRoomHost(ctx, vote.RoomId)
	if err != nil {
		return err
	}
	if hostId != userId {
		return apperrors.ErrForbidden
	}

	return s.voteRepo.DeleteVote(ctx, voteId)
}

func (s *Service) enrichVote(
	ctx context.Context, vote *models.Vote, userId uuid.UUID,
) (*VoteDetail, error) {
	choices, err := s.voteRepo.GetChoicesByVoteId(ctx, vote.Id)
	if err != nil {
		return nil, err
	}

	userChoiceId, err := s.voteRepo.GetUserVoteForVote(
		ctx, vote.Id, userId,
	)
	if err != nil {
		return nil, err
	}

	expiresAt := vote.CreatedAt.Add(
		time.Duration(vote.DurationSeconds) * time.Second,
	)
	isExpired := time.Now().After(expiresAt)

	totalVotes := 0
	for _, c := range choices {
		totalVotes += c.Votes
	}

	choiceDetails := make([]VoteChoiceDetail, len(choices))
	for i, c := range choices {
		pct := 0.0
		if totalVotes > 0 {
			pct = math.Round(
				float64(c.Votes) / float64(totalVotes) * 100,
			)
		}
		choiceDetails[i] = VoteChoiceDetail{
			Id:          c.Id,
			Name:        c.Name,
			Votes:       c.Votes,
			Percentage:  pct,
			IsCorrect:   c.IsCorrect,
			QueueItemId: c.QueueItemId,
		}
	}

	var winnerId *uuid.UUID
	if isExpired && vote.Type == "next_video" {
		winnerId = s.computeWinner(ctx, vote, choices)
		if winnerId != nil && s.queueItemMover != nil {
			for _, c := range choices {
				if c.Id == *winnerId && c.QueueItemId != nil {
					_ = s.queueItemMover.MoveToTop(
						ctx, *c.QueueItemId,
					)
					break
				}
			}
		}
	}

	return &VoteDetail{
		Id:                vote.Id,
		RoomId:            vote.RoomId,
		Type:              vote.Type,
		VoteText:          vote.VoteText,
		CreatedAt:         vote.CreatedAt,
		DurationSeconds:   vote.DurationSeconds,
		ExpiresAt:         expiresAt,
		IsExpired:         isExpired,
		TotalVotes:        totalVotes,
		UserVotedChoiceId: userChoiceId,
		WinnerId:          winnerId,
		Choices:           choiceDetails,
	}, nil
}

func (s *Service) computeWinner(
	ctx context.Context,
	vote *models.Vote,
	choices []models.VoteChoice,
) *uuid.UUID {
	if len(choices) == 0 {
		return nil
	}

	hostId, err := s.roomHostGetter.GetRoomHost(ctx, vote.RoomId)
	if err != nil {
		maxVotes := -1
		var winner uuid.UUID
		for _, c := range choices {
			if c.Votes > maxVotes {
				maxVotes = c.Votes
				winner = c.Id
			}
		}
		return &winner
	}

	maxVotes := 0
	for _, c := range choices {
		if c.Votes > maxVotes {
			maxVotes = c.Votes
		}
	}

	tied := []models.VoteChoice{}
	for _, c := range choices {
		if c.Votes == maxVotes {
			tied = append(tied, c)
		}
	}

	if len(tied) == 1 {
		return &tied[0].Id
	}

	hostChoiceId, _ := s.voteRepo.GetUserVoteForVote(
		ctx, vote.Id, hostId,
	)
	if hostChoiceId != nil {
		for _, c := range tied {
			if c.Id == *hostChoiceId {
				return &c.Id
			}
		}
	}

	return &tied[0].Id
}

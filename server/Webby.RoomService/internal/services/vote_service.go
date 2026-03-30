package services

import (
	"context"
	"fmt"
	"math"
	"time"
	"webby/internal/apperrors"
	"webby/internal/models"

	"github.com/google/uuid"
)

type VoteRepo interface {
	CreateVote(vote *models.Vote) (uuid.UUID, error)
	CreateChoice(choice *models.VoteChoice) (uuid.UUID, error)
	DeleteVote(id uuid.UUID) error
	GetVoteById(id uuid.UUID) (*models.Vote, error)
	ListByRoom(roomId uuid.UUID) ([]models.Vote, error)
	GetChoicesByVoteId(voteId uuid.UUID) ([]models.VoteChoice, error)
	CastVote(choiceId, userId uuid.UUID) error
	RemoveUserVote(voteId, userId uuid.UUID) error
	GetUserVoteForVote(voteId, userId uuid.UUID) (*uuid.UUID, error)
	GetChoiceById(choiceId uuid.UUID) (*models.VoteChoice, error)
}

type RoomGetter interface {
	GetById(id uuid.UUID) (*models.Room, error)
}

type QueueItemMover interface {
	MoveToTop(id uuid.UUID) error
}

type VoteService struct {
	voteRepo       VoteRepo
	roomGetter     RoomGetter
	memberChecker  MemberChecker
	queueItemMover QueueItemMover
}

func NewVoteService(voteRepo VoteRepo, roomGetter RoomGetter, memberChecker MemberChecker, queueItemMover QueueItemMover) *VoteService {
	return &VoteService{
		voteRepo:       voteRepo,
		roomGetter:     roomGetter,
		memberChecker:  memberChecker,
		queueItemMover: queueItemMover,
	}
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

func (s *VoteService) enrichVote(vote *models.Vote, userId uuid.UUID) (*VoteDetail, error) {
	choices, err := s.voteRepo.GetChoicesByVoteId(vote.Id)
	if err != nil {
		return nil, err
	}

	userChoiceId, err := s.voteRepo.GetUserVoteForVote(vote.Id, userId)
	if err != nil {
		return nil, err
	}

	expiresAt := vote.CreatedAt.Add(time.Duration(vote.DurationSeconds) * time.Second)
	isExpired := time.Now().After(expiresAt)

	totalVotes := 0
	for _, c := range choices {
		totalVotes += c.Votes
	}

	choiceDetails := make([]VoteChoiceDetail, len(choices))
	for i, c := range choices {
		pct := 0.0
		if totalVotes > 0 {
			pct = math.Round(float64(c.Votes) / float64(totalVotes) * 100)
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
		winnerId = s.computeWinner(vote, choices)
		if winnerId != nil {
			for _, c := range choices {
				if c.Id == *winnerId && c.QueueItemId != nil {
					_ = s.queueItemMover.MoveToTop(*c.QueueItemId)
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

func (s *VoteService) computeWinner(vote *models.Vote, choices []models.VoteChoice) *uuid.UUID {
	if len(choices) == 0 {
		return nil
	}

	room, err := s.roomGetter.GetById(vote.RoomId)
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

	hostChoiceId, _ := s.voteRepo.GetUserVoteForVote(vote.Id, room.HostId)
	if hostChoiceId != nil {
		for _, c := range tied {
			if c.Id == *hostChoiceId {
				return &c.Id
			}
		}
	}

	return &tied[0].Id
}

func (s *VoteService) CreateVote(ctx context.Context, roomId, userId uuid.UUID, voteType, voteText string, durationSeconds int, choices []CreateChoiceInput) (*VoteDetail, error) {
	room, err := s.roomGetter.GetById(roomId)
	if err != nil {
		return nil, err
	}
	if room.HostId != userId {
		return nil, apperrors.ErrForbidden
	}

	if voteType != "poll" && voteType != "next_video" {
		return nil, fmt.Errorf("%w: type must be 'poll' or 'next_video'", apperrors.ErrInvalidInput)
	}
	if len(choices) < 2 {
		return nil, fmt.Errorf("%w: at least 2 choices required", apperrors.ErrInvalidInput)
	}

	vote := &models.Vote{
		RoomId:          roomId,
		Type:            voteType,
		VoteText:        voteText,
		DurationSeconds: durationSeconds,
	}

	_, err = s.voteRepo.CreateVote(vote)
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
		_, err := s.voteRepo.CreateChoice(choice)
		if err != nil {
			return nil, err
		}
	}

	return s.enrichVote(vote, userId)
}

func (s *VoteService) ListVotes(ctx context.Context, roomId, userId uuid.UUID) ([]VoteDetail, error) {
	exists, err := s.memberChecker.Exists(roomId, userId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, apperrors.ErrForbidden
	}

	votes, err := s.voteRepo.ListByRoom(roomId)
	if err != nil {
		return nil, err
	}

	details := make([]VoteDetail, 0, len(votes))
	for _, v := range votes {
		d, err := s.enrichVote(&v, userId)
		if err != nil {
			return nil, err
		}
		details = append(details, *d)
	}

	return details, nil
}

func (s *VoteService) GetVote(ctx context.Context, voteId, userId uuid.UUID) (*VoteDetail, error) {
	vote, err := s.voteRepo.GetVoteById(voteId)
	if err != nil {
		return nil, err
	}

	exists, err := s.memberChecker.Exists(vote.RoomId, userId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, apperrors.ErrForbidden
	}

	return s.enrichVote(vote, userId)
}

func (s *VoteService) CastVote(ctx context.Context, voteId, choiceId, userId uuid.UUID) (*VoteDetail, error) {
	vote, err := s.voteRepo.GetVoteById(voteId)
	if err != nil {
		return nil, err
	}

	exists, err := s.memberChecker.Exists(vote.RoomId, userId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, apperrors.ErrForbidden
	}

	expiresAt := vote.CreatedAt.Add(time.Duration(vote.DurationSeconds) * time.Second)
	if time.Now().After(expiresAt) {
		return nil, fmt.Errorf("%w: vote has expired", apperrors.ErrInvalidInput)
	}

	choice, err := s.voteRepo.GetChoiceById(choiceId)
	if err != nil {
		return nil, err
	}
	if choice.VoteId != voteId {
		return nil, fmt.Errorf("%w: choice does not belong to this vote", apperrors.ErrInvalidInput)
	}

	existing, err := s.voteRepo.GetUserVoteForVote(voteId, userId)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: user has already voted", apperrors.ErrConflict)
	}

	if err := s.voteRepo.CastVote(choiceId, userId); err != nil {
		return nil, err
	}

	return s.enrichVote(vote, userId)
}

func (s *VoteService) RemoveVote(ctx context.Context, voteId, userId uuid.UUID) (*VoteDetail, error) {
	vote, err := s.voteRepo.GetVoteById(voteId)
	if err != nil {
		return nil, err
	}

	exists, err := s.memberChecker.Exists(vote.RoomId, userId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, apperrors.ErrForbidden
	}

	expiresAt := vote.CreatedAt.Add(time.Duration(vote.DurationSeconds) * time.Second)
	if time.Now().After(expiresAt) {
		return nil, fmt.Errorf("%w: vote has expired", apperrors.ErrInvalidInput)
	}

	if err := s.voteRepo.RemoveUserVote(voteId, userId); err != nil {
		return nil, err
	}

	return s.enrichVote(vote, userId)
}

func (s *VoteService) DeleteVote(ctx context.Context, voteId, userId uuid.UUID) error {
	vote, err := s.voteRepo.GetVoteById(voteId)
	if err != nil {
		return err
	}

	room, err := s.roomGetter.GetById(vote.RoomId)
	if err != nil {
		return err
	}
	if room.HostId != userId {
		return apperrors.ErrForbidden
	}

	return s.voteRepo.DeleteVote(voteId)
}

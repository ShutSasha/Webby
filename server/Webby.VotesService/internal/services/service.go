package services

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"time"
	"webby/vote-service/internal/apperrors"
	"webby/vote-service/internal/models"
	"webby/vote-service/pkg/logger"

	"github.com/google/uuid"
)

type eventEnvelope struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

const EventTypeVotingStarted = "VOTING_STARTED"

type repository interface {
	CreateVoteWithRightChoice(ctx context.Context, vote *models.Vote) error
	SaveChoices(ctx context.Context, vodeID uuid.UUID, choices []string) error
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

type roomHostGetter interface {
	GetRoomHost(ctx context.Context, roomId uuid.UUID) (uuid.UUID, error)
}

type chatRetriever interface {
	GetChatIDByRoomID(ctx context.Context, roomID, userID uuid.UUID) (uuid.UUID, error)
}

type memberChecker interface {
	Exists(ctx context.Context, roomId, userId uuid.UUID) (bool, error)
}

type queueItemMover interface {
	MoveToTop(ctx context.Context, id uuid.UUID) error
}

type publisher interface {
	Publish(ctx context.Context, channel string, payload any) error
}

type service struct {
	repository     repository
	roomHostGetter roomHostGetter
	chatRetriever  chatRetriever
	memberChecker  memberChecker
	queueItemMover queueItemMover
	publisher      publisher
}

func New(repository repository, roomHostGetter roomHostGetter, chatRetriever chatRetriever, memberChecker memberChecker, queueItemMover queueItemMover, publisher publisher) *service {
	return &service{
		repository:     repository,
		memberChecker:  memberChecker,
		roomHostGetter: roomHostGetter,
		queueItemMover: queueItemMover,
		publisher:      publisher,
	}
}

func (s *service) CreateWithRightChoice(ctx context.Context, roomID, userID uuid.UUID, voteText string, duration int, choices []string) error {
	const op = "service.CreateWithRightChoice"
	log := logger.FromContext(ctx).With("op", op)

	hostID, err := s.roomHostGetter.GetRoomHost(ctx, roomID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if hostID != userID {
		return fmt.Errorf("%s: %w", op, apperrors.ErrNotHost)
	}

	chatID, err := s.chatRetriever.GetChatIDByRoomID(ctx, roomID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	vodeID := uuid.New()
	vote := &models.Vote{
		ID:        vodeID,
		RoomID:    roomID,
		VoteText:  voteText,
		CreatedAt: time.Now(),
		Duration:  duration,
	}
	err = s.repository.CreateVoteWithRightChoice(ctx, vote)
	if err != nil {
		return fmt.Errorf("%s: craete vote %w", op, err)
	}

	err = s.repository.SaveChoices(ctx, vodeID, choices)
	if err != nil {
		return fmt.Errorf("%s: save choices %w", op, err)
	}

	envelope := eventEnvelope{
		Type: EventTypeVotingStarted,
		Payload: models.EnrichedVoting{
			ID:        vote.ID,
			VoteText:  voteText,
			Duration:  duration,
			CreatedAt: vote.CreatedAt,
			Choices:   choices,
		},
	}
	topic := fmt.Sprintf("chat:%s", chatID.String())
	err = s.publisher.Publish(ctx, topic, envelope)
	if err != nil {
		log.Error("Could not notify about voting", slog.String("err", err.Error()))
	}

	return nil
}

func (s *service) ListVotes(
	ctx context.Context, roomId, userId uuid.UUID,
) ([]VoteDetail, error) {
	exists, err := s.memberChecker.Exists(ctx, roomId, userId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, apperrors.ErrForbidden
	}

	votes, err := s.repository.ListByRoom(ctx, roomId)
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

func (s *service) GetVote(
	ctx context.Context, voteId, userId uuid.UUID,
) (*VoteDetail, error) {
	vote, err := s.repository.GetVoteById(ctx, voteId)
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

func (s *service) CastVote(
	ctx context.Context, voteId, choiceId, userId uuid.UUID,
) (*VoteDetail, error) {
	vote, err := s.repository.GetVoteById(ctx, voteId)
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

	choice, err := s.repository.GetChoiceById(ctx, choiceId)
	if err != nil {
		return nil, err
	}
	if choice.VoteId != voteId {
		return nil, fmt.Errorf(
			"%w: choice does not belong to this vote",
			apperrors.ErrInvalidInput,
		)
	}

	existing, err := s.repository.GetUserVoteForVote(ctx, voteId, userId)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf(
			"%w: user has already voted", apperrors.ErrConflict,
		)
	}

	if err := s.repository.CastVote(ctx, choiceId, userId); err != nil {
		return nil, err
	}

	return s.enrichVote(ctx, vote, userId)
}

func (s *service) RemoveVote(
	ctx context.Context, voteId, userId uuid.UUID,
) (*VoteDetail, error) {
	vote, err := s.repository.GetVoteById(ctx, voteId)
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

	if err := s.repository.RemoveUserVote(ctx, voteId, userId); err != nil {
		return nil, err
	}

	return s.enrichVote(ctx, vote, userId)
}

func (s *service) DeleteVote(
	ctx context.Context, voteId, userId uuid.UUID,
) error {
	vote, err := s.repository.GetVoteById(ctx, voteId)
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

	return s.repository.DeleteVote(ctx, voteId)
}

func (s *service) enrichVote(
	ctx context.Context, vote *models.Vote, userId uuid.UUID,
) (*VoteDetail, error) {
	choices, err := s.repository.GetChoicesByVoteId(ctx, vote.Id)
	if err != nil {
		return nil, err
	}

	userChoiceId, err := s.repository.GetUserVoteForVote(
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

func (s *service) computeWinner(
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

	hostChoiceId, _ := s.repository.GetUserVoteForVote(
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

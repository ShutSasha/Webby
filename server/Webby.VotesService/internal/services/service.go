package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"
	"webby/vote-service/internal/apperrors"
	"webby/vote-service/internal/models"
	"webby/vote-service/pkg/logger"

	"github.com/google/uuid"
)

type votingResults struct {
	VotingID    uuid.UUID   `json:"votingId"`
	RightChoice string      `json:"rightChoice"`
	Winners     []uuid.UUID `json:"winners"`
}

type eventEnvelope struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

const eventTypeVotingStarted = "VOTING_STARTED"
const eventTypeVotingResults = "VOTING_RESULTS"
const eventTypeVotingLocked = "VOTING_LOCKED"

type repository interface {
	CreateVoteWithRightChoice(ctx context.Context, vote *models.Vote) error
	SaveChoices(ctx context.Context, vodeID uuid.UUID, choices []string) error
	GetVotingUserWinners(ctx context.Context, voteID uuid.UUID, rightChoice string) ([]uuid.UUID, error)
	ListByRoom(ctx context.Context, roomID uuid.UUID) ([]models.Vote, error)
	GetChoicesForVoting(ctx context.Context, voteID uuid.UUID) ([]string, error)
	IsChoiceValid(ctx context.Context, voteID uuid.UUID, rightChoice string) (bool, error)
	SetVotingRightOption(ctx context.Context, roomID, voteID uuid.UUID, rightChoice string) error
	MarkVoteAsLocked(ctx context.Context, voteID uuid.UUID) error
	CastVote(ctx context.Context, voteID, userID uuid.UUID, choice string) error
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

	activeTimers sync.Map
}

func New(repository repository, roomHostGetter roomHostGetter, chatRetriever chatRetriever, memberChecker memberChecker, queueItemMover queueItemMover, publisher publisher) *service {
	return &service{
		repository:     repository,
		memberChecker:  memberChecker,
		chatRetriever:  chatRetriever,
		roomHostGetter: roomHostGetter,
		queueItemMover: queueItemMover,
		publisher:      publisher,

		activeTimers: sync.Map{},
	}
}

func (s *service) CreateWithRightChoice(ctx context.Context, roomID, userID uuid.UUID, voteText string, duration int, choices []string) error {
	const op = "service.CreateWithRightChoice"
	log := logger.FromContext(ctx).With("op", op)

	hostID, err := s.roomHostGetter.GetRoomHost(ctx, roomID)
	if err != nil {
		return fmt.Errorf("%s: get room host %w", op, err)
	}
	if hostID != userID {
		return fmt.Errorf("%s: %w", op, apperrors.ErrNotHost)
	}

	chatID, err := s.chatRetriever.GetChatIDByRoomID(ctx, roomID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	voteID := uuid.New()
	vote := &models.Vote{
		ID:        voteID,
		RoomID:    roomID,
		VoteText:  voteText,
		CreatedAt: time.Now(),
		Duration:  duration,
	}
	err = s.repository.CreateVoteWithRightChoice(ctx, vote)
	if err != nil {
		return fmt.Errorf("%s: create vote %w", op, err)
	}

	err = s.repository.SaveChoices(ctx, voteID, choices)
	if err != nil {
		return fmt.Errorf("%s: save choices %w", op, err)
	}

	envelope := eventEnvelope{
		Type: eventTypeVotingStarted,
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

	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)

	s.activeTimers.Store(voteID, cancel)
	go func(asyncCtx context.Context, vID uuid.UUID, chatTopic string) {
		defer s.activeTimers.Delete(vID)
		bgLog := logger.FromContext(asyncCtx).With(slog.String("op", op+"_async"))

		<-asyncCtx.Done()

		if errors.Is(asyncCtx.Err(), context.Canceled) {
			bgLog.Info("voting timer cancelled manually")
			return
		}

		err := s.repository.MarkVoteAsLocked(context.Background(), vID)
		if err == nil {
			lockEnvelope := eventEnvelope{
				Type:    eventTypeVotingLocked,
				Payload: map[string]string{"votingId": vID.String()},
			}
			s.publisher.Publish(context.Background(), chatTopic, lockEnvelope)
		}
	}(timeoutCtx, voteID, topic)

	return nil
}

func (s *service) ResolveVoting(ctx context.Context, roomID, userID, voteID uuid.UUID, rightChoice string) error {
	const op = "service.ResolveVoting"
	log := logger.FromContext(ctx).With("op", op)

	hostID, err := s.roomHostGetter.GetRoomHost(ctx, roomID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if hostID != userID {
		return fmt.Errorf("%s: %w", op, apperrors.ErrNotHost)
	}

	isValid, err := s.repository.IsChoiceValid(ctx, voteID, rightChoice)
	if err != nil {
		return fmt.Errorf("%s: validate choice: %w", op, err)
	}
	if !isValid {
		return fmt.Errorf("%s: %w", op, apperrors.ErrInvalidChoice)
	}

	chatID, err := s.chatRetriever.GetChatIDByRoomID(ctx, roomID, userID)
	if err != nil {
		return fmt.Errorf("%s: get chat: %w", op, err)
	}

	if cancelInf, ok := s.activeTimers.LoadAndDelete(voteID); ok {
		cancel := cancelInf.(context.CancelFunc)
		cancel()
	}

	err = s.repository.SetVotingRightOption(ctx, roomID, voteID, rightChoice)
	if err != nil {
		return fmt.Errorf("%s: set right option: %w", op, err)
	}

	userIDs, err := s.repository.GetVotingUserWinners(ctx, voteID, rightChoice)
	if err != nil {
		return fmt.Errorf("%s: get winners: %w", op, err)
	}

	resultsEnvelope := eventEnvelope{
		Type: eventTypeVotingResults,
		Payload: votingResults{
			VotingID:    voteID,
			RightChoice: rightChoice,
			Winners:     userIDs,
		},
	}
	topic := fmt.Sprintf("chat:%s", chatID.String())
	if err := s.publisher.Publish(ctx, topic, resultsEnvelope); err != nil {
		log.Error("failed to publish votingResults", slog.String("err", err.Error()))
	}

	return nil
}

func (s *service) ListVotings(ctx context.Context, roomID, userID uuid.UUID) ([]models.EnrichedVoting, error) {
	const op = "service.ListVotes"
	log := logger.FromContext(ctx).With("op", op)

	exists, err := s.memberChecker.Exists(ctx, roomID, userID)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	if !exists {
		return nil, fmt.Errorf("%s, %w", op, apperrors.ErrMemberNotFound)
	}

	votes, err := s.repository.ListByRoom(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	log.Debug("Retrieving votes", "votes", votes)

	enrichedVotings := make([]models.EnrichedVoting, 0, len(votes))
	for _, v := range votes {
		choices, err := s.repository.GetChoicesForVoting(ctx, v.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		enrichedVoting := models.EnrichedVoting{
			ID:        v.ID,
			VoteText:  v.VoteText,
			Duration:  v.Duration,
			CreatedAt: v.CreatedAt,
			Choices:   choices,
		}

		log.Debug("List enriched votings", "voting", enrichedVoting)

		enrichedVotings = append(enrichedVotings, enrichedVoting)
	}

	return enrichedVotings, nil
}

func (s *service) CastVote(ctx context.Context, roomID, voteID, userID uuid.UUID, choice string) error {
	const op = "service.CastVote"

	exists, err := s.memberChecker.Exists(ctx, roomID, userID)
	if err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}
	if !exists {
		return fmt.Errorf("%s, %w", op, apperrors.ErrMemberNotFound)
	}

	isValid, err := s.repository.IsChoiceValid(ctx, voteID, choice)
	if err != nil {
		return fmt.Errorf("%s: validate choice: %w", op, err)
	}
	if !isValid {
		return fmt.Errorf("%s: %w", op, apperrors.ErrInvalidChoice)
	}

	err = s.repository.CastVote(ctx, voteID, userID, choice)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

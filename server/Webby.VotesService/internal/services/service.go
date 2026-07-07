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
	VotingID    uuid.UUID `json:"votingId"`
	RightChoice string    `json:"rightChoice"`
}

type eventEnvelope struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

const (
	eventTypeVotingStarted          = "VOTING_STARTED"
	eventTypeVotingResults          = "VOTING_RESULTS"
	eventTypeVotingLocked           = "VOTING_LOCKED"
	eventTypeNextVideoVotingStarted = "NEXT_VIDEO_VOTING_STARTED"
	eventTypeNextVideoVotingResult  = "NEXT_VIDEO_VOTING_RESULTS"
)

type repository interface {
	CreateVoteWithRightChoice(ctx context.Context, vote *models.Vote) error
	SaveChoices(ctx context.Context, vodeID uuid.UUID, choices []string) error
	ListByRoom(ctx context.Context, roomID uuid.UUID) ([]models.Vote, error)
	GetChoicesForVoting(ctx context.Context, voteID uuid.UUID) ([]string, error)
	IsChoiceValid(ctx context.Context, voteID uuid.UUID, rightChoice string) (bool, error)
	SetVotingRightOption(ctx context.Context, roomID, voteID uuid.UUID, rightChoice string) error
	MarkVoteAsLocked(ctx context.Context, voteID uuid.UUID) error
	CastVote(ctx context.Context, voteID, userID uuid.UUID, choice string) error
	CreateVotingForNextVideo(ctx context.Context, roomID uuid.UUID, duration int) error
	GetNextVideoResults(ctx context.Context, roomID uuid.UUID) (map[uuid.UUID]int, error)
	VoteForNextVideo(ctx context.Context, roomID, userID, queueItemID uuid.UUID) error
	GetNextVideoVoting(ctx context.Context, roomID uuid.UUID) (*models.NextVideoInfo, error)
	GetUserVote(ctx context.Context, voteID, userID uuid.UUID) (*string, error)
}

type roomHostGetter interface {
	GetRoomHost(ctx context.Context, roomID uuid.UUID) (uuid.UUID, error)
}

type chatRetriever interface {
	GetChatIDByRoomID(ctx context.Context, roomID uuid.UUID) (uuid.UUID, error)
}

type memberChecker interface {
	Exists(ctx context.Context, roomID, userID uuid.UUID) (bool, error)
}

type queueItemMover interface {
	MakeNext(ctx context.Context, roomID, queueItemID uuid.UUID) error
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

	chatID, err := s.chatRetriever.GetChatIDByRoomID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	vote := &models.Vote{
		ID:        uuid.New(),
		RoomID:    roomID,
		VoteText:  voteText,
		CreatedAt: time.Now(),
		Duration:  duration,
	}
	err = s.repository.CreateVoteWithRightChoice(ctx, vote)
	if err != nil {
		return fmt.Errorf("%s: create vote %w", op, err)
	}

	err = s.repository.SaveChoices(ctx, vote.ID, choices)
	if err != nil {
		return fmt.Errorf("%s: save choices %w", op, err)
	}

	envelope := eventEnvelope{
		Type: eventTypeVotingStarted,
		Payload: models.EnrichedVoting{
			ID:        vote.ID,
			VoteText:  voteText,
			Duration:  duration,
			ExpiresAt: vote.CreatedAt.Add(time.Duration(duration) * time.Second),
			Choices:   choices,
		},
	}
	topic := fmt.Sprintf("chat:%s", chatID.String())
	err = s.publisher.Publish(ctx, topic, envelope)
	if err != nil {
		log.Error("Could not notify about voting", slog.String("err", err.Error()))
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)

	s.activeTimers.Store(vote.ID, cancel)
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
	}(timeoutCtx, vote.ID, topic)

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

	chatID, err := s.chatRetriever.GetChatIDByRoomID(ctx, roomID)
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

	resultsEnvelope := eventEnvelope{
		Type: eventTypeVotingResults,
		Payload: votingResults{
			VotingID:    voteID,
			RightChoice: rightChoice,
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
		return nil, fmt.Errorf("%s, %w", op, apperrors.ErrNotMember)
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
			return nil, fmt.Errorf("%s: get choices: %w", op, err)
		}

		myVote, err := s.repository.GetUserVote(ctx, v.ID, userID)
		if err != nil {
			return nil, fmt.Errorf("%s: get user vote: %w", op, err)
		}

		isLocked := v.Status != "active"

		enrichedVoting := models.EnrichedVoting{
			ID:        v.ID,
			VoteText:  v.VoteText,
			Duration:  v.Duration,
			ExpiresAt: v.CreatedAt.Add(time.Duration(v.Duration) * time.Second),
			Choices:   choices,
			IsLocked:  isLocked,
			MyVote:    myVote,
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
		return fmt.Errorf("%s, %w", op, apperrors.ErrNotMember)
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

func (s *service) CreateVotingForNextVideo(ctx context.Context, roomID, userID uuid.UUID) error {
	const op = "service.CreateVotingForNextVideo"
	log := logger.FromContext(ctx).With("op", op)

	hostID, err := s.roomHostGetter.GetRoomHost(ctx, roomID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if userID != hostID {
		return fmt.Errorf("%s: %w", op, apperrors.ErrNotHost)
	}

	nextVoting, err := s.repository.GetNextVideoVoting(ctx, roomID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if nextVoting.Exists {
		return fmt.Errorf("%s: %w", op, apperrors.ErrVotingAlreadyExist)
	}

	chatID, err := s.chatRetriever.GetChatIDByRoomID(ctx, roomID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.repository.CreateVotingForNextVideo(ctx, roomID, 20)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	envelope := eventEnvelope{
		Type: eventTypeNextVideoVotingStarted,
		Payload: map[string]any{
			"duration":  20,
			"expiresAt": time.Now().Add(time.Duration(20) * time.Second),
		},
	}
	topic := fmt.Sprintf("chat:%s", chatID.String())
	if err := s.publisher.Publish(ctx, topic, envelope); err != nil {
		log.Error("failed to publish next video voting", slog.String("err", err.Error()))
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(20*time.Second))
	s.activeTimers.Store(roomID, cancel)
	go func(asyncCtx context.Context, rID, cID uuid.UUID) {
		defer cancel()
		defer s.activeTimers.Delete(rID)

		bgLog := logger.FromContext(asyncCtx).With("op", op+"_async")

		<-asyncCtx.Done()

		videoResults, err := s.repository.GetNextVideoResults(context.Background(), rID)
		if err != nil {
			bgLog.Error("Failed to get video results", "err", err.Error())
			return
		}

		if len(videoResults) == 0 {
			bgLog.Info("No votes cast, skipping next video selection")

			resultsEnvelope := eventEnvelope{
				Type:    eventTypeNextVideoVotingResult,
				Payload: map[string]string{"winnerId": "none"},
			}
			topic := fmt.Sprintf("chat:%s", cID.String())
			if err := s.publisher.Publish(context.Background(), topic, resultsEnvelope); err != nil {
				log.Error("failed to publish empty voting results", slog.String("err", err.Error()))
			}

			return
		}

		winnerID := uuid.Nil
		votes := 0
		for id, count := range videoResults {
			if count > votes {
				winnerID = id
				votes = count
			}
		}

		err = s.queueItemMover.MakeNext(context.Background(), rID, winnerID)
		if err != nil {
			bgLog.Error("Failed to make video next", "err", err.Error())
			return
		}

		resultsEnvelope := eventEnvelope{
			Type:    eventTypeNextVideoVotingResult,
			Payload: map[string]string{"winnerId": winnerID.String()},
		}
		topic := fmt.Sprintf("chat:%s", cID.String())
		if err := s.publisher.Publish(context.Background(), topic, resultsEnvelope); err != nil {
			log.Error("failed to publish next video voting results", slog.String("err", err.Error()))
		}
	}(timeoutCtx, roomID, chatID)

	return nil
}

func (s *service) VoteForNextVideo(ctx context.Context, roomID, userID, queueItemID uuid.UUID) error {
	const op = "service.VoteForNextVideo"

	exists, err := s.memberChecker.Exists(ctx, roomID, userID)
	if err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}
	if !exists {
		return fmt.Errorf("%s, %w", op, apperrors.ErrNotMember)
	}

	voting, err := s.repository.GetNextVideoVoting(ctx, roomID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !voting.Exists {
		return fmt.Errorf("%s: %w", op, apperrors.ErrNoVoting)
	}

	err = s.repository.VoteForNextVideo(ctx, roomID, userID, queueItemID)
	if err != nil {
		return fmt.Errorf("%s, %w", op, err)
	}

	return nil
}

func (s *service) HasNextVideoVoting(ctx context.Context, roomID, userID uuid.UUID) (*models.NextVideoInfo, error) {
	const op = "service.HasNextVideoVoting"

	exists, err := s.memberChecker.Exists(ctx, roomID, userID)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}
	if !exists {
		return nil, fmt.Errorf("%s, %w", op, apperrors.ErrNotMember)
	}

	voting, err := s.repository.GetNextVideoVoting(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	return voting, nil
}

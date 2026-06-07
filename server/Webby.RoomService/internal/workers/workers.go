package workers

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type ActiveRoom struct {
	ChatID  string
	UserIDs []uuid.UUID
}

type presenceRetriever interface {
	GetActiveRooms(ctx context.Context, minMembers int, zombieTTL time.Duration) ([]ActiveRoom, error)
}

type pointsRepository interface {
	AddPointsBulk(ctx context.Context, userIDs []uuid.UUID, points int) (map[uuid.UUID]int, error)
}

type eventPublisher interface {
	Publish(ctx context.Context, channel string, payload any) error
}

type pointsWorker struct {
	presenceRepo  presenceRetriever
	repo          pointsRepository
	pub           eventPublisher
	logger        *slog.Logger
	interval      time.Duration
	pointsPerTick int
	zombieTTL     time.Duration
}

func NewPointsWorker(
	presenceRepo presenceRetriever,
	repo pointsRepository,
	pub eventPublisher,
	logger *slog.Logger,
	interval time.Duration,
	pointsPerTick int,
	zombieTTL time.Duration,
) *pointsWorker {
	return &pointsWorker{
		presenceRepo:  presenceRepo,
		repo:          repo,
		pub:           pub,
		logger:        logger,
		interval:      interval,
		pointsPerTick: pointsPerTick,
		zombieTTL:     zombieTTL,
	}
}

func (w *pointsWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("points worker stopped")
			return
		case <-ticker.C:
			w.processPoints(ctx)
		}
	}
}

func (w *pointsWorker) processPoints(ctx context.Context) {
	const op = "workers.pointsWorker.processPoints"
	log := w.logger.With("op", op)

	rooms, err := w.presenceRepo.GetActiveRooms(ctx, 2, w.zombieTTL)
	if err != nil {
		log.Error("failed to get active rooms", slog.String("err", err.Error()))
		return
	}

	for _, room := range rooms {
		updatedTotals, err := w.repo.AddPointsBulk(ctx, room.UserIDs, w.pointsPerTick)
		if err != nil {
			log.Error("failed to add points", slog.String("err", err.Error()), slog.String("chat_id", room.ChatID))
			continue
		}

		totalsPayload := make(map[string]int, len(updatedTotals))
		for uid, total := range updatedTotals {
			totalsPayload[uid.String()] = total
		}

		w.pub.Publish(ctx, "chat:"+room.ChatID, map[string]any{
			"type": "ROOM_POINTS_UPDATED",
			"payload": map[string]any{
				"added_points": w.pointsPerTick,
				"totals":       totalsPayload,
			},
		})
	}
}

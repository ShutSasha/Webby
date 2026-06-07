package workers

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type pointsRepository interface {
	AddPointsBulk(ctx context.Context, userIDs []uuid.UUID, points int) (map[uuid.UUID]int, error)
}

type eventPublisher interface {
	Publish(ctx context.Context, channel string, payload any) error
}

type pointsWorker struct {
	rdb           *redis.Client
	repo          pointsRepository
	pub           eventPublisher
	logger        *slog.Logger
	interval      time.Duration
	pointsPerTick int
}

func NewPointsWorker(rdb *redis.Client, repo pointsRepository, pub eventPublisher, logger *slog.Logger, interval time.Duration, pointsPerTick int) *pointsWorker {
	return &pointsWorker{
		rdb:           rdb,
		repo:          repo,
		pub:           pub,
		logger:        logger,
		interval:      interval,
		pointsPerTick: pointsPerTick,
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
	chats, err := w.rdb.SMembers(ctx, "active_chats").Result()
	if err != nil {
		log.Error("failed to get active chats", slog.String("err", err.Error()))
		return
	}

	for _, chatIDStr := range chats {
		presenceKey := "chat:" + chatIDStr + ":presence"

		count, err := w.rdb.SCard(ctx, presenceKey).Result()
		if err != nil || count < 2 {
			continue
		}

		userIDsStr, err := w.rdb.SMembers(ctx, presenceKey).Result()
		if err != nil {
			continue
		}

		userIDs := make([]uuid.UUID, len(userIDsStr))
		for i, idStr := range userIDsStr {
			userID, _ := uuid.Parse(idStr)
			userIDs[i] = userID
		}

		updatedTotals, err := w.repo.AddPointsBulk(ctx, userIDs, w.pointsPerTick)
		if err != nil {
			log.Error("failed to add points", slog.String("err", err.Error()))
			continue
		}

		totalsPayload := make(map[string]int, len(updatedTotals))
		for uid, total := range updatedTotals {
			totalsPayload[uid.String()] = total
		}

		w.pub.Publish(ctx, "chat:"+chatIDStr, map[string]any{
			"type": "ROOM_POINTS_UPDATED",
			"payload": map[string]any{
				"added_points": w.pointsPerTick,
				"totals":       totalsPayload,
			},
		})
	}
}

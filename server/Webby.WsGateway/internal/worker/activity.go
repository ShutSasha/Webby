package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	clients "webby-wsgateway/internal/grpc"
	"webby-wsgateway/internal/ws"
)

type Activity struct {
	logger   *slog.Logger
	ws       *ws.Server
	room     *clients.RoomClient
	interval time.Duration
}

func NewActivity(logger *slog.Logger, wsSrv *ws.Server, room *clients.RoomClient, interval time.Duration) *Activity {
	return &Activity{logger: logger, ws: wsSrv, room: room, interval: interval}
}

func (a *Activity) Run(ctx context.Context) error {
	t := time.NewTicker(a.interval)
	defer t.Stop()

	a.logger.Info("activity worker started", slog.Duration("interval", a.interval))

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			a.tick(ctx)
		}
	}
}

func (a *Activity) tick(ctx context.Context) {
	snap := a.ws.SnapshotActivity()
	if len(snap) == 0 {
		return
	}

	batch := make([]clients.RoomActivity, 0, len(snap))
	for roomID, users := range snap {
		batch = append(batch, clients.RoomActivity{
			RoomID:  roomID.String(),
			UserIDs: toStrings(users),
		})
	}

	callCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := a.room.AwardActiveUsers(callCtx, batch); err != nil {
		a.logger.Error("AwardActiveUsers failed", slog.String("err", err.Error()))
		return
	}
	a.logger.Debug("activity reported", slog.Int("rooms", len(batch)))
}

func toStrings(ids []uuid.UUID) []string {
	out := make([]string, len(ids))
	for i, id := range ids {
		out[i] = id.String()
	}
	return out
}

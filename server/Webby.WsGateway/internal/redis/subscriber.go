package redisbus

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"webby-wsgateway/internal/ws"
)

type EventType string

const (
	EventMessageCreated EventType = "MESSAGE_CREATED"
	EventVoteUpdated    EventType = "VOTE_UPDATED"
	EventPointsAwarded  EventType = "POINTS_AWARDED"
)

type Envelope struct {
	Type    EventType       `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Subscrib2er struct {
	logger *slog.Logger
	rdb    *redis.Client
	ws     *ws.Server
	patt   string
}

func NewSubscriber(logger *slog.Logger, rdb *redis.Client, wsSrv *ws.Server, pattern string) *Subscriber {
	return &Subscriber{logger: logger, rdb: rdb, ws: wsSrv, patt: pattern}
}

func (s *Subscriber) Run(ctx context.Context) error {
	pubsub := s.rdb.PSubscribe(ctx, s.patt)
	defer pubsub.Close()

	s.logger.Info("redis subscriber started", slog.String("pattern", s.patt))

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-ch:
			if !ok {
				return nil
			}
			s.handle(msg.Channel, msg.Payload)
		}
	}
}

func (s *Subscriber) handle(channel, payload string) {
	roomID, err := roomIDFromChannel(channel)
	if err != nil {
		s.logger.Warn("redis: bad channel", slog.String("channel", channel))
		return
	}

	var env Envelope
	if err := json.Unmarshal([]byte(payload), &env); err != nil {
		s.logger.Warn("redis: bad envelope", slog.String("err", err.Error()))
		return
	}

	switch env.Type {
	case EventMessageCreated, EventVoteUpdated, EventPointsAwarded:
		var raw any
		if len(env.Payload) > 0 {
			_ = json.Unmarshal(env.Payload, &raw)
		}
		s.ws.BroadcastToRoom(roomID, string(env.Type), raw)
	default:
		s.logger.Warn("redis: unknown event type", slog.String("type", string(env.Type)))
	}
}

func roomIDFromChannel(ch string) (uuid.UUID, error) {
	const prefix = "room:"
	if !strings.HasPrefix(ch, prefix) {
		return uuid.Nil, errors.New("channel must start with 'room:'")
	}
	return uuid.Parse(ch[len(prefix):])
}

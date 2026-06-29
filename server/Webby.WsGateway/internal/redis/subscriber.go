package redisbus

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"webby/wsgateway/internal/ws"
)

type broadcaster interface {
	Broadcast(userID uuid.UUID, payload []byte)
}

type envelope struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type subscriber struct {
	logger      *slog.Logger
	redisClient *redis.Client
	wsServer    *ws.Server
	broadcaster broadcaster
	pattern     string
}

func NewSubscriber(logger *slog.Logger, redisClient *redis.Client, wsServer *ws.Server, broadcaster broadcaster, pattern string) *subscriber {
	return &subscriber{
		logger:      logger,
		redisClient: redisClient,
		wsServer:    wsServer,
		broadcaster: broadcaster,
		pattern:     pattern,
	}
}

func (s *subscriber) Run(ctx context.Context) error {
	pubsub := s.redisClient.PSubscribe(ctx, s.pattern)
	defer pubsub.Close()

	s.logger.Info("redis subscriber started", slog.String("pattern", s.pattern))

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-ch:
			if !ok {
				return nil
			}

			if strings.HasPrefix(msg.Channel, "user:") {
				s.handleUser(msg.Channel, msg.Payload)
			} else if strings.HasPrefix(msg.Channel, "chat:") {
				s.handleChat(msg.Channel, msg.Payload)
			} else {
				s.logger.Warn("redis: unknown channel prefix", slog.String("channel", msg.Channel))
			}
		}
	}
}

func (s *subscriber) handleUser(channel, payload string) {
	userID, err := userIDFromChannel(channel)
	if err != nil {
		s.logger.Warn("redis: bad channel", slog.String("channel", channel))
		return
	}

	var env envelope
	if err := json.Unmarshal([]byte(payload), &env); err != nil {
		s.logger.Warn("redis: bad envelope", slog.String("err", err.Error()))
		return
	}

	if strings.TrimSpace(env.Type) == "" {
		s.logger.Warn("redis: empty event type")
		return
	}

	s.broadcaster.Broadcast(userID, []byte(env.Type))
}

func (s *subscriber) handleChat(channel, payload string) {
	chatID, err := chatIDFromChannel(channel)
	if err != nil {
		s.logger.Warn("redis: bad channel", slog.String("channel", channel))
		return
	}

	var env envelope
	if err := json.Unmarshal([]byte(payload), &env); err != nil {
		s.logger.Warn("redis: bad envelope", slog.String("err", err.Error()))
		return
	}

	if strings.TrimSpace(env.Type) == "" {
		s.logger.Warn("redis: empty event type")
		return
	}

	var raw any
	if len(env.Payload) > 0 {
		if err := json.Unmarshal(env.Payload, &raw); err != nil {
			s.logger.Warn("redis: bad payload json", slog.String("err", err.Error()))
			return
		}
	}

	s.wsServer.BroadcastToRoom(chatID, env.Type, raw)
}

func chatIDFromChannel(ch string) (uuid.UUID, error) {
	const prefix = "chat:"
	if !strings.HasPrefix(ch, prefix) {
		return uuid.Nil, errors.New("channel must start with 'chat:'")
	}
	return uuid.Parse(ch[len(prefix):])
}

func userIDFromChannel(ch string) (uuid.UUID, error) {
	const prefix = "user:"
	if !strings.HasPrefix(ch, prefix) {
		return uuid.Nil, errors.New("channel must start with 'user:'")
	}
	return uuid.Parse(ch[len(prefix):])
}

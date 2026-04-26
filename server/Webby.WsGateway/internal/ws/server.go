package ws

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	socketio "github.com/googollee/go-socket.io"

	clients "webby-wsgateway/internal/grpc"
)

const grpcCallTimeout = 5 * time.Second

type session struct {
	UserID uuid.UUID
	RoomID uuid.UUID
}

type Server struct {
	io *socketio.Server

	logger    *slog.Logger
	jwtSecret []byte
	chat      *clients.ChatClient

	mu sync.RWMutex
	rooms map[uuid.UUID]map[uuid.UUID]int
}

func NewServer(logger *slog.Logger, jwtSecret []byte, chat *clients.ChatClient) *Server {
	s := &Server{
		io:        socketio.NewServer(nil),
		logger:    logger,
		jwtSecret: jwtSecret,
		chat:      chat,
		rooms:     make(map[uuid.UUID]map[uuid.UUID]int),
	}
	s.registerHandlers()
	return s
}

func (s *Server) IO() *socketio.Server { return s.io }

func (s *Server) Serve() error { return s.io.Serve() }
func (s *Server) Close() error { return s.io.Close() }

func (s *Server) BroadcastToRoom(roomID uuid.UUID, event string, payload any) {
	s.io.BroadcastToRoom("/", roomID.String(), event, payload)
}

func (s *Server) SnapshotActivity() map[uuid.UUID][]uuid.UUID {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make(map[uuid.UUID][]uuid.UUID, len(s.rooms))
	for roomID, users := range s.rooms {
		if len(users) == 0 {
			continue
		}
		ids := make([]uuid.UUID, 0, len(users))
		for userID := range users {
			ids = append(ids, userID)
		}
		out[roomID] = ids
	}
	return out
}

// --- handlers -------------------------------------------------------------

func (s *Server) registerHandlers() {
	s.io.OnConnect("/", s.onConnect)
	s.io.OnDisconnect("/", s.onDisconnect)
	s.io.OnEvent("/", "send_message", s.onSendMessage)
}

func (s *Server) onConnect(c socketio.Conn) error {
	u := c.URL()
	q := u.Query()
	token := q.Get("token")
	roomIDStr := q.Get("room_id")

	if token == "" || roomIDStr == "" {
		return errors.New("token and room_id are required")
	}

	userID, err := s.parseUserID(token)
	if err != nil {
		s.logger.Warn("ws auth failed", slog.String("err", err.Error()))
		return errors.New("unauthorized")
	}

	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		return errors.New("invalid room_id")
	}

	c.SetContext(session{UserID: userID, RoomID: roomID})
	c.Join(roomID.String())
	s.track(roomID, userID, +1)

	s.logger.Info("ws connected",
		slog.String("user_id", userID.String()),
		slog.String("room_id", roomID.String()),
	)
	return nil
}

func (s *Server) onDisconnect(c socketio.Conn, reason string) {
	sess, ok := c.Context().(session)
	if !ok {
		return
	}
	s.track(sess.RoomID, sess.UserID, -1)
	s.logger.Info("ws disconnected",
		slog.String("user_id", sess.UserID.String()),
		slog.String("room_id", sess.RoomID.String()),
		slog.String("reason", reason),
	)
}

func (s *Server) onSendMessage(c socketio.Conn, data map[string]string) string {
	sess, ok := c.Context().(session)
	if !ok {
		return `{"error":"unauthorized"}`
	}

	content := data["content"]
	if content == "" {
		return `{"error":"content is required"}`
	}

	ctx, cancel := context.WithTimeout(context.Background(), grpcCallTimeout)
	defer cancel()

	msgID, err := s.chat.SaveMessage(ctx, sess.RoomID.String(), sess.UserID.String(), content)
	if err != nil {
		s.logger.Error("SaveMessage failed", slog.String("err", err.Error()))
		return `{"error":"failed to send"}`
	}
	return `{"ok":true,"message_id":"` + msgID + `"}`
}

// --- internals ------------------------------------------------------------

func (s *Server) track(roomID, userID uuid.UUID, delta int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	users := s.rooms[roomID]
	if users == nil {
		if delta <= 0 {
			return
		}
		users = make(map[uuid.UUID]int)
		s.rooms[roomID] = users
	}
	users[userID] += delta
	if users[userID] <= 0 {
		delete(users, userID)
	}
	if len(users) == 0 {
		delete(s.rooms, roomID)
	}
}

func (s *Server) parseUserID(tokenStr string) (uuid.UUID, error) {
	tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil || !tok.Valid {
		return uuid.Nil, errors.New("invalid token")
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("invalid claims")
	}
	
	for _, key := range []string{"sub", "user_id", "userId", "uid"} {
		if v, ok := claims[key].(string); ok && v != "" {
			id, err := uuid.Parse(v)
			if err != nil {
				return uuid.Nil, err
			}
			return id, nil
		}
	}
	return uuid.Nil, errors.New("user id claim not found")
}

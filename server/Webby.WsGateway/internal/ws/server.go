package ws

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

type SessionInfo struct {
	ChatID string
	UserID string
}

type presenceManager interface {
	AddUser(ctx context.Context, chatID, userID string) error
	RemoveUser(ctx context.Context, chatID, userID string) error
	UpdateHeartbeats(ctx context.Context, sessions []SessionInfo) error
}

type session struct {
	ChatID uuid.UUID
	UserID uuid.UUID
}

type userIDRetriever interface {
	GetUserID(ctx context.Context, token string) (uuid.UUID, error)
}

type Server struct {
	io *socketio.Server

	service         userIDRetriever
	presenceManager presenceManager
	logger          *slog.Logger
	callTimeout     time.Duration

	activeSessions sync.Map
}

func NewServer(service userIDRetriever, presenceManager presenceManager, logger *slog.Logger, callTimeout time.Duration) *Server {
	s := &Server{
		io: socketio.NewServer(&engineio.Options{
			Transports: []transport.Transport{
				&polling.Transport{
					CheckOrigin: func(r *http.Request) bool {
						return true
					},
				},
				&websocket.Transport{
					CheckOrigin: func(r *http.Request) bool {
						return true
					},
				},
			},
		}),

		service:         service,
		presenceManager: presenceManager,
		logger:          logger,
		callTimeout:     callTimeout,
	}
	s.registerHandlers()
	return s
}

func (server *Server) IO() *socketio.Server { return server.io }

func (server *Server) Serve() error { return server.io.Serve() }
func (server *Server) Close() error { return server.io.Close() }

func (server *Server) BroadcastToRoom(chatID uuid.UUID, event string, payload any) {
	server.io.BroadcastToRoom("/", chatID.String(), event, payload)
}

func (server *Server) registerHandlers() {
	server.io.OnConnect("/", server.onConnect)
	server.io.OnDisconnect("/", server.onDisconnect)
}

func (server *Server) onConnect(c socketio.Conn) error {
	const op = "ws.server.onConnect"
	log := server.logger.With("op", op)
	ctx, cancel := context.WithTimeout(context.Background(), server.callTimeout)
	defer cancel()

	url := c.URL()
	query := url.Query()

	token := query.Get("token")
	if token == "" {
		return errors.New("token is required")
	}

	chatID, err := uuid.Parse(query.Get("chat_id"))
	if err != nil {
		return errors.New("invalid chat_id")
	}

	userID, err := server.service.GetUserID(ctx, token)
	if err != nil {
		log.Warn("auth failed", slog.String("err", err.Error()))

		return errors.New("internal server error")
	}

	c.SetContext(session{UserID: userID, ChatID: chatID})
	c.Join(chatID.String())

	server.activeSessions.Store(c.ID(), SessionInfo{
		ChatID: chatID.String(),
		UserID: userID.String(),
	})

	err = server.presenceManager.AddUser(ctx, chatID.String(), userID.String())
	if err != nil {
		log.Error("failed to add presence", slog.String("err", err.Error()))
	}

	log.Info("ws connected",
		slog.String("user_id", userID.String()),
		slog.String("chat_id", chatID.String()),
	)
	return nil
}

func (server *Server) onDisconnect(c socketio.Conn, reason string) {
	const op = "ws.server.onDisconnect"
	log := server.logger.With("op", op)

	sess, ok := c.Context().(session)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), server.callTimeout)
	defer cancel()

	if val, ok := server.activeSessions.LoadAndDelete(c.ID()); ok {
		info := val.(SessionInfo)
		err := server.presenceManager.RemoveUser(ctx, info.ChatID, info.UserID)
		if err != nil {
			log.Error("failed to remove presence", slog.String("err", err.Error()))
		}
	}

	log.Info("ws disconnected",
		slog.String("user_id", sess.UserID.String()),
		slog.String("chat_id", sess.ChatID.String()),
		slog.String("reason", reason),
	)
}

func (server *Server) RunHeartbeat(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			server.syncPresence(ctx)
		}
	}
}

func (server *Server) syncPresence(ctx context.Context) {
	var sessions []SessionInfo

	server.activeSessions.Range(func(key, value any) bool {
		sessions = append(sessions, value.(SessionInfo))
		return true
	})

	if len(sessions) > 0 {
		err := server.presenceManager.UpdateHeartbeats(ctx, sessions)
		if err != nil {
			server.logger.Error("failed to bulk update presence", slog.String("err", err.Error()))
		}
	}
}

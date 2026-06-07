package ws

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"

	clients "webby/wsgateway/internal/grpc"

	"webby/wsgateway/internal/domain"
)

type session struct {
	UserID uuid.UUID
	ChatID uuid.UUID
}

type Service interface {
	GetUserID(ctx context.Context, token string) (uuid.UUID, error)
}

type Server struct {
	io *socketio.Server

	service     Service
	logger      *slog.Logger
	jwtSecret   []byte
	chat        *clients.ChatClient
	callTimeout time.Duration
}

func NewServer(service Service, logger *slog.Logger, jwtSecret []byte, chat *clients.ChatClient, callTimeout time.Duration) *Server {
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

		service:     service,
		logger:      logger,
		jwtSecret:   jwtSecret,
		chat:        chat,
		callTimeout: callTimeout,
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

// --- handlers -------------------------------------------------------------

func (server *Server) registerHandlers() {
	server.io.OnConnect("/", server.onConnect)
	server.io.OnDisconnect("/", server.onDisconnect)
}

func (server *Server) onConnect(c socketio.Conn) error {
	ctx, cancel := context.WithTimeout(context.Background(), server.callTimeout)
	defer cancel()

	url := c.URL()
	query := url.Query()

	token := query.Get("token")
	chatIDStr := query.Get("chat_id")

	if token == "" || chatIDStr == "" {
		return errors.New("token and chat_id are required")
	}

	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		return errors.New("invalid chat_id")
	}

	userID, err := server.service.GetUserID(ctx, token)
	if err != nil {
		server.logger.Warn("auth failed", slog.String("err", err.Error()))

		if errors.Is(err, domain.ErrUserNotFound) {
			return errors.New("unauthorized")
		}

		return errors.New("internal server error")
	}

	c.SetContext(session{UserID: userID, ChatID: chatID})
	c.Join(chatID.String())

	server.logger.Info("ws connected",
		slog.String("user_id", userID.String()),
		slog.String("chat_id", chatID.String()),
	)
	return nil
}

func (server *Server) onDisconnect(c socketio.Conn, reason string) {
	sess, ok := c.Context().(session)
	if !ok {
		return
	}
	server.logger.Info("ws disconnected",
		slog.String("user_id", sess.UserID.String()),
		slog.String("chat_id", sess.ChatID.String()),
		slog.String("reason", reason),
	)
}

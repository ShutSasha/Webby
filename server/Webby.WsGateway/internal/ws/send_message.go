package ws

import (
	"context"
	"log/slog"

	socketio "github.com/googollee/go-socket.io"
)

func (server *Server) onSendMessage(c socketio.Conn, data map[string]string) Response {
	sess, ok := c.Context().(session)
	if !ok {
		return Response{OK: false, Error: "unauthorized"}
	}

	content := data["content"]
	if content == "" {
		return Response{OK: false, Error: "content is required"}
	}

	ctx, cancel := context.WithTimeout(context.Background(), callTimeout)
	defer cancel()

	msgID, err := server.chat.SaveMessage(ctx, sess.ChatID.String(), sess.UserID.String(), content)
	if err != nil {
		server.logger.Error("SaveMessage failed",
			slog.String("err", err.Error()),
			slog.String("user_id", sess.UserID.String()),
		)
		return Response{OK: false, Error: "failed to save message"}
	}

	return Response{OK: true, MessageID: msgID}
}

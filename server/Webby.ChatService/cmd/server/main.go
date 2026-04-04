package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"
	"webby-chat/internal/config"
	"webby-chat/internal/database"
	grpcserver "webby-chat/internal/grpc"
	"webby-chat/internal/grpc/chatpb"
	httpserver "webby-chat/internal/handlers"
	"webby-chat/internal/repository"
	"webby-chat/internal/services"
	"webby-chat/internal/ws"
	"webby-chat/pkg/slogpretty"

	"google.golang.org/grpc"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// @Version 1.0
// @Title Webby.ChatService
// @Description This API provides endpoints for real-time chat in rooms via Socket.IO and REST.
// @Security BearerAuth
// @SecurityScheme BearerAuth http bearer Enter your JWT token
func main() {
	ctx := context.Background()

	if err := run(ctx, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, w io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	cfg := config.MustLoad()
	logger := setupLogger(cfg.Env, w)

	db, err := database.New(cfg.ConnectionString)
	if err != nil {
		logger.Error("database connection failed", slog.String("error", err.Error()))
		return err
	}
	defer db.Close()

	logger.Info("database connected successfully")

	// Repositories
	chatRepo := repository.NewChatRepository(db)
	chatMemberRepo := repository.NewChatMemberRepository(db)
	messageRepo := repository.NewMessageRepository(db)

	// Services
	chatService := services.NewChatService(chatRepo, chatMemberRepo)
	messageService := services.NewMessageService(messageRepo, chatMemberRepo)

	// WebSocket
	hubManager := ws.NewHubManager()
	socketServer := ws.SetupSocketIO(logger, []byte(cfg.JwtSecret), hubManager, chatService, messageService)
	go func() {
		if err := socketServer.Serve(); err != nil {
			logger.Error("socket.io serve error", slog.String("error", err.Error()))
		}
	}()
	defer socketServer.Close()

	// HTTP server
	server := httpserver.NewServer(cfg, logger, chatService, socketServer)
	httpServer := &http.Server{
		Addr:         net.JoinHostPort(cfg.Http.Host, strconv.Itoa(cfg.Http.Port)),
		ReadTimeout:  cfg.Http.Timeout,
		WriteTimeout: cfg.Http.Timeout,
		Handler:      server,
	}

	go func() {
		logger.Info("HTTP server listening",
			slog.String("host", cfg.Http.Host),
			slog.Int("port", cfg.Http.Port),
		)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("error listening and serving", slog.Any("error", err))
		}
	}()

	// gRPC server
	grpcSrv := grpc.NewServer()
	chatGrpcServer := grpcserver.NewChatGrpcServer(chatService, logger)
	chatpb.RegisterChatGrpcServiceServer(grpcSrv, chatGrpcServer)

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Grpc.Port))
	if err != nil {
		logger.Error("failed to listen for gRPC", slog.String("error", err.Error()))
		return err
	}

	go func() {
		logger.Info("gRPC server listening", slog.Int("port", cfg.Grpc.Port))
		if err := grpcSrv.Serve(grpcListener); err != nil {
			logger.Error("gRPC server error", slog.String("error", err.Error()))
		}
	}()

	// Graceful shutdown
	var wg sync.WaitGroup
	wg.Go(func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		grpcSrv.GracefulStop()
		logger.Info("gRPC server stopped")

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error("error shutting down http server", slog.Any("error", err))
		}
		logger.Info("HTTP server stopped gracefully")
	})
	wg.Wait()

	return nil
}

func setupLogger(env string, w io.Writer) *slog.Logger {
	var l *slog.Logger

	switch env {
	case envLocal:
		l = setupPrettySlog()
	case envDev:
		l = slog.New(
			slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		l = slog.New(
			slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		l = slog.New(
			slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return l
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}

// suppress unused import warning
var _ = log.Println

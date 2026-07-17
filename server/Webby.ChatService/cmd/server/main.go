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
	"webby/chat-service/internal/config"
	"webby/chat-service/internal/database"
	grpcserver "webby/chat-service/internal/grpc"
	"webby/chat-service/internal/grpc/chatpb"
	handlers "webby/chat-service/internal/handlers"
	"webby/chat-service/internal/repository"
	"webby/chat-service/internal/services"
	"webby/chat-service/pkg/slogpretty"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

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

	roomMemberClient, err := grpcserver.NewMemberClient(cfg.Grpc.RoomServiceAddress)
	if err != nil {
		logger.Error(
			"room service gRPC connection failed",
			slog.String("error", err.Error()),
		)
		return err
	}
	defer roomMemberClient.Close()

	userClient, err := grpcserver.NewUserClient(cfg.Grpc.UserServiceAddress)
	if err != nil {
		logger.Error(
			"user service gRPC connection failed",
			slog.String("error", err.Error()),
		)
		return err
	}
	defer userClient.Close()

	// Redis publisher for event dispatch
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer rdb.Close()

	// Repositories
	chatRepo := repository.NewChatRepository(db)
	chatMemberRepo := repository.NewChatMemberRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	redisPublisher := repository.NewRedisPublisher(rdb)

	// Services
	chatService := services.NewChatService(chatRepo, chatMemberRepo, roomMemberClient, userClient, redisPublisher)
	chatMemberService := services.NewChatMemberService(chatRepo, chatMemberRepo)
	messageService := services.NewMessageService(messageRepo, chatMemberRepo, chatRepo, userClient, redisPublisher)

	// HTTP server
	server := handlers.NewServer(cfg, logger, chatService, messageService)
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
	chatGrpcServer := grpcserver.NewChatGrpcServer(chatService, chatMemberService, logger)
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

package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"
	"webby/room-queue-service/internal/config"
	"webby/room-queue-service/internal/database"
	grpcserver "webby/room-queue-service/internal/grpc"
	"webby/room-queue-service/internal/grpc/queuepb"
	httpserver "webby/room-queue-service/internal/handlers"
	"webby/room-queue-service/internal/publisher"
	"webby/room-queue-service/internal/repository"
	"webby/room-queue-service/internal/services"
	"webby/room-queue-service/pkg/slogpretty"

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
		logger.Error(
			"database connection failed",
			slog.String("error", err.Error()),
		)
		return err
	}
	defer db.Close()

	logger.Info("database connected successfully")

	queueItemRepo := repository.NewQueueItemRepository(db)

	mediaClient, err := grpcserver.NewMediaClient(
		cfg.Grpc.MediaServiceAddress,
	)
	if err != nil {
		logger.Error(
			"media service gRPC connection failed",
			slog.String("error", err.Error()),
		)
		return err
	}
	defer mediaClient.Close()

	chatClient, err := grpcserver.NewChatClient(cfg.Grpc.ChatServiceAddress)
	if err != nil {
		logger.Error(
			"chat service gRPC connection failed",
			slog.String("error", err.Error()),
		)
		return err
	}
	defer chatClient.Close()

	memberClient, err := grpcserver.NewMemberClient(cfg.Grpc.RoomServiceAddress)
	if err != nil {
		logger.Error("room service gRPC connection failed", slog.String("error", err.Error()))
		return err
	}
	defer memberClient.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer rdb.Close()

	redisPublisher := publisher.New(rdb)
	queueService := services.New(
		queueItemRepo,
		mediaClient, chatClient, memberClient,
		redisPublisher,
	)

	server := httpserver.NewServer(cfg, logger, queueService)
	httpServer := &http.Server{
		Addr: net.JoinHostPort(
			cfg.Http.Host, strconv.Itoa(cfg.Http.Port),
		),
		ReadTimeout:  cfg.Http.Timeout,
		WriteTimeout: cfg.Http.Timeout,
		Handler:      server,
	}

	go func() {
		logger.Info(
			"Server listening",
			slog.String("host", cfg.Http.Host),
			slog.Int("port", cfg.Http.Port),
		)
		if err := httpServer.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			logger.Error(
				"error listening and serving",
				slog.Any("error", err),
			)
		}
	}()

	grpcListener, err := net.Listen(
		"tcp",
		net.JoinHostPort(cfg.Grpc.Host, strconv.Itoa(cfg.Grpc.Port)),
	)
	if err != nil {
		logger.Error("grpc listen failed", slog.Any("error", err))
		return err
	}

	grpcSrv := grpc.NewServer()
	queuepb.RegisterQueueGrpcServiceServer(
		grpcSrv,
		grpcserver.NewQueueServer(queueService),
	)

	go func() {
		logger.Info(
			"gRPC server listening",
			slog.String("host", cfg.Grpc.Host),
			slog.Int("port", cfg.Grpc.Port),
		)
		if err := grpcSrv.Serve(grpcListener); err != nil {
			logger.Error(
				"error serving grpc",
				slog.Any("error", err),
			)
		}
	}()

	var wg sync.WaitGroup
	wg.Go(func() {
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(
			shutdownCtx, 10*time.Second,
		)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error(
				"error shutting down http server",
				slog.Any("error", err),
			)
		}
		grpcSrv.GracefulStop()
		logger.Info("server stopped gracefully")
	})
	wg.Wait()

	return nil
}

func setupLogger(env string, w io.Writer) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(
				w, &slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(
				w, &slog.HandlerOptions{Level: slog.LevelInfo},
			),
		)
	}

	return log
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

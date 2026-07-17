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
	"webby/room-service/internal/config"
	"webby/room-service/internal/database"
	grpcClient "webby/room-service/internal/grpc"
	"webby/room-service/internal/grpc/memberpb"
	"webby/room-service/internal/grpc/roompb"
	"webby/room-service/internal/handlers"
	"webby/room-service/internal/publisher"
	"webby/room-service/internal/repository"
	"webby/room-service/internal/services"
	"webby/room-service/internal/workers"
	"webby/room-service/pkg/slogpretty"

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

	roomRepository := repository.NewRoomRepository(db)
	roomMemberRepository := repository.NewRoomMemberRepository(db)
	reactionRepository := repository.NewReactionRepository(db)
	fileStorage := repository.NewFileStorage(cfg)

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer rdb.Close()

	publisher := publisher.New(rdb)
	timecodesRepository := repository.NewTimecodesRepo(rdb)
	roomPresenceRepository := repository.NewRedisPresenceRepository(rdb)

	chatClient, err := grpcClient.NewChatClient(cfg.Grpc.ChatServiceAddress)
	if err != nil {
		logger.Error("chat service gRPC connection failed", slog.String("error", err.Error()))
		return err

	}
	defer chatClient.Close()

	categoryClient, err := grpcClient.NewCategoryClient(cfg.Grpc.CategoryServiceAddress)
	if err != nil {
		logger.Error("category service gRPC connection failed", slog.String("error", err.Error()))
		return err
	}
	defer categoryClient.Close()

	notificationClient, err := grpcClient.NewNotificationClient(cfg.Grpc.NotificationServiceAddress)
	if err != nil {
		logger.Error("notification service gRPC connection failed", slog.String("error", err.Error()))
		return err
	}
	defer notificationClient.Close()

	logger.Info("repositories initialized")

	presenceSerivce := services.NewPresenceService(roomPresenceRepository, chatClient)
	roomService := services.NewRoomService(roomRepository, roomMemberRepository, fileStorage, categoryClient, chatClient)
	roomMemberService := services.NewRoomMemberService(roomRepository, roomMemberRepository, chatClient, notificationClient, chatClient, publisher)
	reactionService := services.NewReactionService(reactionRepository, roomMemberRepository, chatClient, roomMemberRepository, publisher)
	syncService := services.NewSynchronizeService(chatClient, publisher, timecodesRepository, roomMemberRepository)

	server := handlers.NewServer(cfg, logger, roomService, roomMemberService, syncService, reactionService)
	httpServer := &http.Server{
		Addr:         net.JoinHostPort(cfg.Http.Host, strconv.Itoa(cfg.Http.Port)),
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
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("error listening and serving", slog.Any("error", err))
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
	memberpb.RegisterMemberGrpcServiceServer(
		grpcSrv,
		grpcClient.NewMemberServer(roomMemberRepository),
	)
	roompb.RegisterRoomGrpcServiceServer(
		grpcSrv,
		grpcClient.NewRoomServer(roomRepository),
	)

	go func() {
		logger.Info(
			"gRPC server listening",
			slog.String("host", cfg.Grpc.Host),
			slog.Int("port", cfg.Grpc.Port),
		)
		if err := grpcSrv.Serve(grpcListener); err != nil {
			logger.Error("error serving grpc", slog.Any("error", err))
		}
	}()

	worker := workers.NewPointsWorker(presenceSerivce, roomMemberRepository, publisher, logger, cfg.Worker.Interval, cfg.Worker.PointsPerTick, cfg.Worker.ZombieTTL)

	go worker.Run(ctx)

	var wg sync.WaitGroup
	wg.Go(func() {
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error("error shutting down http server", slog.Any("error", err))
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
			slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo}),
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

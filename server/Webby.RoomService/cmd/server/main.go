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
	"webby/internal/config"
	"webby/internal/database"
	grpcClient "webby/internal/grpc"
	"webby/internal/grpc/memberpb"
	httpserver "webby/internal/handlers"
	"webby/internal/repository"
	"webby/internal/services"
	"webby/pkg/slogpretty"

	"google.golang.org/grpc"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// @Version 1.0
// @Title Webby.RoomService
// @Description This API provides endpoints for managing rooms and categories.
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

	config := config.MustLoad()
	logger := setupLogger(config.Env, w)

	db, err := database.New(config.ConnectionString)
	if err != nil {
		logger.Error("database connection failed", slog.String("error", err.Error()))
		return err
	}

	defer db.Close()

	logger.Info("database connected successfully")

	roomRepository := repository.NewRoomRepository(db)
	roomMemberRepository := repository.NewRoomMemberRepository(db)
	fileStorage := repository.NewFileStorage(config)

	mediaClient, err := grpcClient.NewMediaClient(config.Grpc.MediaServiceAddress)
	if err != nil {
		logger.Error("media service gRPC connection failed", slog.String("error", err.Error()))
		return err
	}
	defer mediaClient.Close()

	chatClient, err := grpcClient.NewChatClient(config.Grpc.ChatServiceAddress)
	if err != nil {
		logger.Warn("chat service gRPC connection failed — chat features disabled", slog.String("error", err.Error()))
		chatClient = nil
	} else {
		defer chatClient.Close()
	}

	categoryClient, err := grpcClient.NewCategoryClient(config.Grpc.CategoryServiceAddress)
	if err != nil {
		logger.Error("category service gRPC connection failed", slog.String("error", err.Error()))
		return err
	}
	defer categoryClient.Close()

	queueClient, err := grpcClient.NewQueueClient(config.Grpc.QueueServiceAddress)
	if err != nil {
		logger.Warn("queue service gRPC connection failed — vote queue features disabled", slog.String("error", err.Error()))
		queueClient = nil
	} else {
		defer queueClient.Close()
	}

	roomService := services.NewRoomService(roomRepository, roomMemberRepository, fileStorage, chatClient, categoryClient)
	voteRepository := repository.NewVoteRepository(db)
	voteService := services.NewVoteService(voteRepository, roomRepository, roomMemberRepository, queueClient)

	logger.Info("repositories initialized")

	server := httpserver.NewServer(
		config,
		logger,
		roomService,
		voteService,
	)
	httpServer := &http.Server{
		Addr:         net.JoinHostPort(config.Http.Host, strconv.Itoa(config.Http.Port)),
		ReadTimeout:  config.Http.Timeout,
		WriteTimeout: config.Http.Timeout,
		Handler:      server,
	}

	go func() {
		logger.Info(
			"Server listening",
			slog.String("host", config.Http.Host),
			slog.Int("port", config.Http.Port),
		)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("error listening and serving", slog.Any("error", err))
		}
	}()

	grpcListener, err := net.Listen(
		"tcp",
		net.JoinHostPort(config.Grpc.Host, strconv.Itoa(config.Grpc.Port)),
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

	go func() {
		logger.Info(
			"gRPC server listening",
			slog.String("host", config.Grpc.Host),
			slog.Int("port", config.Grpc.Port),
		)
		if err := grpcSrv.Serve(grpcListener); err != nil {
			logger.Error("error serving grpc", slog.Any("error", err))
		}
	}()

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

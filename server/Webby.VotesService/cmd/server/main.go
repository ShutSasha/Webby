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
	"webby-vote-service/internal/config"
	"webby-vote-service/internal/database"
	grpcserver "webby-vote-service/internal/grpc"
	"webby-vote-service/internal/grpc/votepb"
	httpserver "webby-vote-service/internal/handlers"
	"webby-vote-service/internal/repository"
	"webby-vote-service/internal/services"
	"webby-vote-service/pkg/slogpretty"

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

	voteRepo := repository.NewVoteRepository(db)

	memberClient, err := grpcserver.NewMemberClient(
		cfg.Grpc.RoomServiceAddress,
	)
	if err != nil {
		logger.Error(
			"room service gRPC connection failed (member)",
			slog.String("error", err.Error()),
		)
		return err
	}
	defer memberClient.Close()

	roomClient, err := grpcserver.NewRoomClient(
		cfg.Grpc.RoomServiceAddress,
	)
	if err != nil {
		logger.Error(
			"room service gRPC connection failed (room)",
			slog.String("error", err.Error()),
		)
		return err
	}
	defer roomClient.Close()

	var queueClient *grpcserver.QueueClient
	if cfg.Grpc.QueueServiceAddress != "" {
		queueClient, err = grpcserver.NewQueueClient(
			cfg.Grpc.QueueServiceAddress,
		)
		if err != nil {
			logger.Warn(
				"queue service gRPC connection failed "+
					"— vote queue features disabled",
				slog.String("error", err.Error()),
			)
			queueClient = nil
		} else {
			defer queueClient.Close()
		}
	}

	voteService := services.New(
		voteRepo, memberClient, roomClient, queueClient,
	)

	logger.Info("services initialized")

	server := httpserver.NewServer(cfg, logger, voteService)
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
		net.JoinHostPort(
			cfg.Grpc.Host, strconv.Itoa(cfg.Grpc.Port),
		),
	)
	if err != nil {
		logger.Error("grpc listen failed", slog.Any("error", err))
		return err
	}

	grpcSrv := grpc.NewServer()
	votepb.RegisterVoteGrpcServiceServer(
		grpcSrv,
		grpcserver.NewVoteServer(voteRepo),
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

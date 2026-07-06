package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"
	"webby/admin-service/internal/config"
	"webby/admin-service/internal/database"
	grpcClient "webby/admin-service/internal/grpc"
	"webby/admin-service/internal/handlers"
	httpserver "webby/admin-service/internal/handlers"
	"webby/admin-service/internal/repository"
	"webby/admin-service/internal/services"
	"webby/admin-service/pkg/slogpretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	ctx := context.Background()

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	cfg := config.MustLoad()

	var logWriter io.Writer = os.Stdout

	if cfg.Env == envDev || cfg.Env == envProd {
		if err := os.MkdirAll("logs", 0755); err != nil {
			return fmt.Errorf("failed to create logs directory: %w", err)
		}

		fileName := fmt.Sprintf("app_%s.log", time.Now().Format("2006-01-02_15-04-05"))
		logPath := filepath.Join("logs", fileName)

		logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}

		defer logFile.Close()
		logWriter = logFile
	}

	logger := setupLogger(cfg.Env, logWriter)

	db, err := database.New(cfg.ConnectionString)
	if err != nil {
		logger.Error("database connection failed", slog.String("error", err.Error()))
		return err
	}

	defer db.Close()

	logger.Info("database connected successfully")

	complaintRepository := repository.New(db)

	complaintClient, err := grpcClient.NewComplaintClient(cfg.Grpc.ComplaintServiceAddress)
	if err != nil {
		logger.Error("complaint service gRPC connection failed", slog.String("error", err.Error()))
		return err
	}
	defer complaintClient.Close()

	mediaClient, err := grpcClient.NewMediaClient(cfg.Grpc.MediaServiceAddress)
	if err != nil {
		logger.Error("media gRPC connection failed", slog.String("error", err.Error()))
		return err
	}
	defer mediaClient.Close()

	videoClient, err := grpcClient.NewVideoClient(cfg.Grpc.MediaServiceAddress)
	if err != nil {
		logger.Error("video service gRPC connection failed", slog.String("error", err.Error()))
		return err
	}
	defer videoClient.Close()

	userClient, err := grpcClient.NewUserClient(cfg.Grpc.ComplaintServiceAddress)
	if err != nil {
		logger.Error("user service gRPC connection failed", slog.String("error", err.Error()))
		return err
	}
	defer userClient.Close()

	notificationClient, err := grpcClient.NewNotificationClient(cfg.Grpc.NotificationServiceAddress)
	if err != nil {
		logger.Error("notification service gRPC connection failed", slog.String("error", err.Error()))
		return err
	}
	defer notificationClient.Close()

	roomClient, err := grpcClient.NewRoomClient(cfg.Grpc.RoomServiceAddress)
	if err != nil {
		logger.Error("room service gRPC connection failed", slog.String("error", err.Error()))
		return err
	}
	defer roomClient.Close()

	complaintService := services.NewComplaintsService(complaintRepository, complaintClient, mediaClient, videoClient, userClient, notificationClient)
	statsService := services.NewStatsService(userClient, roomClient, complaintRepository)
	// Health checkers
	databaseChecker := handlers.NewDatabaseChecker(db)
	server := httpserver.NewServer(cfg, logger, complaintService, statsService, databaseChecker)
	httpServer := &http.Server{
		Addr:         cfg.Http.HostPort,
		ReadTimeout:  cfg.Http.Timeout,
		WriteTimeout: cfg.Http.Timeout,
		Handler:      server,
	}

	go func() {
		logger.Info("Server listening", "host", cfg.Http.HostPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("error listening and serving", slog.Any("error", err))
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
		logger.Info("server stopped gracefully")
	})
	wg.Wait()

	return nil
}

func setupLogger(env string, w io.Writer) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog(w)
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

func setupPrettySlog(w io.Writer) *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(w)

	return slog.New(handler)
}

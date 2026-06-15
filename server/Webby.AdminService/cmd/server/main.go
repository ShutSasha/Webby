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
	"webby/admin-service/internal/config"
	"webby/admin-service/internal/database"
	grpcClient "webby/admin-service/internal/grpc"
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

	repository := repository.New(db)

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

	notificationClient, err := grpcClient.NewNotificationClient(cfg.Grpc.NotificationServiceAddress)
	if err != nil {
		logger.Error("notification service gRPC connection failed", slog.String("error", err.Error()))
		return err
	}
	defer notificationClient.Close()

	categoryService := services.New(repository, complaintClient, mediaClient, notificationClient)

	server := httpserver.NewServer(cfg, logger, categoryService)
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

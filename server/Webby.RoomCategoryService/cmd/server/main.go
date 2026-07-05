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
	"path/filepath"
	"sync"
	"time"
	"webby/room-category-service/internal/config"
	"webby/room-category-service/internal/database"
	grpcserver "webby/room-category-service/internal/grpc"
	"webby/room-category-service/internal/grpc/categorypb"
	httpserver "webby/room-category-service/internal/handlers"
	"webby/room-category-service/internal/repository"
	"webby/room-category-service/internal/services"
	"webby/room-category-service/pkg/slogpretty"

	"google.golang.org/grpc"
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

	categoryRepository := repository.New(db)
	categoryService := services.New(categoryRepository)

	server := httpserver.NewServer(
		cfg,
		logger,
		categoryService,
	)
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

	grpcListener, err := net.Listen("tcp", cfg.Grpc.HostPort)
	if err != nil {
		logger.Error("grpc listen failed", slog.Any("error", err))
		return err
	}

	grpcSrv := grpc.NewServer()
	categorypb.RegisterCategoryGrpcServiceServer(grpcSrv, grpcserver.NewCategoryServer(categoryService))

	go func() {
		logger.Info("gRPC server listening", "host", cfg.Grpc.HostPort)
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

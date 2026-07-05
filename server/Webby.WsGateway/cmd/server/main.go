package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"

	"webby/wsgateway/internal/config"
	grpcClient "webby/wsgateway/internal/grpc"
	"webby/wsgateway/internal/handlers"
	redisbus "webby/wsgateway/internal/redis"
	"webby/wsgateway/internal/repositories"
	"webby/wsgateway/internal/services"
	"webby/wsgateway/internal/sse"
	"webby/wsgateway/internal/ws"
	"webby/wsgateway/pkg/slogpretty"
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer func() {
		logger.Info("closing redis client")
		rdb.Close()
	}()

	tokenRepository := repositories.NewTokenRepository(rdb, cfg.TokenTTL)
	tokenService := services.NewTokenService(tokenRepository)

	presenceRepository := repositories.NewPresenceRepository(rdb)
	presenceService := services.NewPresenceService(presenceRepository)

	chatClient, err := grpcClient.NewChatClient(cfg.Grpc.ChatServiceAddress)
	if err != nil {
		logger.Error("chat service gRPC connection failed", slog.String("error", err.Error()))
		return err
	}
	defer chatClient.Close()

	sseBroker := sse.New()
	wsSrv := ws.NewServer(tokenService, presenceService, chatClient, logger, cfg.Http.CallTimeout)
	apiRouter := handlers.NewServer(cfg, tokenService, logger, wsSrv, sseBroker)

	httpSrv := &http.Server{
		Addr:    cfg.Http.HostPort,
		Handler: apiRouter,
	}

	var wg sync.WaitGroup

	wg.Go(func() {
		logger.Info("socket.io engine starting")
		if err := wsSrv.Serve(); err != nil {
			logger.Error("socket.io serve loop stopped", slog.String("err", err.Error()))
		}
	})

	wg.Go(func() {
		logger.Info("heartbeat worker started")
		wsSrv.RunHeartbeat(ctx)
		logger.Info("heartbeat worker stopped")
	})

	sub := redisbus.NewSubscriber(logger, rdb, wsSrv, sseBroker, cfg.Redis.Pattern)
	wg.Go(func() {
		if err := sub.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			logger.Error("redis subscriber error", slog.String("err", err.Error()))
		}
		logger.Info("redis subscriber stopped")
	})

	wg.Go(func() {
		logger.Info("http server listening", slog.String("addr", httpSrv.Addr))
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http listen and serve error", slog.String("err", err.Error()))
		}
	})

	<-ctx.Done()
	logger.Info("shutdown initiated")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("shutting down http server...")
	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		logger.Error("http server shutdown failed", slog.String("err", err.Error()))
		return err
	}

	logger.Info("closing socket.io server...")
	if err := wsSrv.Close(); err != nil {
		logger.Error("socket.io server close failed", slog.String("err", err.Error()))
		return err
	}

	wg.Wait()
	logger.Info("all systems stopped clean. bye")
	return nil
}

func setupLogger(env string, w io.Writer) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog(w)
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

func setupPrettySlog(w io.Writer) *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(w)

	return slog.New(handler)
}

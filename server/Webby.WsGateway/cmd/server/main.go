package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
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
)

func main() {
	cfg := config.MustLoad()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

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
		logger.Warn("chat service gRPC connection failed — chat features disabled", slog.String("error", err.Error()))
		chatClient = nil
	} else {
		defer chatClient.Close()
	}

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
	}

	logger.Info("closing socket.io server...")
	if err := wsSrv.Close(); err != nil {
		logger.Error("socket.io server close failed", slog.String("err", err.Error()))
	}

	wg.Wait()
	logger.Info("all systems stopped clean. bye")
}

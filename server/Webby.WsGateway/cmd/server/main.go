package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"

	"webby/wsgateway/internal/config"
	"webby/wsgateway/internal/handlers"
	redisbus "webby/wsgateway/internal/redis"
	"webby/wsgateway/internal/repositories"
	"webby/wsgateway/internal/services"
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
	defer rdb.Close()

	repository := repositories.New(rdb, cfg.TokenTTL)
	service := services.New(repository)

	wsSrv := ws.NewServer(service, logger, cfg.Http.CallTimeout)
	go func() {
		if err := wsSrv.Serve(); err != nil {
			logger.Error("socket.io serve", slog.String("err", err.Error()))
		}
	}()
	defer wsSrv.Close()

	logger.Info("database connected successfully")

	server := handlers.NewServer(cfg, service, logger, wsSrv)

	httpSrv := &http.Server{
		Addr:         net.JoinHostPort(cfg.Http.Host, strconv.Itoa(cfg.Http.Port)),
		Handler:      server,
		ReadTimeout:  cfg.Http.Timeout,
		WriteTimeout: cfg.Http.Timeout,
	}

	var wg sync.WaitGroup

	wg.Go(func() {
		logger.Info("http listening", slog.String("addr", httpSrv.Addr))
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http serve", slog.String("err", err.Error()))
		}
	})

	sub := redisbus.NewSubscriber(logger, rdb, wsSrv, cfg.Redis.Pattern)
	wg.Go(func() {
		if err := sub.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			logger.Error("redis subscriber", slog.String("err", err.Error()))
		}
	})

	<-ctx.Done()
	logger.Info("shutdown initiated")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	httpSrv.Shutdown(shutdownCtx)

	wg.Wait()
	logger.Info("bye")
}

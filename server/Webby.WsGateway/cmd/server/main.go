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

	"webby-wsgateway/internal/config"
	clients "webby-wsgateway/internal/grpc"
	redisbus "webby-wsgateway/internal/redis"
	"webby-wsgateway/internal/worker"
	"webby-wsgateway/internal/ws"
)

func main() {
	cfg := config.MustLoad()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	chatClient, err := clients.NewChatClient(cfg.Grpc.Chat)
	if err != nil {
		logger.Error("chat client", slog.String("err", err.Error()))
		os.Exit(1)
	}
	defer chatClient.Close()

	votesClient, err := clients.NewVotesClient(cfg.Grpc.Votes)
	if err != nil {
		logger.Error("votes client", slog.String("err", err.Error()))
		os.Exit(1)
	}
	defer votesClient.Close()
	_ = votesClient

	roomClient, err := clients.NewRoomClient(cfg.Grpc.Room)
	if err != nil {
		logger.Error("room client", slog.String("err", err.Error()))
		os.Exit(1)
	}
	defer roomClient.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	defer rdb.Close()

	wsSrv := ws.NewServer(logger, []byte(cfg.JwtSecret), chatClient)
	go func() {
		if err := wsSrv.Serve(); err != nil {
			logger.Error("socket.io serve", slog.String("err", err.Error()))
		}
	}()
	defer wsSrv.Close()

	mux := http.NewServeMux()
	mux.Handle("/socket.io/", wsSrv.IO())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("GET /swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "./docs/oas.json")
	})

	httpSrv := &http.Server{
		Addr:         net.JoinHostPort(cfg.Http.Host, strconv.Itoa(cfg.Http.Port)),
		Handler:      mux,
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

	act := worker.NewActivity(logger, wsSrv, roomClient, cfg.Worker.ActivityInterval)
	wg.Go(func() {
		if err := act.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			logger.Error("activity worker", slog.String("err", err.Error()))
		}
	})

	<-ctx.Done()
	logger.Info("shutdown initiated")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(shutdownCtx)

	wg.Wait()
	logger.Info("bye")
}

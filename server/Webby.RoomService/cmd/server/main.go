package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"
	"webby/internal/config"
	"webby/internal/database"
	grpcClient "webby/internal/grpc"
	httpserver "webby/internal/handlers"
	"webby/internal/repository"
	"webby/internal/services"
	"webby/pkg/slogpretty"
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
	categoryRepository := repository.NewCategoryRepository(db)
	queueItemRepository := repository.NewQueueItemRepository(db)
	fileStorage := repository.NewFileStorage(config)

	mediaClient, err := grpcClient.NewMediaClient(config.Grpc.MediaServiceAddress)
	if err != nil {
		logger.Error("media service gRPC connection failed", slog.String("error", err.Error()))
		return err
	}
	defer mediaClient.Close()

	roomService := services.NewRoomService(roomRepository, roomMemberRepository, fileStorage)
	categoryService := services.NewCategoryService(categoryRepository)
	queueItemService := services.NewQueueItemService(queueItemRepository, mediaClient, roomMemberRepository)
	voteRepository := repository.NewVoteRepository(db)
	voteService := services.NewVoteService(voteRepository, roomRepository, roomMemberRepository, queueItemRepository)

	logger.Info("repositories initialized")

	server := httpserver.NewServer(
		config,
		logger,
		roomService,
		categoryService,
		queueItemService,
		voteService,
	)
	httpServer := &http.Server{
		Addr:         net.JoinHostPort(config.Http.Host, strconv.Itoa(config.Http.Port)),
		ReadTimeout:  config.Http.Timeout,
		WriteTimeout: config.Http.Timeout,
		Handler:      server,
	}

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

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

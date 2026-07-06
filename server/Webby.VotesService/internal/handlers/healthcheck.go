package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type HealthChecker interface {
	CheckHealth(ctx context.Context) error
}

type healthcheck struct {
	checkers []HealthChecker
}

func NewHealthCheck(checkers ...HealthChecker) *healthcheck {
	return &healthcheck{
		checkers: checkers,
	}
}

func (h *healthcheck) Check(ctx context.Context) error {
	for _, checker := range h.checkers {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		if err := checker.CheckHealth(ctx); err != nil {
			cancel()
			return err
		}
		cancel()
	}
	return nil
}

type DatabaseChecker struct {
	pool *pgxpool.Pool
}

func NewDatabaseChecker(pool *pgxpool.Pool) *DatabaseChecker {
	return &DatabaseChecker{pool: pool}
}

func (dc *DatabaseChecker) CheckHealth(ctx context.Context) error {
	if err := dc.pool.Ping(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}
	return nil
}

type RedisChecker struct {
	client *redis.Client
}

func NewRedisChecker(client *redis.Client) *RedisChecker {
	return &RedisChecker{client: client}
}

func (rc *RedisChecker) CheckHealth(ctx context.Context) error {
	if err := rc.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis health check failed: %w", err)
	}
	return nil
}

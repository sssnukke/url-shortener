package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

func NewPostgresPool(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := os.Getenv("POSTGRES_DSN")
	log.Println("connecting with DSN:", dsn)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool postgres: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return pool, nil
}

func NewRedisPool(ctx context.Context) (*redis.Client, error) {
	addr := os.Getenv("REDIS_ADDR")

	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return client, nil
}

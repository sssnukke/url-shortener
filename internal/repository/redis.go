package repository

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type CacheRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

type redisRepo struct {
	client *redis.Client
}

func NewRedisRepo(client *redis.Client) CacheRepository {
	return &redisRepo{client: client}
}

func (r *redisRepo) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()

	if err == redis.Nil {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *redisRepo) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *redisRepo) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

package storage

import (
	"context"
	"github.com/mhsnrafi/giftassignmentservice/internal/errors"
	"time"

	"github.com/go-redis/redis/v8"
)

type DistributedLocker interface {
	Lock(ctx context.Context, key string, ttl time.Duration) error
	Unlock(ctx context.Context, key string) error
}

type RedisDistributedLocker struct {
	redisClient *redis.Client
}

func NewDistributedLocker(redisClient *redis.Client) DistributedLocker {
	return &RedisDistributedLocker{
		redisClient: redisClient,
	}
}

func (l *RedisDistributedLocker) Lock(ctx context.Context, key string, ttl time.Duration) error {
	ok, err := l.redisClient.SetNX(ctx, key, 1, ttl).Result()
	if err != nil {
		return err
	}
	if !ok {
		return errors.ErrCouldNotObtainLock
	}

	return nil
}

func (l *RedisDistributedLocker) Unlock(ctx context.Context, key string) error {
	return l.redisClient.Del(ctx, key).Err()
}

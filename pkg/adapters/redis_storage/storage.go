package redis_storage

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/paragor/parashort/pkg/domain/storage"
)

const (
	reservedValue  = "-"
	keyPrefix      = "keys_"
	reservedPrefix = "reserve_"
)

type RedisStorage struct {
	client *redis.Client
}

func (s *RedisStorage) Ping(ctx context.Context) error {
	res := s.client.Ping(ctx)
	if nil != res.Err() {
		return fmt.Errorf("redis problems: %w", res.Err())
	}
	return nil
}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{client: client}
}

func (s *RedisStorage) Save(ctx context.Context, key, value string, ttl time.Duration) error {
	status := s.client.Set(ctx, keyPrefix+key, value, ttl)
	return status.Err()
}

func (s *RedisStorage) ListKeys(ctx context.Context) ([]string, error) {
	res, err := s.client.Keys(ctx, keyPrefix+"*").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return []string{}, nil
		}
		return nil, err
	}
	for i, _ := range res {
		res[i] = strings.TrimPrefix(res[i], keyPrefix)
	}

	return res, nil
}

func (s *RedisStorage) Delete(ctx context.Context, key string) error {
	res := s.client.Del(ctx, keyPrefix+key)
	if res.Err() != nil && !errors.Is(res.Err(), redis.Nil) {
		return res.Err()
	}

	return nil
}

func (s *RedisStorage) Get(ctx context.Context, key string) (string, error) {
	status := s.client.Get(ctx, keyPrefix+key)
	res, err := status.Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", storage.ErrNotFound
		}
		return "", err
	}

	if res == "" || res == reservedValue {
		return "", storage.ErrNotFound
	}
	return res, nil
}

func (s *RedisStorage) ReserveKey(ctx context.Context, key string, ttl time.Duration) error {
	// а то мало ли...
	defer s.client.Expire(ctx, reservedPrefix+key, ttl)

	result := s.client.GetSet(ctx, reservedPrefix+key, reservedValue)
	if result.Err() != nil && !errors.Is(result.Err(), redis.Nil) {
		return result.Err()
	}

	if result.Val() == reservedValue {
		return storage.ErrKeyAlreadyReversed
	}

	res := s.client.Exists(ctx, keyPrefix+key)
	if res.Err() != nil {
		if errors.Is(res.Err(), redis.Nil) {
			return storage.ErrKeyAlreadyReversed
		}
		return res.Err()
	}
	if res.Val() > 0 {
		return storage.ErrKeyAlreadyReversed
	}

	return nil
}

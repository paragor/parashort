package storage

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrNotFound           = errors.New("Item not found")
	ErrKeyAlreadyReversed = errors.New("Key already reserved")
)

type StorageEngine interface {
	Save(ctx context.Context, key, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	ListKeys(ctx context.Context) ([]string, error)

	ReserveKey(ctx context.Context, key string, ttl time.Duration) error
}

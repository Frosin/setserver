package storage

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("not found")

type Storage interface {
	Add(ctx context.Context, key, value string) error
	Delete(ctx context.Context, key string) (string, error)
}

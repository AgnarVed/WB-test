package service

import (
	"context"
	"encoding/json"
)

type CacheInterface interface {
	Get(ctx context.Context, key string) (json.RawMessage, bool, error)
}

package pluggy

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	AUTH_CACHE_API_KEY           = "pluggy:api_key"
	AUTH_CACHE_CONNECT_TOKEN_KEY = "pluggy:connect_tokens"

	PLUGGY_CLIENT_ID     = os.Getenv("PLUGGY_CLIENT_ID")
	PLUGGY_CLIENT_SECRET = os.Getenv("PLUGGY_CLIENT_SECRET")
)

type auth struct {
	cache *redis.Client
}

func NewAuth(cache *redis.Client) *auth {
	return &auth{cache}
}

func (a *auth) getApiKey() (string, error) {
	res, err := a.cache.Get(context.TODO(), AUTH_CACHE_API_KEY).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("error getting cached Pluggy.ai api key: %w", err)
	}
	return res, nil
}

func (a *auth) setApiKey(k string) error {
	return a.cache.Set(context.TODO(), AUTH_CACHE_API_KEY, k, time.Hour*2).Err()
}

func (a *auth) getConnectToken(itemID string) (string, error) {
	res, err := a.cache.HGet(context.TODO(), AUTH_CACHE_CONNECT_TOKEN_KEY, itemID).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("error getting cached Pluggy.ai connect token: %w", err)
	}
	return res, nil
}

func (a *auth) setConnectToken(itemID, token string) error {
	return a.cache.HSet(context.TODO(), AUTH_CACHE_CONNECT_TOKEN_KEY, itemID, token).Err()
}

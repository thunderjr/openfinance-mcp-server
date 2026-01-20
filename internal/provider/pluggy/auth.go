package pluggy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
)

var (
	AUTH_CACHE_API_KEY           = "pluggy:api_key"
	AUTH_CACHE_CONNECT_TOKEN_KEY = "pluggy:connect_tokens"
)

type pluggyCredentials struct {
	ClientID     string `json:"PLUGGY_CLIENT_ID"`
	ClientSecret string `json:"PLUGGY_CLIENT_SECRET"`
}

type auth struct {
	cache       *redis.Client
	credentials *pluggyCredentials
	credOnce    sync.Once
}

func NewAuth(cache *redis.Client) *auth {
	return &auth{cache: cache}
}

func (a *auth) loadCredentials() {
	a.credOnce.Do(func() {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			logger.Fatalf("error getting home directory: %v", err)
		}

		configPath := filepath.Join(homeDir, ".config", "finops", "pluggy.json")
		data, err := os.ReadFile(configPath)
		if err != nil {
			logger.Fatalf("error reading pluggy credentials from %s: %v", configPath, err)
		}

		var creds pluggyCredentials
		if err := json.Unmarshal(data, &creds); err != nil {
			logger.Fatalf("error parsing pluggy credentials: %v", err)
		}

		if creds.ClientID == "" || creds.ClientSecret == "" {
			logger.Fatalf("pluggy credentials file missing PLUGGY_CLIENT_ID or PLUGGY_CLIENT_SECRET")
		}

		a.credentials = &creds
	})
}

func (a *auth) getClientID() string {
	a.loadCredentials()
	return a.credentials.ClientID
}

func (a *auth) getClientSecret() string {
	a.loadCredentials()
	return a.credentials.ClientSecret
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

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

const (
	authCacheAPIKey            = "pluggy:api_key"
	authCacheConnectTokenPrefix = "pluggy:connect_token:"

	apiKeyTTL       = 2 * time.Hour
	connectTokenTTL = 30 * time.Minute
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
	ctx := context.Background()
	res, err := a.cache.Get(ctx, authCacheAPIKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("pluggy.auth: error getting cached api key: %w", err)
	}
	return res, nil
}

func (a *auth) setApiKey(k string) error {
	ctx := context.Background()
	return a.cache.Set(ctx, authCacheAPIKey, k, apiKeyTTL).Err()
}

func (a *auth) getConnectToken(itemID string) (string, error) {
	ctx := context.Background()
	key := authCacheConnectTokenPrefix + itemID
	res, err := a.cache.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("pluggy.auth: error getting cached connect token: %w", err)
	}
	return res, nil
}

func (a *auth) setConnectToken(itemID, token string) error {
	ctx := context.Background()
	key := authCacheConnectTokenPrefix + itemID
	return a.cache.Set(ctx, key, token, connectTokenTTL).Err()
}

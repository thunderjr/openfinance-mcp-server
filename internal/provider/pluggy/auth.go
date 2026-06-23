package pluggy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
)

const (
	credClientID     = "pluggy:client_id"
	credClientSecret = "pluggy:client_secret"
	cacheAPIKey      = "pluggy:api_key"
	connectTokenKey  = "pluggy:connect_token:"

	apiKeyTTL       = 2 * time.Hour
	connectTokenTTL = 30 * time.Minute
)

// Store is the encrypted persistence backing auth (api keys, connect tokens,
// and the Pluggy client credentials). It is satisfied by sqlite.Store.
type Store interface {
	Get(key string) (string, bool, error)
	Set(key, value string, ttl time.Duration) error
}

type pluggyCredentials struct {
	ClientID     string `json:"PLUGGY_CLIENT_ID"`
	ClientSecret string `json:"PLUGGY_CLIENT_SECRET"`
}

type auth struct {
	store       Store
	credentials *pluggyCredentials
	credOnce    sync.Once
}

func NewAuth(store Store) *auth {
	return &auth{store: store}
}

// loadCredentials resolves the Pluggy client credentials. The encrypted store
// is authoritative; on first run the credentials are seeded (once) from the
// environment or the legacy ~/.config/finops/pluggy.json file and then
// persisted, so later runs read only from the encrypted database.
func (a *auth) loadCredentials() {
	a.credOnce.Do(func() {
		if creds, ok := a.credentialsFromStore(); ok {
			a.credentials = creds
			return
		}

		creds, source := a.bootstrapCredentials()
		if creds == nil {
			logger.Fatalf("pluggy.auth: no credentials found - set PLUGGY_CLIENT_ID/PLUGGY_CLIENT_SECRET or create ~/.config/finops/pluggy.json")
		}

		if err := a.persistCredentials(creds); err != nil {
			logger.Fatalf("pluggy.auth: error saving credentials to encrypted store: %v", err)
		}
		logger.Infof("pluggy.auth: imported credentials from %s into encrypted store", source)
		a.credentials = creds
	})
}

func (a *auth) credentialsFromStore() (*pluggyCredentials, bool) {
	id, idOK, err := a.store.Get(credClientID)
	if err != nil {
		logger.Fatalf("pluggy.auth: error reading client id from store: %v", err)
	}
	secret, secretOK, err := a.store.Get(credClientSecret)
	if err != nil {
		logger.Fatalf("pluggy.auth: error reading client secret from store: %v", err)
	}
	if idOK && secretOK && id != "" && secret != "" {
		return &pluggyCredentials{ClientID: id, ClientSecret: secret}, true
	}
	return nil, false
}

// bootstrapCredentials reads credentials from the environment first, then the
// legacy plaintext config file. It does not delete the legacy file.
func (a *auth) bootstrapCredentials() (*pluggyCredentials, string) {
	if id, secret := os.Getenv("PLUGGY_CLIENT_ID"), os.Getenv("PLUGGY_CLIENT_SECRET"); id != "" && secret != "" {
		return &pluggyCredentials{ClientID: id, ClientSecret: secret}, "environment"
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Fatalf("pluggy.auth: error getting home directory: %v", err)
	}
	configPath := filepath.Join(homeDir, ".config", "finops", "pluggy.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ""
		}
		logger.Fatalf("pluggy.auth: error reading %s: %v", configPath, err)
	}

	var creds pluggyCredentials
	if err := json.Unmarshal(data, &creds); err != nil {
		logger.Fatalf("pluggy.auth: error parsing %s: %v", configPath, err)
	}
	if creds.ClientID == "" || creds.ClientSecret == "" {
		return nil, ""
	}
	return &creds, configPath
}

func (a *auth) persistCredentials(creds *pluggyCredentials) error {
	if err := a.store.Set(credClientID, creds.ClientID, 0); err != nil {
		return err
	}
	return a.store.Set(credClientSecret, creds.ClientSecret, 0)
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
	key, _, err := a.store.Get(cacheAPIKey)
	if err != nil {
		return "", fmt.Errorf("pluggy.auth: error getting cached api key: %w", err)
	}
	return key, nil
}

func (a *auth) setApiKey(k string) error {
	return a.store.Set(cacheAPIKey, k, apiKeyTTL)
}

func (a *auth) getConnectToken(itemID string) (string, error) {
	token, _, err := a.store.Get(connectTokenKey + itemID)
	if err != nil {
		return "", fmt.Errorf("pluggy.auth: error getting cached connect token: %w", err)
	}
	return token, nil
}

func (a *auth) setConnectToken(itemID, token string) error {
	return a.store.Set(connectTokenKey+itemID, token, connectTokenTTL)
}

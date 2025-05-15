package pluggy

import (
	"net/http"
	"time"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
)

type Client struct {
	http.Client
	apiKey      string
	auth        *auth
	rateLimiter *rateLimiter
}

func NewClient(auth *auth) *Client {
	if PLUGGY_CLIENT_ID == "" || PLUGGY_CLIENT_SECRET == "" {
		logger.Fatal("missing Pluggy.ai credentials")
	}

	apiKey, err := auth.getApiKey()
	if err != nil {
		logger.Fatal(err)
	}

	client := &Client{
		auth:   auth,
		apiKey: apiKey,
		rateLimiter: &rateLimiter{
			timestamp: time.Now(),
		},
	}

	if apiKey == "" {
		apiKey, err = client.ApiKey()
		if err != nil {
			logger.Fatalf("[pluggy] error authorizing: %v", err)
		}
		if err = auth.setApiKey(apiKey); err != nil {
			logger.Fatalf("[redis][pluggy] error saving api key: %v", err)
		}
	}

	client.apiKey = apiKey
	return client
}

package pluggy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
)

func (c *Client) ConnectToken(itemID string) (string, error) {
	data, err := json.Marshal(map[string]any{
		"itemId": itemID,
		"options": map[string]bool{
			"avoidDuplicates": true,
		},
	})
	if err != nil {
		return "", fmt.Errorf("[pluggy.ConnectToken] error marshalling data: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.pluggy.ai/connect_token", bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("[pluggy.ConnectToken] error creating request: %w", err)
	}
	req.Header.Set("X-API-KEY", c.apiKey)

	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("[pluggy.ConnectToken] error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("[pluggy.ConnectToken] unexpected status code: %d", resp.StatusCode)
	}

	var body struct {
		AccessToken string `json:"accessToken"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("[pluggy.ConnectToken] error decoding response: %w", err)
	}

	if err := c.auth.setConnectToken(itemID, body.AccessToken); err != nil {
		logger.Errorf("[pluggy.ConnectToken] error saving connect token to cache: \nToken: %s\n%v", body.AccessToken, err)
	}

	return body.AccessToken, nil
}

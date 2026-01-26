package pluggy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
)

func (c *Client) GetAccounts(itemID string) (*paginatedResponse[account], error) {
	req, err := http.NewRequest("GET", "https://api.pluggy.ai/accounts?itemId="+itemID, nil)
	if err != nil {
		return nil, fmt.Errorf("pluggy.GetAccounts: error creating request: %w", err)
	}
	req.Header.Set("X-API-KEY", c.apiKey)

	res, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("pluggy.GetAccounts: error making request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("pluggy.GetAccounts: error reading response body: %w", err)
		}
		logger.Debug("body: \n", string(bodyBytes))

		return nil, fmt.Errorf("pluggy.GetAccounts: request failed with status %d", res.StatusCode)
	}

	var data paginatedResponse[account]
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("pluggy.GetAccounts: error decoding response: %w", err)
	}

	return &data, nil
}

func (c *Client) GetAccount(accountID string) (*account, error) {
	if accountID == "" {
		return nil, fmt.Errorf("pluggy.GetAccount: accountID is required")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.pluggy.ai/accounts/%s", accountID), nil)
	if err != nil {
		return nil, fmt.Errorf("pluggy.GetAccount: error creating request: %w", err)
	}
	req.Header.Set("X-API-KEY", c.apiKey)

	res, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("pluggy.GetAccount: error making request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("pluggy.GetAccount: error reading response body: %w", err)
		}
		logger.Debug("body: \n", string(bodyBytes))

		return nil, fmt.Errorf("pluggy.GetAccount: error getting account: request failed with status %d", res.StatusCode)
	}

	var data account
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("pluggy.GetAccount: error decoding response: %w", err)
	}

	return &data, nil
}

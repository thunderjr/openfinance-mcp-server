package pluggy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) ApiKey() (string, error) {
	data, err := json.Marshal(map[string]string{
		"clientId":     PLUGGY_CLIENT_ID,
		"clientSecret": PLUGGY_CLIENT_SECRET,
	})
	if err != nil {
		return "", fmt.Errorf("[pluggy.ApiKey] error marshalling data: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.pluggy.ai/auth", bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("[pluggy.ApiKey] error creating request: %w", err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("[pluggy.ApiKey] error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("[pluggy.ApiKey] unexpected status code: %d", resp.StatusCode)
	}

	var body struct {
		ApiKey string `json:"apiKey"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", fmt.Errorf("[pluggy.ApiKey] error decoding response: %w", err)
	}

	return body.ApiKey, nil
}

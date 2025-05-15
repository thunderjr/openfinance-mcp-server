package pluggy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ItemStatus string

const (
	ItemStatusUpdating         ItemStatus = "UPDATING"
	ItemStatusError            ItemStatus = "LOGIN_ERROR"
	ItemStatusOutdated         ItemStatus = "OUTDATED"
	ItemStatusWaitingUserInput ItemStatus = "WAITING_USER_INPUT"
	ItemStatusUpdated          ItemStatus = "UPDATED"
)

func (c *Client) WaitUpdated(itemID string) error {
	status := ItemStatusUpdating
	for status == ItemStatusUpdating {
		res, err := c.GetItem(itemID)
		if err != nil {
			return fmt.Errorf("pluggy.GetItem: error waiting updated: %w", err)
		}

		status = ItemStatus(res.Status)
		time.Sleep(5 * time.Second)
	}

	if status != ItemStatusUpdated {
		return fmt.Errorf("pluggy.WaitUpdated: item status is %s", status)
	}

	return nil
}

func (c *Client) GetItem(id string) (*itemResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.pluggy.ai/items/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("pluggy_client: error creating request: %w", err)
	}
	req.Header.Set("X-API-KEY", c.apiKey)

	res, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("pluggy_client: error item making request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pluggy_client: error item request failed with status %d", res.StatusCode)
	}

	var data itemResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("pluggy_client: error item decoding response: %w", err)
	}

	return &data, nil
}

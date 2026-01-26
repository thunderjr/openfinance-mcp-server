package pluggy

import (
	"context"
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

const (
	waitUpdatedInitialBackoff = 2 * time.Second
	waitUpdatedMaxBackoff     = 30 * time.Second
	waitUpdatedMaxRetries     = 30
)

// WaitUpdated polls the item status until it's no longer updating.
// Supports context cancellation and uses exponential backoff.
func (c *Client) WaitUpdated(ctx context.Context, itemID string) error {
	backoff := waitUpdatedInitialBackoff
	retries := 0

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("pluggy.WaitUpdated: context cancelled: %w", ctx.Err())
		default:
		}

		res, err := c.GetItem(itemID)
		if err != nil {
			return fmt.Errorf("pluggy.WaitUpdated: error getting item: %w", err)
		}

		status := ItemStatus(res.Status)
		if status != ItemStatusUpdating {
			if status != ItemStatusUpdated {
				return fmt.Errorf("pluggy.WaitUpdated: item status is %s", status)
			}
			return nil
		}

		retries++
		if retries >= waitUpdatedMaxRetries {
			return fmt.Errorf("pluggy.WaitUpdated: max retries (%d) exceeded", waitUpdatedMaxRetries)
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("pluggy.WaitUpdated: context cancelled: %w", ctx.Err())
		case <-time.After(backoff):
		}

		// Exponential backoff with max cap
		backoff = min(backoff*2, waitUpdatedMaxBackoff)
	}
}

func (c *Client) GetItem(id string) (*itemResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.pluggy.ai/items/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("pluggy.GetItem: error creating request: %w", err)
	}
	req.Header.Set("X-API-KEY", c.apiKey)

	res, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("pluggy.GetItem: error making request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pluggy.GetItem: request failed with status %d", res.StatusCode)
	}

	var data itemResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("pluggy.GetItem: error decoding response: %w", err)
	}

	return &data, nil
}

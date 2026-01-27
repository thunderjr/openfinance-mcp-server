package tools

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	mcp "github.com/metoro-io/mcp-golang"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
	internalMcp "github.com/thunderjr/openfinance-mcp-server/internal/provider/mcp"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/pluggy"
)

const refreshConnectionsTimeout = 10 * time.Minute

type RefreshConnectionsTool struct {
	client *pluggy.Client
}

func NewRefreshConnectionsTool(client *pluggy.Client) *RefreshConnectionsTool {
	return &RefreshConnectionsTool{client}
}

func (t *RefreshConnectionsTool) Name() string {
	return "refresh_connections"
}

func (t *RefreshConnectionsTool) Description() string {
	return `Triggers a sync update for all bank connections (or a specific one) and waits for completion.

Use this tool when:
- The user wants to refresh their bank data
- You need up-to-date account balances or transactions
- A single call to update and wait instead of manually calling update + wait_for_sync

Returns: List of refresh results with item_id, associated connection names, final status, and any errors.`
}

func (t *RefreshConnectionsTool) Handle() internalMcp.ToolHandlerFunc {
	return t.handleRefreshConnections
}

type RefreshConnectionsArgs struct {
	ConnectionID string `json:"connection_id" jsonschema:"description=Optional connection ID to refresh. If omitted all connections are refreshed."`
}

type refreshResult struct {
	ItemID      string   `json:"item_id"`
	Connections []string `json:"connections"`
	Status      string   `json:"status"`
	Error       string   `json:"error,omitempty"`
}

func (t *RefreshConnectionsTool) handleRefreshConnections(args RefreshConnectionsArgs) (*mcp.ToolResponse, error) {
	logger.Info("Refreshing bank connections")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return toolError("Error getting home directory", err)
	}

	configPath := filepath.Join(homeDir, ".config", "finops", "accounts.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return mcp.NewToolResponse(mcp.NewTextContent("[]")), nil
		}
		return toolError("Error reading accounts config", err)
	}

	var rawConnections []struct {
		Name        string `json:"name"`
		ItemID      string `json:"item_id"`
		AccountID   string `json:"account_id"`
		AccountType string `json:"account_type"`
	}

	if err := json.Unmarshal(data, &rawConnections); err != nil {
		return toolError("Error parsing accounts config", err)
	}

	// Deduplicate by item_id and collect connection names
	type itemInfo struct {
		itemID      string
		connections []string
	}

	seen := make(map[string]*itemInfo)
	var order []string

	for _, raw := range rawConnections {
		if args.ConnectionID != "" && raw.ItemID != args.ConnectionID {
			continue
		}

		if info, ok := seen[raw.ItemID]; ok {
			info.connections = append(info.connections, raw.Name)
		} else {
			seen[raw.ItemID] = &itemInfo{
				itemID:      raw.ItemID,
				connections: []string{raw.Name},
			}
			order = append(order, raw.ItemID)
		}
	}

	if len(order) == 0 {
		return mcp.NewToolResponse(mcp.NewTextContent("[]")), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), refreshConnectionsTimeout)
	defer cancel()

	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		results []refreshResult
	)

	for _, itemID := range order {
		wg.Add(1)
		go func(info *itemInfo) {
			defer wg.Done()

			result := refreshResult{
				ItemID:      info.itemID,
				Connections: info.connections,
			}

			logger.Info("Triggering sync for item:", info.itemID)

			_, err := t.client.UpdateItem(info.itemID)
			if err != nil {
				result.Status = "ERROR"
				result.Error = err.Error()
				mu.Lock()
				results = append(results, result)
				mu.Unlock()
				return
			}

			err = t.client.WaitUpdated(ctx, info.itemID)
			if err != nil {
				result.Status = "ERROR"
				result.Error = err.Error()
				mu.Lock()
				results = append(results, result)
				mu.Unlock()
				return
			}

			item, err := t.client.GetItem(info.itemID)
			if err != nil {
				result.Status = "ERROR"
				result.Error = err.Error()
				mu.Lock()
				results = append(results, result)
				mu.Unlock()
				return
			}

			result.Status = item.Status
			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(seen[itemID])
	}

	wg.Wait()

	return jsonResponse(results)
}

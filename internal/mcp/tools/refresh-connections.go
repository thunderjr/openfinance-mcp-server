package tools

import (
	"context"
	"sync"
	"time"

	mcp "github.com/metoro-io/mcp-golang"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/connections"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
	internalMcp "github.com/thunderjr/openfinance-mcp-server/internal/provider/mcp"
)

const refreshConnectionsTimeout = 10 * time.Minute

type RefreshConnectionsTool struct {
	registry *connections.Registry
	gateway  connectionGateway
}

func NewRefreshConnectionsTool(registry *connections.Registry, gateway connectionGateway) *RefreshConnectionsTool {
	return &RefreshConnectionsTool{registry: registry, gateway: gateway}
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
	Connections []string `json:"connections,omitempty"`
	Status      string   `json:"status"`
	Error       string   `json:"error,omitempty"`
}

func (t *RefreshConnectionsTool) handleRefreshConnections(args RefreshConnectionsArgs) (*mcp.ToolResponse, error) {
	logger.Info("Refreshing bank connections")

	itemIDs, err := t.registry.List()
	if err != nil {
		return toolError("Error reading saved connections", err)
	}

	// Optionally narrow to a single connection.
	if args.ConnectionID != "" {
		filtered := itemIDs[:0:0]
		for _, id := range itemIDs {
			if id == args.ConnectionID {
				filtered = append(filtered, id)
			}
		}
		itemIDs = filtered
	}

	if len(itemIDs) == 0 {
		return mcp.NewToolResponse(mcp.NewTextContent("[]")), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), refreshConnectionsTimeout)
	defer cancel()

	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		results []refreshResult
	)

	for _, itemID := range itemIDs {
		wg.Add(1)
		go func(itemID string) {
			defer wg.Done()

			result := refreshResult{ItemID: itemID}

			logger.Info("Triggering sync for item:", itemID)
			info, err := t.gateway.RefreshItem(ctx, itemID)
			if err != nil {
				result.Status = "ERROR"
				result.Error = err.Error()
			} else {
				result.Status = info.Status
				if info.ConnectorName != "" {
					result.Connections = []string{info.ConnectorName}
				}
			}

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(itemID)
	}

	wg.Wait()

	return jsonResponse(results)
}

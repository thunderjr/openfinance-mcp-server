package tools

import (
	"context"
	"time"

	mcp "github.com/metoro-io/mcp-golang"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
	internalMcp "github.com/thunderjr/openfinance-mcp-server/internal/provider/mcp"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/pluggy"
)

const waitForSyncTimeout = 5 * time.Minute

type WaitForSyncTool struct {
	client *pluggy.Client
}

func NewWaitForSyncTool(client *pluggy.Client) *WaitForSyncTool {
	return &WaitForSyncTool{client}
}

func (t *WaitForSyncTool) Name() string {
	return "wait_for_sync"
}

func (t *WaitForSyncTool) Description() string {
	return `Waits for a bank connection to finish syncing and returns the final status.

Use this tool when:
- After creating or updating a bank connection
- The connection status shows it's still syncing
- You need to wait for fresh data before querying accounts/transactions

Returns: Connection status with sync completion details.`
}

func (t *WaitForSyncTool) Handle() internalMcp.ToolHandlerFunc {
	return t.handleWaitForSync
}

type WaitForSyncArgs struct {
	ConnectionID string `json:"connection_id" jsonschema:"required,description=The connection ID to wait for sync completion"`
}

func (t *WaitForSyncTool) handleWaitForSync(args WaitForSyncArgs) (*mcp.ToolResponse, error) {
	if args.ConnectionID == "" {
		return validationError("Connection ID is required")
	}

	logger.Info("Waiting for sync to complete for connection:", args.ConnectionID)

	ctx, cancel := context.WithTimeout(context.Background(), waitForSyncTimeout)
	defer cancel()

	err := t.client.WaitUpdated(ctx, args.ConnectionID)
	if err != nil {
		return toolError("Error waiting for sync", err)
	}

	item, err := t.client.GetItem(args.ConnectionID)
	if err != nil {
		return toolError("Error getting connection status", err)
	}

	return jsonResponse(item)
}

package tools

import (
	"encoding/json"
	"fmt"

	mcp "github.com/metoro-io/mcp-golang"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
	internalMcp "github.com/thunderjr/openfinance-mcp-server/internal/provider/mcp"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/pluggy"
)

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
		errorMessage := "Invalid connection_id parameter: must be a non-empty string"
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	logger.Info("Waiting for sync to complete for connection:", args.ConnectionID)

	err := t.client.WaitUpdated(args.ConnectionID)
	if err != nil {
		errorMessage := fmt.Sprintf("Error waiting for sync: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	item, err := t.client.GetItem(args.ConnectionID)
	if err != nil {
		errorMessage := fmt.Sprintf("Error getting connection status: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	strItem, err := json.Marshal(item)
	if err != nil {
		errorMessage := fmt.Sprintf("Error marshalling connection status: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	return mcp.NewToolResponse(mcp.NewTextContent(string(strItem))), nil
}

package tools

import (
	mcp "github.com/metoro-io/mcp-golang"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
	internalMcp "github.com/thunderjr/openfinance-mcp-server/internal/provider/mcp"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/pluggy"
)

type GetConnectionStatusArgs struct {
	ConnectionID string `json:"connection_id" jsonschema:"required,description=The connection ID to check status for"`
}

type GetConnectionStatusTool struct {
	client *pluggy.Client
}

func NewGetConnectionStatusTool(client *pluggy.Client) *GetConnectionStatusTool {
	return &GetConnectionStatusTool{client}
}

func (t *GetConnectionStatusTool) Name() string {
	return "get_connection_status"
}

func (t *GetConnectionStatusTool) Description() string {
	return `Checks the health and sync status of a bank connection.

Use this tool when:
- You need to verify if a connection is still active
- Checking if a connection needs re-authentication
- Debugging connection issues before querying data

Returns: Connection details including status, last sync time, and any errors.`
}

func (t *GetConnectionStatusTool) Handle() internalMcp.ToolHandlerFunc {
	return t.handleGetConnectionStatus
}

func (t *GetConnectionStatusTool) handleGetConnectionStatus(args GetConnectionStatusArgs) (*mcp.ToolResponse, error) {
	if args.ConnectionID == "" {
		return validationError("Connection ID is required")
	}

	logger.Info("Getting connection status for:", args.ConnectionID)

	item, err := t.client.GetItem(args.ConnectionID)
	if err != nil {
		return toolError("Error getting connection status", err)
	}

	return jsonResponse(item)
}

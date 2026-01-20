package tools

import (
	"encoding/json"
	"fmt"

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
		errorMessage := "Connection ID is required"
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	logger.Info("Getting connection status for:", args.ConnectionID)

	item, err := t.client.GetItem(args.ConnectionID)
	if err != nil {
		errorMessage := fmt.Sprintf("Error getting connection status: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	itemJSON, err := json.Marshal(item)
	if err != nil {
		errorMessage := fmt.Sprintf("Error marshalling connection status: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	return mcp.NewToolResponse(mcp.NewTextContent(string(itemJSON))), nil
}

package tools

import (
	"fmt"

	mcp "github.com/metoro-io/mcp-golang"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
	internalMcp "github.com/thunderjr/openfinance-mcp-server/internal/provider/mcp"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/pluggy"
)

type PluggyConnectTokenTool struct {
	client *pluggy.Client
}

func NewPluggyConnectTokenTool(client *pluggy.Client) *PluggyConnectTokenTool {
	return &PluggyConnectTokenTool{client}
}

func (t *PluggyConnectTokenTool) Name() string {
	return "pluggy_connect_token"
}

func (t *PluggyConnectTokenTool) Description() string {
	return "Generates a new connect token for a specific Pluggy item"
}

func (t *PluggyConnectTokenTool) Handle() internalMcp.ToolHandlerFunc {
	return t.handleConnectToken
}

type ConnectTokenArgs struct {
	ItemID string `json:"item_id" jsonschema:"required,description=The Pluggy item ID to generate a connect token for"`
}

func (t *PluggyConnectTokenTool) handleConnectToken(args ConnectTokenArgs) (*mcp.ToolResponse, error) {
	if args.ItemID == "" {
		errorMessage := "Invalid item_id parameter: must be a non-empty string"
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	logger.Info("Generating connect token for item:", args.ItemID)

	token, err := t.client.ConnectToken(args.ItemID)
	if err != nil {
		errorMessage := fmt.Sprintf("Error generating connect token: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	return mcp.NewToolResponse(mcp.NewTextContent(token)), nil
}

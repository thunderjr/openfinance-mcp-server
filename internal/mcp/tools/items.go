package tools

import (
	"encoding/json"
	"fmt"

	mcp "github.com/metoro-io/mcp-golang"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
	internalMcp "github.com/thunderjr/openfinance-mcp-server/internal/provider/mcp"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/pluggy"
)

type ItemArgs struct {
	ItemID string `json:"item_id" jsonschema:"required,description=The ID of the item to retrieve details for"`
}

type PluggyItemTool struct {
	client *pluggy.Client
}

func NewPluggyItemTool(client *pluggy.Client) *PluggyItemTool {
	return &PluggyItemTool{client}
}

func (t *PluggyItemTool) Name() string {
	return "get_item_details"
}

func (t *PluggyItemTool) Description() string {
	return "Retrieves detailed information for a specific item"
}

func (t *PluggyItemTool) Handle() internalMcp.ToolHandlerFunc {
	return t.handleGetItem
}

func (t *PluggyItemTool) handleGetItem(args ItemArgs) (*mcp.ToolResponse, error) {
	if args.ItemID == "" {
		errorMessage := "Item ID is required"
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	logger.Info("Getting item details for:", args.ItemID)

	item, err := t.client.GetItem(args.ItemID)
	if err != nil {
		errorMessage := fmt.Sprintf("Error getting item: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	itemJSON, err := json.Marshal(item)
	if err != nil {
		errorMessage := fmt.Sprintf("Error marshalling item: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	return mcp.NewToolResponse(mcp.NewTextContent(string(itemJSON))), nil
}

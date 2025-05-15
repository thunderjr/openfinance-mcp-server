package tools

import (
	"encoding/json"
	"fmt"

	mcp "github.com/metoro-io/mcp-golang"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
	internalMcp "github.com/thunderjr/openfinance-mcp-server/internal/provider/mcp"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/pluggy"
)

type PluggyWaitItemUpdatedTool struct {
	client *pluggy.Client
}

func NewPluggyWaitItemUpdatedTool(client *pluggy.Client) *PluggyWaitItemUpdatedTool {
	return &PluggyWaitItemUpdatedTool{client}
}

func (t *PluggyWaitItemUpdatedTool) Name() string {
	return "pluggy_wait_item_updated"
}

func (t *PluggyWaitItemUpdatedTool) Description() string {
	return "Waits for an item to complete updating and returns final item status"
}

func (t *PluggyWaitItemUpdatedTool) Handle() internalMcp.ToolHandlerFunc {
	return t.handleWaitItemUpdated
}

type WaitItemUpdatedArgs struct {
	ItemID string `json:"item_id" jsonschema:"required,description=The Pluggy item ID to wait for update completion"`
}

func (t *PluggyWaitItemUpdatedTool) handleWaitItemUpdated(args WaitItemUpdatedArgs) (*mcp.ToolResponse, error) {
	if args.ItemID == "" {
		errorMessage := "Invalid item_id parameter: must be a non-empty string"
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	logger.Info("Waiting for item to update:", args.ItemID)

	err := t.client.WaitUpdated(args.ItemID)
	if err != nil {
		errorMessage := fmt.Sprintf("Error waiting for item update: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	item, err := t.client.GetItem(args.ItemID)
	if err != nil {
		errorMessage := fmt.Sprintf("Error getting updated item: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	strItem, err := json.Marshal(item)
	if err != nil {
		errorMessage := fmt.Sprintf("Error marshalling updated item: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	return mcp.NewToolResponse(mcp.NewTextContent(string(strItem))), nil
}

package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/thunderjr/openfinance-mcp-server/internal/infra/logger"
	"github.com/thunderjr/openfinance-mcp-server/internal/infra/pluggy"
)

type PluggyWaitItemUpdatedTool struct {
	client *pluggy.Client
}

func NewPluggyWaitItemUpdatedTool(client *pluggy.Client) *PluggyWaitItemUpdatedTool {
	return &PluggyWaitItemUpdatedTool{client}
}

func (t *PluggyWaitItemUpdatedTool) Tool() mcp.Tool {
	return mcp.Tool{
		Name:        "pluggy_wait_item_updated",
		Description: "Waits for an item to complete updating and returns final item status",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"item_id": map[string]interface{}{
					"type":        "string",
					"description": "The Pluggy item ID to wait for update completion",
				},
			},
			Required: []string{"item_id"},
		},
		Annotations: mcp.ToolAnnotation{
			Title:          "Wait for Item Update",
			ReadOnlyHint:   true,
			OpenWorldHint:  true,
			IdempotentHint: true,
		},
	}
}

func (t *PluggyWaitItemUpdatedTool) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	itemID, ok := args["item_id"].(string)
	if !ok {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Type: "text",
					Text: "Invalid item_id parameter: must be a string",
				},
			},
			IsError: true,
		}, nil
	}

	logger.Info("Waiting for item to update:", itemID)

	err := t.client.WaitUpdated(itemID)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error waiting for item update: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	item, err := t.client.GetItem(itemID)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error getting updated item: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	strItem, err := json.Marshal(item)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error marshalling updated item: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Type: "text",
				Text: string(strItem),
			},
		},
	}, nil
}

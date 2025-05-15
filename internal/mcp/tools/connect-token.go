package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/thunderjr/openfinance-mcp-server/internal/infra/logger"
	"github.com/thunderjr/openfinance-mcp-server/internal/infra/pluggy"
)

type PluggyConnectTokenTool struct {
	client *pluggy.Client
}

func NewPluggyConnectTokenTool(client *pluggy.Client) *PluggyConnectTokenTool {
	return &PluggyConnectTokenTool{client}
}

func (t *PluggyConnectTokenTool) Tool() mcp.Tool {
	return mcp.Tool{
		Name:        "pluggy_connect_token",
		Description: "Generates a new connect token for a specific Pluggy item",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]interface{}{
				"item_id": map[string]interface{}{
					"type":        "string",
					"description": "The Pluggy item ID to generate a connect token for",
				},
			},
			Required: []string{"item_id"},
		},
		Annotations: mcp.ToolAnnotation{
			Title:          "Generate Connect Token",
			ReadOnlyHint:   false,
			OpenWorldHint:  true,
			IdempotentHint: true,
		},
	}
}

func (t *PluggyConnectTokenTool) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	itemID, ok := args["item_id"].(string)
	if !ok {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Type: "text",
					Text: "Invalid item_id parameter: must be a string",
				},
			},
		}, nil
	}

	logger.Info("Generating connect token for item:", itemID)

	token, err := t.client.ConnectToken(itemID)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error generating connect token: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Type: "text",
				Text: token,
			},
		},
	}, nil
}

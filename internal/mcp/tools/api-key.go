package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/thunderjr/openfinance-mcp-server/internal/infra/logger"
	internalMcp "github.com/thunderjr/openfinance-mcp-server/internal/infra/mcp"
	"github.com/thunderjr/openfinance-mcp-server/internal/infra/pluggy"
)

var _ internalMcp.ToolHandler = (*PluggyApiKeyTool)(nil)

type PluggyApiKeyTool struct {
	client *pluggy.Client
}

func NewPluggyApiKeyTool(client *pluggy.Client) *PluggyApiKeyTool {
	return &PluggyApiKeyTool{client}
}

func (t *PluggyApiKeyTool) Tool() mcp.Tool {
	return mcp.Tool{
		Name:        "pluggy_api_key",
		Description: "Generates a new Pluggy API key using the client ID and client secret",
		Annotations: mcp.ToolAnnotation{
			Title:          "Generate Pluggy API Key",
			OpenWorldHint:  true,
			IdempotentHint: true,
		},
	}
}

func (t *PluggyApiKeyTool) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	logger.Info("Generating new Pluggy API key")

	apiKey, err := t.client.ApiKey()
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Error generating Pluggy API key: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Type: "text",
				Text: apiKey,
			},
		},
	}, nil
}

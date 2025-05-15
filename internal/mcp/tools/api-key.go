package tools

import (
	"fmt"

	mcp "github.com/metoro-io/mcp-golang"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
	internalMcp "github.com/thunderjr/openfinance-mcp-server/internal/provider/mcp"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/pluggy"
)

type PluggyApiKeyTool struct {
	client *pluggy.Client
}

func NewPluggyApiKeyTool(client *pluggy.Client) *PluggyApiKeyTool {
	return &PluggyApiKeyTool{client}
}

func (t *PluggyApiKeyTool) Name() string {
	return "pluggy_api_key"
}

func (t *PluggyApiKeyTool) Description() string {
	return "Generates a new Pluggy API key using the client ID and client secret"
}

func (t *PluggyApiKeyTool) Handle() internalMcp.ToolHandlerFunc {
	return t.handleApiKey
}

type ApiKeyArgs struct{}

func (t *PluggyApiKeyTool) handleApiKey(args ApiKeyArgs) (*mcp.ToolResponse, error) {
	logger.Info("Generating new Pluggy API key")

	apiKey, err := t.client.ApiKey()
	if err != nil {
		errorMessage := fmt.Sprintf("Error generating Pluggy API key: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	return mcp.NewToolResponse(mcp.NewTextContent(apiKey)), nil
}

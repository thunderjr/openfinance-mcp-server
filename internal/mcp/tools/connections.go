package tools

import (
	"encoding/json"
	"os"
	"path/filepath"

	mcp "github.com/metoro-io/mcp-golang"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
	internalMcp "github.com/thunderjr/openfinance-mcp-server/internal/provider/mcp"
)

type Connection struct {
	Name         string `json:"name"`
	ConnectionID string `json:"connection_id"`
	AccountID    string `json:"account_id"`
}

type ListConnectionsTool struct{}

func NewListConnectionsTool() *ListConnectionsTool {
	return &ListConnectionsTool{}
}

func (t *ListConnectionsTool) Name() string {
	return "list_connections"
}

func (t *ListConnectionsTool) Description() string {
	return `Lists all saved bank connections from local configuration.

Use this tool when:
- Starting a conversation to see available accounts
- The user asks about their connected banks
- You need connection_id values for other tools

Returns: List of connections with name, connection_id, and account_id.`
}

func (t *ListConnectionsTool) Handle() internalMcp.ToolHandlerFunc {
	return t.handleListConnections
}

type ListConnectionsArgs struct{}

func (t *ListConnectionsTool) handleListConnections(args ListConnectionsArgs) (*mcp.ToolResponse, error) {
	logger.Info("Listing saved bank connections")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return toolError("Error getting home directory", err)
	}

	configPath := filepath.Join(homeDir, ".config", "finops", "accounts.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return mcp.NewToolResponse(mcp.NewTextContent("[]")), nil
		}
		return toolError("Error reading accounts config", err)
	}

	var rawConnections []struct {
		Name      string `json:"name"`
		ItemID    string `json:"item_id"`
		AccountID string `json:"account_id"`
	}

	if err := json.Unmarshal(data, &rawConnections); err != nil {
		return toolError("Error parsing accounts config", err)
	}

	connections := make([]Connection, len(rawConnections))
	for i, raw := range rawConnections {
		connections[i] = Connection{
			Name:         raw.Name,
			ConnectionID: raw.ItemID,
			AccountID:    raw.AccountID,
		}
	}

	return jsonResponse(connections)
}

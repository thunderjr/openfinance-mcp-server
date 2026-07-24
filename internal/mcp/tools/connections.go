package tools

import (
	"context"

	mcp "github.com/metoro-io/mcp-golang"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/connections"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
	internalMcp "github.com/thunderjr/openfinance-mcp-server/internal/provider/mcp"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/pluggy"
)

// connectionGateway is the narrow seam over the Pluggy item subsystem the
// connection tools depend on. *pluggy.Client satisfies it; tests fake it at the
// network boundary.
type connectionGateway interface {
	Connection(itemID string) (pluggy.ConnectionInfo, error)
	RefreshItem(ctx context.Context, itemID string) (pluggy.ConnectionInfo, error)
}

// registerItem records a resolved item_id, logging (not failing) on error so a
// registry hiccup never breaks the tool the user actually called.
func registerItem(registry *connections.Registry, itemID string) {
	if registry == nil || itemID == "" {
		return
	}
	if err := registry.Register(itemID); err != nil {
		logger.Info("connections: failed to register item:", itemID, err)
	}
}

type Connection struct {
	Name         string `json:"name"`
	ConnectionID string `json:"connection_id"`
	AccountID    string `json:"account_id"`
	AccountType  string `json:"account_type,omitempty"`
}

type ListConnectionsTool struct {
	registry *connections.Registry
	gateway  connectionGateway
}

func NewListConnectionsTool(registry *connections.Registry, gateway connectionGateway) *ListConnectionsTool {
	return &ListConnectionsTool{registry: registry, gateway: gateway}
}

func (t *ListConnectionsTool) Name() string {
	return "list_connections"
}

func (t *ListConnectionsTool) Description() string {
	return `Lists all saved bank connections.

Use this tool when:
- Starting a conversation to see available accounts
- The user asks about their connected banks
- You need connection_id values for other tools

Returns: List of connections with name, connection_id, account_id, and account_type (BANK or CREDIT).`
}

func (t *ListConnectionsTool) Handle() internalMcp.ToolHandlerFunc {
	return t.handleListConnections
}

type ListConnectionsArgs struct{}

func (t *ListConnectionsTool) handleListConnections(ListConnectionsArgs) (*mcp.ToolResponse, error) {
	logger.Info("Listing saved bank connections")

	itemIDs, err := t.registry.List()
	if err != nil {
		return toolError("Error reading saved connections", err)
	}

	rows := make([]Connection, 0)
	for _, itemID := range itemIDs {
		info, err := t.gateway.Connection(itemID)
		if err != nil {
			// A connection Pluggy can no longer resolve (deleted, expired) is
			// skipped, not fatal: the rest of the list still returns.
			logger.Info("Skipping unresolved connection:", itemID, err)
			continue
		}
		for _, acc := range info.Accounts {
			rows = append(rows, Connection{
				Name:         info.ConnectorName + " - " + acc.Name,
				ConnectionID: itemID,
				AccountID:    acc.ID,
				AccountType:  acc.Type,
			})
		}
	}

	return jsonResponse(rows)
}

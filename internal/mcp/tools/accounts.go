package tools

import (
	mcp "github.com/metoro-io/mcp-golang"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
	internalMcp "github.com/thunderjr/openfinance-mcp-server/internal/provider/mcp"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/pluggy"
)

type ListAccountsArgs struct {
	ConnectionID string `json:"connection_id" jsonschema:"required,description=The connection ID to list accounts for"`
}

type ListAccountsTool struct {
	client *pluggy.Client
}

func NewListAccountsTool(client *pluggy.Client) *ListAccountsTool {
	return &ListAccountsTool{client}
}

func (t *ListAccountsTool) Name() string {
	return "list_accounts"
}

func (t *ListAccountsTool) Description() string {
	return `Lists all accounts associated with a bank connection.

Use this tool when:
- You need to see all accounts for a connected bank
- Looking for account IDs to query transactions or balances
- The user asks about their accounts at a specific bank

Returns: List of accounts with IDs, types, names, and current balances.`
}

func (t *ListAccountsTool) Handle() internalMcp.ToolHandlerFunc {
	return t.handleListAccounts
}

type GetAccountArgs struct {
	AccountID string `json:"account_id" jsonschema:"required,description=The account ID to retrieve details for"`
}

type GetAccountTool struct {
	client *pluggy.Client
}

func NewGetAccountTool(client *pluggy.Client) *GetAccountTool {
	return &GetAccountTool{client}
}

func (t *GetAccountTool) Name() string {
	return "get_account"
}

func (t *GetAccountTool) Description() string {
	return `Gets detailed information for a specific account.

Use this tool when:
- You need current balance for a specific account
- Looking for account details like number, type, or owner
- The user asks about a specific account

Returns: Account details including balance, type, name, and metadata.`
}

func (t *GetAccountTool) Handle() internalMcp.ToolHandlerFunc {
	return t.handleGetAccount
}

func (t *ListAccountsTool) handleListAccounts(args ListAccountsArgs) (*mcp.ToolResponse, error) {
	if args.ConnectionID == "" {
		return validationError("Connection ID is required")
	}

	logger.Info("Listing accounts for connection:", args.ConnectionID)

	accounts, err := t.client.GetAccounts(args.ConnectionID)
	if err != nil {
		return toolError("Error listing accounts", err)
	}

	return jsonResponse(accounts)
}

func (t *GetAccountTool) handleGetAccount(args GetAccountArgs) (*mcp.ToolResponse, error) {
	if args.AccountID == "" {
		return validationError("Account ID is required")
	}

	logger.Info("Getting account details for:", args.AccountID)

	account, err := t.client.GetAccount(args.AccountID)
	if err != nil {
		return toolError("Error getting account", err)
	}

	return jsonResponse(account)
}

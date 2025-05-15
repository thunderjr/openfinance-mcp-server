package tools

import (
	"encoding/json"
	"fmt"

	mcp "github.com/metoro-io/mcp-golang"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
	internalMcp "github.com/thunderjr/openfinance-mcp-server/internal/provider/mcp"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/pluggy"
)

type AccountsArgs struct {
	ItemID string `json:"item_id" jsonschema:"required,description=The ID of the item to retrieve accounts for"`
}

type PluggyAccountsTool struct {
	client *pluggy.Client
}

func NewPluggyAccountsTool(client *pluggy.Client) *PluggyAccountsTool {
	return &PluggyAccountsTool{client}
}

func (t *PluggyAccountsTool) Name() string {
	return "get_item_accounts"
}

func (t *PluggyAccountsTool) Description() string {
	return "Retrieves all accounts associated with a specific item"
}

func (t *PluggyAccountsTool) Handle() internalMcp.ToolHandlerFunc {
	return t.handleGetAccounts
}

type AccountArgs struct {
	AccountID string `json:"account_id" jsonschema:"required,description=The ID of the account to retrieve details for"`
}

type PluggyAccountTool struct {
	client *pluggy.Client
}

func NewPluggyAccountTool(client *pluggy.Client) *PluggyAccountTool {
	return &PluggyAccountTool{client}
}

func (t *PluggyAccountTool) Name() string {
	return "get_account_details"
}

func (t *PluggyAccountTool) Description() string {
	return "Retrieves detailed information for a specific account"
}

func (t *PluggyAccountTool) Handle() internalMcp.ToolHandlerFunc {
	return t.handleGetAccount
}

func (t *PluggyAccountsTool) handleGetAccounts(args AccountsArgs) (*mcp.ToolResponse, error) {
	if args.ItemID == "" {
		errorMessage := "Item ID is required"
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	accounts, err := t.client.GetAccounts(args.ItemID)
	if err != nil {
		errorMessage := fmt.Sprintf("Error getting accounts: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	accountsJSON, err := json.Marshal(accounts)
	if err != nil {
		errorMessage := fmt.Sprintf("Error marshalling accounts: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	return mcp.NewToolResponse(mcp.NewTextContent(string(accountsJSON))), nil
}

func (t *PluggyAccountTool) handleGetAccount(args AccountArgs) (*mcp.ToolResponse, error) {
	if args.AccountID == "" {
		errorMessage := "Account ID is required"
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	logger.Info("Getting account details for:", args.AccountID)

	account, err := t.client.GetAccount(args.AccountID)
	if err != nil {
		errorMessage := fmt.Sprintf("Error getting account: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	accountJSON, err := json.Marshal(account)
	if err != nil {
		errorMessage := fmt.Sprintf("Error marshalling account: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	return mcp.NewToolResponse(mcp.NewTextContent(string(accountJSON))), nil
}

package tools

import (
	"encoding/json"
	"fmt"
	"time"

	mcp "github.com/metoro-io/mcp-golang"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
	internalMcp "github.com/thunderjr/openfinance-mcp-server/internal/provider/mcp"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/pluggy"
)

type TransactionsArgs struct {
	AccountID   string    `json:"account_id" jsonschema:"required,description=The ID of the account to retrieve transactions for"`
	From        *string   `json:"from,omitempty" jsonschema:"description=Filter transactions from this date (format: yyyy-mm-dd)"`
	To          *string   `json:"to,omitempty" jsonschema:"description=Filter transactions to this date (format: yyyy-mm-dd)"`
	Page        *int      `json:"page,omitempty" jsonschema:"description=Page number for pagination (default: 1)"`
	PageSize    *int      `json:"page_size,omitempty" jsonschema:"description=Number of results per page (default: 20, max: 500)"`
	CreatedFrom *string   `json:"created_from,omitempty" jsonschema:"description=Filter transactions created from this date (ISO 8601 format)"`
	IDs         *[]string `json:"ids,omitempty" jsonschema:"description=Filter transactions by specific IDs"`
}

type PluggyTransactionsTool struct {
	client *pluggy.Client
}

func NewPluggyTransactionsTool(client *pluggy.Client) *PluggyTransactionsTool {
	return &PluggyTransactionsTool{client}
}

func (t *PluggyTransactionsTool) Name() string {
	return "get_account_transactions"
}

func (t *PluggyTransactionsTool) Description() string {
	return "Retrieves transactions for a specific account with optional filters"
}

func (t *PluggyTransactionsTool) Handle() internalMcp.ToolHandlerFunc {
	return t.handleGetTransactions
}

func (t *PluggyTransactionsTool) handleGetTransactions(args TransactionsArgs) (*mcp.ToolResponse, error) {
	if args.AccountID == "" {
		errorMessage := "Account ID is required"
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	filter := &pluggy.TransactionFilter{}

	if args.From != nil && *args.From != "" {
		from, err := time.Parse("2006-01-02", *args.From)
		if err == nil {
			filter.From = from
			logger.Info("Filter transactions from:", filter.From)
		} else {
			errorMessage := fmt.Sprintf("Invalid 'from' date format: %v", err)
			return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
		}
	}

	if args.To != nil && *args.To != "" {
		to, err := time.Parse("2006-01-02", *args.To)
		if err == nil {
			filter.To = to
			logger.Info("Filter transactions to:", filter.To)
		} else {
			errorMessage := fmt.Sprintf("Invalid 'to' date format: %v", err)
			return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
		}
	}

	if args.Page != nil && *args.Page > 0 {
		filter.Page = *args.Page
		logger.Info("Transaction page:", filter.Page)
	}

	if args.PageSize != nil && *args.PageSize > 0 {
		filter.PageSize = *args.PageSize
		logger.Info("Transaction page size:", filter.PageSize)
	}

	if args.CreatedFrom != nil && *args.CreatedFrom != "" {
		createdFrom, err := time.Parse(time.RFC3339, *args.CreatedFrom)
		if err == nil {
			filter.CreatedAtFrom = createdFrom
			logger.Info("Filter transactions created from:", filter.CreatedAtFrom)
		} else {
			errorMessage := fmt.Sprintf("Invalid 'created_from' date format: %v", err)
			return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
		}
	}

	if args.IDs != nil && len(*args.IDs) > 0 {
		filter.IDs = *args.IDs
		logger.Info("Filter by IDs:", filter.IDs)
	}

	transactions, err := t.client.GetTransactions(args.AccountID, filter)
	if err != nil {
		errorMessage := fmt.Sprintf("Error getting transactions: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	transactionsJSON, err := json.Marshal(transactions)
	if err != nil {
		errorMessage := fmt.Sprintf("Error marshalling transactions: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	return mcp.NewToolResponse(mcp.NewTextContent(string(transactionsJSON))), nil
}

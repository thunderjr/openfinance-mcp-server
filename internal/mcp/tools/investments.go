package tools

import (
	"encoding/json"
	"fmt"

	mcp "github.com/metoro-io/mcp-golang"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
	internalMcp "github.com/thunderjr/openfinance-mcp-server/internal/provider/mcp"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/pluggy"
)

type ListInvestmentsArgs struct {
	ConnectionID string  `json:"connection_id" jsonschema:"required,description=The connection ID to list investments for"`
	Type         *string `json:"type,omitempty" jsonschema:"description=Filter investments by type (COE, EQUITY, ETF, FIXED_INCOME, MUTUAL_FUND, SECURITY, OTHER)"`
	Page         *int    `json:"page,omitempty" jsonschema:"description=Page number for pagination (default: 1)"`
	PageSize     *int    `json:"page_size,omitempty" jsonschema:"description=Number of results per page (default: 20, max: 500)"`
}

type ListInvestmentsTool struct {
	client *pluggy.Client
}

func NewListInvestmentsTool(client *pluggy.Client) *ListInvestmentsTool {
	return &ListInvestmentsTool{client}
}

func (t *ListInvestmentsTool) Name() string {
	return "list_investments"
}

func (t *ListInvestmentsTool) Description() string {
	return `Lists investment positions associated with a bank connection.

Use this tool when:
- The user asks about their investments or portfolio
- Checking investment values or performance
- Looking for specific investment types (stocks, ETFs, fixed income)

Returns: List of investment positions with current values, types, and quantities.`
}

func (t *ListInvestmentsTool) Handle() internalMcp.ToolHandlerFunc {
	return t.handleListInvestments
}

func (t *ListInvestmentsTool) handleListInvestments(args ListInvestmentsArgs) (*mcp.ToolResponse, error) {
	if args.ConnectionID == "" {
		errorMessage := "Connection ID is required"
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	filter := &pluggy.InvestmentsFilter{}

	if args.Type != nil && *args.Type != "" {
		filter.Type = pluggy.InvestmentType(*args.Type)
		logger.Info("Filter investments by type:", filter.Type)
	}

	if args.Page != nil && *args.Page > 0 {
		filter.Page = *args.Page
		logger.Info("Investments page:", filter.Page)
	}

	if args.PageSize != nil && *args.PageSize > 0 {
		filter.PageSize = *args.PageSize
		logger.Info("Investments page size:", filter.PageSize)
	}

	logger.Info("Listing investments for connection:", args.ConnectionID)

	investments, err := t.client.GetInvestments(args.ConnectionID, filter)
	if err != nil {
		errorMessage := fmt.Sprintf("Error listing investments: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	investmentsJSON, err := json.Marshal(investments)
	if err != nil {
		errorMessage := fmt.Sprintf("Error marshalling investments: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	return mcp.NewToolResponse(mcp.NewTextContent(string(investmentsJSON))), nil
}

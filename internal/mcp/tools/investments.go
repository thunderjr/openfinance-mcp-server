package tools

import (
	mcp "github.com/metoro-io/mcp-golang"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
	internalMcp "github.com/thunderjr/openfinance-mcp-server/internal/provider/mcp"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/pluggy"
)

const (
	InvestmentDefaultPageSize = 20
	InvestmentMaxPageSize     = 50
)

type ListInvestmentsArgs struct {
	ConnectionID string  `json:"connection_id" jsonschema:"required,description=The connection ID to list investments for"`
	Type         *string `json:"type,omitempty" jsonschema:"description=Filter investments by type (COE, EQUITY, ETF, FIXED_INCOME, MUTUAL_FUND, SECURITY, OTHER)"`
	Page         *int    `json:"page,omitempty" jsonschema:"description=Page number for pagination (default: 1)"`
	PageSize     *int    `json:"page_size,omitempty" jsonschema:"description=Number of results per page (default: 20, max: 50)"`
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
		return validationError("Connection ID is required")
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

	pageSize := InvestmentDefaultPageSize
	if args.PageSize != nil && *args.PageSize > 0 {
		pageSize = min(*args.PageSize, InvestmentMaxPageSize)
	}
	filter.PageSize = pageSize
	logger.Info("Investments page size:", filter.PageSize)

	logger.Info("Listing investments for connection:", args.ConnectionID)

	investments, err := t.client.GetInvestments(args.ConnectionID, filter)
	if err != nil {
		return toolError("Error listing investments", err)
	}

	return jsonResponse(investments)
}

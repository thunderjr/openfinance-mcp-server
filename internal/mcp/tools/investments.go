package tools

import (
	"encoding/json"
	"fmt"

	mcp "github.com/metoro-io/mcp-golang"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
	internalMcp "github.com/thunderjr/openfinance-mcp-server/internal/provider/mcp"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/pluggy"
)

type InvestmentsArgs struct {
	ItemID   string  `json:"item_id" jsonschema:"required,description=The ID of the item to retrieve investments for"`
	Type     *string `json:"type,omitempty" jsonschema:"description=Filter investments by type (COE, EQUITY, ETF, FIXED_INCOME, MUTUAL_FUND, SECURITY, OTHER)"`
	Page     *int    `json:"page,omitempty" jsonschema:"description=Page number for pagination (default: 1)"`
	PageSize *int    `json:"page_size,omitempty" jsonschema:"description=Number of results per page (default: 20, max: 500)"`
}

type PluggyInvestmentsTool struct {
	client *pluggy.Client
}

func NewPluggyInvestmentsTool(client *pluggy.Client) *PluggyInvestmentsTool {
	return &PluggyInvestmentsTool{client}
}

func (t *PluggyInvestmentsTool) Name() string {
	return "get_item_investments"
}

func (t *PluggyInvestmentsTool) Description() string {
	return "Retrieves investments associated with a specific item with optional filters"
}

func (t *PluggyInvestmentsTool) Handle() internalMcp.ToolHandlerFunc {
	return t.handleGetInvestments
}

func (t *PluggyInvestmentsTool) handleGetInvestments(args InvestmentsArgs) (*mcp.ToolResponse, error) {
	if args.ItemID == "" {
		errorMessage := "Item ID is required"
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

	investments, err := t.client.GetInvestments(args.ItemID, filter)
	if err != nil {
		errorMessage := fmt.Sprintf("Error getting investments: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	investmentsJSON, err := json.Marshal(investments)
	if err != nil {
		errorMessage := fmt.Sprintf("Error marshalling investments: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	return mcp.NewToolResponse(mcp.NewTextContent(string(investmentsJSON))), nil
}

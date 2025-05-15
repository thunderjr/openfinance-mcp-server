package resources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/thunderjr/openfinance-mcp-server/internal/infra/logger"
	"github.com/thunderjr/openfinance-mcp-server/internal/infra/pluggy"
	"github.com/thunderjr/openfinance-mcp-server/internal/infra/uri"
)

type PluggyGetInvestmentsResource struct {
	client *pluggy.Client
}

func NewGetInvestmentsResource(client *pluggy.Client) *PluggyGetInvestmentsResource {
	return &PluggyGetInvestmentsResource{client}
}

func (p *PluggyGetInvestmentsResource) Resource() mcp.ResourceTemplate {
	return mcp.NewResourceTemplate(
		"items://{item_id}/investments?type={type}&page={page}&pageSize={pageSize}",
		"Item Investments",
		mcp.WithTemplateDescription("Retrieves investments associated with a specific item. Supports query parameters: type, page, and pageSize. 'type' filters by investment type (COE, EQUITY, ETF, FIXED_INCOME, MUTUAL_FUND, SECURITY, OTHER), 'page' and 'pageSize' control pagination (default: page 1, size 20, max 500)."),
		mcp.WithTemplateMIMEType("application/json"),
	)
}

func (p *PluggyGetInvestmentsResource) Handle(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	resource, err := uri.Match("items://{item_id}/investments", request.Params.URI)
	if err != nil {
		return nil, fmt.Errorf("error parsing investments uri: err: %w", err)
	}

	itemID, exists := resource.PathParams["item_id"]
	if !exists {
		return nil, fmt.Errorf("item ID not found in URI parameters")
	}

	filter := &pluggy.InvestmentsFilter{}

	if investmentType, exists := resource.QueryParams["type"]; exists && investmentType[0] != "" {
		filter.Type = pluggy.InvestmentType(investmentType[0])
		logger.Info("Filter investments by type:", filter.Type)
	}

	page := uri.GetIntParam(resource.QueryParams, "page")
	if page > 0 {
		filter.Page = page
		logger.Info("Investments page:", filter.Page)
	}

	pageSize := uri.GetIntParam(resource.QueryParams, "pageSize")
	if pageSize > 0 {
		filter.PageSize = pageSize
		logger.Info("Investments page size:", filter.PageSize)
	}

	investments, err := p.client.GetInvestments(itemID, filter)
	if err != nil {
		return nil, fmt.Errorf("error getting investments: err: %w", err)
	}

	strInvestments, err := json.Marshal(investments)
	if err != nil {
		return nil, fmt.Errorf("error marshalling investments: err: %w", err)
	}

	res := []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(strInvestments),
		},
	}

	return res, nil
}

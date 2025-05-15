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

type PluggyGetBillsResource struct {
	client *pluggy.Client
}

func NewGetBillsResource(client *pluggy.Client) *PluggyGetBillsResource {
	return &PluggyGetBillsResource{client}
}

func (p *PluggyGetBillsResource) Resource() mcp.ResourceTemplate {
	return mcp.NewResourceTemplate(
		"accounts://{account_id}/bills",
		"Account Bills",
		mcp.WithTemplateDescription("Retrieves all bills associated with a specific account"),
		mcp.WithTemplateMIMEType("application/json"),
	)
}

func (p *PluggyGetBillsResource) Handle(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	resource, err := uri.Match("accounts://{account_id}/bills", request.Params.URI)
	if err != nil {
		return nil, fmt.Errorf("error parsing bills uri: err: %w", err)
	}

	accountID, exists := resource.PathParams["account_id"]
	if !exists {
		return nil, fmt.Errorf("account ID not found in URI parameters")
	}

	logger.Info("Getting bills for account:", accountID)

	bills, err := p.client.GetBills(accountID)
	if err != nil {
		return nil, fmt.Errorf("error getting bills: err: %w", err)
	}

	strBills, err := json.Marshal(bills)
	if err != nil {
		return nil, fmt.Errorf("error marshalling bills: err: %w", err)
	}

	res := []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(strBills),
		},
	}

	return res, nil
}

type PluggyGetBillResource struct {
	client *pluggy.Client
}

func NewGetBillResource(client *pluggy.Client) *PluggyGetBillResource {
	return &PluggyGetBillResource{client}
}

func (p *PluggyGetBillResource) Resource() mcp.ResourceTemplate {
	return mcp.NewResourceTemplate(
		"bills://{bill_id}",
		"Bill Details",
		mcp.WithTemplateDescription("Retrieves detailed information for a specific bill"),
		mcp.WithTemplateMIMEType("application/json"),
	)
}

func (p *PluggyGetBillResource) Handle(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	resource, err := uri.Match("bills://{bill_id}", request.Params.URI)
	if err != nil {
		return nil, fmt.Errorf("error parsing bill uri: err: %w", err)
	}

	billID, exists := resource.PathParams["bill_id"]
	if !exists {
		return nil, fmt.Errorf("bill ID not found in URI parameters")
	}

	logger.Info("Getting bill details for:", billID)

	bill, err := p.client.GetBill(billID)
	if err != nil {
		return nil, fmt.Errorf("error getting bill: err: %w", err)
	}

	strBill, err := json.Marshal(bill)
	if err != nil {
		return nil, fmt.Errorf("error marshalling bill: err: %w", err)
	}

	res := []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(strBill),
		},
	}

	return res, nil
}

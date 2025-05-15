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

type PluggyGetAccountsResource struct {
	client *pluggy.Client
}

func NewGetAccountsResource(client *pluggy.Client) *PluggyGetAccountsResource {
	return &PluggyGetAccountsResource{client}
}

func (p *PluggyGetAccountsResource) Resource() mcp.ResourceTemplate {
	return mcp.NewResourceTemplate(
		"items://{item_id}/accounts",
		"Item Accounts",
		mcp.WithTemplateDescription("Retrieves all accounts associated with a specific item"),
		mcp.WithTemplateMIMEType("application/json"),
	)
}

func (p *PluggyGetAccountsResource) Handle(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	resource, err := uri.Match("items://{item_id}/accounts", request.Params.URI)
	if err != nil {
		return nil, fmt.Errorf("error parsing accounts uri: err: %w", err)
	}

	itemID, exists := resource.PathParams["item_id"]
	if !exists {
		return nil, fmt.Errorf("item ID not found in URI parameters")
	}

	accounts, err := p.client.GetAccounts(itemID)
	if err != nil {
		return nil, fmt.Errorf("error getting accounts: err: %w", err)
	}

	strAccounts, err := json.Marshal(accounts)
	if err != nil {
		return nil, fmt.Errorf("error marshalling accounts: err: %w", err)
	}

	res := []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(strAccounts),
		},
	}

	return res, nil
}

type PluggyGetAccountResource struct {
	client *pluggy.Client
}

func NewGetAccountResource(client *pluggy.Client) *PluggyGetAccountResource {
	return &PluggyGetAccountResource{client}
}

func (p *PluggyGetAccountResource) Resource() mcp.ResourceTemplate {
	return mcp.NewResourceTemplate(
		"accounts://{account_id}",
		"Account Details",
		mcp.WithTemplateDescription("Retrieves detailed information for a specific account"),
		mcp.WithTemplateMIMEType("application/json"),
	)
}

func (p *PluggyGetAccountResource) Handle(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	resource, err := uri.Match("accounts://{account_id}", request.Params.URI)
	if err != nil {
		return nil, fmt.Errorf("error parsing account uri: err: %w", err)
	}

	accountID, exists := resource.PathParams["account_id"]
	if !exists {
		return nil, fmt.Errorf("account ID not found in URI parameters")
	}

	logger.Info("Getting account details for:", accountID)

	account, err := p.client.GetAccount(accountID)
	if err != nil {
		return nil, fmt.Errorf("error getting account: err: %w", err)
	}

	strAccount, err := json.Marshal(account)
	if err != nil {
		return nil, fmt.Errorf("error marshalling account: err: %w", err)
	}

	res := []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(strAccount),
		},
	}

	return res, nil
}

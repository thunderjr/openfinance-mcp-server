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

type PluggyGetItemResource struct {
	client *pluggy.Client
}

func NewGetItemResource(client *pluggy.Client) *PluggyGetItemResource {
	return &PluggyGetItemResource{client}
}

func (p *PluggyGetItemResource) Resource() mcp.ResourceTemplate {
	return mcp.NewResourceTemplate(
		"items://{item_id}",
		"Item Details",
		mcp.WithTemplateDescription("Retrieves detailed information for a specific item"),
		mcp.WithTemplateMIMEType("application/json"),
	)
}

func (p *PluggyGetItemResource) Handle(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	resource, err := uri.Match("items://{item_id}", request.Params.URI)
	if err != nil {
		return nil, fmt.Errorf("error parsing item uri: err: %w", err)
	}

	itemID, exists := resource.PathParams["item_id"]
	if !exists {
		return nil, fmt.Errorf("item ID not found in URI parameters")
	}

	logger.Info("Getting item details for:", itemID)

	item, err := p.client.GetItem(itemID)
	if err != nil {
		return nil, fmt.Errorf("error getting item: err: %w", err)
	}

	strItem, err := json.Marshal(item)
	if err != nil {
		return nil, fmt.Errorf("error marshalling item: err: %w", err)
	}

	res := []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(strItem),
		},
	}

	return res, nil
}

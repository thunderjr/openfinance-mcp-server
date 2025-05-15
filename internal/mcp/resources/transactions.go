package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/thunderjr/openfinance-mcp-server/internal/infra/logger"
	"github.com/thunderjr/openfinance-mcp-server/internal/infra/pluggy"
	"github.com/thunderjr/openfinance-mcp-server/internal/infra/uri"
)

type PluggyGetTransactionsTool struct {
	client *pluggy.Client
}

func NewGetTransactionsResource(client *pluggy.Client) *PluggyGetTransactionsTool {
	return &PluggyGetTransactionsTool{client}
}

func (p *PluggyGetTransactionsTool) Resource() mcp.ResourceTemplate {
	return mcp.NewResourceTemplate(
		"accounts://{account_id}/transactions?from={from}&to={to}&page={page}&pageSize={pageSize}&createdAtFrom={createdAtFrom}&ids={ids}",
		"Account Transactions",
		mcp.WithTemplateDescription("Queries paginated account transactions by date. Supports query parameters: from, to, page, pageSize, createdAtFrom, and ids. 'from' and 'to' filter by transaction date (format: yyyy-mm-dd), 'page' and 'pageSize' control pagination (default: page 1, size 20, max 500), 'createdAtFrom' filters by creation date (ISO 8601 format), and 'ids' filters by specific transaction IDs."),
		mcp.WithTemplateMIMEType("application/json"),
	)
}

func (p *PluggyGetTransactionsTool) Handle(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	resource, err := uri.Match("accounts://{account_id}/transactions", request.Params.URI)
	if err != nil {
		return nil, fmt.Errorf("error parsing transactions uri: err: %w", err)
	}

	accountID, exists := resource.PathParams["account_id"]
	if !exists {
		return nil, fmt.Errorf("account ID not found in URI parameters")
	}

	filter := &pluggy.TransactionFilter{}

	from := uri.GetTimeParam(resource.QueryParams, "from", "2006-01-02")
	if !from.IsZero() {
		filter.From = from
		logger.Info("Filter transactions from:", filter.From)
	}

	to := uri.GetTimeParam(resource.QueryParams, "to", "2006-01-02")
	if !to.IsZero() {
		filter.To = to
		logger.Info("Filter transactions to:", filter.To)
	}

	page := uri.GetIntParam(resource.QueryParams, "page")
	if page > 0 {
		filter.Page = page
		logger.Info("Transaction page:", filter.Page)
	}

	pageSize := uri.GetIntParam(resource.QueryParams, "pageSize")
	if pageSize > 0 {
		filter.PageSize = pageSize
		logger.Info("Transaction page size:", filter.PageSize)
	}

	createdAtFrom := uri.GetTimeParam(resource.QueryParams, "createdAtFrom", time.RFC3339)
	if !createdAtFrom.IsZero() {
		filter.CreatedAtFrom = createdAtFrom
		logger.Info("Filter transactions created from:", filter.CreatedAtFrom)
	}

	ids := uri.GetStringArrayParam(resource.QueryParams, "ids")
	if len(ids) > 0 {
		filter.IDs = ids
		logger.Info("Filter by IDs:", filter.IDs)
	}

	transactions, err := p.client.GetTransactions(accountID, filter)
	if err != nil {
		return nil, fmt.Errorf("error getting transactions: err: %w", err)
	}

	strTxs, err := json.Marshal(transactions)
	if err != nil {
		return nil, fmt.Errorf("error marshalling transactions: err: %w", err)
	}

	res := []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(strTxs),
		},
	}

	return res, nil
}

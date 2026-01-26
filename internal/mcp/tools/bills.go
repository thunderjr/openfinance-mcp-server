package tools

import (
	mcp "github.com/metoro-io/mcp-golang"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
	internalMcp "github.com/thunderjr/openfinance-mcp-server/internal/provider/mcp"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/pluggy"
)

type ListCreditCardBillsArgs struct {
	AccountID string `json:"account_id" jsonschema:"required,description=The credit card account ID to list bills for"`
}

type ListCreditCardBillsTool struct {
	client *pluggy.Client
}

func NewListCreditCardBillsTool(client *pluggy.Client) *ListCreditCardBillsTool {
	return &ListCreditCardBillsTool{client}
}

func (t *ListCreditCardBillsTool) Name() string {
	return "list_credit_card_bills"
}

func (t *ListCreditCardBillsTool) Description() string {
	return `Lists credit card bills for an account.

Use this tool when:
- The user asks about their credit card bills
- Checking bill due dates or amounts
- Looking at credit card payment history

Returns: List of bills with amounts, due dates, and payment status.`
}

func (t *ListCreditCardBillsTool) Handle() internalMcp.ToolHandlerFunc {
	return t.handleListCreditCardBills
}

func (t *ListCreditCardBillsTool) handleListCreditCardBills(args ListCreditCardBillsArgs) (*mcp.ToolResponse, error) {
	if args.AccountID == "" {
		return validationError("Account ID is required")
	}

	logger.Info("Listing credit card bills for account:", args.AccountID)

	bills, err := t.client.GetBills(args.AccountID)
	if err != nil {
		return toolError("Error listing credit card bills", err)
	}

	return jsonResponse(bills)
}

type GetCreditCardBillArgs struct {
	BillID string `json:"bill_id" jsonschema:"required,description=The bill ID to retrieve details for"`
}

type GetCreditCardBillTool struct {
	client *pluggy.Client
}

func NewGetCreditCardBillTool(client *pluggy.Client) *GetCreditCardBillTool {
	return &GetCreditCardBillTool{client}
}

func (t *GetCreditCardBillTool) Name() string {
	return "get_credit_card_bill"
}

func (t *GetCreditCardBillTool) Description() string {
	return `Gets detailed information for a specific credit card bill.

Use this tool when:
- The user asks about a specific bill's details
- You need to see individual transactions within a bill
- Checking payment details for a specific billing period

Returns: Bill details including total amount, due date, transactions, and payment status.`
}

func (t *GetCreditCardBillTool) Handle() internalMcp.ToolHandlerFunc {
	return t.handleGetCreditCardBill
}

func (t *GetCreditCardBillTool) handleGetCreditCardBill(args GetCreditCardBillArgs) (*mcp.ToolResponse, error) {
	if args.BillID == "" {
		return validationError("Bill ID is required")
	}

	logger.Info("Getting credit card bill details for:", args.BillID)

	bill, err := t.client.GetBill(args.BillID)
	if err != nil {
		return toolError("Error getting credit card bill", err)
	}

	return jsonResponse(bill)
}

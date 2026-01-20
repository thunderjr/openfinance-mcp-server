package tools

import (
	"encoding/json"
	"fmt"

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
		errorMessage := "Account ID is required"
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	logger.Info("Listing credit card bills for account:", args.AccountID)

	bills, err := t.client.GetBills(args.AccountID)
	if err != nil {
		errorMessage := fmt.Sprintf("Error listing credit card bills: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	billsJSON, err := json.Marshal(bills)
	if err != nil {
		errorMessage := fmt.Sprintf("Error marshalling bills: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	return mcp.NewToolResponse(mcp.NewTextContent(string(billsJSON))), nil
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
		errorMessage := "Bill ID is required"
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	logger.Info("Getting credit card bill details for:", args.BillID)

	bill, err := t.client.GetBill(args.BillID)
	if err != nil {
		errorMessage := fmt.Sprintf("Error getting credit card bill: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	billJSON, err := json.Marshal(bill)
	if err != nil {
		errorMessage := fmt.Sprintf("Error marshalling bill: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	return mcp.NewToolResponse(mcp.NewTextContent(string(billJSON))), nil
}

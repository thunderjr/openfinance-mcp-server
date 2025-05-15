package tools

import (
	"encoding/json"
	"fmt"

	mcp "github.com/metoro-io/mcp-golang"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
	internalMcp "github.com/thunderjr/openfinance-mcp-server/internal/provider/mcp"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/pluggy"
)

type BillsArgs struct {
	AccountID string `json:"account_id" jsonschema:"required,description=The ID of the account to retrieve bills for"`
}

type PluggyBillsTool struct {
	client *pluggy.Client
}

func NewPluggyBillsTool(client *pluggy.Client) *PluggyBillsTool {
	return &PluggyBillsTool{client}
}

func (t *PluggyBillsTool) Name() string {
	return "get_account_bills"
}

func (t *PluggyBillsTool) Description() string {
	return "Retrieves all bills associated with a specific account"
}

func (t *PluggyBillsTool) Handle() internalMcp.ToolHandlerFunc {
	return t.handleGetBills
}

func (t *PluggyBillsTool) handleGetBills(args BillsArgs) (*mcp.ToolResponse, error) {
	if args.AccountID == "" {
		errorMessage := "Account ID is required"
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	logger.Info("Getting bills for account:", args.AccountID)

	bills, err := t.client.GetBills(args.AccountID)
	if err != nil {
		errorMessage := fmt.Sprintf("Error getting bills: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	billsJSON, err := json.Marshal(bills)
	if err != nil {
		errorMessage := fmt.Sprintf("Error marshalling bills: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	return mcp.NewToolResponse(mcp.NewTextContent(string(billsJSON))), nil
}

type BillArgs struct {
	BillID string `json:"bill_id" jsonschema:"required,description=The ID of the bill to retrieve details for"`
}

type PluggyBillTool struct {
	client *pluggy.Client
}

func NewPluggyBillTool(client *pluggy.Client) *PluggyBillTool {
	return &PluggyBillTool{client}
}

func (t *PluggyBillTool) Name() string {
	return "get_bill_details"
}

func (t *PluggyBillTool) Description() string {
	return "Retrieves detailed information for a specific bill"
}

func (t *PluggyBillTool) Handle() internalMcp.ToolHandlerFunc {
	return t.handleGetBill
}

func (t *PluggyBillTool) handleGetBill(args BillArgs) (*mcp.ToolResponse, error) {
	if args.BillID == "" {
		errorMessage := "Bill ID is required"
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	logger.Info("Getting bill details for:", args.BillID)

	bill, err := t.client.GetBill(args.BillID)
	if err != nil {
		errorMessage := fmt.Sprintf("Error getting bill: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	billJSON, err := json.Marshal(bill)
	if err != nil {
		errorMessage := fmt.Sprintf("Error marshalling bill: %v", err)
		return mcp.NewToolResponse(mcp.NewTextContent(errorMessage)), fmt.Errorf(errorMessage)
	}

	return mcp.NewToolResponse(mcp.NewTextContent(string(billJSON))), nil
}

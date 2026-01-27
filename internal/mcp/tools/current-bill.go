package tools

import (
	"fmt"
	"sort"
	"time"

	"github.com/shopspring/decimal"

	mcp "github.com/metoro-io/mcp-golang"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
	internalMcp "github.com/thunderjr/openfinance-mcp-server/internal/provider/mcp"
	"github.com/thunderjr/openfinance-mcp-server/internal/provider/pluggy"
)

type GetCurrentBillArgs struct {
	AccountID string `json:"account_id" jsonschema:"required,description=The credit card account ID"`
}

type GetCurrentBillTool struct {
	client *pluggy.Client
}

func NewGetCurrentBillTool(client *pluggy.Client) *GetCurrentBillTool {
	return &GetCurrentBillTool{client}
}

func (t *GetCurrentBillTool) Name() string {
	return "get_current_bill"
}

func (t *GetCurrentBillTool) Description() string {
	return `Calculates the actual current credit card bill amount.

Use this tool when:
- The user asks "how much is my credit card bill?"
- You need the actual amount due (excluding future installments)
- Checking the current billing period summary

This tool automatically:
1. Gets the account balance (total outstanding including future installments)
2. Fetches open/pending transactions for the current billing period
3. Identifies installment purchases and calculates future installment amounts
4. Subtracts future installments from the balance to get the actual bill

Returns: Summary with total_balance, future_installments_deduction, actual_bill_amount, transaction_count, and top transactions.`
}

func (t *GetCurrentBillTool) Handle() internalMcp.ToolHandlerFunc {
	return t.handleGetCurrentBill
}

type currentBillResponse struct {
	AccountID                   string           `json:"account_id"`
	TotalBalance                decimal.Decimal  `json:"total_balance"`
	FutureInstallmentsDeduction decimal.Decimal  `json:"future_installments_deduction"`
	ActualBillAmount            decimal.Decimal  `json:"actual_bill_amount"`
	TransactionCount            int              `json:"transaction_count"`
	InstallmentCount            int              `json:"installment_count"`
	DueDate                     string           `json:"due_date,omitempty"`
	TopTransactions             []topTransaction `json:"top_transactions,omitempty"`
}

type topTransaction struct {
	Description string          `json:"description"`
	Amount      decimal.Decimal `json:"amount"`
	Installment string          `json:"installment,omitempty"`
}

func (t *GetCurrentBillTool) handleGetCurrentBill(args GetCurrentBillArgs) (*mcp.ToolResponse, error) {
	if args.AccountID == "" {
		return validationError("Account ID is required")
	}

	logger.Info("Getting current bill for account:", args.AccountID)

	// 1. Get account to retrieve balance and credit data
	account, err := t.client.GetAccount(args.AccountID)
	if err != nil {
		return toolError("Error getting account", err)
	}

	if account.Type != "CREDIT" {
		return validationError("Account is not a credit card account")
	}

	totalBalance := account.Balance

	// 2. Get due date from creditData
	dueDate := ""
	if account.CreditData != nil {
		dueDate = account.CreditData.BalanceDueDate
	}

	// 3. Fetch open bill transactions (current billing period)
	now := time.Now()
	from := now.AddDate(0, -1, -5) // slightly more than a month back
	txns, err := t.client.GetTransactions(args.AccountID, &pluggy.TransactionFilter{
		From:     from,
		To:       now,
		PageSize: 500,
	})
	if err != nil {
		return toolError("Error fetching transactions", err)
	}

	// 4. Filter to open bill transactions and calculate future installments
	futureDeduction := decimal.Zero
	installmentCount := 0
	var openTxns []pluggy.Transaction

	for _, tx := range txns.Results {
		// Open bill transactions have no billId
		if tx.CreditCardMetadata == nil || tx.CreditCardMetadata.BillID != "" {
			continue
		}
		openTxns = append(openTxns, tx)

		if tx.CreditCardMetadata.TotalInstallments > 1 {
			future := tx.CreditCardMetadata.TotalInstallments - tx.CreditCardMetadata.InstallmentNumber
			if future > 0 {
				installmentCount++
				// Each future installment has the same amount as this one
				futureDeduction = futureDeduction.Add(
					tx.Amount.Abs().Mul(decimal.NewFromInt(int64(future))),
				)
			}
		}
	}

	actualBill := totalBalance.Sub(futureDeduction)

	// 5. Build top transactions (top 5 by absolute amount)
	sort.Slice(openTxns, func(i, j int) bool {
		return openTxns[i].Amount.Abs().GreaterThan(openTxns[j].Amount.Abs())
	})

	topCount := 5
	if len(openTxns) < topCount {
		topCount = len(openTxns)
	}

	topTxns := make([]topTransaction, topCount)
	for i := 0; i < topCount; i++ {
		tx := openTxns[i]
		top := topTransaction{
			Description: tx.Description,
			Amount:      tx.Amount,
		}
		if tx.CreditCardMetadata != nil && tx.CreditCardMetadata.TotalInstallments > 1 {
			top.Installment = fmt.Sprintf("%d/%d",
				tx.CreditCardMetadata.InstallmentNumber,
				tx.CreditCardMetadata.TotalInstallments,
			)
		}
		topTxns[i] = top
	}

	return jsonResponse(currentBillResponse{
		AccountID:                   args.AccountID,
		TotalBalance:                totalBalance,
		FutureInstallmentsDeduction: futureDeduction,
		ActualBillAmount:            actualBill,
		TransactionCount:            len(openTxns),
		InstallmentCount:            installmentCount,
		DueDate:                     dueDate,
		TopTransactions:             topTxns,
	})
}

package pluggy

import (
	"time"

	"github.com/shopspring/decimal"
)

type itemResponse struct {
	ID              string    `json:"id"`
	Status          string    `json:"status"`
	ExecutionStatus string    `json:"executionStatus"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	LastUpdatedAt   time.Time `json:"lastUpdatedAt"`
	NextAutoSyncAt  time.Time `json:"nextAutoSyncAt"`
	Products        []string  `json:"products"`
	Connector       struct {
		ID             int      `json:"id"`
		Name           string   `json:"name"`
		InstitutionURL string   `json:"institutionUrl"`
		ImageURL       string   `json:"imageUrl"`
		PrimaryColor   string   `json:"primaryColor"`
		Type           string   `json:"type"`
		Country        string   `json:"country"`
		Products       []string `json:"products"`
		Oauth          bool     `json:"oauth"`
		Health         struct {
			Status string `json:"status"`
		} `json:"health"`
	} `json:"connector"`
}

type account struct {
	ID           string          `json:"id"`
	Type         string          `json:"type"`    // "BANK" or "CREDIT"
	Subtype      string          `json:"subtype"` // "CHECKING_ACCOUNT", "CREDIT_CARD"
	Number       string          `json:"number"`
	Name         string          `json:"name"`
	Balance      decimal.Decimal `json:"balance"`
	ItemID       string          `json:"itemId"`
	TaxNumber    string          `json:"taxNumber,omitempty"`
	CurrencyCode string          `json:"currencyCode"`
	BankData     *bankData       `json:"bankData,omitempty"`
	CreditData   *creditData     `json:"creditData,omitempty"`
}

type bankData struct {
	TransferNumber               string          `json:"transferNumber"`
	ClosingBalance               decimal.Decimal `json:"closingBalance"`
	AutomaticallyInvestedBalance decimal.Decimal `json:"automaticallyInvestedBalance"`
}

type creditData struct {
	Brand                string          `json:"brand"`          // "MASTERCARD"
	BalanceDueDate       string          `json:"balanceDueDate"` // "2025-01-09"
	AvailableCreditLimit decimal.Decimal `json:"availableCreditLimit"`
	MinimumPayment       decimal.Decimal `json:"minimumPayment"`
	CreditLimit          decimal.Decimal `json:"creditLimit"`
	Status               string          `json:"status"` // "ACTIVE"
}

type Transaction struct {
	ID                      string              `json:"id"`
	Description             string              `json:"description"`
	DescriptionRaw          *string             `json:"descriptionRaw,omitempty"`
	CurrencyCode            string              `json:"currencyCode"`
	Amount                  decimal.Decimal     `json:"amount"`
	Date                    time.Time           `json:"date"`
	Balance                 decimal.Decimal     `json:"balance"`
	Category                string              `json:"category"`
	AccountID               string              `json:"accountId"`
	ProviderCode            interface{}         `json:"providerCode"`
	Status                  string              `json:"status"`
	PaymentData             *paymentData        `json:"paymentData,omitempty"`
	Type                    string              `json:"type"`
	AmountInAccountCurrency *decimal.Decimal    `json:"amountInAccountCurrency,omitempty"`
	Merchant                *merchant           `json:"merchant,omitempty"`
	CreditCardMetadata      *creditCardMetadata `json:"creditCardMetadata,omitempty"`
	OperationType           *string             `json:"operationType,omitempty"`
}

type paymentData struct {
	Payer           payer    `json:"payer"`
	Reason          string   `json:"reason"`
	Receiver        receiver `json:"receiver"`
	PaymentMethod   string   `json:"paymentMethod"`
	ReferenceNumber string   `json:"referenceNumber"`
}

type payer struct {
	Name           string         `json:"name"`
	BranchNumber   string         `json:"branchNumber"`
	AccountNumber  string         `json:"accountNumber"`
	RoutingNumber  string         `json:"routingNumber"`
	DocumentNumber documentNumber `json:"documentNumber"`
}

type receiver struct {
	Name           string         `json:"name"`
	BranchNumber   string         `json:"branchNumber"`
	AccountNumber  string         `json:"accountNumber"`
	RoutingNumber  string         `json:"routingNumber"`
	DocumentNumber documentNumber `json:"documentNumber"`
}

type documentNumber struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type merchant struct {
	Name         string `json:"name"`
	BusinessName string `json:"businessName"`
	Cnpj         string `json:"cnpj"`
	Category     string `json:"category"`
	Cnae         string `json:"cnae"`
}

type creditCardMetadata struct {
	PayeeMCC          int       `json:"payeeMCC"`
	CardNumber        string    `json:"cardNumber"`
	TotalAmount       int       `json:"totalAmount"`
	TotalInstallments int       `json:"totalInstallments"`
	InstallmentNumber int       `json:"installmentNumber"`
	PurchaseDate      time.Time `json:"purchaseDate"`
	BillID            string    `json:"billId"`
}

type Investment struct {
	ID                   string          `json:"id"`
	Code                 string          `json:"code"`
	Name                 string          `json:"name"`
	Balance              decimal.Decimal `json:"balance"`
	CurrencyCode         string          `json:"currencyCode"`
	Type                 string          `json:"type"`
	Subtype              string          `json:"subtype"`
	LastMonthRate        decimal.Decimal `json:"lastMonthRate,omitempty"`
	AnnualRate           decimal.Decimal `json:"annualRate,omitempty"`
	LastTwelveMonthsRate decimal.Decimal `json:"lastTwelveMonthsRate,omitempty"`
	ItemID               string          `json:"itemId"`
	Value                decimal.Decimal `json:"value,omitempty"`
	Quantity             decimal.Decimal `json:"quantity,omitempty"`
	Amount               decimal.Decimal `json:"amount"`
	Taxes                decimal.Decimal `json:"taxes"`
	Taxes2               decimal.Decimal `json:"taxes2"`
	Date                 time.Time       `json:"date"`
	Owner                string          `json:"owner"`
	Number               interface{}     `json:"number"`
	AmountProfit         decimal.Decimal `json:"amountProfit"`
	AmountWithdrawal     decimal.Decimal `json:"amountWithdrawal"`
	AmountOriginal       decimal.Decimal `json:"amountOriginal,omitempty"`
	Status               string          `json:"status"`
	Transactions         []struct {
		TradeDate    time.Time       `json:"tradeDate"`
		Date         time.Time       `json:"date"`
		Description  string          `json:"description"`
		Quantity     decimal.Decimal `json:"quantity"`
		Value        decimal.Decimal `json:"value"`
		Amount       decimal.Decimal `json:"amount"`
		Type         string          `json:"type"`
		MovementType string          `json:"movementType"`
	} `json:"transactions,omitempty"`
	Rate            decimal.Decimal `json:"rate,omitempty"`
	RateType        string          `json:"rateType,omitempty"`
	FixedAnnualRate decimal.Decimal `json:"fixedAnnualRate,omitempty"`
	Issuer          string          `json:"issuer,omitempty"`
	IssuerCNPJ      string          `json:"issuerCNPJ,omitempty"`
	IssueDate       time.Time       `json:"issueDate,omitempty"`
	Institution     struct {
		Name   string `json:"name"`
		Number string `json:"number"`
	} `json:"institution,omitempty"`
}

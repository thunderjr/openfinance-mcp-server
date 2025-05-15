package pluggy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/thunderjr/openfinance-mcp-server/internal/provider/logger"
)

type FinanceCharge struct {
	ID             string  `json:"id"`
	Type           string  `json:"type"`
	Amount         float64 `json:"amount"`
	CurrencyCode   string  `json:"currencyCode"`
	AdditionalInfo string  `json:"additionalInfo,omitempty"`
}

type Bill struct {
	ID                      string          `json:"id"`
	DueDate                 time.Time       `json:"dueDate"`
	TotalAmount             float64         `json:"totalAmount"`
	TotalAmountCurrencyCode string          `json:"totalAmountCurrencyCode"`
	MinimumPaymentAmount    float64         `json:"minimumPaymentAmount,omitempty"`
	AllowsInstallments      bool            `json:"allowsInstallments,omitempty"`
	FinanceCharges          []FinanceCharge `json:"financeCharges"`
}

func (c *Client) GetBills(accountID string) (*paginatedResponse[Bill], error) {
	if accountID == "" {
		return nil, fmt.Errorf("pluggyClient.GetBills: accountID is required")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.pluggy.ai/bills?accountId=%s", accountID), nil)
	if err != nil {
		return nil, fmt.Errorf("pluggyClient.GetBills: error creating request: %w", err)
	}
	req.Header.Set("X-API-KEY", c.apiKey)

	res, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("pluggyClient.GetBills: error making request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("pluggyClient.GetBills: error reading response body: %w", err)
		}
		logger.Debug("body: \n", string(bodyBytes))

		return nil, fmt.Errorf("pluggyClient.GetBills: request failed with status %d", res.StatusCode)
	}

	var data paginatedResponse[Bill]
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("pluggyClient.GetBills: error decoding response: %w", err)
	}

	return &data, nil
}

func (c *Client) GetBill(billID string) (*Bill, error) {
	if billID == "" {
		return nil, fmt.Errorf("pluggyClient.GetBill: billID is required")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.pluggy.ai/bills/%s", billID), nil)
	if err != nil {
		return nil, fmt.Errorf("pluggyClient.GetBill: error creating request: %w", err)
	}
	req.Header.Set("X-API-KEY", c.apiKey)

	res, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("pluggyClient.GetBill: error making request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("pluggyClient.GetBill: error reading response body: %w", err)
		}
		logger.Debug("body: \n", string(bodyBytes))

		return nil, fmt.Errorf("pluggyClient.GetBill: request failed with status %d", res.StatusCode)
	}

	var data Bill
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("pluggyClient.GetBill: error decoding response: %w", err)
	}

	return &data, nil
}

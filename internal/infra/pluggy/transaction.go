package pluggy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/thunderjr/openfinance-mcp-server/internal/core"
)

type TransactionFilter struct {
	IDs           []string  `json:"ids,omitempty"`
	Page          int       `json:"page,omitempty"`          // default: 1
	PageSize      int       `json:"pageSize,omitempty"`      // default: 20 max. 500
	To            time.Time `json:"to,omitempty"`            // yyyy-mm-dd
	From          time.Time `json:"from,omitempty"`          // yyyy-mm-dd
	CreatedAtFrom time.Time `json:"createdAtFrom,omitempty"` // ISO 8601
}

func (c *Client) GetTransactions(accountID string, query *TransactionFilter) (*paginatedResponse[Transaction], error) {
	c.rateLimiter.wait()

	q := url.Values{}
	url := "https://api.pluggy.ai/transactions"
	q.Set("accountId", accountID)

	if query != nil {
		if len(query.IDs) > 0 {
			q.Set("ids", strings.Join(query.IDs, ","))
		} else {
			if !query.From.IsZero() {
				q.Set("from", query.From.Format("2006-01-02"))
			}
			if !query.To.IsZero() {
				q.Set("to", query.To.Format("2006-01-02"))
			}
		}
		if query.PageSize > 0 {
			q.Set("pageSize", fmt.Sprintf("%d", query.PageSize))
		}
		if query.Page > 0 {
			q.Set("page", fmt.Sprintf("%d", query.Page))
		}
	}

	url = fmt.Sprintf("%s?%s", url, q.Encode())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("pluggy_client: error creating request: %w", err)
	}
	req.Header.Set("X-API-KEY", c.apiKey)

	res, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("pluggy_client: error item making request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pluggy_client: error item request failed with status %d", res.StatusCode)
	}

	var data paginatedResponse[Transaction]
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("pluggy_client: error item decoding response: %w", err)
	}

	return &data, nil
}

func (tx *Transaction) ToStatement() *core.Statement {
	var (
		currentInstallment int
		totalInstallments  int
	)

	name := tx.Description
	if tx.CreditCardMetadata != nil {
		currentInstallment = tx.CreditCardMetadata.InstallmentNumber
		totalInstallments = tx.CreditCardMetadata.TotalInstallments

		re := regexp.MustCompile(`\d{1,2}\/\d{1,2}`)
		name = strings.TrimSpace(re.ReplaceAllString(tx.Description, ""))
	}

	return &core.Statement{
		Name:               name,
		CorrelationID:      tx.ID,
		Timestamp:          tx.Date,
		Amount:             tx.Amount,
		CurrentInstallment: &currentInstallment,
		TotalInstallments:  &totalInstallments,
		MonthKey:           fmt.Sprintf("%d-%02d", tx.Date.Year(), tx.Date.Month()),
	}
}

package pluggy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/thunderjr/openfinance-mcp-server/internal/infra/logger"
)

type InvestmentType string

const (
	InvestmentTypeCOE         InvestmentType = "COE"
	InvestmentTypeEquity      InvestmentType = "EQUITY"
	InvestmentTypeETF         InvestmentType = "ETF"
	InvestmentTypeFixedIncome InvestmentType = "FIXED_INCOME"
	InvestmentTypeMutualFund  InvestmentType = "MUTUAL_FUND"
	InvestmentTypeSecurity    InvestmentType = "SECURITY"
	InvestmentTypeOther       InvestmentType = "OTHER"
)

type InvestmentsFilter struct {
	Type     InvestmentType `json:"type,omitempty"`
	Page     int            `json:"page,omitempty"`     // default: 1
	PageSize int            `json:"pageSize,omitempty"` // default: 20 max. 500

}

func (c *Client) GetInvestments(itemID string, query *InvestmentsFilter) (*paginatedResponse[Investment], error) {
	c.rateLimiter.wait()

	q := url.Values{}
	url := "https://api.pluggy.ai/investments"
	q.Set("itemId", itemID)

	if query != nil {
		if query.Type != "" {
			q.Set("type", string(query.Type))
		}

		if query.PageSize > 0 {
			q.Set("pageSize", fmt.Sprintf("%d", query.PageSize))
		}
		if query.Page > 0 {
			q.Set("page", fmt.Sprintf("%d", query.Page))
		}
	}

	url = fmt.Sprintf("%s?%s", url, q.Encode())
	logger.Debug("url", url)
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

	var data paginatedResponse[Investment]
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("pluggy_client: error item decoding response: %w", err)
	}

	return &data, nil
}

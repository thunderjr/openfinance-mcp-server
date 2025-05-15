package core

import (
	"time"

	"github.com/shopspring/decimal"
)

type Statement struct {
	Name               string          `json:"name"`
	Tag                string          `json:"tag"`
	CorrelationID      string          `json:"correlation_id"`
	MonthKey           string          `json:"month_key"` // YYYY-MM
	Timestamp          time.Time       `json:"timestamp"` // Purchase date
	Amount             decimal.Decimal `json:"amount" gorm:"type:float"`
	CurrentInstallment *int            `json:"current_installment,omitempty" gorm:"null"`
	TotalInstallments  *int            `json:"total_installments,omitempty" gorm:"null"`
}

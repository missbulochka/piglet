package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type Bill struct {
	ID             string
	BillType       bool
	BillStatus     bool
	Name           string
	CurrentSum     decimal.Decimal
	GoalSum        decimal.Decimal
	Date           time.Time
	MonthlyPayment decimal.Decimal
}

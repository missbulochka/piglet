package models

import "github.com/shopspring/decimal"

type Bill struct {
	ID             string
	BillType       bool
	BillStatus     bool
	Name           string
	CurrentSum     decimal.Decimal
	Date           string
	MonthlyPayment decimal.Decimal
}

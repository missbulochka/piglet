package models

import "github.com/shopspring/decimal"

type Bill struct {
	ID             string
	BillType       bool `validate:"boolean"`
	BillStatus     bool
	Name           string `validate:"required"`
	CurrentSum     decimal.Decimal
	Date           string `validate:"required,datetime=02-01-2006"`
	MonthlyPayment decimal.Decimal
}

package models

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/shopspring/decimal"
)

type Bill struct {
	ID             string
	BillType       bool
	BillStatus     bool
	Name           string
	CurrentSum     decimal.Decimal
	Date           *timestamp.Timestamp
	MonthlyPayment decimal.Decimal
}

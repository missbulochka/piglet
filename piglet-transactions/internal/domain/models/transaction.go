package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type Transaction struct {
	Id         uuid.UUID
	Date       time.Time
	TransType  uint8
	Sum        decimal.Decimal
	Comment    string
	IdCategory uuid.UUID
	DebtType   bool
	IdBillTo   uuid.UUID
	IdBillFrom uuid.UUID
	Person     string
	Repeat     bool
}

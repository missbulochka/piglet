package models

type Bill struct {
	ID             string
	BillType       bool `validate:"required"`
	BillStatus     bool
	Name           string `validate:"required"`
	CurrentSum     float32
	Date           string `validate:"required,datetime=2006-01-02"`
	MonthlyPayment float32
}

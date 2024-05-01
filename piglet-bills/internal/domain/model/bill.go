package models

import "time"

type Bill struct {
	ID             string
	billType       bool
	billStatus     bool
	Name           string
	currentSum     float32
	date           time.Time
	monthlyPayment float32
}

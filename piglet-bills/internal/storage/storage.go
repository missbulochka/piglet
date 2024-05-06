package storage

import "errors"

var (
	ErrBillExists   = errors.New("bill already exists")
	ErrBillNotFound = errors.New("bill not found")
)

const (
	CreateBill = `
		INSERT INTO bills (
		    id,
			"billName",
			"currentSum"
		) VALUES (
			$1, $2, $3
		) RETURNING id, "billName", "currentSum"
	`

	CreateAccount = `
		INSERT INTO accounts (
		    bill_id, 
		    "billStatus"
		) VALUES (
		    $1, $2
		) RETURNING "billStatus"
	`

	CreateGoals = `
		INSERT INTO goals (
		    bill_id, 
		    date,
		    "monthlyPayment"
		) VALUES (
		    $1, $2, $3
		) RETURNING date, "monthlyPayment"
		`
)

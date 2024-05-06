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
			"currentSum",
		    "billType"
		) VALUES (
			$1, $2, $3, $4
		) RETURNING id, "billName", "currentSum", "billType"
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

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
			bill_name,
			current_sum,
		    bill_type
		) VALUES (
			$1, $2, $3, $4
		) RETURNING id, bill_name, current_sum, bill_type
	`

	CreateAccount = `
		INSERT INTO accounts (
		    bill_id, 
		    bill_status
		) VALUES (
		    $1, $2
		) RETURNING bill_status
	`

	CreateGoals = `
		INSERT INTO goals (
		    bill_id, 
		    goal_sum,
		    date,
		    monthly_payment
		) VALUES (
		    $1, $2, $3, $4
		) RETURNING goal_sum, date, monthly_payment
		`
)

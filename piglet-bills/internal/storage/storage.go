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

const (
	// HACK: подумать над аккуратностью запроса
	GetOneBill = `
		SELECT id, bill_name, current_sum, bill_type
		FROM bills 
		WHERE ($1 = '' OR id::text = $1) OR ($2 = '' OR bill_name = $2)
		LIMIT 1;
	`

	GetOneAccount = `
		SELECT bill_status
		FROM accounts
		WHERE bill_id::text = $1
		LIMIT 1 
	`

	GetOneGoal = `
		SELECT goal_sum, date, monthly_payment
		FROM goals
		WHERE bill_id::text = $1
		LIMIT 1
	`

	GetSomeBills = `
		SELECT id, bill_name, current_sum, bill_type
		FROM bills
		WHERE bill_type = $1
		ORDER BY id
	`
	GetAllAccounts = `SELECT bill_status FROM accounts ORDER BY bill_id`
	GetAllGoals    = `
		SELECT goal_sum, date, monthly_payment
		FROM goals
		ORDER BY bill_id
	`
)

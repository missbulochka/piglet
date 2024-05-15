package storage

const (
	InsertTransaction = `INSERT INTO transactions 
    		(id, trans_date, type, sum, comment)
    		VALUES ($1, $2, $3, $4, $5)`

	InsertIncome = `INSERT INTO income
    	(trans_id, id_category, id_bill_to, sender, repeat)
		VALUES ($1, $2, $3, $4, $5)`

	InsertExpense = `INSERT INTO expense
    	(trans_id, id_category, id_bill_from, recipient, repeat)
		VALUES ($1, $2, $3, $4, $5)`

	InsertDebt = `INSERT INTO debt
    	(trans_id, type, id_bill_from, id_bill_to, creditor_debtor) 
		VALUES ($1, $2, $3, $4, $5)`

	InsertTransfer = `INSERT INTO transfer
    	(trans_id, id_bill_from, id_bill_to)
		VALUES ($1, $2, $3)`
)

const (
	GetOneTransaction = `SELECT trans_date, type, sum, comment
			FROM transactions
			WHERE id = $1`
	GetOneIncome = `SELECT id_category, id_bill_to, sender, repeat
			FROM income
			WHERE trans_id = $1`
	GetOneExpense = `SELECT id_category, id_bill_from, recipient, repeat
			FROM expense
			WHERE trans_id = $1`
	GetOneDebt = `SELECT type, id_bill_from, id_bill_to, creditor_debtor
			FROM debt
			WHERE trans_id = $1`
	GetOneTransfer = `SELECT id_bill_from, id_bill_to
			FROM transfer
			WHERE trans_id = $1`
)

const (
	DeleteTransaction = `DELETE FROM transactions WHERE id = $1`
	DeleteIncome      = `DELETE FROM income WHERE trans_id = $1`
	DeleteExpenses    = `DELETE FROM expense WHERE trans_id = $1`
	DeleteDebt        = `DELETE FROM debt WHERE trans_id = $1`
	DeleteTransfer    = `DELETE FROM transfer WHERE trans_id = $1`
)

const (
	GetCategory = `SELECT id, type, name, mandatory
		FROM categories WHERE id = $1`
)

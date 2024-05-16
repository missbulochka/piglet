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
	UpdateTransaction = `UPDATE transactions
		SET trans_date = $2, type = $3, sum = $4, comment = $5
		WHERE id = $1`
	UpdateIncome = `UPDATE income
    	SET id_category = $2, id_bill_to = $3, sender = $4, repeat = $5
		WHERE trans_id = $1`
	UpdateExpense = `UPDATE expense
    	SET id_category = $2, id_bill_from = $3, recipient = $4, repeat = $5
		WHERE trans_id = $1`
	UpdateDebt = `UPDATE debt
    	SET  type = $2, id_bill_from = $3, id_bill_to = $4, creditor_debtor = $5 
		WHERE trans_id = $1`
	UpdateTransfer = `UPDATE transfer
    	SET id_bill_from = $2, id_bill_to = $3
		WHERE trans_id = $1`
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
	GetSomeTransactions = `SELECT id, trans_date, type, sum, comment
			FROM "transactions"
			ORDER BY "trans_date" DESC
			LIMIT $1`
)

const (
	DeleteTransaction = `DELETE FROM transactions WHERE id = $1`
	DeleteIncome      = `DELETE FROM income WHERE trans_id = $1`
	DeleteExpenses    = `DELETE FROM expense WHERE trans_id = $1`
	DeleteDebt        = `DELETE FROM debt WHERE trans_id = $1`
	DeleteTransfer    = `DELETE FROM transfer WHERE trans_id = $1`
)

const (
	InsertCategory = `INSERT INTO categories 
    	(id, type, name, mandatory)
    	VALUES ($1, $2, $3, $4)`
	UpdateCategory = `
		UPDATE categories
		SET type = $2, name = $3, mandatory = $4
		WHERE id = $1`
	GetCategory = `SELECT id, type, name, mandatory
		FROM categories WHERE id::text = $1 or name = $1`
	GetAllCategories = `SELECT * FROM categories ORDER BY id`
	DeleteCategory   = `DELETE FROM categories WHERE id = $1`
)

const (
	InsertBill = `INSERT INTO bills (id, status) VALUES ($1, $2)`
	UpdateBill = `UPDATE bills SET status = $2 WHERE id = $1`
	GetBill    = `SELECT status FROM bills WHERE id = $1`
	DeleteBill = `DELETE FROM bills WHERE id = $1`
)

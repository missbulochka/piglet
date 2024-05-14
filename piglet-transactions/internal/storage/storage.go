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
	GetCategory = `SELECT id, type, name, mandatory
		FROM categories WHERE id = $1`
)

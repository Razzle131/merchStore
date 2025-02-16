package transactionRepo

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/Razzle131/merchStore/internal/consts"
	"github.com/Razzle131/merchStore/internal/db"
	"github.com/Razzle131/merchStore/internal/model"
	"github.com/Razzle131/merchStore/internal/serverErrors"
)

type TransactionRepoPg struct {
	db *db.DB
}

func NewTransactionRepoPg(db *db.DB) *TransactionRepoPg {
	return &TransactionRepoPg{
		db: db,
	}
}

// mode 1 -> income transactions
//
// mode 2 -> sent transactions
func (r *TransactionRepoPg) GetTransactions(login string, mode int) ([]model.Transaction, error) {
	field := ""
	if mode == 1 {
		field = "user_to"
	} else {
		field = "user_from"
	}

	builder := sq.Select("user_from", "user_to", "amount").From("transaction").Where(sq.Eq{field: login})

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, serverErrors.ErrInternal
	}

	transactions := make([]model.Transaction, 0, consts.MinSliceCap)
	rows, err := r.db.Pool.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, serverErrors.ErrInternal
	}

	for rows.Next() {
		var t model.Transaction
		err := rows.Scan(
			&t.LoginFrom,
			&t.LoginTo,
			&t.Amount,
		)
		if err != nil {
			return nil, serverErrors.ErrInternal
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (r *TransactionRepoPg) AddTransaction(loginFrom, loginTo string, amount int) error {
	sql, args, err := sq.Insert("transaction").Values(loginFrom, loginTo, amount).
		Columns("user_from", "user_to", "amount").ToSql()
	if err != nil {
		return serverErrors.ErrInternal
	}

	_, err = r.db.Pool.Exec(context.Background(), sql, args...)
	if err != nil {
		return serverErrors.ErrInternal
	}

	return nil
}

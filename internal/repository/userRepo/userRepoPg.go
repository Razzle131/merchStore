package userRepo

import (
	"context"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/Razzle131/merchStore/internal/consts"
	"github.com/Razzle131/merchStore/internal/db"
	"github.com/Razzle131/merchStore/internal/model"
	"github.com/Razzle131/merchStore/internal/serverErrors"
)

type UserRepoPg struct {
	db *db.DB
}

func NewUserRepoPg(db *db.DB) *UserRepoPg {
	return &UserRepoPg{
		db: db,
	}
}

func (r *UserRepoPg) AddUser(login, password string) error {
	tx, err := r.db.Pool.Begin(context.Background())
	if err != nil {
		slog.Error("failed to start db transaction")
		return serverErrors.ErrInternal
	}

	sql, args, err := sq.Insert("users").Values(login, password).Columns("login", "password").ToSql()
	if err != nil {
		return serverErrors.ErrInternal
	}

	_, err = r.db.Pool.Exec(context.Background(), sql, args...)
	if err != nil {
		tx.Rollback(context.Background())
		return serverErrors.ErrInternal
	}

	sql, args, err = sq.Insert("wallet").Values(1000).Columns("cash").ToSql()
	if err != nil {
		tx.Rollback(context.Background())
		return serverErrors.ErrInternal
	}

	_, err = r.db.Pool.Exec(context.Background(), sql, args...)
	if err != nil {
		tx.Rollback(context.Background())
		return serverErrors.ErrInternal
	}

	sql, args, err = sq.Insert("inventory").Values().Columns("cash").ToSql()
	if err != nil {
		tx.Rollback(context.Background())
		return serverErrors.ErrInternal
	}

	_, err = r.db.Pool.Exec(context.Background(), sql, args...)
	if err != nil {
		tx.Rollback(context.Background())
		return serverErrors.ErrInternal
	}
	return nil
}

func (r *UserRepoPg) getUserWalletByLogin(login string) (*model.WalletInfo, error) {
	builder := sq.Select("user_from", "user_to", "amount").From("transaction")

	inTransactions := builder.Where(sq.Eq{"user_to": login})
	outTransactions := builder.Where(sq.Eq{"user_from": login})

	inTransactionsRes := make([]model.Transaction, 0, consts.MinSliceCap)
	outTransactionsRes := make([]model.Transaction, 0, consts.MinSliceCap)

	inSql, inArgs, err := inTransactions.ToSql()
	if err != nil {
		return nil, serverErrors.ErrInternal
	}

	outSql, outArgs, err := outTransactions.ToSql()
	if err != nil {
		return nil, serverErrors.ErrInternal
	}

	slog.Debug("get in wallet sql: " + inSql)
	slog.Debug("get out wallet sql: " + outSql)

	inRows, err := r.db.Pool.Query(context.Background(), inSql, inArgs...)
	for inRows.Next() {
		var t model.Transaction
		err := inRows.Scan(
			&t.LoginFrom,
			&t.LoginTo,
			&t.Amount,
		)
		if err != nil {
			return nil, serverErrors.ErrInternal
		}
		inTransactionsRes = append(inTransactionsRes, t)
	}

	outRows, err := r.db.Pool.Query(context.Background(), outSql, outArgs...)
	for outRows.Next() {
		var t model.Transaction
		err := outRows.Scan(
			&t.LoginFrom,
			&t.LoginTo,
			&t.Amount,
		)
		if err != nil {
			return nil, serverErrors.ErrInternal
		}
		outTransactionsRes = append(outTransactionsRes, t)
	}

	return &model.WalletInfo{
		In:   inTransactionsRes,
		Out:  outTransactionsRes,
		Cash: 0,
	}, nil
}

func (r *UserRepoPg) getInventoryById(id int) (*model.Inventory, error) {
	builder := sq.Select("item", "quantity").From("inventory").Where(sq.Eq{"id": id})

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, serverErrors.ErrInternal
	}

	rows, err := r.db.Pool.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, serverErrors.ErrInternal
	}

	res := make(map[string]int, consts.MinSliceCap)
	for rows.Next() {
		var item string
		var quantity int

		err := rows.Scan(
			&item,
			&quantity,
		)
		if err != nil {
			return nil, serverErrors.ErrInternal
		}

		res[item] = quantity
	}

	return &model.Inventory{
		Items: res,
	}, nil
}

func (r *UserRepoPg) GetUserById(id int) (model.User, error) {
	builder := sq.Select("u.id", "u.login", "u.password", "w.cash")
	builder.From("users").Join("wallet USING (id)")

	sql, args, err := builder.ToSql()
	if err != nil {
		return model.User{}, serverErrors.ErrInternal
	}

	rows, err := r.db.Pool.Query(context.Background(), sql, args...)
	if err != nil {
		return model.User{}, serverErrors.ErrInternal
	}

	if !rows.Next() {
		return model.User{}, serverErrors.ErrUserNotFound
	}

	res := model.User{}
	var cash int
	err = rows.Scan(
		&res.Id,
		&res.Login,
		&res.Password,
		&cash,
	)
	if err != nil {
		return model.User{}, serverErrors.ErrInternal
	}

	wal, err := r.getUserWalletByLogin(res.Login)
	if err != nil {
		return model.User{}, serverErrors.ErrInternal
	}

	items, err := r.getInventoryById(res.Id)
	if err != nil {
		return model.User{}, serverErrors.ErrInternal
	}

	return model.User{
		Id:       res.Id,
		Login:    res.Login,
		Password: res.Password,
		Wallet:   *wal,
		Items:    *items,
	}, nil
}

func (r *UserRepoPg) GetUserByLogin(login string) (model.User, error) {
	builder := sq.Select("id").From("users").Where(sq.Eq{"login": login})

	sql, args, err := builder.ToSql()
	if err != nil {
		return model.User{}, serverErrors.ErrInternal
	}

	rows, err := r.db.Pool.Query(context.Background(), sql, args...)
	if err != nil {
		return model.User{}, serverErrors.ErrInternal
	}

	if !rows.Next() {
		return model.User{}, serverErrors.ErrUserNotFound
	}

	var id int
	err = rows.Scan(&id)
	if err != nil {
		return model.User{}, serverErrors.ErrInternal
	}

	return r.GetUserById(id)
}

func (r *UserRepoPg) UpdateUser(newUser model.User) error {
	if r.users[newUser.Id].Id != newUser.Id {
		return serverErrors.ErrNotAllowed
	}

	r.users[newUser.Id] = newUser
	return nil
}

func (r *UserRepoPg) GetNewId() int {
	return len(r.users)
}

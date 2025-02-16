package merchRepo

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/Razzle131/merchStore/internal/db"
	"github.com/Razzle131/merchStore/internal/serverErrors"
)

type MerchRepoPg struct {
	db *db.DB
}

func NewMerchRepoPg(db *db.DB) *MerchRepoPg {
	return &MerchRepoPg{
		db: db,
	}
}

func (r *MerchRepoPg) GetMerchPrice(item string) (int, error) {
	sql, args, err := sq.Select("price").From("items").Where(sq.Eq{"item": item}).ToSql()
	if err != nil {
		return 0, serverErrors.ErrInternal
	}

	rows, err := r.db.Pool.Query(context.Background(), sql, args...)
	if err != nil {
		return 0, serverErrors.ErrInternal
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, serverErrors.ErrItemNotFound
	}

	var price int
	err = rows.Scan(&price)
	if err != nil {
		return 0, serverErrors.ErrInternal
	}

	return price, nil
}

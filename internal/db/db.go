package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, conn string) (*DB, error) {
	pool, err := pgxpool.Connect(ctx, conn)
	if err != nil {
		return nil, err
	}

	return &DB{
		Pool: pool,
	}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}

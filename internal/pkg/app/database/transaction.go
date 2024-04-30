package database

import (
	"context"
	"database/sql"
)

type contextKey int

const txKey contextKey = iota

func BeginTx(ctx context.Context, db *sql.DB) (context.Context, *sql.Tx, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, nil, err
	}

	txCtx := context.WithValue(ctx, txKey, tx)

	return txCtx, tx, nil
}

func TxOrDB(ctx context.Context, db Connection) Connection {
	tx, ok := ctx.Value(txKey).(*sql.Tx)
	if ok {
		return tx
	}

	return db
}

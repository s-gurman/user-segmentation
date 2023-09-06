package segmentrepo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxTx struct {
	pgx.Tx
}

func newPgxTx(ctx context.Context, db *pgxpool.Pool, opts pgx.TxOptions) (pgxTx, error) {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return pgxTx{}, fmt.Errorf("segmentrepo - tx begin err: %w", err)
	}
	return pgxTx{Tx: tx}, nil
}

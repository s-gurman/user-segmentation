package segmentrepo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// mockgen -destination ./mocks/pgx_tx.go -package mocks github.com/jackc/pgx/v5 Tx
// mockgen -source ./pgx_pool.go -destination ./mocks/pgx_pool.go -package mocks

type PgxPool interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

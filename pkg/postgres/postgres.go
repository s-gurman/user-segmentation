package postgres

import (
	"context"
	"fmt"

	"github.com/s-gurman/user-segmentation/config"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Postgres interface {
	GetPool() *pgxpool.Pool
	Close()
}

type postgres struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, cfg config.PGConfig) (Postgres, error) {
	pgURL := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Address, cfg.DBName, cfg.SSLMode,
	)
	db, err := pgxpool.New(ctx, pgURL)
	if err != nil {
		return nil, fmt.Errorf("pg - open err: %w", err)
	}
	if err = db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pg - ping err: %w", err)
	}
	return postgres{db: db}, nil
}

func (pg postgres) GetPool() *pgxpool.Pool {
	return pg.db
}

func (pg postgres) Close() {
	pg.db.Close()
}

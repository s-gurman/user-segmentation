package segmentrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/s-gurman/user-segmentation/internal/domain"
	"github.com/s-gurman/user-segmentation/internal/e"
	"github.com/s-gurman/user-segmentation/pkg/postgres"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (tx pgxTx) insertSegment(ctx context.Context, slug domain.Slug) (int, error) {
	query := `INSERT INTO segments (slug) VALUES ($1) RETURNING id`

	var id int
	if err := tx.QueryRow(ctx, query, slug).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			msg := "attempt to create already created segment"
			return 0, e.NewBadRequest(msg, "segmentrepo" /*from*/)
		}
		return 0, fmt.Errorf("segmentrepo - tx insert segment err: %w", err)
	}

	return id, nil
}

func (tx pgxTx) initSegmentByRandomUsers(
	ctx context.Context,
	segID int,
	autoaddPercent float32,
) (int64, error) {

	query := `
INSERT INTO experiments (user_id, segment_id)
SELECT id, $1 FROM users TABLESAMPLE BERNOULLI ($2)`

	tag, err := tx.Exec(ctx, query, segID, autoaddPercent)
	if err != nil {
		return 0, fmt.Errorf("segmentrepo - tx insert segment users err: %w", err)
	}

	return tag.RowsAffected(), nil
}

type SegmentRepo struct {
	db *pgxpool.Pool
}

func NewSegmentRepository(pg postgres.Postgres) SegmentRepo {
	return SegmentRepo{db: pg.GetPool()}
}

func (repo SegmentRepo) CreateSegment(
	ctx context.Context,
	slug domain.Slug,
	autoaddPercent float32,
) (int, int64, error) {

	txOpts := pgx.TxOptions{IsoLevel: pgx.RepeatableRead}
	tx, err := newPgxTx(ctx, repo.db, txOpts)
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback(ctx) // nolint:errcheck

	segID, err := tx.insertSegment(ctx, slug)
	if err != nil {
		return 0, 0, err
	}
	if autoaddPercent == 0 {
		if err = tx.Commit(ctx); err != nil {
			return 0, 0, fmt.Errorf("segmentrepo - tx commit err: %w", err)
		}
		return segID, 0, nil
	}

	insertedCount, err := tx.initSegmentByRandomUsers(ctx, segID, autoaddPercent)
	if err != nil {
		return 0, 0, err
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, 0, fmt.Errorf("segmentrepo - tx commit err: %w", err)
	}

	return segID, insertedCount, nil
}

func (repo SegmentRepo) DeleteSegment(
	ctx context.Context,
	slug domain.Slug,
) error {

	query := `DELETE FROM segments WHERE slug = $1`

	tag, err := repo.db.Exec(ctx, query, slug)
	if err != nil {
		return fmt.Errorf("segmentrepo - segment delete exec err: %w", err)
	}
	if tag.RowsAffected() != 1 {
		msg := "attempt to delete unknown segment"
		return e.NewNotFound(msg, "segmentrepo" /*from*/)
	}

	return nil
}

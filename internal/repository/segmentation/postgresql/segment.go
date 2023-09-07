package segmentrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/s-gurman/user-segmentation/internal/domain"
	"github.com/s-gurman/user-segmentation/internal/e"
	"github.com/s-gurman/user-segmentation/pkg/postgres"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	segSlugAttr = "slug"
	segIDAttr   = "id"
	segTable    = "segments"
	usersIDAttr = "id"
	usersTable  = "users"
)

func (tx pgxTx) createSegment(ctx context.Context, slug domain.Slug) (int, error) {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES ($1) RETURNING %s",
		segTable, segSlugAttr, segIDAttr,
	)

	var id int
	if err := tx.QueryRow(ctx, query, slug).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			msg := "attempt to create existing segment"
			return 0, e.NewBadRequest(msg, "segmentrepo" /*from*/)
		}
		return 0, fmt.Errorf("segmentrepo - tx id scan err: %w", err)
	}

	return id, nil
}

func (tx pgxTx) getUserBatch(ctx context.Context, autoaddPercent float32) ([]int, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s TABLESAMPLE BERNOULLI ($1)",
		usersIDAttr, usersTable,
	)

	rows, err := tx.Query(ctx, query, autoaddPercent)
	if err != nil {
		return nil, fmt.Errorf("segmentrepo - tx user batch query err: %w", err)
	}

	var userIDs []int
	if err := pgxscan.ScanAll(&userIDs, rows); err != nil {
		return nil, fmt.Errorf("segmentrepo - tx user batch scan err: %w", err)
	}

	return userIDs, nil
}

func (tx pgxTx) addSegmentToUsers(ctx context.Context, userIDs []int, segID int) error {
	experiments := make([][]any, len(userIDs))
	for i, userID := range userIDs {
		experiments[i] = []any{userID, segID}
	}

	affected, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{expTable},
		[]string{expUserAttr, expSegAttr},
		pgx.CopyFromRows(experiments),
	)
	if err != nil {
		return fmt.Errorf("segmentrepo - tx copy experiments err: %w", err)
	}
	if affected != int64(len(userIDs)) {
		return errors.New("segmentrepo err: autoadd segment to users failed")
	}

	return nil
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
) (int, int, error) {

	txOpts := pgx.TxOptions{IsoLevel: pgx.RepeatableRead}
	tx, err := newPgxTx(ctx, repo.db, txOpts)
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback(ctx) // nolint:errcheck

	segID, err := tx.createSegment(ctx, slug)
	if err != nil {
		return 0, 0, err
	}
	if autoaddPercent == 0 {
		if err = tx.Commit(ctx); err != nil {
			return 0, 0, fmt.Errorf("segmentrepo - tx commit err: %w", err)
		}
		return segID, 0, nil
	}

	userIDs, err := tx.getUserBatch(ctx, autoaddPercent)
	if err != nil {
		return 0, 0, err
	}
	if err = tx.addSegmentToUsers(ctx, userIDs, segID); err != nil {
		return 0, 0, err
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, 0, fmt.Errorf("segmentrepo - tx commit err: %w", err)
	}

	return segID, len(userIDs), nil
}

func (repo SegmentRepo) DeleteSegment(
	ctx context.Context,
	slug domain.Slug,
) error {

	query := fmt.Sprintf(
		"DELETE FROM %s WHERE %s = $1",
		segTable, segSlugAttr,
	)

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

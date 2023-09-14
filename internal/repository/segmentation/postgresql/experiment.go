package segmentrepo

import (
	"context"
	"fmt"
	"strings"

	"github.com/s-gurman/user-segmentation/internal/domain"
	"github.com/s-gurman/user-segmentation/internal/e"
	"github.com/s-gurman/user-segmentation/internal/t"
	"github.com/s-gurman/user-segmentation/pkg/postgres"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type ExperimentRepo struct {
	db PgxPool
}

func NewExperimentRepository(pg postgres.Postgres) ExperimentRepo {
	return ExperimentRepo{db: pg.GetPool()}
}

func getSegmentIDs(ctx context.Context, tx pgx.Tx, slugs []domain.Slug) ([]int, error) {
	query := `SELECT id FROM segments WHERE slug = ANY($1)`

	rows, err := tx.Query(ctx, query, slugs)
	if err != nil {
		return nil, fmt.Errorf("segmentrepo - tx slugs query err: %w", err)
	}

	var segIDs []int
	if err := pgxscan.ScanAll(&segIDs, rows); err != nil {
		return nil, fmt.Errorf("segmentrepo - tx ids scan err: %w", err)
	}
	if len(segIDs) != len(slugs) {
		msg := "attempt to work with unknown segment, create it first"
		return nil, e.NewNotFound(msg, "segmentrepo" /*from*/)
	}

	return segIDs, nil
}

func softDeleteUserSegments(
	ctx context.Context,
	tx pgx.Tx,
	userID int,
	segIDs []int,
) error {

	query := `
UPDATE experiments SET expired_at = NOW()::timestamp(0)
FROM segments
WHERE segments.id = ANY($2)
	AND user_id = $1 AND segment_id = segments.id
	AND (expired_at IS NULL OR expired_at > NOW()::timestamp(0))`

	tag, err := tx.Exec(ctx, query, userID, segIDs)
	if err != nil {
		return fmt.Errorf("segmentrepo - tx user segments delete err: %w", err)
	}
	if tag.RowsAffected() != int64(len(segIDs)) {
		msg := "attempt to delete user's inactive segment"
		return e.NewBadRequest(msg, "segmentrepo" /*from*/)
	}

	return nil
}

func addUserSegments(
	ctx context.Context,
	tx pgx.Tx,
	userID int,
	segIDs []int,
	expired *t.CustomTime,
) error {

	query := `
INSERT INTO experiments (user_id, segment_id, expired_at)
SELECT $1, id, NULL FROM segments WHERE id = ANY($2)
ON CONFLICT ON CONSTRAINT experiments_user_segment_unique
	DO UPDATE SET started_at = NOW()::timestamp(0), expired_at = NULL
WHERE experiments.expired_at IS NOT NULL
	AND experiments.expired_at <= NOW()::timestamp(0)`

	args := []any{userID, segIDs}
	if expired != nil {
		query = strings.Replace(query, "NULL", "$3", 2)
		args = append(args, expired.Time)
	}

	tag, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("segmentrepo - tx user segments add err: %w", err)
	}
	if tag.RowsAffected() != int64(len(segIDs)) {
		msg := "attempt to add already active segment"
		return e.NewBadRequest(msg, "segmentrepo" /*from*/)
	}

	return nil
}

func (repo ExperimentRepo) UpdateUserSegments(
	ctx context.Context,
	userID int,
	toDel, toAdd []domain.Slug,
	expired *t.CustomTime,
) error {

	tx, err := repo.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		return fmt.Errorf("segmentrepo - tx begin err: %w", err)
	}
	defer tx.Rollback(ctx) // nolint:errcheck

	idsToDel, err := getSegmentIDs(ctx, tx, toDel)
	if err != nil {
		return err
	}
	idsToAdd, err := getSegmentIDs(ctx, tx, toAdd)
	if err != nil {
		return err
	}

	if len(idsToDel) > 0 {
		if err = softDeleteUserSegments(ctx, tx, userID, idsToDel); err != nil {
			return err
		}
	}
	if len(idsToAdd) > 0 {
		if err = addUserSegments(ctx, tx, userID, idsToAdd, expired); err != nil {
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("segmentrepo - tx commit err: %w", err)
	}

	return nil
}

func (repo ExperimentRepo) GetUserSegments(
	ctx context.Context,
	userID int,
) ([]string, error) {

	query := `
SELECT segments.slug FROM experiments
JOIN segments ON segment_id = segments.id
WHERE user_id = $1
	AND (expired_at IS NULL OR expired_at > NOW()::timestamp(0))`

	rows, err := repo.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("segmentrepo - user segments query err: %w", err)
	}

	slugs := make([]string, 0)
	if err := pgxscan.ScanAll(&slugs, rows); err != nil {
		return nil, fmt.Errorf("segmentrepo - user segments scan err: %w", err)
	}

	return slugs, nil
}

package segmentrepo

import (
	"context"
	"fmt"

	"github.com/s-gurman/user-segmentation/internal/domain"
	"github.com/s-gurman/user-segmentation/internal/e"
	"github.com/s-gurman/user-segmentation/pkg/postgres"
	"golang.org/x/sync/errgroup"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	expExpiredAttr = "expired_at"
	expStartedAttr = "started_at"
	expSegmentAttr = "segment_id"
	expUserAttr    = "user_id"
	expTable       = "experiments"
)

type ExperimentRepo struct {
	db *pgxpool.Pool
}

func NewExperimentRepository(pg postgres.Postgres) ExperimentRepo {
	return ExperimentRepo{db: pg.GetPool()}
}

type pgxTx struct {
	pgx.Tx
}

func (tx pgxTx) getSegmentsID(ctx context.Context, slugs []domain.Slug) ([]int32, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s=ANY($1)",
		segmIDAttr, segmTable, segmSlugAttr,
	)

	rows, err := tx.Query(ctx, query, slugs)
	if err != nil {
		return nil, fmt.Errorf("segmentrepo - tx segments query err: %w", err)
	}

	var result []int32
	if err := pgxscan.ScanAll(&result, rows); err != nil {
		return nil, fmt.Errorf("segmentrepo - tx rows scan err: %w", err)
	}
	if len(result) != len(slugs) {
		return nil, e.NewNotFound("attempt to work with unknown segments" /*from*/, "segmentrepo" /*from*/)
	}

	return result, nil
}

func (tx pgxTx) softDeleteUserSegments(ctx context.Context, userID int, segmentsID []int32) error {
	query := fmt.Sprintf(
		"UPDATE %s SET %s=NOW() WHERE %s=$1 AND %s=ANY($2) AND (%s IS NULL OR %s > NOW())",
		expTable, expExpiredAttr, expUserAttr, expSegmentAttr, expExpiredAttr, expExpiredAttr,
	)

	tag, err := tx.Exec(ctx, query, userID, segmentsID)
	if err != nil {
		return fmt.Errorf("segmentrepo - user segments exec err: %w", err)
	}
	if tag.RowsAffected() != int64(len(segmentsID)) {
		return e.NewBadRequest("attempt to delete inactive user segments" /*msg*/, "segmentrepo" /*from*/)
	}

	return nil
}

func (repo ExperimentRepo) UpdateUserSegments(ctx context.Context, userID int, toDel, toAdd []domain.Slug) error {
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("segmentrepo - tx begin err: %w", err)
	}
	defer tx.Rollback(ctx) // nolint:errcheck

	var idsToDel, idsToAdd []int32
	pgxtx := pgxTx{tx}

	g, groupCtx := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		idsToDel, err = pgxtx.getSegmentsID(groupCtx, toDel)
		return
	})
	g.Go(func() (err error) {
		idsToAdd, err = pgxtx.getSegmentsID(groupCtx, toAdd)
		return
	})
	if err = g.Wait(); err != nil {
		return err
	}

	err = pgxtx.softDeleteUserSegments(ctx, userID, idsToDel)
	if err != nil {
		return err
	}

	_ = idsToAdd

	return nil
}
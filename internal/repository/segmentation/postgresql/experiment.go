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
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	expUniqueConstr = "experiments_user_segment_unique"
	expStartedAttr  = "started_at"
	expExpiredAttr  = "expired_at"
	expUserAttr     = "user_id"
	expSegAttr      = "segment_id"
	expTable        = "experiments"
)

func (tx pgxTx) getSegmentIDs(ctx context.Context, slugs []domain.Slug) ([]int, error) {
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s = ANY($1)",
		segIDAttr, segTable, segSlugAttr,
	)

	rows, err := tx.Query(ctx, query, slugs)
	if err != nil {
		return nil, fmt.Errorf("segmentrepo - tx slugs query err: %w", err)
	}

	var segIDs []int
	if err := pgxscan.ScanAll(&segIDs, rows); err != nil {
		return nil, fmt.Errorf("segmentrepo - tx ids scan err: %w", err)
	}
	if len(segIDs) != len(slugs) {
		msg := "attempt to work with unknown segments, create them first"
		return nil, e.NewNotFound(msg, "segmentrepo" /*from*/)
	}

	return segIDs, nil
}

func (tx pgxTx) softDeleteUserSegments(
	ctx context.Context,
	userID int, segIDs []int,
) error {

	query := fmt.Sprintf(
		`UPDATE %s SET %s = NOW() WHERE %s = $1
AND %s = ANY($2) AND (%s IS NULL OR %s > NOW())`,
		expTable, expExpiredAttr, expUserAttr,
		expSegAttr, expExpiredAttr, expExpiredAttr,
	)

	tag, err := tx.Exec(ctx, query, userID, segIDs)
	if err != nil {
		return fmt.Errorf("segmentrepo - tx user segments delete err: %w", err)
	}
	if tag.RowsAffected() != int64(len(segIDs)) {
		msg := "attempt to delete user's inactive segments"
		return e.NewBadRequest(msg, "segmentrepo" /*from*/)
	}

	return nil
}

func (tx pgxTx) addUserSegments(
	ctx context.Context,
	userID int, segIDs []int,
	expired *t.CustomTime,
) error {

	query := fmt.Sprintf(
		`INSERT INTO %s (%s, %s, %s, %s) VALUES ($1, $2, NOW(), DEFAULT)
ON CONFLICT ON CONSTRAINT %s DO UPDATE SET %s = EXCLUDED.%s, %s = DEFAULT
WHERE %s.%s IS NOT NULL AND %s.%s <= EXCLUDED.%s`,
		expTable, expUserAttr, expSegAttr, expStartedAttr, expExpiredAttr,
		expUniqueConstr, expStartedAttr, expStartedAttr, expExpiredAttr,
		expTable, expExpiredAttr, expTable, expExpiredAttr, expStartedAttr,
	)
	if expired != nil {
		query = strings.ReplaceAll(query, "DEFAULT", "$3")
	}

	for _, segID := range segIDs {
		args := []any{userID, segID}
		if expired != nil {
			args = append(args, expired.Time)
		}

		tag, err := tx.Exec(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("segmentrepo - tx user segments add err: %w", err)
		}
		if tag.RowsAffected() != 1 {
			msg := "attempt to add already active segment"
			return e.NewBadRequest(msg, "segmentrepo" /*from*/)
		}
	}

	return nil
}

type ExperimentRepo struct {
	db *pgxpool.Pool
}

func NewExperimentRepository(pg postgres.Postgres) ExperimentRepo {
	return ExperimentRepo{db: pg.GetPool()}
}

func (repo ExperimentRepo) UpdateUserSegments(
	ctx context.Context,
	userID int,
	toDel, toAdd []domain.Slug,
	expired *t.CustomTime,
) error {

	txOpts := pgx.TxOptions{IsoLevel: pgx.RepeatableRead}
	tx, err := newPgxTx(ctx, repo.db, txOpts)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // nolint:errcheck

	idsToDel, err := tx.getSegmentIDs(ctx, toDel)
	if err != nil {
		return err
	}
	idsToAdd, err := tx.getSegmentIDs(ctx, toAdd)
	if err != nil {
		return err
	}

	if err = tx.softDeleteUserSegments(ctx, userID, idsToDel); err != nil {
		return err
	}
	if err = tx.addUserSegments(ctx, userID, idsToAdd, expired); err != nil {
		return err
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

	query := fmt.Sprintf(
		`SELECT %s.%s FROM %s JOIN %s ON %s = %s.%s
WHERE %s = $1 AND (%s IS NULL OR %s > NOW())`,
		segTable, segSlugAttr, expTable, segTable, expSegAttr, segTable, segIDAttr,
		expUserAttr, expExpiredAttr, expExpiredAttr,
	)

	rows, err := repo.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("segmentrepo - user segments query err: %w", err)
	}

	var slugs []string
	if err := pgxscan.ScanAll(&slugs, rows); err != nil {
		return nil, fmt.Errorf("segmentrepo - tx slugs scan err: %w", err)
	}

	return slugs, nil
}

package segmentrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/s-gurman/user-segmentation/internal/domain"
	"github.com/s-gurman/user-segmentation/internal/e"
	"github.com/s-gurman/user-segmentation/pkg/postgres"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	segSlugAttr = "slug"
	segIDAttr   = "id"
	segTable    = "segments"
)

type SegmentRepo struct {
	db *pgxpool.Pool
}

func NewSegmentRepository(pg postgres.Postgres) SegmentRepo {
	return SegmentRepo{db: pg.GetPool()}
}

func (repo SegmentRepo) CreateSegment(ctx context.Context, slug domain.Slug) (int, error) {
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES ($1) RETURNING %s",
		segTable, segSlugAttr, segIDAttr,
	)

	var id int
	if err := repo.db.QueryRow(ctx, query, slug).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			msg := "attempt to create existing segment"
			return 0, e.NewBadRequest(msg, "segmentrepo" /*from*/)
		}
		return 0, fmt.Errorf("segmentrepo - segment row scan err: %w", err)
	}

	return id, nil
}

func (repo SegmentRepo) DeleteSegment(ctx context.Context, slug domain.Slug) error {
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

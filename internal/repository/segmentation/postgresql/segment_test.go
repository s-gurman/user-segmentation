package segmentrepo

import (
	"context"
	"errors"
	"testing"

	"github.com/s-gurman/user-segmentation/internal/domain"
	"github.com/s-gurman/user-segmentation/internal/repository/segmentation/postgresql/mocks"

	"github.com/pashagolub/pgxmock/v2"
	"go.uber.org/mock/gomock"
)

func TestDeleteSegment_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	query := `DELETE FROM segments WHERE slug = $1`
	slug := domain.Slug("some slug")

	var expected error

	pool := mocks.NewMockPgxPool(ctrl)
	pool.
		EXPECT().
		Exec(ctx, query, slug).
		Return(pgxmock.NewResult("DELETE", 1), expected)

	repo := SegmentRepo{db: pool}

	err := repo.DeleteSegment(ctx, slug)

	if !errors.Is(err, expected) {
		t.Errorf("unexpected error: %s", err)
	}
}

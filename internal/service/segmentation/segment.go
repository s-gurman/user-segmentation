package segmentsvc

import (
	"context"
	"fmt"

	"github.com/s-gurman/user-segmentation/internal/domain"
	"github.com/s-gurman/user-segmentation/internal/e"
)

func (svc SegmentationSvc) CreateSegment(ctx context.Context, name string) (int, error) {
	slug, err := domain.NewSlug(name)
	if err != nil {
		return 0, e.NewBadRequest(err.Error(), "segmentsvc" /*from*/)
	}

	id, err := svc.segrepo.InsertOne(ctx, slug)
	if err != nil {
		return 0, fmt.Errorf("segmentsvc - insert one segment: %w", err)
	}

	return id, nil
}

func (svc SegmentationSvc) DeleteSegment(ctx context.Context, name string) error {
	slug, err := domain.NewSlug(name)
	if err != nil {
		return e.NewBadRequest(err.Error(), "segmentsvc" /*from*/)
	}
	if err := svc.segrepo.DeleteOne(ctx, slug); err != nil {
		return fmt.Errorf("segmentsvc - delete one segment: %w", err)
	}
	return nil
}

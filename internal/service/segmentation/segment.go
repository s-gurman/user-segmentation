package segmentsvc

import (
	"context"
	"fmt"

	"github.com/s-gurman/user-segmentation/internal/domain"
	"github.com/s-gurman/user-segmentation/internal/e"
)

func (svc SegmentationSvc) CreateSegment(
	ctx context.Context,
	name string,
	autoaddPercent float32,
) (int, int, error) {

	slug, err := domain.NewSlug(name)
	if err != nil {
		return 0, 0, e.NewBadRequest(err.Error(), "segmentsvc" /*from*/)
	}
	if autoaddPercent < 0 || autoaddPercent > 100 {
		msg := "autoadd_percent option should be in range [0, 100]"
		return 0, 0, e.NewBadRequest(msg, "segmentsvc")
	}

	id, autoaddCount, err := svc.segrepo.CreateSegment(ctx, slug, autoaddPercent)
	if err != nil {
		return 0, 0, fmt.Errorf("segmentsvc - create segment: %w", err)
	}

	return id, autoaddCount, nil
}

func (svc SegmentationSvc) DeleteSegment(ctx context.Context, name string) error {
	slug, err := domain.NewSlug(name)
	if err != nil {
		return e.NewBadRequest(err.Error(), "segmentsvc" /*from*/)
	}

	if err := svc.segrepo.DeleteSegment(ctx, slug); err != nil {
		return fmt.Errorf("segmentsvc - delete segment: %w", err)
	}

	return nil
}

package segmentsvc

import (
	"context"

	"github.com/s-gurman/user-segmentation/internal/domain"
	"github.com/s-gurman/user-segmentation/internal/t"
)

type (
	ExperimentStorage interface {
		UpdateUserSegments(_ context.Context, userID int, toDel, toAdd []domain.Slug, delTime *t.CustomTime) error
		GetUserSegments(_ context.Context, userID int) ([]string, error)
	}
	SegmentStorage interface {
		CreateSegment(_ context.Context, slug domain.Slug, autoaddPercent float32) (int, int, error)
		DeleteSegment(_ context.Context, slug domain.Slug) error
	}
)

type SegmentationSvc struct {
	segrepo SegmentStorage
	exprepo ExperimentStorage
}

func New(segment SegmentStorage, experiment ExperimentStorage) SegmentationSvc {
	return SegmentationSvc{
		segrepo: segment,
		exprepo: experiment,
	}
}

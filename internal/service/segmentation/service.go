package segmentsvc

import (
	"context"

	"github.com/s-gurman/user-segmentation/internal/domain"
)

type (
	ExperimentStorage interface {
		UpdateUserSegments(ctx context.Context, userID int, toDel, toAdd []domain.Slug) error
	}

	SegmentStorage interface {
		InsertOne(ctx context.Context, slug domain.Slug) (int, error)
		DeleteOne(ctx context.Context, slug domain.Slug) error
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

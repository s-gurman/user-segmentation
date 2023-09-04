package segmentsvc

import (
	"context"
	"fmt"

	"github.com/s-gurman/user-segmentation/internal/domain"
	"github.com/s-gurman/user-segmentation/internal/e"

	"golang.org/x/sync/errgroup"
)

func GetSegmentsSlug(toDel, toAdd []string) ([]domain.Slug, []domain.Slug, error) {
	if len(toDel) == 0 && len(toAdd) == 0 {
		return nil, nil, e.NewBadRequest("empty experiment update lists" /*msg*/, "segmentsvc" /*from*/)
	}
	nameToSlug := func(names []string, slugs []domain.Slug) error {
		for i, name := range names {
			slug, err := domain.NewSlug(name)
			if err != nil {
				return err
			}
			slugs[i] = slug
		}
		return nil
	}
	slugsToDel := make([]domain.Slug, len(toDel))
	slugsToAdd := make([]domain.Slug, len(toAdd))

	g := new(errgroup.Group)
	g.Go(func() error {
		return nameToSlug(toDel, slugsToDel)
	})
	g.Go(func() error {
		return nameToSlug(toAdd, slugsToAdd)
	})
	if err := g.Wait(); err != nil {
		return nil, nil, e.NewBadRequest(err.Error(), "segmentsvc" /*from*/)
	}

	return slugsToDel, slugsToAdd, nil
}

func (svc SegmentationSvc) UpdateExperiments(ctx context.Context, userID int, toDel, toAdd []string) error {
	slugsToDel, slugsToAdd, err := GetSegmentsSlug(toDel, toAdd)
	if err != nil {
		return err
	}
	if err = svc.exprepo.UpdateUserSegments(ctx, userID, slugsToDel, slugsToAdd); err != nil {
		return fmt.Errorf("segmentsvc - update user segments: %w", err)
	}
	return nil
}

func (svc SegmentationSvc) GetUserExperiments(ctx context.Context, userID int) error {
	return nil
}

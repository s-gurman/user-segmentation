package segmentsvc

import (
	"context"
	"fmt"
	"time"

	"github.com/s-gurman/user-segmentation/internal/domain"
	"github.com/s-gurman/user-segmentation/internal/e"
	"github.com/s-gurman/user-segmentation/internal/t"

	"golang.org/x/sync/errgroup"
)

const _defaultTimeTruncate = time.Second

func namesToSlugs(names []string, slugs []domain.Slug, list string) error {
	// no need to use sync Map (it will be created and used by only one goroutine)
	namesMap := make(map[string]struct{})

	for i, name := range names {
		if _, found := namesMap[name]; found {
			msg := fmt.Sprintf("%s list contains non-unique segments", list)
			return e.NewBadRequest(msg, "segmentsvc")
		}
		namesMap[name] = struct{}{}

		slug, err := domain.NewSlug(name)
		if err != nil {
			msg := fmt.Sprintf("%s list contains invalid segment name", list)
			return e.NewBadRequest(msg, "segmentsvc")
		}
		slugs[i] = slug
	}

	return nil
}

func getSlugsToUpdate(toDel, toAdd []string) ([]domain.Slug, []domain.Slug, error) {
	if len(toDel) == 0 && len(toAdd) == 0 {
		msg := "empty experiment update lists, must add or delete at least one segment"
		return nil, nil, e.NewBadRequest(msg, "segmentsvc" /*from*/)
	}
	slugsToDel := make([]domain.Slug, len(toDel))
	slugsToAdd := make([]domain.Slug, len(toAdd))

	g := new(errgroup.Group)
	g.Go(func() error {
		return namesToSlugs(toDel, slugsToDel, "delete")
	})
	g.Go(func() error {
		return namesToSlugs(toAdd, slugsToAdd, "add")
	})
	if err := g.Wait(); err != nil {
		return nil, nil, err
	}

	return slugsToDel, slugsToAdd, nil
}

func (svc SegmentationSvc) UpdateExperiments(
	ctx context.Context,
	userID int,
	toDel, toAdd []string,
	delTime *t.CustomTime,
) error {

	if delTime != nil && !delTime.After(time.Now().Truncate(_defaultTimeTruncate)) {
		msg := "segment deletion time should be in the future"
		return e.NewBadRequest(msg, "segmentsvc" /*from*/)
	}

	slugsToDel, slugsToAdd, err := getSlugsToUpdate(toDel, toAdd)
	if err != nil {
		return err
	}
	if err = svc.exprepo.UpdateUserSegments(ctx, userID, slugsToDel, slugsToAdd, delTime); err != nil {
		return fmt.Errorf("segmentsvc - update user segments: %w", err)
	}
	return nil
}

func (svc SegmentationSvc) GetUserExperiments(
	ctx context.Context,
	userID int,
) ([]string, error) {

	slugs, err := svc.exprepo.GetUserSegments(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("segmentsvc - get user segments: %w", err)
	}

	return slugs, nil
}

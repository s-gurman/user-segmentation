package segmentsvc

import (
	"context"
	"errors"
	"testing"

	"github.com/s-gurman/user-segmentation/internal/domain"
	"github.com/s-gurman/user-segmentation/internal/service/segmentation/mocks"

	"go.uber.org/mock/gomock"
)

func TestDeleteSegment_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	name := "some slug"
	slug := domain.Slug(name)

	var expected error

	segrepo := mocks.NewMockSegmentStorage(ctrl)
	segrepo.
		EXPECT().
		DeleteSegment(ctx, slug).
		Return(expected)

	svc := SegmentationSvc{segrepo: segrepo}

	err := svc.DeleteSegment(ctx, name)

	if !errors.Is(err, expected) {
		t.Errorf("unexpected error: %s", err)
	}
}

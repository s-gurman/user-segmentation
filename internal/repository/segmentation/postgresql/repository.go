package segmentrepo

import "github.com/s-gurman/user-segmentation/pkg/postgres"

type Repository struct {
	Experiment ExperimentRepo
	Segment    SegmentRepo
}

func NewPostgreSQL(pg postgres.Postgres) Repository {
	return Repository{
		Experiment: NewExperimentRepository(pg),
		Segment:    NewSegmentRepository(pg),
	}
}

package domain

import "time"

type Experiment struct {
	ID        int
	UserID    int
	SegmentID int
	StartedAt time.Time
	ExpiredAt time.Time
}

package repository

import (
	"context"
	"template/internal/models"
	"time"
)

type Regions interface {
	Create(ctx context.Context, region models.RegionBase) (int, error)
}

type UsersSegments interface {
	Create(ctx context.Context, userSegment models.UserSegmentBase) (int, error)
}

type Users interface {
	Create(ctx context.Context, user models.UserBase) (int, error)
}

type Matrix interface {
	GetHistory(ctx context.Context, timeStart time.Time, timeEnd time.Time, matrixType string)
}

package repository

import (
	"context"
	"template/internal/models"
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

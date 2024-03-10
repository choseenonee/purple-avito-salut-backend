package repository

import (
	"context"
	"github.com/guregu/null"
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

type Microcategories interface {
	Create(ctx context.Context, microcategory models.MicrocategoryBase) (int, error)
}

type Matrix interface {
	CreateMatrix(ctx context.Context, matrix models.MatrixBase) error
	GetHistory(ctx context.Context, timeStart time.Time, timeEnd time.Time, isBaseline null.Bool) ([]models.Matrix, error)
}

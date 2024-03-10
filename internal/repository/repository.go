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

type Microcategories interface {
	Create(ctx context.Context, microсategory models.MicrocategoryBase) (int, error)
}

type Matrix interface {
	CreateMatrix(ctx context.Context, matrix models.MatrixBase) (string, error)
	GetMatrix(ctx context.Context, matrixName string, page int) (models.Matrix, error)
	GetHistory(ctx context.Context, matrix models.GetHistoryMatrix) ([]models.Matrix, error)
	GetDifference(ctx context.Context, matrixName1, matrixName2 string) (models.MatrixDifference, error)
}

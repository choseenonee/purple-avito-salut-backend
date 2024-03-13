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

type Microcategories interface {
	Create(ctx context.Context, microсategory models.MicrocategoryBase) (int, error)
}

type Matrix interface {
	CreateMatrix(ctx context.Context, matrix models.MatrixDifferenceRequest) (string, error)
	CreateMatrixWithoutParent(ctx context.Context, matrix models.MatrixBase) (string, error)
	GetMatrix(ctx context.Context, matrixName string, page int) (models.Matrix, error)
	GetMatrixPages(ctx context.Context, matrixName string) (int, error)
	GetMatricesByDuration(ctx context.Context, timeStart, timeEnd time.Time) ([]models.Matrix, error)
	GetHistory(ctx context.Context, matrix models.GetHistoryMatrix) ([]models.ResponseHistoryMatrix, error)
	GetPriceTendency(ctx context.Context, data models.GetTendencyNode) ([]models.ResponseTendencyNode, error)
	GetDifference(ctx context.Context, matrixName1, matrixName2 string) (models.MatrixDifferenceResponse, error)
	GetRelationsWithPrice(ctx context.Context, matrixName string) ([][4]int, [][4]int, error)
}

type Nodes interface {
	GetMicrocategoryTree() ([]models.NodeRaw, error)
}

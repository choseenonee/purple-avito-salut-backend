package service

import (
	"context"
	"template/internal/models"
	"time"
)

type Matrix interface {
	Create(ctx context.Context, matrix models.MatrixBase) (string, error)
	GetMatrix(ctx context.Context, matrixName string, page int) (models.Matrix, error)
	GetMatricesByDuration(ctx context.Context, timeStart, timeEnd time.Time) ([]models.Matrix, error)
	GetHistory(ctx context.Context, matrix models.GetHistoryMatrix) ([]models.ResponseHistoryMatrix, error)
	GetDifference(ctx context.Context, matrixName1, matrixName2 string) (models.MatrixDifference, error)
}

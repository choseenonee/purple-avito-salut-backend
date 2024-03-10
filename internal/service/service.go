package service

import (
	"context"
	"template/internal/models"
)

type Matrix interface {
	Create(ctx context.Context, matrix models.MatrixBase) (string, error)
	GetHistory(ctx context.Context, matrix models.GetHistoryMatrix) ([]models.ResponseHistoryMatrix, error)
	GetDifference(ctx context.Context, matrixName1, matrixName2 string) (models.MatrixDifference, error)
	GetTendency(ctx context.Context, data models.GetTendencyNode) ([]models.ResponseTendencyNode, error)
}

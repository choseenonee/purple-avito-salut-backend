package service

import (
	"context"
	"template/internal/models"
)

type Matrix interface {
	Create(ctx context.Context, matrix models.MatrixBase) (string, error)
	GetHistory(ctx context.Context, matrix models.GetHistoryMatrix) ([]models.Matrix, error)
}

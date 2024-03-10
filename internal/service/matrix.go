package service

import (
	"context"
	"template/internal/models"
	"template/internal/repository"
)

type matrixService struct {
	matrixRepo repository.Matrix
}

func (m matrixService) Create(ctx context.Context, matrix models.MatrixBase) (string, error) {
	// TODO: implement validation maybe???
	return m.matrixRepo.CreateMatrix(ctx, matrix)
}

func (m matrixService) GetHistory(ctx context.Context, data models.GetHistoryMatrix) ([]models.ResponseHistoryMatrix, error) {
	return m.matrixRepo.GetHistory(ctx, data)
}

func InitMatrixService(matrixRepo repository.Matrix) Matrix {
	return matrixService{matrixRepo: matrixRepo}
}

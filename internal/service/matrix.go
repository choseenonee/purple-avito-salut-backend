package service

import (
	"context"
	"template/internal/models"
	"template/internal/repository"
)

type matrixService struct {
	matrixRepo repository.Matrix
}

func InitMatrixService(matrixRepo repository.Matrix) Matrix {
	return matrixService{matrixRepo: matrixRepo}
}

func (m matrixService) Create(ctx context.Context, matrix models.MatrixBase) (string, error) {
	// TODO: implement validation maybe???
	return m.matrixRepo.CreateMatrix(ctx, matrix)
}

func (m matrixService) GetHistory(ctx context.Context, data models.GetHistoryMatrix) ([]models.Matrix, error) {
	return m.matrixRepo.GetHistory(ctx, data)
}

func (m matrixService) GetDifference(ctx context.Context, matrixName1, matrixName2 string) (models.MatrixDifference, error) {
	return m.matrixRepo.GetDifference(ctx, matrixName1, matrixName2)
}

func (m matrixService) GetMatrix(ctx context.Context, matrixName string, page int) (models.Matrix, error) {
	return m.matrixRepo.GetMatrix(ctx, matrixName, page)
}

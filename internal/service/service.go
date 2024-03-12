package service

import (
	"context"
	"errors"
	"template/internal/models"
	"template/internal/repository"
)

type Service interface {
	GetMicroCategoryPath(ctx context.Context, microCategoryID int) ([]int, error)
	GetRegionPath(ctx context.Context, microCategoryID int) ([]int, error)
	GetPrice(ctx context.Context, inData models.InData) (models.OutData, error)
}

type serviceStruct struct {
	repo                      repository.Repository
	currentBaseLineMatrixName string
}

func InitService(repo repository.Repository, currentBaseLineMatrixName string) Service {
	return serviceStruct{repo: repo, currentBaseLineMatrixName: currentBaseLineMatrixName}
}

func (s serviceStruct) GetPrice(ctx context.Context, inData models.InData) (models.OutData, error) {
	microCategoryPath, err := s.repo.GetMicroCategoryPath(ctx, inData.MicroCategoryID)
	if err != nil {
		return models.OutData{}, err
	}

	regionPath, err := s.repo.GetRegionPath(ctx, inData.RegionID)
	if err != nil {
		return models.OutData{}, err
	}

	// TODO: посмотреть в прыжках и обрезать path

	for _, regionID := range regionPath {
		for _, microcategoryID := range microCategoryPath {
			price, err := s.repo.GetPriceFromBaseLine(ctx, microcategoryID, regionID, s.currentBaseLineMatrixName)
			if err != nil {
				return models.OutData{}, err
			}
			if price != 0 {
				return models.OutData{
					MatrixName: s.currentBaseLineMatrixName,
					Price:      price,
					InData:     inData,
				}, nil
			}
		}
	}

	return models.OutData{}, errors.New("wtf how i did not find price?)))")
}

func (s serviceStruct) GetMicroCategoryPath(ctx context.Context, microCategoryID int) ([]int, error) {
	return s.repo.GetMicroCategoryPath(ctx, microCategoryID)
}

func (s serviceStruct) GetRegionPath(ctx context.Context, microCategoryID int) ([]int, error) {
	return s.repo.GetRegionPath(ctx, microCategoryID)
}

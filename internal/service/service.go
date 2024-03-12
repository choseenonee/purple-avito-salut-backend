package service

import (
	"context"
	"template/internal/repository"
)

type Service interface {
	GetMicroCategoryPath(ctx context.Context, microCategoryID int) ([]int, error)
}

type serviceStruct struct {
	repo repository.Repository
}

func InitService(repo repository.Repository) Service {
	return serviceStruct{repo: repo}
}

func (s serviceStruct) GetMicroCategoryPath(ctx context.Context, microCategoryID int) ([]int, error) {
	return s.repo.GetMicroCategoryPath(ctx, microCategoryID)
}

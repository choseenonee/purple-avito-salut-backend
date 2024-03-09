package repository

import (
	"context"
	"template/internal/model/entities"
)

type User interface {
	Create(ctx context.Context, userCreate entities.UserCreate) (int, error)
	Get(ctx context.Context, userID int) (entities.User, error)
	GetHashedPassword(ctx context.Context, name string) (int, string, error)
	Delete(ctx context.Context, userID int) error
	//	CREATE, READ, UPDATE, DELETE - the basics
}

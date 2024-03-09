package service

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"template/internal/model/entities"
)

type Public interface {
	CreateUser(ctx context.Context, userCreate entities.UserCreate) (string, string, error)
	LoginUser(ctx context.Context, userLogin entities.UserCreate) (string, string, error)
	Refresh(ctx context.Context, sessionID string, span trace.Span) (string, string, error)
}

type User interface {
	GetMe(ctx context.Context, userID int, span trace.Span) (entities.User, error)
	Delete(ctx context.Context, userID int, sessionID string) error
}
